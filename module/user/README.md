# 宝藏文档

```text

\app
\config.toml
\files
    uploads
\web

```

## 数据更新

服务启动时通过 GORM `AutoMigrate` 初始化 PostgreSQL 表结构。该过程不会迁移已有 MySQL 数据；本项目当前也不提供 MySQL 到 PostgreSQL 的数据搬迁脚本。

以下语句用于修复已经导入 PostgreSQL 的历史数据：

```sql
UPDATE td_doc SET group_id = 'root'  WHERE group_id = '' OR group_id ='0';

UPDATE td_doc_group SET p_id = 'root' WHERE p_id = '' OR p_id ='0';
```
