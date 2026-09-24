# 宝藏文档 (Treasure Doc) — API 后端服务

> 基于 Gin 框架构建的高性能文档管理系统后端，支持文档管理、空间隔离、团队协作、版本历史等功能。
>
> Backend API service for Treasure Doc — a high-performance document management system built on Gin, featuring document CRUD, room-based isolation, team collaboration, and version history.

---

## 技术栈 Tech Stack

| 类别 | 技术 | 版本 |
| ------ | ------ | ------ |
| **语言** | Go | 1.26.8（`go.mod` 声明） |
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
  │     └── api.*Handler       ← Handler 层：参数解析 → 调用 Service → 组装响应
  │
  │   （跨域 CORS 由前置反向代理统一处理，Go 侧不设置 Access-Control-* 头）
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
| POST | `/api/user/change-pwd` | 修改当前用户密码（需原密码） | ✅ |

> 重置他人密码不提供 HTTP 接口，只能通过服务端二进制的 `resetpwd` 子命令执行。

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
| GET | `/` | 前端页面入口（静态资源同样挂在根路径，如 `/assets/...`） |

### Blog 公开接口

匿名只读接口挂载在 `/api/blog`，覆盖文章、日记、作品、工具、书签、个人资料、站点信息和内容统计。完整请求参数与响应契约见 [doc/blog-api.md](doc/blog-api.md)，后端模块说明见 [module/blog/README.md](module/blog/README.md)。

Blog 管理接口挂载在 `/api/blog-mgr`，仅 admin/root 可访问，提供分类、标签、全部内容资源以及 Profile、Site 的管理能力。管理接口说明见 [module/blog_mgr/README.md](module/blog_mgr/README.md)。

---

## 快速开始 Quick Start

### 环境要求 Requirements

- **Go** 1.26+（`go.mod` 声明 go 1.26.8，与 Dockerfile 镜像一致）
- **SQLite**（默认，零依赖单文件，开箱即用）或 **PostgreSQL** 12+（可选）
- **Redis** (可选，用于缓存/验证码)

### 配置文件 Configuration

```bash
cp config.example.toml config.toml
```

主要配置项：

```toml
[app]
port = 2026
runMode = "dev"       # dev-开发模式 release-生产模式
trustedProxies = []   # 信任的反代 IP/CIDR，命中时才采信 X-Forwarded-For（空=不信任代理头）

[database]
# driver：sqlite（默认，零依赖单文件）或 postgres
driver = "sqlite"
dsn = "file:./treasure_doc.db?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on"
tablePrefix = "td_"

# PostgreSQL 示例：将 driver 改为 postgres 并填写下面 dsn
# driver = "postgres"
# dsn = "host=127.0.0.1 user=postgres password=postgres dbname=treasure_doc port=5432 sslmode=disable TimeZone=Asia/Shanghai"

[backup]
# SQLite 定时备份（仅 driver=sqlite 生效；postgres 不启用）
enable = false
interval = 86400 # 备份周期（秒），86400 = 每天
dir = "backup"   # 备份文件存放目录
compress = true  # 是否 gzip 压缩（.db.gz）
keepDays = 7     # 旧备份保留天数，0 = 不清理
autoPack = false # 定时备份是否同时生成「完整备份包」（含图片）
apiToken = ""    # NAS 免登录拉取完整备份的机器令牌（X-Backup-Token）
allowIPs = []    # 备份接口来源 IP 白名单（精确 IP 或 CIDR，空=不限）

[publish]
allowIPs = []    # 发布接口来源 IP 白名单（精确 IP 或 CIDR，空=不限）

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
2. 注册默认管理员账号：`treasuredocmgr / treasuredocmgr`（首次运行，请尽快修改密码）
3. 服务监听 `:2026`

> AutoMigrate 仅用于在空数据库中初始化或调整表结构，不会迁移已有数据；切换数据库驱动需重启服务。

> ⚠️ **首次登录后请立即修改 root 密码。**

### 备份与恢复 Backup & Restore

SQLite 场景下，数据 = **数据库文件 + 上传目录（`files/`，含图片/视频）**。系统提供两种备份形态：

| 形态 | 内容 | 触发方式 |
| --- | --- | --- |
| 数据库快照（`.db` / `.db.gz`） | 仅数据库 | 定时任务或 `POST /api/blog-mgr/backups` |
| 完整备份包（`.tar.gz`） | `treasure_doc.db` + `files/` 图片 + `manifest.json` | 后台「导出完整备份」/ NAS 接口 |

完整备份包是自包含归档，可跨环境恢复。

**后台管理**（`[backup]` 相关接口均在 `config.example.toml` 的 `[backup]` 段配置）：

- `GET /api/blog-mgr/backups` — 备份列表
- `POST /api/blog-mgr/backups` — 立即生成数据库快照
- `GET /api/blog-mgr/backups/:name` — 下载指定备份
- `POST /api/blog-mgr/backups/export` — 导出「完整备份包」（含图片，落盘到备份目录）
- `POST /api/blog-mgr/backups/restore` — 导入并恢复（上传 `.tar.gz`；先校验完整性/版本再替换数据库并补齐图片，仅管理员可操作，导入前建议先开站点维护模式）

**NAS 定时拉取**（机器令牌，免登录、只读）：

```bash
# 每日拉取最新完整备份到 NAS
30 2 * * * curl -fsSL -H "X-Backup-Token: <你的令牌>" \
  -o "$HOME/backup/treasure-$(date +\%Y\%m\%d).tar.gz" \
  "https://<你的域名>/api/backup/export"
```

要点：
- NAS 令牌在 `[backup].apiToken` 配置；该接口仅开放「下载」，不提供导入/删除。
- 想让定时任务/NAS 得到**含图片**的完整包，需同时开启 `[backup].autoPack = true`（定时任务产出完整备份包）。
- 恢复会把「当前库」替换成备份内容，属于破坏性操作：建议恢复前先对现有库导出一次留底。
- 可用 `[backup].allowIPs` 限制来源 IP（精确 IP 或 CIDR，空=不限）。

### 内容发布接口 Publish API

供脚本 / CI / NAS 用机器令牌**免登录直接发布内容**（文章/日记/作品/工具/收藏集），与后台手动发布走同一套校验与默认分类逻辑。鉴权：

- 令牌为 `X-Publish-Token`，在后台“备份”页或 `PUT /api/blog-mgr/publish/token` 设置（写入数据库，即时生效，留空=关闭接口）。
- 可用 `[publish].allowIPs` 限制来源 IP；接口自带**每 IP 限流**（`/api/publish`）。
- **反代下注意**：IP 白名单/限流按真实客户端 IP 判定。部署在 nginx 后需在 `[app].trustedProxies` 填入反代地址（如 `["127.0.0.1"]`），后端才会采信 `X-Forwarded-For` 拿到真实 IP，否则取到的是反代出口 IP。此项为启动期配置，改动需重启。

发布语义：`publishStatus` 恒为 `published`（入参忽略）；未指定分类自动落到“默认分类”；`slug`/`publicId` 未填自动按标题生成。

**调用示例**

```bash
# 1) 配置发布令牌
curl -sSL -X PUT http://<域名>/api/blog-mgr/publish/token \
  -H "X-Token: <管理员token>" -H "Content-Type: application/json" \
  -d '{"token":"your-publish-secret"}'

# 2) 直接发布一篇已发布状态的文章（未给分类 → 默认分类；未给 slug → 按标题生成）
curl -sSL -X POST http://<域名>/api/publish/posts \
  -H "X-Publish-Token: your-publish-secret" -H "Content-Type: application/json" \
  -d '{"title":"Hello API 文章","content":"正文…"}'

# 3) 其他资源：/api/publish/diaries | portfolio-items | tools | bookmarks
```

---

## 项目结构 Project Structure

```
treasure-doc/
│
├── module/                          # 业务模块
│   ├── user/                        # 核心入口：用户/鉴权/上传/备份、路由汇总与前端托管
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
│       │       ├── admin.go         # 管理员鉴权
│       │       └── ...               # 其它中间件（CORS 由反代处理，不在 Go 侧）
│       │
│       ├── utils/                   # 工具函数
│       │   ├── datetime.go / file.go / slice.go / user.go
│       │
│       ├── cli/                     # 仅含 reset-pwd/README.md；密码重置用主程序 resetpwd 子命令
│       │
│       ├── web/                     # 前端静态文件
│       └── files/                   # 用户上传文件
│
│   ├── blog/                        # 公开只读博客 API（含 seed 演示数据）
│   ├── blog_mgr/                    # 后台管理通用 CRUD（由 user 模块路由挂载）
│   └── common/                      # 统一响应 {code,msg,data}
│
├── doc/                             # 设计文档
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
  -p 2026:2026 \
  -v /path/to/web:/app/web \
  -v /path/to/files:/app/files \
  -v /path/to/config.toml:/app/config.toml \
  treasure-doc

# 调试模式
docker run --rm --name treasure-doc -it \
  -p 2026:2026 \
  -v /path/to/web:/app/web \
  -v /path/to/files:/app/files \
  -v /path/to/config.toml:/app/config.toml \
  treasure-doc /bin/sh
```

---

## 开发工具 Development Tools

### CLI 密码重置

重置默认管理员账号（`treasuredocmgr`）的密码，新密码需为 8-16 位。必须从 `module/user` 运行（`FilesPath`/`WebPath` 为相对路径）：

```bash
cd module/user
go run . -c config.toml resetpwd <新密码>
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

> 以下为仍需改进的已知问题：

1. **缺少 DAO 层** — Service 直接操作 `global.Db`，查询逻辑分散。建议抽取 `data/repository/` 封装。另 `module/user/api/file_api.go`/`media_api.go` 在 Handler 内直接查库，属遗留例外，新代码不要效仿。
2. **工作目录约束** — `FilesPath`/`WebPath` 为相对路径，服务与 resetpwd 都必须从 `module/user` 运行，从仓库根启动会解析错路径。
3. **GORM SQL 日志** — 当前 GORM Logger 使用 Silent 级别，不打印 SQL。排查时可临时切换为 Info：

   ```go
   &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
   ```

4. **测试覆盖** — 核心 Service 缺少单元测试。
5. **用户管理占位代码** — `module/admin/` 并不存在；`api/user_manage_api.go` 存在但未注册任何路由，其 DTO 内嵌的 `request.Sort` 也因此暂无调用方；实际的后台管理能力在 `module/blog_mgr`。

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
