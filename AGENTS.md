# AGENTS.md

本文件为 AI 编码代理在本仓库工作时提供指引。

## 项目范围

- 本仓库是 Go 1.26 单模块项目（`fastduck/treasure-doc`，`go.mod` 声明 `go 1.26.8`，与 Dockerfile 的 `golang:1.26.8-alpine3.23` 一致），核心服务位于 `module/user`，技术栈为 Gin、GORM、SQLite（默认）/ PostgreSQL（可选）；Redis 可选。
- 先阅读根目录 [README.md](README.md) 了解产品与 API 概览，公开接口明细见 [doc/blog-api.md](doc/blog-api.md)，跨域反代示例见 [doc/nginx-cors.example.conf](doc/nginx-cors.example.conf)，部署目录和历史数据修复见 [module/user/README.md](module/user/README.md)。设计文档描述的是目标状态，实际行为以代码和 `module/user/router/router.go` 为准。
- 代码与文档主要使用中文。保持现有命名、分层和错误响应风格，不做与任务无关的架构重构。

## 常用命令

在仓库根目录运行：

```bash
go test ./...
go test ./module/user/global/gid/...
go fmt ./...
go build -o treasure_user.exe ./module/user
docker build -t treasure-doc .
```

本地启动必须从 `module/user` 运行，因为配置和静态目录使用相对路径：

```bash
cd module/user
go run . -c config.toml
```

- 首次运行前按需从 `config.example.toml` 创建本地 `config.toml`；默认 SQLite 零依赖开箱即用，也可在 `[database]` 改为 PostgreSQL。
- 启动服务会连接数据库、执行 `AutoMigrate`，并尝试注册默认 root 用户；不要把启动服务当作无副作用的验证步骤。
- 重置密码是主程序的子命令（见 `module/user/main.go` 的 `runResetPwd`）：`cd module/user && go run . -c config.toml resetpwd <新密码>`，仅重置默认管理员账号，新密码须满足 8–16 位规则。`module/user/cli/reset-pwd/` 目录下只有 README，没有可执行代码。
- 仓库没有 CI、Makefile 或 lint 配置。门禁检查（`go test ./...`、`gofmt` / `go fmt`）**仅在重要改动时执行**：新增大模块、大块业务逻辑调整、重构或大规模跨文件改动必须跑通；单文件小改（样式微调、文案/提示语、参数调整、单行 bug 修复）可直接提交。

## 代码边界

仓库含四个业务模块，共用同一 Gin Engine：`module/user`（进程入口、用户/鉴权/上传/备份、路由汇总与前端托管）、`module/blog`（公开只读博客 API）、`module/blog_mgr`（后台管理 CRUD）、`module/common`（统一响应）。

请求链路为 `main.go` → `global.InitModule()` → `router.InitRouter()` → middleware → `api` → `internal/service` → `global.Db` → `data/model`。`module/blog` 与 `module/blog_mgr` 沿用同样的分层，各自持有 `api / internal/service / data / router`，由 `module/user/router/router.go` 统一挂载。

各层职责：

- `api/`：绑定请求、从 Gin context 获取当前用户、调用 Service，并使用 `data/response` 输出；不要在 Handler 中新增数据库查询。
- `internal/service/`：业务规则、事务和 GORM 查询。项目没有 DAO/Repository 层；除非任务明确要求架构调整，否则沿用这一结构。
- `data/request/` 与 `data/response/`：请求 DTO、分页/排序参数、响应 DTO 和业务错误码。
- `data/model/`：GORM 模型、表名、软删除和创建钩子。新增模型通常嵌入 `BaseModel` 并实现 `TableName()`。
- `router/`：实际生效的路由清单。存在 API 或 Service 文件不代表端点已经暴露。
- `global/`（仅 user 模块）：配置、数据库、Redis、日志、迁移和 validator 等进程级状态。

## 已有功能实现地图（功能 → 代码位置）

- **启动与初始化**：`module/user/main.go` → `global.InitModule()`（`global/global.go`，顺序为配置 → 日志 → 可选 Redis → 数据库 → 可选 SQLite 定时备份调度器 → validator（`InitTrans`）→ 配置热更新监听 → `migrateDbTable`（AutoMigrate）→ `seedBlogData`，返回清理函数）→ `router.InitRouter()`。默认 root 用户不在 `InitModule` 里：由 `registerAPI` 中 `api.NewUserApi()` → `service.NewUserService()` → `RegisterRootUser()` 注册，失败直接 `log.Fatalf`（排查“默认账号没出现”时先看这里）。访问日志/gzip 等在 `main.go` 与 `router/middleware/` 注册。
- **登录与鉴权**：`api/user_api.go`（`GET /api/user/captcha` 图形验证码 → `internal/service/captcha_service.go`；`POST /api/user/login` → `user_service.go` 签发 token，模型 `data/model/user_token.go`）。后续请求经 `router/middleware/auth.go` 校验 `X-Token` 并注入 `global.UserInfoKey`；后台再叠加 `middleware/admin.go`（`RequireAdmin`，`userType` ∈ {2, 100}）。dev 模式下 `debug.enableMockLogin` 可跳过真实鉴权（release 下永不生效）。登录与上传限流在 `middleware/ratelimit.go`，规则表在 `router/router.go`。
- **公开博客只读 API（`/api/blog/*`）**：`module/blog/router/router.go` → `api/handler.go` + `api/feed.go` → `internal/service/`（`service.go` 内容与列表、`catalog.go` 分类/标签/归档/统计、`feed.go` RSS）。可见性规则集中在 `service.published()`：仅 `publish_status = published` 且 `published_at <= now`，草稿/未来文章对公众 404。`robots.txt`/`sitemap.xml`/`rss.xml` 由 `RegisterSiteFiles` 挂在站点根。演示数据种子在 `module/blog/seed/`，开关读 `config/blog_seed.go`。
- **后台管理 CRUD（`/api/blog-mgr/*`）**：`module/blog_mgr/router/router.go` 按 `api.ResourceNames()`（categories/tags/posts/diaries/portfolio-items/tools/bookmarks）循环注册同一套泛化 Handler：List/Detail/Create/Update/`PATCH /:id/fields` 快捷字段/Delete/DeleteMany/`POST /:id/restore`。业务在 `internal/service/service.go`：快捷修改白名单仅 `pinned`（文章/日记）与 `publishStatus`；`version` 乐观锁（`requiresVersion`/`modelVersion`，版本过期返回冲突码）；删除为 GORM 软删除，回收站与恢复依赖 `Unscoped()`。请求 DTO 与字段级校验错误（`request.Field(...)`）在 `data/request/request.go`。
- **站点设置与个人资料**：`GET/PUT /api/blog-mgr/profile|/site`，逻辑在 `blog_mgr/internal/service/site.go`；site 的模块可见性/里程碑等存 JSON `settings` 字段。
- **上传 / 媒体库 / 备份**：Handler 在 user 模块——`api/file_api.go`（`UploadBlogImage` ≤ 8MB、`UploadBlogMedias` ≤ 50MB，文件按内容 sha256 命名存入 `files/blog/`，返回 `/files/blog/<hash>.<ext>`，相同内容去重）、`api/media_api.go`（列表、`references` 扫描内容的封面/图标/相册/JSON 设置与 Markdown 正文统计引用数、单删/批删）、`api/backup_api.go` + `global/db_backup.go`（SQLite 在线备份下载）。注意：这几个 Handler 直接查 `global.Db`，是「Handler 不做 DB 查询」规则的遗留例外，新代码不要效仿，也不要顺手重构。这些路由统一在 `blog_mgr/router/router.go` 注册，走后台鉴权链。
- **统一响应**：外层 `{code,msg,data}` 由 `module/common/response` 输出；业务错误码在各模块 `data/response`。
- **前端托管**：`router/router.go` `registerFrontend` 服务 `module/user/web/`（Vite 构建产物）：注入 `<base href="/">`、`/web` 301 收敛、SPA history 兜底；带扩展名的静态资源缺失时返回 404（不回退 index.html，避免 MIME 伪错误）。上传文件经 `r.Static("/files", config.FilesPath)` 暴露。

## 实现约束

- 用户资源查询必须显式包含 `user_id` 所有权条件，并沿用相邻 Service 的鉴权方式；Room/Team 空间模型尚未实现，不要为其预留路由或中间件。
- 配置文件监听只允许热更新 `app.registerEnabled`。数据库（driver/dsn）、Redis、日志、端口、运行模式和 Debug 配置运行中变更会被忽略并记录警告，修改后必须重启服务。
- 模型的 `TableName()` 当前硬编码为 `td_*`。不要假设修改 `database.tablePrefix` 会自动改变已有模型表名。
- 普通业务成功和失败通常通过 HTTP 200 响应体中的 `code` 区分，但认证中间件会返回 HTTP 401。新增响应时遵循相邻端点。
- 博客内容资源更新用模型上的 `version` 字段做乐观锁（`blog_mgr` 的 `requiresVersion`/`modelVersion`）：初始值来自列标签 `default:1` 与 `buildModel` 里的 `max(version, 1)` 归一化（`BaseModel.BeforeCreate` 只生成雪花 ID，不设 version），更新时 `gorm.Expr("version + 1")`，冲突返回 409 业务码；修改更新流程必须保留版本冲突检查。快捷 PATCH 只允许白名单字段，新增可快捷修改字段须同步服务端白名单与前端表格。
- 内容删除一律 GORM 软删除；回收站列表、彻底查询和恢复依赖 `Unscoped()`（见 `blog_mgr/internal/service/service.go`），不要改成物理删除。媒体文件被批删也不回写引用它的内容行——内容中保留失效的 `/files/...` 字符串，由前端降级占位图兜底。
- 列表排序目前按模块硬编码：公开博客 `orderByDate()`（`module/blog/internal/service/service.go`，`pinned DESC` + `published_on` 方向由 query `sort` 参数控制），后台列表 `created_at ASC/DESC`（`module/blog_mgr/internal/service/service.go`）。`module/user/data/request/request_req.go` 的 `request.Sort` 仅被尚未注册路由的 user-manage DTO 内嵌，暂无调用方。任何新排序必须服务端校验字段白名单和 `asc`/`desc` 方向，禁止把请求值直接拼进 `ORDER BY`。
- 跨域（CORS）由前置反向代理（如 nginx）统一处理，Go 侧不设置任何 `Access-Control-*` 头，也不要在路由链里再加 CORS 中间件。Service 构造方式不完全统一，新增代码时参考同类、相邻模块，不要强制套用单例或中间件模板。
- GORM `AutoMigrate` 只「加列不删列」：模型**新增**字段会自动 `ALTER TABLE ADD COLUMN`、删除字段不会 DROP，旧库会残留列。**删除模型字段时，必须同步在 `module/user/global/db.go` 的 `migrateDbTable()` 里补 `Db.Migrator().DropColumn(...)` 兼容代码**（参考已移除的 `td_blog_bookmark.public_id`），否则 `NOT NULL` 且无默认值的残留列会让新插入记录直接失败。

## 新增业务模块

1. 在 `data/model` 增加模型；若有新表，将模型加入 `global/db.go` 的 `TableMigrate`。
2. 在 `data/request/<module>` 增加请求 DTO，仅在确有专用输出结构时增加 response DTO。
3. 在 `internal/service` 实现业务、事务和 GORM 查询，并补齐用户所有权条件。
4. 在 `api` 增加 Handler，沿用相邻模块的绑定、校验、认证用户读取和统一响应方式。
5. 在 `router/router.go` 注册路由，并明确选择所需的 Auth 中间件（CORS 由反代处理，不在 Go 侧添加）。
6. 为纯逻辑优先添加单元测试；涉及数据库（PostgreSQL/SQLite）的流程若无法自动测试，至少保证 `go test ./...` 和构建通过，并说明未做集成验证。
7. **同步 API 规格**：在 `router/router.go` 注册路由后，必须按下面的「API 变更登记」补齐 `doc/api_docs/` 里的端点文件与根引用。

## API 变更登记

`doc/api_docs/` 是接口的**唯一权威契约**（OpenAPI 3.1）。它采用「**一个端点一个文件**」的组织方式，与 `module/user/router/router.go`、`module/blog/router/router.go`、`module/blog_mgr/router/router.go` 里实际注册的路由一一对应。

目录结构：

- [doc/api_docs/openapi.yaml](doc/api_docs/openapi.yaml) —— 聚合根，只声明 `info` / `servers` / `tags` / `securitySchemes` 与全部 `paths`；每个 operation 用操作级 `$ref` 指向端点文件（如 `get: {$ref: './blog/list-posts.yaml'}`）。
- `doc/api_docs/blog/`（公开只读博客，21 个）、`blog-mgr/`（后台管理，80 个，资源 CRUD 在 `blog-mgr/<资源>/` 子目录）、`publish/`（机器令牌发布，12 个）、`user/`（鉴权与 `/ping`，5 个）、`backup/`（NAS 导出，1 个）—— 共 **119 个端点文件**。
- `doc/api_docs/components/{schemas,responses,parameters}.yaml` —— 跨端点共享的模型、具名响应与参数。端点文件用相对路径引用（`../components/...`，资源子目录下是 `../../components/...`）。

**任何新增、修改或删除 API 的改动，都必须同步这里的规格**（"每次修改和新增了 API 都要在这里登记"）：

1. **新增端点**：在对应模块目录下新建端点文件（内容是一个 operation 对象，**不含** `openapi:` / `paths:` 包装；首行注释写明「方法 + 路径 + 所属模块 + 注册位置」），并在 `openapi.yaml` 的 `paths` 里加上该路径与该方法的 `$ref`。两处缺一不可——只加文件不会被聚合，只加引用会指向不存在的文件。
2. **修改端点**：路径、方法、参数、请求体、data 结构、错误码或鉴权方式任一变化，都要改对应端点文件；新增/修改的字段结构优先加进 `components/schemas.yaml` 并以 `$ref` 复用，不要在端点文件里内联大段 schema。
3. **删除端点**：删除端点文件，并移除 `openapi.yaml` 里对应的 operation 引用。
4. **状态码只能用一次的键**：YAML 不允许同一状态码出现两次。同一 400 若有多个业务码（如参数错误与关联引用不存在），要在 `components/responses.yaml` 里定义一个合并响应对象（参考 `ManageBadRequest` / `ManageToolBadRequest`），而不是写两个 `'400'`。
5. **响应约定要写准**：`module/blog` 与 `module/blog_mgr` 走真实 HTTP 状态码（400/404/409/500）；`module/user` 与上传/媒体/备份类 Handler **失败也返回 HTTP 200**，靠 `code` 区分（失败 `code: 1`、`data: {}`），这类端点的 `data` 要声明成「成功结构 或 `EmptyData`」的 `oneOf`。不要把两套约定混用。
6. **YAML 书写陷阱**：纯文本标量里不能出现「**冒号 + 空格**」（如 `javascript: 这类`、`{"field": "x"}`），也不能以反引号/`@` 开头，否则解析直接失败；这类值要用单引号整体括起来。
7. **落地前自检**：
   ```bash
   cd doc/api_docs
   npx --yes @redocly/cli@latest lint openapi.yaml            # 应无 error
   npx --yes @redocly/cli@latest bundle openapi.yaml -o /tmp/openapi.bundled.yaml
   ```
   另外确认 `openapi.yaml` 的 operation 数、`operationId` 唯一性，与端点文件数量三者一致（当前 119）。
   `lint` 会固定报出一批 `operation-4xx-response` 警告（`/ping`、站点根三个文件与若干只读列表接口确实没有 4xx 分支），这是该规则的固有意见、不是缺陷；出现**其它**规则名的 error/warning 才需要处理。

`doc/blog-api.md` 是面向人的接口说明，可能滞后于实现；契约不明确时以 `doc/api_docs/` 与路由源码为准。

## 已知差异

- 根目录 [README.md](README.md) 的部分配置示例和架构描述可能早于当前实现，例如运行模式、初始化顺序、CLI 生成器和模型字段；实现任务先核对源码。
- 启动期 `AutoMigrate` 只负责当前所选数据库的表结构，不会迁移已有数据；切换 driver 须重启服务。
- GORM 日志默认是 Silent；排查 SQL 时可临时调整日志级别，但不要把调试配置作为无关改动提交。
