# 宝藏文档

```text

\app
\config.toml
\files
    uploads
\web

```

## 数据更新
```sql
UPDATE td_doc SET group_id = 'root'  WHERE group_id = '' OR group_id ='0';

UPDATE td_doc_group SET p_id = 'root' WHERE p_id = '' OR p_id ='0';
```