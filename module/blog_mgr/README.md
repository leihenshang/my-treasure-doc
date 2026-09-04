# Blog 管理 API

本模块管理 `module/blog` 使用的 PostgreSQL 数据，挂载在 `/api/blog-mgr`。所有请求必须携带有效 `X-Token`，且当前用户类型必须为 admin 或 root。

开发环境可通过配置启用固定 root Mock：

```toml
[app]
runMode = "dev"

[debug]
enableMockLogin = true
mockUserId = "9999999999"
```

启用后管理请求可以不带 `X-Token`。`mockUserId` 只设置内存 Mock 用户 ID，不查询数据库，角色固定为 root。release 模式禁止启用 Mock，配置修改后必须重启。

## 资源路由

以下资源使用统一 CRUD：

- `categories`
- `tags`
- `posts`
- `diaries`
- `portfolio-items`
- `tools`
- `bookmarks`

```text
GET    /api/blog-mgr/{resource}
POST   /api/blog-mgr/{resource}
GET    /api/blog-mgr/{resource}/{id}
PATCH  /api/blog-mgr/{resource}/{id}
DELETE /api/blog-mgr/{resource}/{id}
POST   /api/blog-mgr/{resource}/{id}/restore
```

管理路径中的 `id` 是数据库 ID，不是文章 slug 或日记 publicId。删除使用软删除，不提供物理删除。

文章、日记和书签的管理列表、详情、创建和更新响应均包含 `tagIds`。公开 `/api/blog/tags` 和 `/api/blog/diary/tags` 返回包含 `id`、`name` 的标签对象数组，可直接用于管理端选择器。

列表参数包括 `page`、`pageSize`、`keyword`、`status`、`deleted`、`scope`、`categoryId` 和 `sort`。`deleted` 可取 `exclude`、`only`、`all`；`sort` 可取 `asc`、`desc`。

Profile 和 Site 使用单例接口：

```text
GET /api/blog-mgr/profile
PUT /api/blog-mgr/profile
GET /api/blog-mgr/site
PUT /api/blog-mgr/site
```

## 内容更新

文章、日记、作品、工具和书签更新时必须提交当前 `version`。更新成功后版本递增；版本过期或唯一标识冲突返回 HTTP 409。

文章创建示例：

```json
{
  "slug": "hello-world",
  "title": "Hello World",
  "summary": "第一篇文章",
  "categoryId": "essay",
  "author": "Treasure Doc",
  "content": "# Hello World",
  "publishStatus": "draft",
  "publishedOn": "2026-09-04",
  "pinned": false,
  "tagIds": []
}
```

状态可取 `draft`、`published`、`archived`。首次发布且未提供 `publishedAt` 时使用服务器当前时间；未来时间表示预约发布。

## 校验规则

- 公开 URL 仅允许 `https://`；Profile 联系方式额外允许 `mailto:`。
- 工具 `kind` 仅允许 `own` 或 `link`。外链工具必须提供 HTTPS URL。
- 分类 scope 仅允许 `post`、`portfolio`、`bookmark`，创建后不可修改。
- Profile 技能 level 必须在 0 至 100。
- Site 模块路径必须以 `/Blog` 开头。
- 文章、日记和书签的标签关系与主体更新处于同一事务。

## 验证

```bash
go test ./module/blog_mgr/...
go test ./...
go build -o treasure_user.exe ./module/user
```

PostgreSQL JSONB、事务回滚和写入后公开可见性的完整验证需要独立测试数据库，不应使用项目运行配置指向的数据库。
