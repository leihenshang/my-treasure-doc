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
UPDATE td_doc SET group_id = ''  WHERE group_id = '0';
UPDATE td_doc_group SET p_id = '' WHERE p_id = '0';


# 统一设置
UPDATE td_doc SET group_id = '0'  WHERE group_id = '';

UPDATE td_doc_group SET p_id = '0' WHERE p_id = '';
```