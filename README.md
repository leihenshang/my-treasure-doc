# 宝藏文档 (Treasure Doc) — API 后端服务

> 基于 Gin 框架构建的高性能文档管理系统后端，支持文档管理、空间隔离、团队协作、版本历史等功能。
>
> Backend API service for Treasure Doc — a high-performance document management system built on Gin, featuring document CRUD, room-based isolation, team collaboration, and version history.

---

## 技术栈 Tech Stack

| 类别 | 技术 | 版本 |
| ------ | ------ | ------ |
| **语言** | Go | 1.22+ |
| **Web 框架** | Gin | v1.9.1 |
| **ORM** | GORM | v1.24 |
| **数据库** | SQLite（默认，零依赖）/ PostgreSQL（可选） | SQLite 单文件 / PG 12+ |
| **缓存** | Redis | 可选 |
| **配置管理** | Viper + TOML | v1.8.1 |
| **日志** | Zap + Lumberjack | v1.17.0 |
| **ID 生成** | Sonyflake | v1.2.0 |
| **密码加密** | golang.org/x/crypto | v0.23.0 |
| **验证码** | base64Captcha | v1.3.6 |
| **参数校验** | go-playground/validator | v10.14.0 |
| **部署** | Docker 多阶段构建 | Alpine 3.20 |

---

## 架构分层 Architecture

```
main.go
  │
  ├─ global.InitModule()        ← 统一初始化 Config → Logger → Redis → 数据库（按 driver 选择 SQLite/PostgreSQL）→ Validator
  │                               (仅注册开关支持热更新，基础设施配置变更需重启)
  │
  ├─ router.InitRouter(r)       ← 路由注册 + 中间件
  │     ├── middleware.Auth()   ← X-Token 认证（支持 Mock 开发模式）
  │     ├── middleware.Cors()   ← 跨域
  │     └── api.*Handler       ← Handler 层：参数解析 → 调用 Service → 组装响应
  │
  ├─ api/                       ← Handler 层（请求绑定、响应格式化）
  │     └── internal/service/   ← Service 层（业务逻辑、DB 操作）
  │
  └─ data/model/                ← DO 层（GORM 模型 + BeforeCreate 钩子）
```

---

## API 概览（≈30 个端点）

所有 API 挂载在 `/api` 下，认证统一通过 Header `X-Token` 传递。

### 用户模块

| 方法 | 路径 | 说明 | 认证 |
| ------ | ------ | ------ | ------ |
| GET | `/api/user/captcha` | 获取图形验证码（Base64 PNG） | ❌ |
| POST | `/api/user/reg` | 注册（自动创建默认空间） | ❌ |
| POST | `/api/user/login` | 登录（返回 token，7 天有效；需图形验证码） | ❌ |
| POST | `/api/user/logout` | 退出登录 | ✅ |
| POST | `/api/user/update-profile` | 更新个人资料 | ✅ |
| POST | `/api/user-manage/create` | 创建用户（管理） | ✅ |
| GET | `/api/user-manage/detail` | 用户详情（管理） | ✅ |
| GET | `/api/user-manage/list` | 用户列表（管理） | ✅ |
| POST | `/api/user-manage/update` | 更新用户（管理） | ✅ |
| POST | `/api/user-manage/delete` | 删除用户（管理） | ✅ |
| POST | `/api/user-manage/reset-pwd` | 重置密码（管理） | ✅ |

### 文档模块

| 方法 | 路径 | 说明 | 认证 |
| ------ | ------ | ------ | ------ |
| POST | `/api/doc/create` | 创建文档 | ✅ |
| GET | `/api/doc/detail` | 文档详情 | ✅ |
| GET | `/api/doc/list` | 文档列表 | ✅ |
| POST | `/api/doc/update` | 更新文档 | ✅ |
| POST | `/api/doc/delete` | 删除文档（进回收站） | ✅ |
| POST | `/api/doc/recover` | 恢复文档 | ✅ |
| GET | `/api/doc-history/detail` | 历史版本详情 | ✅ |
| GET | `/api/doc-history/list` | 历史版本列表 | ✅ |
| POST | `/api/doc-history/recover` | 恢复历史版本 | ✅ |

### 文档分组

| 方法 | 路径 | 说明 | 认证 |
| ------ | ------ | ------ | ------ |
| POST | `/api/doc-group/create` | 创建分组 | ✅ |
| GET | `/api/doc-group/list` | 分组列表 | ✅ |
| GET | `/api/doc-group/tree` | 分组树形结构 | ✅ |
| POST | `/api/doc-group/update` | 重命名/移动分组 | ✅ |
| POST | `/api/doc-group/delete` | 删除分组 | ✅ |
| GET | `/api/doc-group/detail` | 分组详情 | ✅ |

### 笔记 & 文件 & 空间

| 方法 | 路径 | 说明 | 认证 |
| ------ | ------ | ------ | ------ |
| POST | `/api/note/create` | 创建笔记 | ✅ |
| GET | `/api/note/detail` | 笔记详情 | ✅ |
| GET | `/api/note/list` | 笔记列表 | ✅ |
| POST | `/api/note/update` | 更新笔记 | ✅ |
| POST | `/api/note/delete` | 删除笔记 | ✅ |
| POST | `/api/file/upload` | 文件上传 | ✅ |
| POST | `/api/room/create` | 创建空间 | ✅ |
| GET | `/api/room/detail` | 空间详情 | ✅ |
| GET | `/api/room/list` | 空间列表 | ✅ |
| POST | `/api/room/update` | 更新空间 | ✅ |
| POST | `/api/room/delete` | 删除空间 | ✅ |

### 健康检查

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ping` | 返回 `{"msg": "pong!"}` |
| GET | `/` | 重定向到 `/web` 前端页面 |

### Blog 公开接口

匿名只读接口挂载在 `/api/blog`，覆盖文章、日记、作品、工具、书签、个人资料、站点信息和内容统计。完整请求参数与响应契约见 [doc/blog-api.md](doc/blog-api.md)，后端模块说明见 [module/blog/README.md](module/blog/README.md)。

Blog 管理接口挂载在 `/api/blog-mgr`，仅 admin/root 可访问，提供分类、标签、全部内容资源以及 Profile、Site 的管理能力。管理接口说明见 [module/blog_mgr/README.md](module/blog_mgr/README.md)。

---

## 快速开始 Quick Start

### 环境要求 Requirements

- **Go** 1.22+
- **SQLite**（默认，零依赖单文件，开箱即用）或 **PostgreSQL** 12+（可选）
- **Redis** (可选，用于缓存/验证码)

### 配置文件 Configuration

```bash
cp config.example.toml config.toml
```

主要配置项：

```toml
[app]
port = 2021
runMode = "dev"       # dev-开发模式 release-生产模式

[database]
# driver：sqlite（默认，零依赖单文件）或 postgres
driver = "sqlite"
dsn = "file:./treasure_doc.db?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on"
tablePrefix = "td_"

# PostgreSQL 示例：将 driver 改为 postgres 并填写下面 dsn
# driver = "postgres"
# dsn = "host=127.0.0.1 user=postgres password=postgres dbname=treasure_doc port=5432 sslmode=disable TimeZone=Asia/Shanghai"

[redis]
enable = false        # 不启用 Redis 可跳过
host = "127.0.0.1"
port = 6379

[log]
level = "info"
```

### 本地运行 Local Development

```bash
# 进入用户模块
cd module/user

# 运行服务（默认加载 ./config.toml）
go run . -c config.toml

# 指定配置文件路径
go run . -c /path/to/config.toml
```

启动后自动完成：

1. 初始化数据库连接，自动建表（GORM AutoMigrate）
2. 注册 root 账号：`treasure-root / treasure-root`（首次运行）
3. 服务监听 `:2021`

> AutoMigrate 仅用于在空数据库中初始化或调整表结构，不会迁移已有数据；切换数据库驱动需重启服务。

> ⚠️ **首次登录后请立即修改 root 密码。**

---

## 项目结构 Project Structure

```
treasure-doc/
│
├── module/                          # 业务模块
│   └── user/                        # 核心用户文档模块
│       ├── main.go                  # 程序入口
│       ├── config.example.toml      # 配置示例
│       ├── config.toml              # 运行时配置（gitignore 建议）
│       │
│       ├── api/                     # Handler 层：请求解析 & 响应组装
│       │   ├── user_api.go
│       │   ├── user_manage_api.go
│       │   ├── doc_api.go
│       │   ├── doc_group_api.go
│       │   ├── doc_history_api.go
│       │   ├── note_api.go
│       │   ├── file_api.go
│       │   └── room_api.go
│       │
│       ├── config/                  # Config 结构体 + Viper 封装
│       │   ├── config.go            # 配置加载与受限热更新
│       │   ├── app.go / database.go / redis.go / log.go / debug.go
│       │
│       ├── data/                    # 数据层
│       │   ├── model/               # GORM 模型（DO）
│       │   │   ├── base.go          # BaseModel + 雪花 ID 自动生成
│       │   │   ├── user.go / user_token.go / user_conf.go
│       │   │   ├── doc.go / doc_group.go / doc_history.go
│       │   │   ├── note.go / room.go
│       │   │   ├── team.go / team_user.go
│       │   │   ├── global_conf.go / verify_code.go
│       │   ├── request/             # 请求 DTO
│       │   └── response/            # 响应 DTO + 错误码
│       │
│       ├── global/                  # 全局单例 & 初始/销毁
│       │   ├── global.go            # InitModule / 业务配置热更新策略
│       │   ├── constant.go
│       │   ├── db.go                # 数据库初始化（按 driver 分支 SQLite/PostgreSQL）+ 优雅关闭
│       │   ├── logger.go            # Zap 初始化
│       │   ├── trans.go             # 校验器翻译
│       │   └── gid/                 # Sonyflake ID 生成
│       │
│       ├── internal/                # 内部逻辑
│       │   ├── service/             # Service 层（业务逻辑）
│       │   │   ├── user_service.go / doc_service.go / ...
│       │   │   ├── room_service.go / team_service.go
│       │   │   ├── captcha_service.go
│       │   │   └── doc_history_service.go
│       │   └── auth/                # 认证逻辑
│       │
│       ├── router/                  # 路由 & 中间件
│       │   ├── router.go
│       │   └── middleware/
│       │       ├── auth.go          # X-Token 认证（含 Mock 模式）
│       │       └── cors.go          # 跨域
│       │
│       ├── utils/                   # 工具函数
│       │   ├── datetime.go / file.go / slice.go / user.go
│       │
│       ├── cli/                     # CLI 工具
│       │   ├── cli.go               # -gen 生成模型
│       │   └── reset-pwd/           # 密码重置工具
│       │
│       ├── web/                     # 前端静态文件
│       └── files/                   # 用户上传文件
│
├── module/admin/                    # 管理后台（开发中）
├── list_sort/                       # 列表排序算法（独立工具包）
├── log/                             # 运行时日志
├── doc/                             # 设计文档
├── template/                        # 模板文件
│
├── Dockerfile                       # 多阶段 Docker 构建
├── build.sh                         # Docker 构建脚本
├── go.mod / go.sum
└── README.md
```

---

## 数据模型 Database Models

所有模型嵌入 `BaseModel`，使用 **Sonyflake 雪花 ID**（字符串类型，19 位数字），支持软删除。

| 模型 | 对应表 | 关键字段 | 说明 |
| ------ | -------- | ---------- | ------ |
| `User` | `td_user` | Account, Password, UserType, UserStatus, CurrentRoomId | 支持 root/普通用户 |
| `UserToken` | `td_user_token` | Token, TokenExpire, LoginIp, LoginTime | 每用户最多 3 个并发会话 |
| `Doc` | `td_doc` | Title, Content, GroupId, RoomId, UserId | 软删除 → 回收站 |
| `DocGroup` | `td_doc_group` | Name, PId, RoomId | 树形分组 |
| `DocHistory` | `td_doc_history` | DocId, Content, Version | 版本快照，支持恢复 |
| `Note` | `td_note` | Title, Content, RoomId | 轻量笔记 |
| `Room` | `td_room` | Name, UserId, IsDefault | 空间/文档隔离单元 |
| `Team` / `TeamUser` | `td_team` / `td_team_user` | — | 团队及成员关系 |
| `GlobalConf` | `td_global_conf` | Key, Value | 系统级 KV |
| `UserConf` | `td_user_conf` | Key, Value, UserId | 用户级 KV |
| `VerifyCode` | `td_verify_code` | Code, Type, Target, ExpireAt | 验证码（待流程接入） |

---

## 部署部署 Build & Deployment

### 跨平台编译 Cross-platform Build

```bash
# Linux
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o treasure_user ./module/user

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o treasure_user.exe ./module/user

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o treasure_user ./module/user
```

### Docker 部署

```bash
# 构建镜像
docker build -t treasure-doc .

# 后台运行
docker run -d --name treasure-doc \
  --restart=always \
  -p 2021:2021 \
  -v /path/to/web:/app/web \
  -v /path/to/files:/app/files \
  -v /path/to/config.toml:/app/config.toml \
  treasure-doc

# 调试模式
docker run --rm --name treasure-doc -it \
  -p 2021:2021 \
  -v /path/to/web:/app/web \
  -v /path/to/files:/app/files \
  -v /path/to/config.toml:/app/config.toml \
  treasure-doc /bin/sh
```

---

## 开发工具 Development Tools

### CLI 数据库模型生成

```bash
cd module/user/cli
go run . -gen
```

自动读取数据库表结构，生成 GORM 模型到 `data/model/`。

### CLI 密码重置

```bash
go run ./module/user/cli/reset-pwd -u <账号> -p <新密码> -cfg <config.toml 绝对路径>
```

---

## 关键设计 Key Design Decisions

| 决策 | 实现 |
| ------ | ------ |
| **ID 生成** | Sonyflake → 19 位数字字符串，比 UUID 更短，分布式友好 |
| **密码加密** | `golang.org/x/crypto/bcrypt` |
| **会话管理** | 每用户最多 3 个 token，超限自动剔除最早的会话 |
| **空间隔离** | 注册时自动创建默认 Room，文档按 Room 隔离（多租户基础） |
| **配置热更新** | 仅 `app.registerEnabled` 可热更新；数据库（driver/dsn）、Redis、日志、端口、运行模式和 Debug 变更需重启 |
| **Mock 认证** | 仅 dev 模式可通过 `debug.enableMockLogin` 注入固定 root，修改后需重启 |
| **Blog Seed** | 默认关闭；迁移后幂等填充公开演示数据，远程库需显式 `allowRemote=true` |

---

## 已知问题 & 待改进 Known Issues

> 以下对应 `question.md` 中的记录：

1. **缺少 DAO 层** — Service 直接操作 `global.Db`，查询逻辑分散。建议抽取 `data/repository/` 封装。
2. **CLI 配置路径** — cli 子工具应通过 `-cfg` 参数传递配置文件绝对路径。
3. **GORM SQL 日志** — 当前 GORM Logger 使用 Silent 级别，不打印 SQL。排查时可临时切换为 Info：

   ```go
   &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
   ```

4. **测试覆盖** — 仅在 `list_sort/` 和 `gid/` 有测试，核心 Service 缺少单元测试。
5. **Admin 模块** — `module/admin/` 仅占位，用户管理功能暂在 user 模块 `user-manage` 路由下。

---

## 数据库维护 Database Maintenance

```sql
-- 修复文档分组关联
UPDATE td_doc SET group_id = 'root' WHERE group_id = '' OR group_id = '0';

-- 修复文档组父子关系
UPDATE td_doc_group SET p_id = 'root' WHERE p_id = '' OR p_id = '0';
```

---

## 许可证 License

MIT License
