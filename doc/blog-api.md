# Treasure Blog API 接口文档

> 版本：v1.0
> 适用范围：当前前端 `/Blog` 模块
> 后端实现建议：Go 1.23+、Gin/Echo/Chi 均可
> 文档状态：公开端与管理端联调契约

## 1. 概述

本文档定义 Treasure Doc 前端 `/Blog` 公开模块所需的只读 API，以及对应内容的后台管理 API，包括文章、日记、作品、工具、书签、分类、标签、个人资料和站点信息。

公开接口允许匿名访问，不要求登录，不校验 `X-Token`。管理接口要求登录且仅允许 admin/root 用户访问。图片上传及其他文档管理接口不在本文档范围内。

### 1.1 Base URL

```text
/api/blog
```

本地开发时由 Vite 将 `/api` 代理至 Go 服务。当前 `vite.config.ts` 的目标为 `http://localhost:2021`。

### 1.2 通用约定

- 协议：HTTPS（本地开发可使用 HTTP）
- 数据格式：`application/json; charset=utf-8`
- 字符编码：UTF-8
- 日期格式：`YYYY-MM-DD`
- ID 类型：字符串。ID 应被视为不透明标识，不应在前端解析为数据库主键
- Markdown：详情接口的 `content` 返回 Markdown 原文
- 字段命名：lower camelCase
- 空数组：返回 `[]`，不要返回 `null`
- 可选字段：无值时可省略；如返回则类型必须符合定义
- 所有 URL 字段只允许安全协议，公开网页使用 `https://`，个人联系方式可使用 `mailto:`

## 2. 统一响应格式

为兼容项目已有的 `TreasureResponse<T>`，所有接口使用统一响应信封。

### 2.1 普通响应

```json
{
  "code": 0,
  "msg": "",
  "data": {}
}
```

| 字段 | 类型 | 必填 | 说明 |
| -------- | ------- | ---- | -------------------------------------- |
| `code` | integer | 是 | 业务状态码，`0` 表示成功 |
| `msg` | string | 是 | 成功时为空字符串，失败时为可读错误信息 |
| `data` | any | 是 | 业务数据；错误时为`null` |

### 2.2 分页响应

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "list": [],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 0,
      "orderBy": "date_desc"
    }
  }
}
```

| 字段 | 类型 | 必填 | 说明 |
| ---------------------------- | ------- | ---- | ----------------------------- |
| `data.list` | array | 是 | 当前页数据 |
| `data.pagination.page` | integer | 是 | 当前页，从`1` 开始 |
| `data.pagination.pageSize` | integer | 是 | 每页数量 |
| `data.pagination.total` | integer | 是 | 筛选后的总数 |
| `data.pagination.orderBy` | string | 是 | 实际排序规则，如`date_desc` |

### 2.3 HTTP 状态和业务错误码

HTTP 状态码与业务码应同时正确返回。

| HTTP | 业务码    | 含义                             |
| ---- | --------- | -------------------------------- |
| 200  | `0`     | 成功                             |
| 400  | `40001` | 查询参数格式错误                 |
| 400  | `40002` | 不支持的排序方式                 |
| 404  | `40401` | 文章不存在或已删除               |
| 404  | `40402` | 日记不存在或已删除               |
| 404  | `40403` | 作品不存在或已移除               |
| 404  | `40404` | 工具不存在，或该条目不是自研工具 |
| 429  | `42900` | 请求过于频繁                     |
| 500  | `50000` | 服务内部错误                     |

错误示例：

```json
{
  "code": 40401,
  "msg": "文章不存在或已被删除",
  "data": null
}
```

列表无匹配项不是错误，应返回 HTTP 200、空数组和 `total: 0`。

## 3. 接口总览

| 模块 | 方法 | Endpoint                           | 说明                   |
| ---- | ---- | ---------------------------------- | ---------------------- |
| 文章 | GET  | `/api/blog/categories`           | 文章分类               |
| 文章 | GET  | `/api/blog/tags`                 | 文章标签               |
| 文章 | GET  | `/api/blog/posts`                | 文章分页列表           |
| 文章 | GET  | `/api/blog/posts/{id}`           | 文章详情               |
| 日记 | GET  | `/api/blog/diary/tags`           | 日记标签               |
| 日记 | GET  | `/api/blog/diaries`              | 日记分页列表           |
| 日记 | GET  | `/api/blog/diaries/{id}`         | 日记详情               |
| 作品 | GET  | `/api/blog/portfolio/categories` | 作品分类               |
| 作品 | GET  | `/api/blog/portfolio/items`      | 作品列表               |
| 作品 | GET  | `/api/blog/portfolio/items/{id}` | 作品详情               |
| 工具 | GET  | `/api/blog/tools`                | 工具与外链列表         |
| 工具 | GET  | `/api/blog/tools/{id}`           | 自研工具详情           |
| 书签 | GET  | `/api/blog/bookmark/categories`  | 书签分类               |
| 书签 | GET  | `/api/blog/bookmarks`            | 书签列表               |
| 关于 | GET  | `/api/blog/profile`              | 个人资料               |
| 关于 | GET  | `/api/blog/site`                 | 站点信息               |
| 统计 | GET  | `/api/blog/stats`                | 内容数量，可选优化接口 |

## 4. 文章接口

### 4.1 数据结构

#### BlogCategory

| 字段 | 类型 | 必填 | 说明 |
| -------- | ------ | ---- | ------------------------ |
| `id` | string | 是 | 分类稳定标识，如`tech` |
| `name` | string | 是 | 分类展示名称，如`技术` |

#### BlogTag

| 字段 | 类型 | 必填 | 说明 |
| -------- | ------ | ---- | --------------------------------------- |
| `id` | string | 是 | 标签数据库 ID，可用于管理端的`tagIds` |
| `name` | string | 是 | 标签展示名称 |

#### BlogPostSummary

| 字段         | 类型     | 必填 | 说明                                 |
| ------------ | -------- | ---- | ------------------------------------ |
| `id`       | string   | 是   | 文章 slug，如`hello-world`         |
| `title`    | string   | 是   | 标题                                 |
| `summary`  | string   | 是   | 列表摘要，建议纯文本                 |
| `category` | string   | 是   | 分类 ID，必须对应`BlogCategory.id` |
| `tags`     | string[] | 是   | 标签列表                             |
| `author`   | string   | 是   | 作者展示名                           |
| `date`     | string   | 是   | 发布日期，`YYYY-MM-DD`             |
| `pinned`   | boolean  | 是   | 是否置顶                             |

#### BlogPost

包含 `BlogPostSummary` 的全部字段，并增加：

| 字段        | 类型   | 必填 | 说明          |
| ----------- | ------ | ---- | ------------- |
| `content` | string | 是   | Markdown 正文 |

### 4.2 获取文章分类

```http
GET /api/blog/categories
```

成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": [
    { "id": "tech", "name": "技术" },
    { "id": "essay", "name": "随笔" },
    { "id": "reading", "name": "阅读" }
  ]
}
```

分类应使用后台配置顺序返回。

### 4.3 获取文章标签

```http
GET /api/blog/tags
```

成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": [
    {"id": "10001", "name": "CSS"},
    {"id": "10002", "name": "TypeScript"},
    {"id": "10003", "name": "Vue"}
  ]
}
```

要求：从已发布文章中去重，并按稳定的名称升序返回。内容列表和详情中的 `tags` 仍为标签名称字符串数组。

### 4.4 获取文章列表

```http
GET /api/blog/posts
```

查询参数：

| 参数 | 类型 | 必填 | 默认值 | 约束与说明 |
| -------------- | ------- | ---- | -------- | -------------------------- |
| `categoryId` | string | 否 | - | 分类精确匹配 |
| `tag` | string | 否 | - | 标签精确匹配 |
| `keyword` | string | 否 | - | 去除首尾空格，大小写不敏感 |
| `sort` | string | 否 | `desc` | 仅允许`desc`、`asc` |
| `page` | integer | 否 | `1` | 必须大于等于 1 |
| `pageSize` | integer | 否 | `10` | 1 至 100；前端当前传入 3 |

关键词至少搜索：`title`、`summary`、`tags`、`content`。不要求匹配作者和分类名称。

排序规则：

1. `pinned DESC`，置顶文章永远排在普通文章之前。
2. 同一置顶状态内，`sort=desc` 按 `date DESC`；`sort=asc` 按 `date ASC`。
3. 日期相同时建议使用 `id ASC` 保证稳定顺序。

请求示例：

```http
GET /api/blog/posts?categorncyId=tech&tag=Vue&keyword=性能&sort=desc&page=1&pageSize=3
```

成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "list": [
      {
        "id": "vue-rendering-performance",
        "title": "Vue 页面渲染性能排查清单",
        "summary": "从组件更新、长列表、资源加载到性能面板，整理一套可以重复使用的排查顺序。",
        "category": "tech",
        "tags": ["Vue", "性能优化"],
        "author": "Treasure Doc",
        "date": "2026-07-02",
        "pinned": false
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 3,
      "total": 1,
      "orderBy": "date_desc"
    }
  }
}
```

列表接口不得返回 `content`。

### 4.5 获取文章详情

```http
GET /api/blog/posts/{id}
```

路径参数：

| 参数   | 类型   | 说明                                            |
| ------ | ------ | ----------------------------------------------- |
| `id` | string | 文章稳定 ID/slug，必须进行 URL 解码和参数化查询 |

成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "id": "hello-world",
    "title": "Hello World，我的博客开张了",
    "summary": "这是我的第一篇文章。",
    "category": "essay",
    "tags": ["随笔", "公告"],
    "author": "Treasure Doc",
    "date": "2026-08-20",
    "pinned": true,
    "content": "# Hello World，我的博客开张了\n\n欢迎来到我的个人博客。"
  }
}
```

不存在时返回 HTTP 404 和业务码 `40401`。

## 5. 日记接口

### 5.1 数据结构

#### DiarySummary

| 字段 | 类型 | 必填 | 说明 |
| ----------- | -------- | ---- | ---------------------------------------------- |
| `id` | string | 是 | 不透明 ID；即使形似日期也不得按日期主键解析 |
| `title` | string | 是 | 标题 |
| `summary` | string | 是 | 列表摘要 |
| `tags` | string[] | 是 | 标签 |
| `date` | string | 是 | 必须为`YYYY-MM-DD`，前端按固定位置拆分年月日 |
| `mood` | string | 是 | 心情展示文本 |
| `weather` | string | 是 | 天气展示文本 |
| `pinned` | boolean | 是 | 是否置顶 |

#### Diary

包含 `DiarySummary` 全部字段，并增加：

| 字段        | 类型   | 必填 | 说明          |
| ----------- | ------ | ---- | ------------- |
| `content` | string | 是   | Markdown 正文 |

### 5.2 获取日记标签

```http
GET /api/blog/diary/tags
```

响应 `data` 为去重并稳定排序后的 `BlogTag[]`。

```json
{
  "code": 0,
  "msg": "",
  "data": [
    {"id": "20001", "name": "博客"},
    {"id": "20002", "name": "工作"},
    {"id": "20003", "name": "技术"}
  ]
}
```

### 5.3 获取日记列表

```http
GET /api/blog/diaries
```

| 参数 | 类型 | 必填 | 默认值 | 约束与说明 |
| ------------ | ------- | ---- | -------- | ---------------------------------------- |
| `tag` | string | 否 | - | 标签精确匹配 |
| `keyword` | string | 否 | - | 匹配标题、摘要、标签和正文，大小写不敏感 |
| `sort` | string | 否 | `desc` | `desc` 或 `asc` |
| `page` | integer | 否 | `1` | 大于等于 1 |
| `pageSize` | integer | 否 | `10` | 1 至 100；前端当前传入 4 |

排序规则与文章一致：先 `pinned DESC`，再按日期升序或降序。

成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "list": [
      {
        "id": "2026-09-02",
        "title": "把一个想法做成可用页面的一整天",
        "summary": "从早上的草图到晚上的移动端检查，记录一次完整而缓慢的产品迭代。",
        "tags": ["工作", "博客", "复盘"],
        "date": "2026-09-02",
        "mood": "充实",
        "weather": "晴转雨",
        "pinned": false
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 4,
      "total": 1,
      "orderBy": "date_desc"
    }
  }
}
```

### 5.4 获取日记详情

```http
GET /api/blog/diaries/{id}
```

成功响应 `data` 为完整 `Diary`，包含 Markdown `content`。不存在时返回 HTTP 404 和业务码 `40402`。

## 6. 作品接口

### 6.1 数据结构

#### PortfolioCategory

| 字段 | 类型 | 必填 | 说明 |
| -------- | ------ | ---- | -------- |
| `id` | string | 是 | 分类 ID |
| `name` | string | 是 | 展示名称 |

#### PortfolioLink

| 字段 | 类型 | 必填 | 说明 |
| --------- | ------ | ---- | ---------------------- |
| `label` | string | 是 | 链接名称，如`GitHub` |
| `url` | string | 是 | `https://` URL |

#### PortfolioSummary

| 字段          | 类型     | 必填 | 说明                                   |
| ------------- | -------- | ---- | -------------------------------------- |
| `id`        | string   | 是   | 作品 slug                              |
| `title`     | string   | 是   | 标题                                   |
| `summary`   | string   | 是   | 摘要                                   |
| `category`  | string   | 是   | 分类 ID                                |
| `cover`     | string   | 是   | 当前可返回 emoji；未来可扩展为图片 URL |
| `techStack` | string[] | 是   | 技术栈                                 |
| `date`      | string   | 是   | `YYYY-MM-DD`                         |

#### PortfolioItem

包含 `PortfolioSummary` 全部字段，并增加：

| 字段        | 类型            | 必填 | 说明          |
| ----------- | --------------- | ---- | ------------- |
| `links`   | PortfolioLink[] | 是   | 项目外链      |
| `content` | string          | 是   | Markdown 详情 |

### 6.2 获取作品分类

```http
GET /api/blog/portfolio/categories
```

`data` 返回 `PortfolioCategory[]`，按后台配置顺序排列。

### 6.3 获取作品列表

```http
GET /api/blog/portfolio/items?categoryId=website
```

| 参数           | 类型   | 必填 | 说明                       |
| -------------- | ------ | ---- | -------------------------- |
| `categoryId` | string | 否   | 分类精确匹配；缺省返回全部 |

当前前端不分页、不搜索。建议按 `date DESC, id ASC` 返回，以保证稳定顺序。

```json
{
  "code": 0,
  "msg": "",
  "data": [
    {
      "id": "treasure-doc",
      "title": "Treasure Doc 文档管理系统",
      "summary": "基于 Vue 3 与 Vditor 的 Markdown 文档管理系统。",
      "category": "website",
      "cover": "📚",
      "techStack": ["Vue 3", "TypeScript", "Naive UI", "Vditor", "Go"],
      "date": "2026-06-01"
    }
  ]
}
```

### 6.4 获取作品详情

```http
GET /api/blog/portfolio/items/{id}
```

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "id": "treasure-doc",
    "title": "Treasure Doc 文档管理系统",
    "summary": "基于 Vue 3 与 Vditor 的 Markdown 文档管理系统。",
    "category": "website",
    "cover": "📚",
    "techStack": ["Vue 3", "TypeScript", "Go"],
    "links": [
      { "label": "GitHub", "url": "https://github.com/leihenshang/my-treasure-doc" }
    ],
    "date": "2026-06-01",
    "content": "# Treasure Doc 文档管理系统\n\n一个以 Markdown 为中心的文档管理应用。"
  }
}
```

不存在时返回 HTTP 404 和业务码 `40403`。

## 7. 工具接口

工具列表是判别联合类型，通过 `type` 区分自研工具和外部链接。

### 7.1 数据结构

#### ExternalLink

| 字段 | 类型 | 必填 | 说明 |
| -------- | ------ | ---- | --------------------- |
| `id` | string | 是 | 唯一 ID |
| `type` | string | 是 | 固定为`link` |
| `name` | string | 是 | 名称 |
| `desc` | string | 是 | 描述 |
| `url` | string | 是 | 外部`https://` 地址 |

#### OwnTool

| 字段 | 类型 | 必填 | 说明 |
| ----------- | ------ | ---- | ------------------------------ |
| `id` | string | 是 | 工具 slug |
| `type` | string | 是 | 固定为`own` |
| `name` | string | 是 | 名称 |
| `desc` | string | 是 | 描述 |
| `cover` | string | 是 | 当前可返回 emoji |
| `status` | string | 是 | 当前允许`开发中`、`规划中` |
| `content` | string | 是 | Markdown 说明 |

### 7.2 获取工具列表

```http
GET /api/blog/tools
```

```json
{
  "code": 0,
  "msg": "",
  "data": [
    {
      "id": "json-formatter",
      "type": "own",
      "name": "JSON 格式化",
      "desc": "格式化与压缩 JSON，支持语法校验。",
      "cover": "🧾",
      "status": "开发中",
      "content": "# JSON 格式化\n\n该工具正在开发中。"
    },
    {
      "id": "mdn",
      "type": "link",
      "name": "MDN Web 文档",
      "desc": "Web 标准与浏览器 API 的权威参考文档",
      "url": "https://developer.mozilla.org/zh-CN/"
    }
  ]
}
```

当前前端要求列表中的自研工具包含 `content`。后续如数据量增大，可增加不含正文的 `OwnToolSummary`，但需要同步修改前端类型。

### 7.3 获取自研工具详情

```http
GET /api/blog/tools/{id}
```

仅 `type=own` 的条目存在详情。成功返回完整 `OwnTool`。

以下情况均返回 HTTP 404 和业务码 `40404`：

- ID 不存在
- ID 对应 `type=link` 的外部链接

## 8. 书签接口

### 8.1 数据结构

#### BookmarkCategory

| 字段 | 类型 | 必填 | 说明 |
| -------- | ------ | ---- | -------- |
| `id` | string | 是 | 分类 ID |
| `name` | string | 是 | 分类名称 |

#### Bookmark

| 字段         | 类型     | 必填 | 说明                                     |
| ------------ | -------- | ---- | ---------------------------------------- |
| `id`       | string   | 是   | 唯一 ID                                  |
| `title`    | string   | 是   | 标题                                     |
| `url`      | string   | 是   | 外部`https://` 地址                    |
| `desc`     | string   | 是   | 描述                                     |
| `category` | string   | 是   | 分类 ID                                  |
| `tags`     | string[] | 是   | 标签                                     |
| `icon`     | string   | 是   | 当前可返回 emoji；未来可使用 favicon URL |

### 8.2 获取书签分类

```http
GET /api/blog/bookmark/categories
```

`data` 返回 `BookmarkCategory[]`，顺序同时用于前端分组展示。

### 8.3 获取书签列表

```http
GET /api/blog/bookmarks
```

| 参数 | 类型 | 必填 | 说明 |
| -------------- | ------ | ---- | ---------------------------------------- |
| `categoryId` | string | 否 | 分类精确匹配 |
| `keyword` | string | 否 | 匹配标题、描述、标签和 URL，大小写不敏感 |

当前接口不分页。未知分类的数据仍可返回，前端会将其归入“其他”。

```json
{
  "code": 0,
  "msg": "",
  "data": [
    {
      "id": "github",
      "title": "GitHub",
      "url": "https://github.com/",
      "desc": "代码托管与开源社区",
      "category": "dev",
      "tags": ["代码托管", "开源"],
      "icon": "🐙"
    }
  ]
}
```

## 9. 个人资料与站点信息

### 9.1 获取个人资料

```http
GET /api/blog/profile
```

#### ProfileLink

| 字段 | 类型 | 必填 | 说明 |
| --------- | ------ | ---- | --------------------------------------------- |
| `id` | string | 是 | 链接 ID |
| `label` | string | 是 | 展示名称 |
| `value` | string | 是 | 展示值 |
| `url` | string | 否 | 为空时仅展示；允许`https://` 或 `mailto:` |
| `icon` | string | 是 | 当前可返回 emoji |

#### ProfileSkill

| 字段 | 类型 | 必填 | 说明 |
| --------- | ------- | ---- | ------------------------------------------ |
| `name` | string | 是 | 技能名称 |
| `level` | integer | 是 | 熟练度，范围 0 至 100 |
| `group` | string | 是 | 分组名称，如`前端`、`后端`、`工程化` |

#### Profile

| 字段         | 类型           | 必填 | 说明                               |
| ------------ | -------------- | ---- | ---------------------------------- |
| `name`     | string         | 是   | 姓名/昵称                          |
| `avatar`   | string         | 是   | 当前可返回 emoji、首字母或图片 URL |
| `role`     | string         | 是   | 职位或身份                         |
| `location` | string         | 是   | 所在地                             |
| `motto`    | string         | 是   | 一句话签名                         |
| `bio`      | string         | 是   | 个人简介                           |
| `links`    | ProfileLink[]  | 是   | 联系方式和外链                     |
| `skills`   | ProfileSkill[] | 是   | 技能列表                           |

成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "name": "Treasure",
    "avatar": "🧑‍💻",
    "role": "前端工程师 · 折腾爱好者",
    "location": "中国 · 杭州",
    "motto": "写代码，写生活。",
    "bio": "个人简介。",
    "links": [
      {
        "id": "github",
        "label": "GitHub",
        "value": "github.com/leihenshang",
        "url": "https://github.com/leihenshang",
        "icon": "🐙"
      }
    ],
    "skills": [
      { "name": "Vue 3", "level": 95, "group": "前端" }
    ]
  }
}
```

### 9.2 获取站点信息

```http
GET /api/blog/site
```

#### SiteModule

| 字段        | 类型    | 必填 | 说明                                                              |
| ----------- | ------- | ---- | ----------------------------------------------------------------- |
| `id`      | string  | 是   | 模块 ID                                                           |
| `icon`    | string  | 是   | 当前可返回 emoji                                                  |
| `name`    | string  | 是   | 模块名称                                                          |
| `desc`    | string  | 是   | 模块描述                                                          |
| `path`    | string  | 是   | 前端站内路径，必须以`/Blog` 开头                                |
| `visible` | boolean | 是   | 是否在公开站点显示并允许前端路由访问；旧数据缺失时按`true` 迁移 |

站点模块为固定集合，后端必须使用以下 ID 与路径，不允许由管理请求修改 ID 或路径：

| ID | 模块 | 固定路径 |
| ------------- | ---- | ------------------- |
| `blog` | 文章 | `/Blog` |
| `diary` | 日记 | `/Blog/Diary` |
| `portfolio` | 作品 | `/Blog/Portfolio` |
| `tools` | 工具 | `/Blog/Tools` |
| `bookmark` | 书签 | `/Blog/Bookmark` |
| `about` | 关于 | `/Blog/About` |

`GET /api/blog/site` 应返回全部六个模块，包括 `visible: false` 的模块，并保持上表顺序。前端使用该字段控制导航、关于页模块卡片和直接 URL 路由访问。允许六个模块全部为 `false`，此时前端展示站点维护页。

`visible` 是前台展示配置，不是安全权限边界。如果关闭模块还需要禁止公开 API 获取数据，后端必须在对应公开内容接口中额外校验模块状态并返回明确错误；不能只依赖前端路由守卫。

#### SiteMilestone

| 字段 | 类型 | 必填 | 说明 |
| --------- | ------ | ---- | -------------- |
| `date` | string | 是 | `YYYY-MM-DD` |
| `title` | string | 是 | 里程碑标题 |
| `desc` | string | 是 | 说明 |

#### SiteInfo

| 字段           | 类型            | 必填 | 说明                   |
| -------------- | --------------- | ---- | ---------------------- |
| `name`       | string          | 是   | 站点名称               |
| `slogan`     | string          | 是   | 标语                   |
| `intro`      | string          | 是   | 站点介绍               |
| `techStack`  | string[]        | 是   | 技术栈                 |
| `modules`    | SiteModule[]    | 是   | 功能模块               |
| `milestones` | SiteMilestone[] | 是   | 更新记录，建议日期降序 |

成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "name": "Treasure Blog",
    "slogan": "records · thinking · life",
    "intro": "个人站点介绍。",
    "techStack": ["Vue 3", "TypeScript", "Vite", "Go"],
    "modules": [
      {
        "id": "blog",
        "icon": "📝",
        "name": "文章",
        "desc": "技术笔记与长文",
        "path": "/Blog",
        "visible": true
      },
      {
        "id": "diary",
        "icon": "📔",
        "name": "日记",
        "desc": "日常碎片与随想",
        "path": "/Blog/Diary",
        "visible": false
      },
      {
        "id": "portfolio",
        "icon": "🎨",
        "name": "作品",
        "desc": "网站、应用与开源项目",
        "path": "/Blog/Portfolio",
        "visible": true
      },
      {
        "id": "tools",
        "icon": "🧰",
        "name": "工具",
        "desc": "自研工具与常用链接",
        "path": "/Blog/Tools",
        "visible": true
      },
      {
        "id": "bookmark",
        "icon": "🔖",
        "name": "书签",
        "desc": "值得反复访问的资源",
        "path": "/Blog/Bookmark",
        "visible": true
      },
      {
        "id": "about",
        "icon": "👤",
        "name": "关于",
        "desc": "个人资料与站点记录",
        "path": "/Blog/About",
        "visible": true
      }
    ],
    "milestones": [
      {
        "date": "2026-09-01",
        "title": "日记模块上线",
        "desc": "新增日记板块与关键词搜索。"
      }
    ]
  }
}
```

## 10. 内容统计接口（建议实现）

关于页当前分别请求文章、日记和作品列表来计算数量。后端建议增加聚合接口，减少请求和无用数据传输。

```http
GET /api/blog/stats
```

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "posts": 21,
    "diaries": 21,
    "works": 6
  }
}
```

| 字段 | 类型 | 说明 |
| ----------- | ------- | -------------- |
| `posts` | integer | 已发布文章数量 |
| `diaries` | integer | 已发布日记数量 |
| `works` | integer | 已发布作品数量 |

该接口属于优化项。若暂不实现，前端仍可通过列表接口的 `pagination.total` 和作品列表长度计算。

## 11. Go 后端实现约束

### 11.1 数据可见性

所有列表、标签、分类和统计只应包含公开且已发布的数据。建议内容表至少保留：

- `status`：draft/published/archived
- `published_at` 或业务日期
- `created_at`
- `updated_at`

草稿和归档内容不得通过公开详情接口访问。

### 11.2 查询安全

- 所有 SQL 使用参数化查询
- 对 `%`、`_` 等 LIKE 特殊字符明确处理
- `sort` 使用白名单映射，禁止直接拼接用户输入
- `pageSize` 最大限制为 100
- Path ID 应限制合理长度，例如 1 至 128 字符
- URL 字段写入数据库前校验协议

### 11.3 搜索语义

为保持与当前 mock 行为一致：

- 去除关键词首尾空格
- 英文字母大小写不敏感
- 文章：搜索标题、摘要、标签、Markdown 正文
- 日记：搜索标题、摘要、标签、Markdown 正文
- 书签：搜索标题、描述、标签、URL

MySQL 可根据字符集和 collation 实现大小写不敏感搜索；数据量增大后建议使用全文索引。

### 11.4 Markdown 安全

后端返回 Markdown 原文，前端使用 Vditor 渲染。仍应将 Markdown 视为不可信输入：

- 后台保存时限制危险 URL 协议
- 禁止或清理原始 HTML 中的脚本和事件属性
- 外链统一增加 `rel="noopener noreferrer"`
- 图片代理、文件上传和 CSP 策略应独立设计

### 11.5 缓存建议

这些接口为公开只读接口，可以设置短时缓存：

```http
Cache-Control: public, max-age=60, stale-while-revalidate=300
```

分类、标签、Profile 和 Site 可缓存更久；详情接口可使用 `ETag` 或 `Last-Modified`。内容发布、更新或删除后应主动失效相关列表和统计缓存。

### 11.6 CORS

同域部署时无需额外 CORS。前后端分域部署时，只允许明确的前端 Origin，并允许 `GET`、`OPTIONS`，不应使用任意 Origin 搭配凭据。

## 12. 前端迁移说明

当前 `src/api/blog.ts`、`diary.ts`、`portfolio.ts`、`tools.ts`、`bookmark.ts`、`profile.ts` 返回裸 mock 数据。接入 Go 后端时需要新增公开 HTTP 客户端或为这些模块单独适配响应信封。

注意：现有 `src/api/index.ts` 会自动携带 `X-Token`，并在 HTTP 401 时清理登录状态、跳转 `/LogIn`。公开博客接口建议使用独立 Axios 实例，不携带 Token，也不复用 401 跳转逻辑。

分页适配关系：

| 当前 mock 返回   | 正式 API 返回                                                |
| ---------------- | ------------------------------------------------------------ |
| `res.list`     | `res.data.list` 或由 API 模块解包后继续暴露为 `res.list` |
| `res.total`    | `res.data.pagination.total`                                |
| `res.page`     | `res.data.pagination.page`                                 |
| `res.pageSize` | `res.data.pagination.pageSize`                             |

推荐在 API 模块内部完成解包，让 Vue 页面保持当前调用方式不变。

## 13. 验收清单

### 13.1 文章

- 分类、标签接口只统计已发布文章
- 分类、标签、关键词可以组合筛选
- 置顶文章始终位于普通文章之前
- 升序和降序仅影响同一置顶分组内部
- 分页无重复、无遗漏
- 列表不返回正文，详情返回完整 Markdown
- 未知 ID 返回 404

### 13.2 日记

- 日期严格返回 `YYYY-MM-DD`
- 搜索覆盖标题、摘要、标签和正文
- 置顶与分页规则和文章一致
- 日记 ID 按字符串处理

### 13.3 作品、工具与书签

- 作品分类筛选正确
- 作品详情包含链接和正文
- 工具联合类型由 `type` 正确区分
- 外部工具 ID 不允许访问自研工具详情接口
- 书签搜索覆盖完整 URL
- 所有外链均通过安全协议校验

### 13.4 关于页

- 技能 `level` 始终位于 0 至 100
- `SiteModule.path` 指向有效 `/Blog` 路由
- Site 返回六个固定模块，ID、路径和顺序符合 9.2 节定义
- `visible=false` 的模块仍包含在 Site 响应中，以便后台重新开启
- 旧数据缺少 `visible` 时迁移为 `true`
- 统计数字只包含已发布内容

## 14. 后续可选扩展

以下能力不属于当前前端的硬性依赖，可在后端基础接口稳定后扩展：

- OpenAPI 3.1 文档：`GET /api/blog/openapi.json`
- RSS/Atom：文章与日记订阅
- 归档接口：按年月聚合内容数量
- 相邻内容：详情页上一篇/下一篇
- 搜索建议与热门标签
- 作品和书签分页
- 备案号、版权年份等站点配置接口

## 15. Blog 管理 API

管理 API 用于维护公开 Blog 使用的 PostgreSQL 数据。

### 15.1 Base URL

```text
/api/blog-mgr
```

管理 API 与公开 API 使用相同响应信封，但管理路径中的资源 `id` 是数据库内部 ID，不是公开 slug/publicId。

标识对应关系：

| 资源 | 管理路径 ID | 公开 API ID |
| ---- | ------------ | --------------------------- |
| 分类 | 数据库`id` | `slug`，公开响应为 `id` |
| 标签 | 数据库`id` | 标签名称 |
| 文章 | 数据库`id` | `slug` |
| 日记 | 数据库`id` | `publicId` |
| 作品 | 数据库`id` | `slug` |
| 工具 | 数据库`id` | `slug` |
| 书签 | 数据库`id` | `publicId` |

### 15.2 认证与权限

除 CORS 预检请求外，所有管理请求必须携带：

```http
X-Token: <login-token>
```

仅以下用户类型允许访问：

| userType | 角色  |
| -------- | ----- |
| `2`    | admin |
| `100`  | root  |

认证和授权顺序为 `CORS -> Auth -> RequireAdmin -> Handler`。`OPTIONS` 预检请求由 CORS 中间件直接返回 HTTP 204。

认证相关响应：

| 场景             | HTTP | 业务码    | 说明                     |
| ---------------- | ---- | --------- | ------------------------ |
| 缺少或无效 Token | 401  | `1`     | 沿用用户模块现有认证响应 |
| 认证上下文缺失   | 401  | `40100` | 未登录或登录已失效       |
| 普通用户访问     | 403  | `40300` | 无管理权限               |

开发模式启用 `debug.enableMockLogin` 时，认证中间件使用 Mock root 用户，不检查 `X-Token`。该配置仅在服务重启后生效。

### 15.3 管理响应与错误码

普通成功响应：

```json
{
  "code": 0,
  "msg": "",
  "data": {}
}
```

创建成功使用 HTTP 201；查询、更新、删除、恢复和单例保存成功使用 HTTP 200。

| HTTP | 业务码    | 含义                   |
| ---- | --------- | ---------------------- |
| 200  | `0`     | 操作成功               |
| 201  | `0`     | 创建成功               |
| 400  | `40001` | 请求参数格式错误       |
| 404  | `40410` | 资源不存在             |
| 409  | `40900` | 版本冲突或唯一标识重复 |
| 500  | `50000` | 服务内部错误           |

错误时 `data` 固定为 `null`：

```json
{
  "code": 40900,
  "msg": "数据已变更或标识重复",
  "data": null
}
```

## 16. 管理路由总览

以下资源使用统一的 CRUD、软删除和恢复路由：

- `categories`
- `tags`
- `posts`
- `diaries`
- `portfolio-items`
- `tools`
- `bookmarks`

每种资源均注册以下 6 个端点，共 42 个：

| 方法   | Endpoint                                  | 说明                         |
| ------ | ----------------------------------------- | ---------------------------- |
| GET    | `/api/blog-mgr/{resource}`              | 分页列表                     |
| POST   | `/api/blog-mgr/{resource}`              | 创建资源                     |
| GET    | `/api/blog-mgr/{resource}/{id}`         | 资源详情，包括已软删除记录   |
| PATCH  | `/api/blog-mgr/{resource}/{id}`         | 更新资源，当前为全量覆盖语义 |
| DELETE | `/api/blog-mgr/{resource}/{id}`         | 软删除资源                   |
| POST   | `/api/blog-mgr/{resource}/{id}/restore` | 恢复资源                     |

Profile 和 Site 额外提供 4 个单例端点：

| 方法 | Endpoint                  | 说明                   |
| ---- | ------------------------- | ---------------------- |
| GET  | `/api/blog-mgr/profile` | 获取个人资料配置       |
| PUT  | `/api/blog-mgr/profile` | 创建或覆盖个人资料配置 |
| GET  | `/api/blog-mgr/site`    | 获取站点配置           |
| PUT  | `/api/blog-mgr/site`    | 创建或覆盖站点配置     |

管理端共计 46 个端点。

## 17. 管理列表

```http
GET /api/blog-mgr/{resource}
```

### 17.1 查询参数

| 参数           | 默认值      | 约束与说明                                                         |
| -------------- | ----------- | ------------------------------------------------------------------ |
| `page`       | `1`       | 大于等于 1                                                         |
| `pageSize`   | `20`      | 1 至 100                                                           |
| `keyword`    | -           | 去除首尾空格；分类、标签、工具匹配`name`，其他资源匹配 `title` |
| `status`     | -           | `draft`、`published` 或 `archived`；仅内容资源使用           |
| `deleted`    | `exclude` | `exclude`、`only` 或 `all`                                   |
| `scope`      | -           | `post`、`portfolio` 或 `bookmark`；仅分类使用                |
| `categoryId` | -           | 分类 slug；仅文章、作品、书签使用                                  |
| `sort`       | `desc`    | `asc` 或 `desc`，按 `created_at` 排序                        |

请求示例：

```http
GET /api/blog-mgr/posts?page=1&pageSize=20&status=draft&deleted=exclude&sort=desc
```

响应示例：

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "list": [],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 0,
      "orderBy": "created_at_desc"
    }
  }
}
```

管理列表和详情响应包含资源业务字段以及 `id`、`createdAt`、`updatedAt`、`deletedAt`。文章、日记和书签的列表、详情、创建和更新响应均包含 `tagIds`，用于编辑回显。

## 18. 分类与标签管理

### 18.1 分类 Category

请求字段：

| 字段 | 类型 | 必填 | 说明 |
| ------------- | ------- | ---- | ----------------------------------------------------- |
| `scope` | string | 是 | `post`、`portfolio`、`bookmark`；创建后不可修改 |
| `slug` | string | 是 | 长度 1 至 128；同一 scope 内唯一 |
| `name` | string | 是 | 展示名称，去除首尾空格后不能为空 |
| `sortOrder` | integer | 否 | 排序值，默认 0 |
| `enabled` | boolean | 否 | 是否启用，缺省为 true |

```http
POST /api/blog-mgr/categories
Content-Type: application/json
```

```json
{
  "scope": "post",
  "slug": "tech",
  "name": "技术",
  "sortOrder": 10,
  "enabled": true
}
```

修改分类 slug 时，服务会在同一事务中同步该 scope 对应内容的 `categoryId`。仍有未删除内容引用分类时，删除返回 HTTP 409。

### 18.2 标签 Tag

```json
{
  "name": "Go"
}
```

`name` 去除首尾空格后不能为空。服务端生成小写 `normalizedName`，因此标签名称按大小写不敏感语义保持唯一。删除标签不会删除关系记录；公开接口会过滤已删除标签。

## 19. 内容资源管理

### 19.1 文章 Post

| 字段              | 类型     | 说明                                        |
| ----------------- | -------- | ------------------------------------------- |
| `slug`          | string   | 公开文章 ID，长度 1 至 128                  |
| `title`         | string   | 必填                                        |
| `summary`       | string   | 摘要                                        |
| `categoryId`    | string   | 文章分类 slug，非空时必须属于`post` scope |
| `author`        | string   | 作者展示名                                  |
| `content`       | string   | Markdown 正文                               |
| `publishStatus` | string   | `draft`、`published`、`archived`      |
| `publishedOn`   | string   | `YYYY-MM-DD`                              |
| `publishedAt`   | string   | 可选，RFC 3339 时间                         |
| `pinned`        | boolean  | 是否置顶                                    |
| `version`       | integer  | 更新时必填且大于等于 1                      |
| `tagIds`        | string[] | 标签数据库 ID；去重后整体替换               |

创建示例：

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

更新使用相同字段，并携带当前版本：

```http
PATCH /api/blog-mgr/posts/123456789
```

```json
{
  "slug": "hello-world",
  "title": "Hello World（修订）",
  "summary": "第一篇文章",
  "categoryId": "essay",
  "author": "Treasure Doc",
  "content": "# Hello World\n\n更新后的正文。",
  "publishStatus": "published",
  "publishedOn": "2026-09-04",
  "pinned": true,
  "version": 1,
  "tagIds": ["10001", "10002"]
}
```

### 19.2 日记 Diary

请求字段为：`publicId`、`title`、`summary`、`content`、`mood`、`weather`、`publishStatus`、`publishedOn`、`publishedAt`、`pinned`、`version`、`tagIds`。

```json
{
  "publicId": "2026-09-04",
  "title": "九月四日",
  "summary": "今天的简短记录。",
  "content": "# 九月四日",
  "mood": "充实",
  "weather": "晴",
  "publishStatus": "published",
  "publishedOn": "2026-09-04",
  "pinned": false,
  "version": 1,
  "tagIds": []
}
```

`publicId` 是公开日记 ID，但管理详情和更新路径仍使用数据库 `id`。

### 19.3 作品 PortfolioItem

请求字段为：`slug`、`title`、`summary`、`categoryId`、`cover`、`techStack`、`links`、`content`、`publishStatus`、`publishedOn`、`publishedAt`、`version`。

```json
{
  "slug": "treasure-doc",
  "title": "Treasure Doc",
  "summary": "Markdown 文档管理系统",
  "categoryId": "website",
  "cover": "📚",
  "techStack": ["Vue 3", "Go"],
  "links": [
    {"label": "GitHub", "url": "https://github.com/leihenshang/my-treasure-doc"}
  ],
  "content": "# Treasure Doc",
  "publishStatus": "published",
  "publishedOn": "2026-09-04",
  "version": 1
}
```

`techStack` 和 `links` 以 JSONB 保存并保持数组顺序。非空 `categoryId` 必须属于 `portfolio` scope。

### 19.4 工具 Tool

请求字段为：`slug`、`kind`、`name`、`description`、`url`、`cover`、`developmentStatus`、`content`、`publishStatus`、`publishedAt`、`sortOrder`、`version`。

外链工具示例：

```json
{
  "slug": "mdn",
  "kind": "link",
  "name": "MDN Web 文档",
  "description": "Web 标准参考",
  "url": "https://developer.mozilla.org/zh-CN/",
  "publishStatus": "published",
  "sortOrder": 10,
  "version": 1
}
```

`kind=link` 时 URL 必须使用 HTTPS，服务端清空 `cover`、`developmentStatus`、`content`。`kind=own` 时 `developmentStatus` 必填，服务端清空 `url`。

### 19.5 书签 Bookmark

请求字段为：`publicId`、`title`、`url`、`description`、`categoryId`、`icon`、`publishStatus`、`publishedAt`、`sortOrder`、`version`、`tagIds`。

```json
{
  "publicId": "github",
  "title": "GitHub",
  "url": "https://github.com/",
  "description": "代码托管与开源社区",
  "categoryId": "dev",
  "icon": "🐙",
  "publishStatus": "published",
  "sortOrder": 10,
  "version": 1,
  "tagIds": []
}
```

书签 URL 必须使用 HTTPS。非空 `categoryId` 必须属于 `bookmark` scope。

## 20. Profile 与 Site 管理

### 20.1 Profile

```http
PUT /api/blog-mgr/profile
Content-Type: application/json
```

```json
{
  "name": "Treasure",
  "avatar": "🧑‍💻",
  "role": "前端工程师",
  "location": "中国 · 杭州",
  "motto": "写代码，写生活。",
  "bio": "个人简介。",
  "links": [
    {
      "id": "github",
      "label": "GitHub",
      "value": "github.com/leihenshang",
      "url": "https://github.com/leihenshang",
      "icon": "🐙"
    }
  ],
  "skills": [
    {"name": "Vue 3", "level": 95, "group": "前端"}
  ]
}
```

`name` 必填；链接 ID 不可重复，非空 URL 仅允许 HTTPS 或 `mailto:`；技能名称和分组必填，`level` 必须在 0 至 100。

### 20.2 Site

```http
PUT /api/blog-mgr/site
Content-Type: application/json
```

```json
{
  "name": "Treasure Blog",
  "slogan": "records · thinking · life",
  "intro": "个人站点介绍。",
  "techStack": ["Vue 3", "TypeScript", "Go"],
  "modules": [
    {
      "id": "blog",
      "icon": "📝",
      "name": "文章",
      "desc": "技术笔记与长文",
      "path": "/Blog",
      "visible": true
    },
    {
      "id": "diary",
      "icon": "📔",
      "name": "日记",
      "desc": "日常碎片与随想",
      "path": "/Blog/Diary",
      "visible": false
    },
    {
      "id": "portfolio",
      "icon": "🎨",
      "name": "作品",
      "desc": "网站、应用与开源项目",
      "path": "/Blog/Portfolio",
      "visible": true
    },
    {
      "id": "tools",
      "icon": "🧰",
      "name": "工具",
      "desc": "自研工具与常用链接",
      "path": "/Blog/Tools",
      "visible": true
    },
    {
      "id": "bookmark",
      "icon": "🔖",
      "name": "书签",
      "desc": "值得反复访问的资源",
      "path": "/Blog/Bookmark",
      "visible": true
    },
    {
      "id": "about",
      "icon": "👤",
      "name": "关于",
      "desc": "个人资料与站点记录",
      "path": "/Blog/About",
      "visible": true
    }
  ],
  "milestones": [
    {"date": "2026-09-04", "title": "管理 API 上线", "desc": "支持内容维护。"}
  ]
}
```

`name` 必填；里程碑日期必须为 `YYYY-MM-DD` 且标题不能为空。

Site 模块保存约束：

1. `modules` 必须包含 9.2 节定义的全部六个固定模块，每个 ID 恰好出现一次。
2. 模块 ID 和路径必须与固定映射一致；后端不得接受未知 ID、重复 ID 或被修改的路径。
3. `icon`、`name`、`desc` 和 `visible` 可由后台修改，其中 `name` 不能为空，`visible` 必须为布尔值。
4. `visible=false` 的模块必须正常持久化，管理端 GET 和公开端 GET 均不得过滤这些模块。
5. 允许全部模块为 `visible=false`。
6. 为兼容已有数据，数据库记录缺少 `visible` 时应在迁移或读取时补为 `true`；已有配置缺少固定模块时应补齐默认模块，不得将缺失项解释为隐藏。
7. `PUT /api/blog-mgr/site` 成功响应应返回标准化后的完整 Site 对象，供管理前端立即更新公开配置缓存。

Profile 和 Site 使用固定单例 key `default`。记录不存在时，Profile GET 返回结构完整的空对象且数组为 `[]`；Site GET 返回结构完整的默认对象，其中 `modules` 包含六个固定模块且 `visible` 均为 `true`，其他数组可为 `[]`。PUT 执行创建或覆盖，并恢复已软删除的单例记录。

## 21. 内容生命周期

### 21.1 发布状态和时间

- `publishStatus` 仅允许 `draft`、`published`、`archived`。
- 首次保存为 `published` 且未提供 `publishedAt` 时，使用服务器当前时间。
- 文章、日记、作品发布时若未提供 `publishedOn`，使用 `publishedAt` 对应日期。
- 未来的 `publishedAt` 表示预约发布；到达该时间前公开接口不可见。
- 草稿和归档内容不会出现在公开 API。

### 21.2 乐观锁

文章、日记、作品、工具和书签更新必须携带当前 `version >= 1`。更新 SQL 同时匹配数据库 ID 和版本号，成功后版本原子递增。

版本过期、唯一 slug/publicId 冲突或唯一标签名称冲突统一返回 HTTP 409 和业务码 `40900`。

### 21.3 标签事务

文章、日记和书签的主体更新与标签关系替换位于同一事务：

1. 去除 `tagIds` 中的首尾空格和重复值。
2. 验证所有标签存在且未删除。
3. 创建或更新主体。
4. 删除旧关系并写入请求中的完整标签集合。
5. 任一步失败则整体回滚。

### 21.4 软删除与恢复

普通资源使用 GORM 软删除。删除成功响应：

```json
{"code": 0, "msg": "", "data": {"deleted": true}}
```

恢复成功响应：

```json
{"code": 0, "msg": "", "data": {"restored": true}}
```

管理详情使用 `Unscoped()`，因此已删除资源仍可通过数据库 ID 查询。管理列表通过 `deleted=only` 查看回收站，通过 `deleted=all` 查看全部数据。

软删除记录仍占用唯一 slug/publicId/标签名称，不能用相同标识重新创建，只能恢复原记录。

## 22. 管理 API 当前实现限制

以下内容描述当前实现边界，调用方应避免依赖未实现语义：

1. `PATCH` 当前是全量覆盖，不是真正的部分更新。应提交资源全部字段；省略字段会写入零值，分类省略 `enabled` 会重置为 true。
2. 管理列表过滤器没有按资源自动裁剪。不要给不含相应列的资源传 `status`、`scope` 或 `categoryId`，否则当前实现可能返回 HTTP 500。
3. `keyword` 只搜索一个主字段，且 `%`、`_` 当前会按 SQL LIKE 通配符解释。
4. 作品 `links[].url` 当前未执行 HTTPS 协议校验；书签、外链工具和 Profile 链接已校验。
5. 创建版本化内容时请求中的 `version` 可以大于 1；调用方应省略或传 1。
6. 非版本资源更新不存在与版本资源冲突的错误语义并不完全一致；调用方应同时处理 404 和 409。
7. 禁用或删除分类不会自动隐藏引用该分类的公开内容，只会使分类不再出现在分类列表。
8. 标签精确筛选在 PostgreSQL 中通常区分大小写；关键词搜索才显式不区分大小写。
9. 当前单元测试未连接真实 PostgreSQL，JSONB、事务回滚、分类 slug 级联和并发乐观锁仍需集成测试验证。
10. 模块 `visible` 默认只控制前端入口与路由；若后端未在内容接口中校验该字段，关闭模块后其公开 API 仍可被直接调用。

## 23. 管理端验收清单

- 无 Token 返回 401，普通用户返回 403，admin/root 可以访问。
- 七类普通资源的列表、创建、详情、更新、删除、恢复均可用。
- 创建返回 HTTP 201，其余成功返回 HTTP 200。
- 内容更新使用当前 version，过期版本返回 409。
- 草稿、归档、预约发布时间未到的内容不出现在公开 API。
- 删除后公开详情返回 404，管理详情仍可读取，恢复后按状态重新决定是否公开。
- 分类 scope 校验正确，修改 slug 后引用内容同步更新。
- 标签整体替换无重复，任一步失败时主体和关系同时回滚。
- 文章、日记、书签的管理响应包含可直接用于编辑回显的 `tagIds`。
- Profile/Site GET 在空库返回空数组而不是 null，PUT 后公开接口可读取新配置。
- Site 管理 GET/PUT 和公开 GET 均返回六个固定模块及显式 `visible` 布尔值。
- 将任一模块设为 `visible=false` 后重新读取仍保持关闭，且响应中没有删除该模块。
- 全部模块可同时关闭；重新开启任一模块后配置可正常恢复。
- 旧 Site 数据缺少 `visible` 或缺少固定模块时，迁移结果默认开启并补齐，不影响升级前已有站点。
