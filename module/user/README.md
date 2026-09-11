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

服务启动时自动确保超级管理员账号 `treasuredocmgr` 存在（初始密码同账号，见 `internal/service/user_service.go` 的
`RegisterRootUser`）。该账号创建时会被标记「强制修改密码」（`td_user.require_pwd_reset = 1`），并贯穿整条链路：

1. 登录响应返回 `requirePwdReset: true`；
2. 管理端登录后**强制弹出**修改密码弹窗：不可关闭、不可 ESC、没有「取消」按钮；
3. 该标记随登录态写入本地缓存，刷新页面或直接输入后台地址同样会被拦截；
4. 改密完成前，后端 `middleware.RequireAdmin()` 会拒绝所有 `/api/blog-mgr/*` 请求并返回 `40301`；
5. 改密成功后在同一个事务里清除标记（`ChangePassword`），并使其它设备的登录态失效。

启动日志不再打印明文密码。

### 修改默认管理员密码

- **管理端修改（推荐）**：按强制弹窗提示输入原密码与新密码即可；非强制场景可在右上角头像下拉菜单选择「修改密码」。
  修改成功后其它设备的登录态会失效，当前登录保持有效，且强制改密标记自动清除。
- **命令行重置**（忘记密码或无法登录时，密码规则与登录入口一致，8-16 位）：

  ```bash
  # 仓库根目录执行；-c 是全局 flag，必须写在子命令 resetpwd 之前
  go run ./module/user -c module/user/config.toml resetpwd '新密码'
  ```

  重置后该账号所有登录态立即失效，下次登录需使用新密码。
  注意两点：`module/user/cli/reset-pwd/` 目录下只有说明文档、**没有独立的 cli 包**，重置统一走上表这个子命令；
  该子命令不清除 `require_pwd_reset` 标记，因此重置后首次登录仍会被要求改一次密码。

### 容器部署下的密码初始化

镜像（`Dockerfile`）**没有**任何自动初始化或注入密码的机制：不存在 `ADMIN_PASSWORD` 之类的环境变量，
也没有 entrypoint 脚本执行种子逻辑，`CMD` 只是直接启动服务。容器首次启动时行为与本地一致 ——
用挂载进来的 `config.toml` 建库，并创建默认管理员 `treasuredocmgr / treasuredocmgr`，随后由强制改密流程接管。

因此在容器里只有两种方式：

- **首次部署（推荐）**：用默认账号登录，按强制弹窗改密码。
- **忘记密码**：在容器内执行同一个二进制的子命令。容器 `WORKDIR` 为 `/app`、配置挂载在 `/app/config.toml`，
  默认路径即可命中，无需额外参数：

  ```bash
  docker exec -it treasure-doc /app/treasure-doc resetpwd '新密码'
  ```

  该命令连接配置指向的同一个数据库（SQLite 单文件或 PostgreSQL），服务**无需停机**；
  SQLite 为 WAL 模式，允许多进程并发读写。执行完进程自行退出。

> 若希望"用环境变量在首次启动时注入管理员密码"，需要新增启动期读取环境变量的逻辑，当前版本未提供。

## 数据更新

服务启动时通过 GORM `AutoMigrate` 初始化数据库表结构（仅 `td_user`、`td_user_token`
及博客相关表）。当前默认使用 SQLite（单文件，零依赖），也可在 `[database]` 配置为 PostgreSQL；
该过程不会迁移已有数据，本项目当前也不提供跨数据库的数据搬迁脚本。
