# AGENTS.md

## 项目范围

- 本仓库是 Go 1.22 单模块项目（`fastduck/treasure-doc`），核心服务位于 `module/user`，技术栈为 Gin、GORM、PostgreSQL；Redis 可选。
- 先阅读根目录 [README.md](README.md) 了解产品与 API 概览。部署目录和历史数据修复见 [module/user/README.md](module/user/README.md)，空间与权限的规划见 [doc/space.md](doc/space.md)。设计文档描述的是目标状态，实际行为以代码和 `router/router.go` 为准。
- 代码与文档主要使用中文。保持现有命名、分层和错误响应风格，不做与任务无关的架构重构。

## 常用命令

在仓库根目录运行：

```bash
go test ./...
go test ./list_sort/... -run '^Test_ParseSortParams$'
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

- 首次运行前按需从 `config.example.toml` 创建本地 `config.toml` 并配置可用的 PostgreSQL。
- 启动服务会连接数据库、执行 `AutoMigrate`，并尝试注册默认 root 用户；不要把启动服务当作无副作用的验证步骤。
- 重置密码使用 `go run ./module/user/cli/reset-pwd -u <账号> -p <新密码> -cfg <配置文件绝对路径>`。`module/user/cli/cli.go` 目前没有可用的代码生成入口。
- 仓库没有 CI、Makefile 或 lint 配置。修改 Go 代码后至少运行相关包测试和 `go test ./...`；提交前对改动文件执行 `gofmt` 或 `go fmt`。

## 代码边界

请求链路为 `main.go` → `global.InitModule()` → `router.InitRouter()` → middleware → `api` → `internal/service` → `global.Db` → `data/model`。

- `api/`：绑定请求、从 Gin context 获取当前用户、调用 Service，并使用 `data/response` 输出；不要在 Handler 中新增数据库查询。
- `internal/service/`：业务规则、事务和 GORM 查询。项目没有 DAO/Repository 层；除非任务明确要求架构调整，否则沿用这一结构。
- `data/request/` 与 `data/response/`：请求 DTO、分页/排序参数、响应 DTO 和业务错误码。
- `data/model/`：GORM 模型、表名、软删除和创建钩子。新增模型通常嵌入 `BaseModel` 并实现 `TableName()`。
- `router/router.go`：实际生效的路由清单。存在 API 或 Service 文件不代表端点已经暴露。
- `global/`：配置、数据库、Redis、日志、迁移和 validator 等进程级状态。

## 实现约束

- 用户资源查询必须显式包含 `user_id` 所有权条件，并沿用相邻 Service 的鉴权方式。当前 Room/Team 权限仍不完整，不要把 [doc/space.md](doc/space.md) 的规划当成已实现规则。
- 配置热更新会替换 `global.Db` 和 `global.Redis`；不要把这两个全局连接复制到长期存活的包级变量或结构体字段中。
- 模型的 `TableName()` 当前硬编码为 `td_*`。不要假设修改 `database.tablePrefix` 会自动改变已有模型表名。
- 普通业务成功和失败通常通过 HTTP 200 响应体中的 `code` 区分，但认证中间件会返回 HTTP 401。新增响应时遵循相邻端点。
- 文档更新使用 `Doc.Version` 做乐观锁，并在事务中写入历史快照；修改更新或恢复流程时必须保留并发冲突检查和事务边界。
- 文档删除使用 GORM 软删除；回收站查询和恢复依赖 `Unscoped()`，不要改成物理删除。
- 业务列表使用 `data/request/request_req.go` 中的排序逻辑；根目录 `list_sort` 是独立包，目前未接入业务接口。拼接 SQL 排序前必须同时校验字段白名单和 `asc`/`desc` 方向。
- CORS 不是全局中间件，Service 构造方式也不完全统一。新增代码时参考同类、相邻模块，不要强制套用单例或中间件模板。

## 新增业务模块

1. 在 `data/model` 增加模型；若有新表，将模型加入 `global/db.go` 的 `TableMigrate`。
2. 在 `data/request/<module>` 增加请求 DTO，仅在确有专用输出结构时增加 response DTO。
3. 在 `internal/service` 实现业务、事务和 GORM 查询，并补齐用户所有权条件。
4. 在 `api` 增加 Handler，沿用相邻模块的绑定、校验、认证用户读取和统一响应方式。
5. 在 `router/router.go` 注册路由，并明确选择所需的 Auth/CORS 中间件。
6. 为纯逻辑优先添加单元测试；涉及 PostgreSQL 的流程若无法自动测试，至少保证 `go test ./...` 和构建通过，并说明未做集成验证。

## 已知差异

- 根目录 [README.md](README.md) 的部分配置示例和架构描述可能早于当前实现，例如运行模式、初始化顺序、CLI 生成器和模型字段；实现任务先核对源码。
- 启动期 `AutoMigrate` 只负责 PostgreSQL 表结构，不会迁移已有 MySQL 数据。
- `team` 已有 API/Service 代码但尚未在路由中注册。
- GORM 日志默认是 Silent；排查 SQL 时可临时调整日志级别，但不要把调试配置作为无关改动提交。
- 历史问题记录见 [question.md](question.md)，独立排序包说明见 [list_sort/README.md](list_sort/README.md)。
