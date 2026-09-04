# 重置用户密码

配置文件必须使用 `[database]` PostgreSQL 配置，并能从当前环境访问数据库。

在仓库根目录运行：

```bash
go run ./module/user/cli/reset-pwd -u <账号> -p <新密码> -cfg <config.toml 绝对路径>
```
