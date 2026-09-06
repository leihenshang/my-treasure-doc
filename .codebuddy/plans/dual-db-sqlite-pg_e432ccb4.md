---
name: dual-db-sqlite-pg
overview: 让服务同时支持 SQLite（默认）和 PostgreSQL，通过配置选择驱动。核心是新增纯 Go 的 SQLite 驱动、把 Database 配置重构为 Driver+Dsn、把 PG 专属的 jsonb/timestamptz 字段类型改为跨方言写法，并同步文档。现有 PG 部署需手动更新 config.toml（按用户选择允许重构）。
todos:
  - id: deps-config
    content: 新增 glebarez/sqlite 依赖，重构 Database 为 Driver+Dsn 并更新 config.example.toml
    status: completed
  - id: db-open
    content: 改造 global/db.go 按 driver 分支打开 sqlite/postgres，设置 sqlite 连接池与 WAL/busy_timeout
    status: completed
    dependencies:
      - deps-config
  - id: model-types
    content: 将 jsonb 与 timestamptz 改为跨方言写法（datatypes.JSON / timestamp），调整自定义 JSON 类型
    status: completed
    dependencies:
      - deps-config
  - id: docs
    content: 同步更新 README、各模块 README、doc/blog-api.md、AGENTS.md、reset-pwd/README 的数据库说明
    status: completed
    dependencies:
      - db-open
      - model-types
  - id: verify
    content: CGO_ENABLED=0 构建、go test ./...、本地 sqlite 起服验证并补充集成冒烟测试
    status: completed
    dependencies:
      - db-open
      - model-types
      - docs
---

## 用户需求

让 treasure-doc 服务同时支持 SQLite 与 PostgreSQL 两种数据库，通过配置项选择驱动；SQLite 为默认驱动，PostgreSQL 仍保留支持。

## 产品概述

在不改变现有业务功能与 API 行为的前提下，将数据库层从"强绑定 PostgreSQL"改造为"配置驱动的跨数据库"架构。新部署开箱即用（SQLite 零依赖、单文件），原有 PostgreSQL 部署只需手动更新 config.toml 即可继续运行。

## 核心特性

- 配置新增 `driver`（sqlite/postgres）+ `dsn` 字段，按驱动打开不同 GORM 方言。
- SQLite 默认启用 WAL、busy_timeout 与单写者连接池，避免 `database is locked`。
- 将 PG 专属字段类型（jsonb / timestamptz）改为跨方言写法，确保同一套模型在两种库均可 AutoMigrate 与读写。
- 文档与示例配置同步，说明两种库的配置方式与"切换驱动需重启"的约束。

## 技术栈

- 语言/框架：Go 1.22 + Gin + GORM（沿用现状）
- SQLite 驱动：`github.com/glebarez/sqlite`（基于 modernc.org/sqlite，纯 Go、零 CGO）。必须用纯 Go 驱动，因为 `Dockerfile:20` 为 `CGO_ENABLED=0 go build`，标准 `gorm.io/driver/sqlite`（mattn/go-sqlite3）需 CGO 无法编译。
- PostgreSQL 驱动：保留 `gorm.io/driver/postgres`（PG 仍需）。
- JSON 类型：`gorm.io/datatypes` 的 `datatypes.JSON`（自动按方言映射 jsonb/json）。

## 实现方案

### 总体策略

在 `global/db.go` 的 `openDatabaseWithConfig` 中按 `cfg.Database.Driver` 分支构造 GORM Dialector；配置结构从 PG 专属字段重构为通用的 `Driver + Dsn`；模型层把 PG 专属列类型改为跨方言声明。其余业务代码（Service/API/事务/乐观锁/软删除）全部走 GORM 泛型能力，无需改动。

### 关键技术决策与权衡

1. **纯 Go SQLite 驱动（glebarez）**：受 Dockerfile CGO_ENABLED=0 约束，是唯一可行选择；现代 c 版本对常见 SQL 支持完整，满足博客低并发读写。
2. **配置统一为 Driver + Dsn**：用户已确认允许重构。PG 的 Dsn 直接复用原拼接串（`host=... user=... dbname=... sslmode=... TimeZone=...`），不再在代码里拆分字段；sqlite 的 Dsn 为 `file:./treasure_doc.db?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on`。代价是现有 PG 的 config.toml 需手动改为新格式（已在文档注明）。
3. **自定义 JSON 类型改为 `datatypes.JSON`**：`module/blog/data/model/model.go` 当前 `type JSON []byte`（自带 Value/Scan）+ `type:jsonb` 标签。SQLite 无 jsonb 列类型，AutoMigrate 会建 NUMERIC 亲和列。方案：将 `JSON` 改为 `type JSON = datatypes.JSON`（类型别名），删除自带的 Value/Scan 方法，字段标签去掉 `type:jsonb`（让 GORM 自动选择 jsonb/json）。调用方 `blog_mgr/.../service.go:172 marshal`、`blog/seed/seed.go:274 toJSON`、`blog/internal/service/service.go:231 decodeJSON` 因别名不变，零改动。
4. **timestamptz 改 timestamp**：`base.go`、`user_token.go`、`model.go` 全量 `type:timestamptz` → `type:timestamp`，两种方言均原生支持，避免 SQLite NUMERIC 亲和导致的时区/排序隐患。
5. **SQLite 并发**：单写者模型，取 `sqlDb.SetMaxOpenConns(1)`（写少读多场景足够），DSN 已含 `_busy_timeout=5000` 防锁；不引入复杂连接池，控制改动面。

### 性能与可靠性

- 博客为低并发只读为主，SQLite WAL + 单写连接足够；无 N+1/额外遍历风险。
- 保留 `DisableForeignKeyConstraintWhenMigrating: true`，AutoMigrate 两种库行为一致。
- 非法 `driver` 在 `openDatabaseWithConfig` 返回明确错误，启动即失败而非运行时panic。
- `dsn` 为空时启动报错，避免静默用空库。

### 避免技术债

沿用现有 `global.Db`、`TableMigrate`、`AutoMigrate` 流程，不新增 DAO/Repository；仅替换 Dialector 构造与列类型声明，符合 AGENTS.md 的分层与命名约定。

## 实现要点（防回归）

- 保留 PG 分支逻辑与原 DSN 内容，仅把字段拼接改为直接读取 `cfg.Database.Dsn`，确保 PG 行为与现状一致。
- `jsonb` 标签移除后，PG 仍能正常使用 jsonb（datatypes.JSON 在 pg 下即 jsonb），博客读写/Seed 不受影响。
- 配置热更新现状：仅 registerEnabled 生效，database 运行期变更本就忽略且需重启；切换驱动同样需重启，符合现状，仅在文档补充一句。
- 不改 `list_sort`、`gid`(sonyflake)、captcha(mem store)、Redis 逻辑。

## 架构设计

数据流不变：`main.go → global.InitModule() → 读取 config → openDatabaseWithConfig(按 driver 选择 Dialector) → AutoMigrate → router → api → service → global.Db → model`。仅在"打开数据库"环节引入 driver 分支，下游完全无感。

## 目录结构与改动文件

```
go.mod                                       # [MODIFY] 新增 github.com/glebarez/sqlite 依赖，保留 gorm.io/driver/postgres
module/user/config/database.go              # [MODIFY] 结构重构为 Driver + Dsn（+ 可选 TablePrefix），Driver 默认 sqlite
module/user/config.example.toml             # [MODIFY] 默认 [database] driver="sqlite" dsn="file:./treasure_doc.db?..."；附 PG 示例注释
module/user/global/db.go                    # [MODIFY] openDatabaseWithConfig 按 driver 分支；sqlite 设置 SetMaxOpenConns(1)；非法 driver/空 dsn 报错
module/blog/data/model/model.go            # [MODIFY] JSON 改为 datatypes.JSON 别名、删自定义 Value/Scan、7 处 jsonb 去类型、timestamptz→timestamp、import datatypes
module/user/data/model/base.go              # [MODIFY] 3 处 timestamptz→timestamp
module/user/data/model/user_token.go        # [MODIFY] 3 处 timestamptz→timestamp
README.md                                    # [MODIFY] 数据库章节改为"SQLite 默认 / 可选 PG"，更新示例与依赖说明
module/user/README.md                        # [MODIFY] 同步数据库初始化说明与示例
module/blog/README.md                        # [MODIFY] JSONB 描述改为跨方言 JSON，更新启动配置示例
doc/blog-api.md                              # [MODIFY] PostgreSQL/JSONB 相关表述（如 1600 行）改为支持 sqlite/pg
AGENTS.md                                    # [MODIFY] 技术栈(第5行)、热更新(第48行)、测试(第63行) 的数据库说明
module/user/cli/reset-pwd/README.md          # [MODIFY] 说明新 config [database] 写法
module/user/global/db_sqlite_test.go         # [NEW] 可选：基于真实 sqlite 文件的集成冒烟测试（AutoMigrate + 一条 JSON 字段写读 + 乐观锁）
```

## 关键代码结构

```
// module/user/config/database.go
type Database struct {
    Driver      string // "sqlite" (默认) | "postgres"
    Dsn         string // sqlite: file:./x.db?... ; postgres: host=... user=... dbname=... sslmode=... TimeZone=...
    TablePrefix string
}

// module/blog/data/model/model.go
type JSON = datatypes.JSON // 替代原 []byte 自定义类型，复用 GORM 跨方言序列化
```