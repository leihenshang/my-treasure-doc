# Blog 公开 API 模块

本模块实现 `doc/blog-api.md` 定义的匿名只读接口，并挂载到现有 `module/user` 服务的 `/api/blog` 路径。模块没有独立的 `main`、配置文件或数据库连接。

## 数据与迁移

- 服务启动时，Blog 模型随 `global.TableMigrate` 通过 GORM `AutoMigrate` 创建或调整表结构（默认 SQLite，可选 PostgreSQL）。
- 内容状态为 `draft`、`published` 或 `archived`。公开查询仅返回 `published`、未软删除且 `published_at` 不晚于当前时间的数据。
- 分类的 `scope` 可取 `post`、`portfolio`、`bookmark`。
- 工具的 `kind` 可取 `own` 或 `link`。只有 `own` 类型可以通过详情接口访问。
- Profile 与 Site 使用 `default` 作为单例 key；空库时返回字段完整的空对象。
- 本模块提供默认关闭的配置驱动 Seed；管理 CRUD 由 `module/blog_mgr` 提供。

## Mock 数据 Seed

Seed 在全部 `AutoMigrate` 成功后执行，默认关闭：

```toml
[blogSeed]
enabled = false
allowRemote = false
restoreDeleted = false
```

- `enabled`：启动时执行 Seed。
- `allowRemote`：允许向非 localhost 数据库写入；远程库必须显式设为 true。
- `restoreDeleted`：恢复与 Seed 固定业务键匹配的软删除记录；默认跳过。

Seed 仅允许在 `runMode = "dev"` 下执行。配置运行中修改不会触发 Seed，必须重启服务。

固定数据包括 20 篇文章、20 篇日记、11 个分类、10 个标签、8 个作品、8 个工具、8 个书签以及 Profile/Site。文章和日记各有 12 条当前公开数据，其余覆盖预约发布、软删除、草稿和归档；公开作品为 4 条。

Seed 使用单一事务并按 slug/publicId/default key 幂等查找。已有记录不会被覆盖，标签关系只补缺失项，不删除管理员后来添加的关系。

## 表

- `td_blog_category`、`td_blog_tag`
- `td_blog_post`、`td_blog_post_tag`
- `td_blog_diary`、`td_blog_diary_tag`
- `td_blog_portfolio_item`
- `td_blog_tool`
- `td_blog_bookmark`、`td_blog_bookmark_tag`
- `td_blog_profile`、`td_blog_site`

标签和分类使用关系表，作品链接、技术栈、Profile 和 Site 的有序集合使用跨方言 JSON 字段（PostgreSQL 下为 jsonb，SQLite 下为 json/text）。

## 本地验证

```bash
go test ./module/blog/...
go test ./...
go build -o treasure_user.exe ./module/user
```

实际启动仍需从 `module/user` 目录运行，并提供可用的数据库配置（默认 SQLite，可选 PostgreSQL）：

```bash
cd module/user
go run . -c config.toml
```

启动会连接数据库并执行迁移，不应把启动命令用作无副作用的测试。
