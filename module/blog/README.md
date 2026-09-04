# Blog 公开 API 模块

本模块实现 `doc/blog-api.md` 定义的匿名只读接口，并挂载到现有 `module/user` 服务的 `/api/blog` 路径。模块没有独立的 `main`、配置文件或数据库连接。

## 数据与迁移

- 服务启动时，Blog 模型随 `global.TableMigrate` 通过 GORM `AutoMigrate` 创建或调整 PostgreSQL 表。
- 内容状态为 `draft`、`published` 或 `archived`。公开查询仅返回 `published`、未软删除且 `published_at` 不晚于当前时间的数据。
- 分类的 `scope` 可取 `post`、`portfolio`、`bookmark`。
- 工具的 `kind` 可取 `own` 或 `link`。只有 `own` 类型可以通过详情接口访问。
- Profile 与 Site 使用 `default` 作为单例 key；空库时返回字段完整的空对象。
- 本模块不自动写入演示数据，也不提供管理端 CRUD。

## 表

- `td_blog_category`、`td_blog_tag`
- `td_blog_post`、`td_blog_post_tag`
- `td_blog_diary`、`td_blog_diary_tag`
- `td_blog_portfolio_item`
- `td_blog_tool`
- `td_blog_bookmark`、`td_blog_bookmark_tag`
- `td_blog_profile`、`td_blog_site`

标签和分类使用关系表，作品链接、技术栈、Profile 和 Site 的有序集合使用 PostgreSQL JSONB。

## 本地验证

```bash
go test ./module/blog/...
go test ./...
go build -o treasure_user.exe ./module/user
```

实际启动仍需从 `module/user` 目录运行，并提供可用 PostgreSQL 配置：

```bash
cd module/user
go run . -c config.toml
```

启动会连接数据库并执行迁移，不应把启动命令用作无副作用的测试。
