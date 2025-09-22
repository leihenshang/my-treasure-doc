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

docker build -f module/user/Dockerfile  -t treasure-doc-go . 

docker run --rm --name treasure-doc-go -it -p 2025:2025 treasure-doc-go:latest 