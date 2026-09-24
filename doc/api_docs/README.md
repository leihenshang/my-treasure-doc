# Treasure Doc API 规格（OpenAPI 3.1）

本目录是后端 HTTP 接口的**唯一权威契约**，采用「**一个端点一个文件**」的 subdocuments 组织方式，与实际注册的路由一一对应：

- `module/blog/router/router.go`（公开只读 `/api/blog/*` + `robots.txt`/`sitemap.xml`/`rss.xml`）
- `module/blog_mgr/router/router.go`（后台管理 `/api/blog-mgr/*`，含资源 CRUD、统计、设置、上传/媒体、备份、机器令牌、编辑历史）
- `module/user/router/router.go`（`registerAPI` 里的鉴权 `/api/user/*` 与 `/ping`；上传/媒体/备份 Handler 也在此注册，但在 `blog_mgr/router/router.go` 挂载路由）

> 想了解某个接口的最快方式：看下面「模块一览」找到对应子目录 → 打开端点文件（首行注释即「方法 + 完整路径 + 模块 + 注册位置」）→ 顺着 `$ref` 跳转 `components/` 拿到参数与字段。嫌逐文件太慢，就 bundle 后用 Redoc 整页浏览（见「渲染浏览」）。

> **给 AI / 代码生成器（对接方）：** 直接读 [openapi.bundled.yaml](openapi.bundled.yaml) 这一个文件即可——所有外部 `$ref` 已展开，93 条路径、119 个接口、全部字段集中在单文件内，无需跳转。它由下面的 `bundle` 命令从拆分源生成，改 API 后重跑覆盖。

---

## 目录结构

```
api_docs/
├── openapi.yaml                     # 聚合根：info / servers / tags / security + 全部 paths（$ref）
├── openapi.bundled.yaml             # 【生成产物】外部 $ref 已展开的单文件，AI/代码生成器对接入口
├── README.md                        # 本文件
├── components/                      # 跨端点共享（用相对路径被端点文件引用）
│   ├── schemas.yaml                 #   模型、请求/响应 DTO、分页、错误信封
│   ├── responses.yaml               #   具名失败响应（按错误码聚合，避免重复）
│   └── parameters.yaml              #   分页/排序/筛选等共享 query/path 参数
├── blog/        (22)            # 公开只读　/api/blog/*　＋ 站点根文件
├── blog-mgr/    (80)            # 后台管理　/api/blog-mgr/*
│   ├── categories/  tags/  posts/  diaries/  portfolio-items/  tools/  bookmarks/
│   │                            #   7 个资源，各 8 个端点（list/create/batch-delete/
│   │                            #   detail/update/update-fields/delete/restore）
│   └── （其余 24 个一层文件：统计/设置/上传/媒体/备份/令牌/历史）
├── memo/        (7)             # 前台速记本　/api/memo/*（登录即博主本人）
├── publish/     (12)            # 机器令牌发布　/api/publish/*
├── user/        (5)             # 鉴权　/api/user/*　＋　/ping
└── backup/      (1)             # NAS 导出　/api/backup/export
```

端点文件共 **127** 个（`blog 22` + `blog-mgr 80` + `memo 7` + `publish 12` + `user 5` + `backup 1`）。

---

## 模块一览（模块 → 位置 → 鉴权）

| 模块 | 目录 | 路径前缀 | 数量 | 注册位置 | 鉴权 | OpenAPI tag |
| --- | --- | --- | --- | --- | --- | --- |
| 公开只读博客 | `blog/` | `/api/blog/*` | 18 | `module/blog/router/router.go` | 匿名（`security: [{}]`） | `public-blog` |
| 站点根文件 | `blog/get-{robots,sitemap,rss}.yaml` | `/robots.txt` `/sitemap.xml` `/rss.xml` | 3 | `RegisterSiteFiles` | 匿名（非 JSON 信封） | `site-files` |
| 内容资源 CRUD | `blog-mgr/<资源>/` | `/api/blog-mgr/<资源>/*` | 56 | `blog_mgr/router/router.go`（泛化循环注册） | `X-Token` + 管理员 | `blog-manage` |
| 统计 / 设置 | `blog-mgr/` 一层 | `/api/blog-mgr/{stats,visitor-stats,profile,site}` | 6 | 同上 | `X-Token` + 管理员 | `blog-manage` |
| 上传 / 媒体库 | `blog-mgr/` 一层 | `/api/blog-mgr/uploads/*`、`/api/blog-mgr/medias/*` | 6 | `blog_mgr/router/router.go`（Handler 在 user 模块） | `X-Token` + 管理员 | `blog-manage` |
| 备份 | `blog-mgr/` 一层 | `/api/blog-mgr/backups*` | 7 | 同上 | `X-Token` + 管理员 | `blog-manage` |
| 机器令牌（发布/备份） | `blog-mgr/` 一层 | `/api/blog-mgr/publish/token`、`/api/blog-mgr/backups/token` | 3 | 同上 | `X-Token` + 管理员 | `blog-manage` |
| 编辑历史 | `blog-mgr/` 一层 | `/api/blog-mgr/history/*` | 3 | 同上 | `X-Token` + 管理员 | `blog-manage` |
| 机器令牌发布 | `publish/` | `/api/publish/*` | 12 | `blog_mgr/router/router.go` 的 `RegisterPublishRoutes` | `X-Publish-Token`（可叠加 IP 白名单） | `publish` |
| 用户鉴权 | `user/` | `/api/user/*` 与 `/ping` | 5 | `module/user/router/router.go`（`registerAPI`） | `login`/`captcha`/`ping` 匿名；`logout`/`change-pwd` 需 `X-Token` | `user` |
| 前台速记本 | `memo/` | `/api/memo/*` | 7 | `module/user/router/router.go`（`registerAPI` → `blogmgrrouter.RegisterMemo`） | `X-Token`（登录即博主本人，不要求管理员） | `memo-manage` |
| NAS 备份导出 | `backup/nas-export.yaml` | `/api/backup/export` | 1 | `module/user/router/router.go` | `X-Backup-Token` | `backup` |

> 鉴权细节：`X-Token` 由登录签发，后台接口还要求 `userType ∈ {2, 100}`；`X-Publish-Token` / `X-Backup-Token` 为后台配置的机器令牌（`/api/publish/*` 创建即发布）。三种 scheme 定义在 `openapi.yaml` 的 `components.securitySchemes`。

---

## 端点文件命名约定

- 单个嵌套资源（categories / tags / posts / diaries / portfolio-items / tools / bookmarks）在 `blog-mgr/<资源>/` 下，固定 8 个文件，`operationId` 为 `get{Resource}List / create{Resource} / ...`（如 `getManagePostsList`）。
- 一层文件用 `动词-名词` 命名，一眼看出用途：`list-*`（列表）、`get-*`（详情/读取）、`create-*`/`update-*`/`delete-*`/`restore-*`、`upload-*`、`batch-delete-*`、`force-update-*`（发布侧强制覆盖）、`lookup-*`（按 slug 查）。
- 每个端点文件首行注释固定为：`# <HTTP 方法> <完整路径> —— <模块>（<一句话说明>）`；第二行 `# 注册位置：...`。**看路径与归属先看这两行注释。**

---

## 如何查找接口（给前端开发者）

三种找法，按需选择：

1. **按功能/路径找文件**：先看上方「模块一览」确定属于哪个前缀，再进对应目录。文件命名已含动词与路径片段（如要「文章的删除」就是 `blog-mgr/posts/delete.yaml`）。
2. **看全量路径与返回结构**：读 `openapi.yaml` 的 `paths` 段（按注释块分模块分组），或 bundle 后用 Redoc 整页搜索。
3. **前端接入参考**：前端仓库 `my-treasure-doc-front` 的实际调用位于 `src/api/*.ts` 与 `src/api/blog-manage/*.ts`，可作为「这个端点前端怎么调」的落点对照。

### 怎么读懂一个端点文件

以一个端点文件为例，它包含三部分信息：

- **首行注释** → 方法、完整路径、模块；
- **`parameters`** → 路径/查询参数，通常是 `$ref` 指向 `components/parameters.yaml`（如 `PublicPage` / `ManageID`）。要看到参数约束（`minimum`、`enum`、默认值、`required`），需点进 `components/parameters.yaml`。
- **`responses`** → 成功响应的 `data` 结构用 `allOf: [Envelope, {data: $ref: '...#/<实体>'}]` 组合，实体定义在 `components/schemas.yaml`；失败响应是具名 `$ref`（`components/responses.yaml`）。

> **注意**：由于共享参数、共享实体都以 `$ref` 拆分，单独打开一个端点文件**看不到完整的参数字段列表**，必须顺着跳 `components/parameters.yaml` 与 `components/schemas.yaml`。这正是「渲染浏览」更省力的原因。

---

## 渲染浏览（整页可检索）

这是 subdocuments 风格，不能直接被 Swagger UI 加载，需先打包：

```bash
# 打包成单文件（入库为 AI/对接方入口，改 API 后重跑覆盖）
npx --yes @redocly/cli@latest bundle doc/api_docs/openapi.yaml -o doc/api_docs/openapi.bundled.yaml

# 本地预览（Redoc 渲染，可搜索接口）
cd doc/api_docs
npx --yes @redocly/cli@latest preview-docs doc/api_docs/openapi.yaml
```

验证规格有效性：

```bash
cd doc/api_docs
npx --yes @redocly/cli@latest lint openapi.yaml      # 应无 error（operation-4xx-response 警告是固有意见）
```

---

## 拿到一个接口后，需要看齐的全局约定

- **统一响应信封**：除 `/ping` 与站点根三个文件外，都返回 `{code, msg, data}`；业务成功 `code: 0`，失败 `code: 非 0`。
- **两种失败约定（务必区分）**：
  - `module/blog` 与 `module/blog_mgr` 走真实 HTTP 状态码（400/404/409/500）；
  - `module/user` 与上传/媒体/备份类 Handler **失败也返回 HTTP 200**，仅靠 `code: 1` + `data: {}` 区分（端点文件里其 `data` 声明为「成功结构 或 `EmptyData`」的 `oneOf`）。
- **乐观锁**：内容资源更新必须回传当前 `version`，过期返回 409 / code 40900。
- **分页**：列表接口的查询参数与分页实体统一在 `components/parameters.yaml` 与 `components/schemas.yaml` 的 `Pagination`。

---

## 变更登记（改 API 时同步规格）

新增 / 修改 / 删除 API，都要同步本目录与 `openapi.yaml`（人读文档 `doc/blog-api.md` 可能滞后，以本目录为准）：

1. **新增端点**：在对应模块目录新建端点文件（内容是 operation 对象，**不含** `openapi:`/`paths:` 包装；首行注释「方法 + 路径 + 模块 + 注册位置」），并在 `openapi.yaml` 的 `paths` 加对应 `$ref`。两者缺一不可。
2. **修改端点**：路径 / 方法 / 参数 / 请求体 / data 结构 / 错误码 / 鉴权任一变化，改对应端点文件；新增字段优先进 `components/schemas.yaml` 并以 `$ref` 复用。
3. **删除端点**：删端点文件，并移除 `openapi.yaml` 里对应引用。
4. **提交前自检**：上面「渲染浏览」的 `lint` + `bundle` 跑通（bundle 会覆盖入库的 `openapi.bundled.yaml`，它作为 AI/对接方入口应一并提交）；`openapi.yaml` 的 operation 数、`operationId` 唯一性、端点文件总数三者一致（当前 119）。改完记得同步更新本 README 的「模块一览」与计数。