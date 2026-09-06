# 宝藏文档（用户与鉴权模块）

本模块当前仅保留博客管理员的登录、退出登录与权限校验等基础能力，
以及博客相关的只读与后台管理路由（由 `module/blog`、`module/blog_mgr` 提供）。
文档、笔记、空间、团队、文件上传、用户注册与用户管理等业务功能已移除。

## 提供的接口

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/user/captcha` | 获取图形验证码 | 匿名 |
| POST | `/api/user/login` | 博客管理员登录 | 匿名（含验证码） |
| POST | `/api/user/logout` | 退出登录 | 登录态（X-Token） |
| 任意 | `/api/blog-mgr/*` | 博客后台管理 | 登录态 + 管理员权限 |

- 通用鉴权中间件：`middleware.Auth()`（校验 `X-Token`，开发模式可开启 mock 登录）。
- 管理员权限校验：`middleware.RequireAdmin()`（要求 `UserType` 为管理员或超级管理员）。

## 默认管理员

服务启动时自动确保超级管理员账号 `treasure-root` 存在（见 `internal/service/user_service.go` 的
`RegisterRootUser`）。该账号的密码在启动日志中输出，请尽快修改。

## 数据更新

服务启动时通过 GORM `AutoMigrate` 初始化数据库表结构（仅 `td_user`、`td_user_token`
及博客相关表）。当前默认使用 SQLite（单文件，零依赖），也可在 `[database]` 配置为 PostgreSQL；
该过程不会迁移已有数据，本项目当前也不提供跨数据库的数据搬迁脚本。
