# 重置默认管理员密码

配置文件必须使用 `[database]` 配置（SQLite 或 PostgreSQL 均可），并能从当前环境访问数据库。

在仓库根目录运行（重置的是默认管理员账号 `treasuredocmgr`，新密码需为 8-16 位）：

```bash
go run ./module/user -c module/user/config.toml resetpwd <新密码>
```

注意：`-c` 必须写在子命令 `resetpwd` 之前，因为它是全局 flag。
