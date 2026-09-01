# CODEBUDDY.md This file provides guidance to CodeBuddy when working with code in this repository.

## 项目概览

宝藏文档 (Treasure Doc) 是一个文档管理系统的后端 API 服务，基于 Gin + GORM + MySQL。仓库是单个 Go module（`fastduck/treasure-doc`，Go 1.22），包含多个 main 包，没有 Makefile / CI / lint 配置。代码注释与设计文档以中文为主。

## 常用命令

- **启动开发服务**（默认读取 `./config.toml`，需先 `cp config.example.toml config.toml` 并按本地环境修改 MySQL 配置）：
  ```bash
  cd module/user && go run main.go
  ```
  用 `-c /绝对路径/config.toml` 指定配置文件。启动后自动 AutoMigrate 建表、注册 root 账号（`treasure-root / treasure-root`），监听 `:2021`。服务需要可用的 MySQL；Redis 在配置 `enable=false` 时跳过。

- **运行测试**（仓库内仅 `list_sort` 和 `global/gid` 有测试）：
  ```bash
  go test ./list_sort/...
  go test ./module/user/global/gid/...
  ```
  运行单个测试：`go test -run TestSortObj ./list_sort/`（`-run` 支持正则匹配测试名）。

- **构建**（CGO 关闭的静态编译）：
  ```bash
  # Linux / macOS
  CGO_ENABLED=0 go build -o treasure_user ./module/user
  # Windows
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o treasure_user.exe ./module/user
  ```

- **Docker 构建**：`docker build -t treasure-doc .`（多阶段构建，见 `Dockerfile`，产物为 `module/user` 编译出的二进制 + `web/` + `files/` + `config.example.toml`）。

- **重置用户密码**（`cli/` 目录下的 CLI 工具）：
  ```bash
  cd module/user/cli/reset-pwd
  go run reset_pwd.go -u <账号> -p <新密码> -cfg <config.toml 绝对路径>
  ```
  注意：CLI 工具默认从当前工作目录加载配置，务必通过 `-cfg` 传入绝对路径，否则读不到配置。

- **代码生成工具**：`module/user/cli/cli.go` 目前只有 `package main`（空实现），README 中提到的 `go run . -gen` 模型生成功能尚未实现，不要依赖它。

## 架构

### 分层结构

请求链路：`main.go` → `global.InitModule()`（启动引导）→ `router.InitRouter()` → `middleware.Auth()` → `api` Handler → `internal/service` → `global.Db`（GORM 直接操作，无 DAO 层）→ `data/model`。

核心目录（均在 `module/user/` 下）：

- **`api/`** — Handler 层。每个业务模块一个 `*Api` 结构体（持有对应 `*Service` 单例）。统一模式：`ShouldBindJSON/ShouldBindQuery` 解析请求 → `auth.GetUserInfoByCtx(c)` 取当前用户 → 调用 Service → `response.Ok*` / `response.FailWithMessage` 输出。
- **`internal/service/`** — 业务逻辑层。`*Service` 结构体 + `sync.Once` 单例（`NewXxxService()`）。直接使用全局 `global.Db` 做 GORM 查询，**没有独立的 DAO/Repository 层**（已知问题，新增查询逻辑时注意保持与现有 Service 风格一致）。错误通过返回 `err` 由 Handler 转成响应；文档更新冲突使用哨兵错误 `service.ErrorDocIsEdited`。
- **`data/model/`** — GORM DO 模型。所有模型内嵌 `BaseModel`（`Id` + 时间戳 + 软删除），`BeforeCreate` 钩子用 Sonyflake 生成 19 位数字字符串 ID（`global/gid.GenId()`）。模型上附带 `TableName()`（形如 `td_doc`）和业务常量/枚举（如 `UserTypeRoot`、`NoteTypeDoc`）。
- **`data/request/`** — 请求 DTO，按模块分子目录（`doc/`、`user/`、`room/` 等）。公共结构在 `request_req.go`：`IDReq`、`Pagination{Page,PageSize,Total}`、`Sort`。
- **`data/response/`** — 统一响应 `Response{code, msg, data}`（`code` 为 `ErrorCode` 枚举）；列表用 `ListResponse{pagination, list}`。所有响应 HTTP 状态码恒为 200，业务错误通过 `code` 字段区分。
- **`global/`** — 包级单例：`Conf`、`Db`、`Zap`、`Log`、`Redis`、`Trans`（校验器翻译）。`InitModule()` 是启动引导序列：配置 → 日志 → Redis（可选）→ MySQL → validator → 配置热更新 → 返回析构函数。`db.go` 里的 `TableMigrate` 列表在启动时逐个 `AutoMigrate`；表名策略为单数 + 配置前缀（`td_`）。`global.constant.go` 定义 `RootGroup = "root"`、`UserInfoKey = "userinfo"` 等关键常量。
- **`router/`** — 路由注册 + 中间件。`middleware.Auth()` 读取 `X-Token` 头，通过 `service.GetUserByToken` 校验（每用户最多 3 个并发会话，超出剔除最早 token），并把 `*model.User` 放入 gin context（key 为 `UserInfoKey`）。`Debug.EnableMockLogin=true` 时跳过校验直接注入 `mockUser`（开发用）。
- **`internal/auth/`** — `GetUserInfoByCtx(c)` 从 context 取回用户。
- **`config/`** — Viper 封装。`InitConf` 加载 TOML；`WatchConf` 监听文件变化触发回调。`App.RunMode` 取值 `release`/`dev`（注意 README 中 `release = false` 的示例已过时，实际配置项是 `runMode` 字符串）。
- **`utils/`** — 工具函数：`PasswordEncrypt`（bcrypt）、日期/文件/slice 辅助。
- **`web/`**（静态前端资源）、**`files/`**（上传文件）通过 `r.Static` 挂载，分别由 `config.WebPath`/`config.FilesPath` 指定。

### 关键设计

- **配置热更新**：`config.WatchConf` 触发后 `global.applyConnectionHotReload` 会比较新旧配置，MySQL/Redis 配置变化时重建连接池并替换 `global.Db`/`global.Redis`（用 `connMu` 互斥）。因此代码中**不要缓存 `global.Db` 引用**，每次通过包级变量读取。
- **多租户/空间隔离（Room）**：用户注册时自动创建默认 Room（`个人空间`），`User.CurrentRoomId` 指向当前空间；`Doc`/`DocGroup`/`Note` 均带 `UserId` 并在查询中强制按 `user_id` 过滤，Service 层入参通常为 `(…, userId string)`。
- **文档版本与历史**：`Doc.Version` 字段做乐观锁（更新时 `WHERE version = ?` 且 `version+1`）；更新在事务里先写 `DocHistory` 快照再更新正文。删除是 GORM 软删除（进回收站），`Recover` 用 `Unscoped().Update("deleted_at", nil)` 恢复，回收站列表用 `Unscoped().Where("deleted_at is not null")` 查询。
- **分组树**：`DocGroup` 用 `PId`（父节点）+ `GroupPath`（逗号分隔的祖先链）组织树形结构，`doc-group/tree` 与 `Detail`（回填父分组链）依赖该设计。
- **笔记与置顶**：`Note` 可关联 `Doc`（`NoteTypeDoc`），文档置顶（`IsPin`）实质是对 Doc-Note 的增删。
- **排序**：列表接口排序参数格式为 `id_asc,name_desc`（或 JSON 数组，`list_sort` 包支持），字段名经白名单映射（如 `createdAt → created_at`）后拼入 SQL，防注入。分页用 `request.Pagination.Offset()`（`(page-1)*pageSize`）。
- **密码与会话**：bcrypt 加密；登录返回 token（7 天有效）；`RegisterRootUser` 在 `NewUserService()` 首次调用时自动注册 root 账号。
- **GORM SQL 日志默认关闭**（`logger.Silent`，见 `global/db.go`）。排查 SQL 问题时需临时把 `LogLevel` 改为 `logger.Info`，这是已知问题（见 `question.md`）。

### 其他包

- **`list_sort/`** — 独立库，解析 `sort` 参数（JSON 或 `field_order` 字符串），带单元测试。可复用于任意列表接口。
- **`module/admin/`** — 仅占位（README），管理端功能暂在 user 模块 `user-manage` 路由下实现。
- **`template/`**、**`doc/`**（`space.md` 空间/团队设计文档）、**`log/`**（运行时日志，gitignore 外）。

### 新增模块的惯例

1. `data/model/` 建模型（内嵌 `BaseModel`，实现 `TableName()`）。
2. `data/request/<module>/` 建请求 DTO，`data/response/` 建专用响应（可选）。
3. `internal/service/` 建 `XxxService` 单例。
4. `api/` 建 `XxxApi`，在 `router/router.go` 挂路由（默认 `middleware.Auth()` + `middleware.Cors()`）。
5. 若需新表，把模型加进 `global/db.go` 的 `TableMigrate`。

## 已知问题（来自 question.md / README）

- 无 DAO 层，Service 直接操作 `global.Db`，查询逻辑分散。
- CLI 子工具（`cli/`、`cli/reset-pwd/`）按 CWD 加载配置，需 `-c`/`-cfg` 传绝对路径。
- GORM 未配置 SQL 日志。
- 核心 Service 缺少单元测试（仅 `list_sort`、`gid` 有测试）。
