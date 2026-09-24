# Memo 速记本模块 · 功能与架构设计

> 状态：设计稿（待评审）
> 目标：在公开站新增一个「个人速记本」入口，博主登录后可以快速记录碎片内容；支持整站点级「公开/仅本人可见」总开关 + 单条速记默认私有、可单独公开。

---

## 1. 背景与目标

现有 Treasure Doc 是「公开只读博客 + 后台管理 CRUD」的架构：前台完全匿名只读，后台管理员用 `X-Token` 管理内容。缺少一个「博主本人、前台随手记」的轻量速记空间。

Memo 模块用于满足：**博主在前台登录后，快速记录灵感/待办/链接等碎片，默认自己可见，可选对访客公开**。

设计原则：完全复用现有登录鉴权与 Markdown 编辑能力，不新建任何用户身份模型，最小侵入地挂到现有前端托管与公开博客体系上。

---

## 2. 关键设计决策（已确认）

| 决策点 | 结论 | 理由 |
| --- | --- | --- |
| 登录身份 | 复用现有管理员登录（`X-Token` + 现有 `User` 模型） | 现系统只有博主一人为管理员，避免新建多用户体系 |
| 可见性模型 | 两级控制：**站点总开关**（公开可浏览 / 仅本人可见）+ **每条速记可见性**（默认私有，可单独公开） | 满足「全局控 + 默认私有、可放开」的诉求 |
| 数据形态 | 个人速记本：多条 Markdown 卡片，可置顶、检索 | 复用现有 `rook` Markdown 预览链 |

---

## 3. 功能需求

### 3.1 前端（公开站 `/Blog` 下新增 memo 入口）

1. **入口与路由**：新增首页板块/导航模块 `memo`（`path: /Blog/Memo`），沿用现有 `constants/blog-modules.ts` 模块化可见性。
2. **列表页**：
   - 未登录访客：仅能看到「站点总开关开启 + 且被标记为公开」的速记（只读）。
   - 已登录博主：额外看到自己的私有速记（带「仅自己可见」角标），并可编辑/删除。
3. **编辑**：复用现有 `MarkdownEditor` + `FloatingSaveButton` 能力；列表页提供「新建速记」。
4. **判断登录态**：读取 `useUserInfoStore` 是否已有 token（现有 `Auth()` 恢复会话逻辑），据此决定是否渲染个人操作区。
5. **站点总开关的入口**：放在后台「站点设置」页新增「Memo 模块」区块，或复用现有 site 模块配置的可见性 + 新增全局开关字段。

### 3.2 后端（公开只读 + 登录管理 两套能力）

1. **公开只读**：返回「站点总开关开启」时「已公开」的速记列表/详情，匿名可访问，不泄漏私有项。
2. **登录管理**（需已登录）：新建 / 更新 / 删除 / 置顶 / 改可见性 / 单条详情。所有操作均为当前登录博主本人，无需 `RequireAdmin` 的二次校验（登录本身就代表博主）。

---

## 4. 数据模型（新增表 `td_blog_memo`）

新建模型 `Memo`，放在 `module/blog/data/model/model.go`（沿用 `BaseModel` + `TableName()` 模式），并加入 `blogmodel.Tables()`。

| 字段 | 列名/类型 | 说明 |
| --- | --- | --- |
| `Title` | `title varchar(200) not null` | 速记标题（可留空则前端取正文前 N 字） |
| `Content` | `content text not null` | Markdown 正文 |
| `Pinned` | `pinned bool not null default false` | 置顶 |
| `Public` | `public bool not null default false` | 是否对访客公开；**默认私有** |
| `PublicAt` | `public_at datetime null` | 公开时间（可选，用于排序） |
| `Tags`（可选） | `JSON` | 预留标签，若需要 |
| `BaseModel` | 继承 | `ID`(雪花)、`CreatedAt/UpdatedAt`、软删除 |

> 采用**单表双用途**：`public=false`（私有，仅博主可见）+ `public=true`（潜在公开）→ 由「站点总开关」决定最终是否公开展示。总开关本身不加列到本表，而是放站点配置（见下文）。

### 站点总开关放哪

推荐放 `td_site` 的 JSON `settings`（现有 `draw.dom` 已支持 JSON settings 存模块可见性/里程碑）：新增 `settings.memoPublicEnabled: boolean`。

- `memoPublicEnabled = false`（默认）：公开列表**永远为空**，任何可达 `/Blog/Memo` 的访客都看不到内容提示「仅博主可见」。
- `memoPublicEnabled = true`：公开列表返回所有 `public = true` 的速记。
- **单条私有（`public=false`）在任何开关下都对访客隐藏**，仅登录博主可见。

---

## 5. 架构（后端分层）

遵循现有四模块分层，不新建独立 module 目录，复用 `module/blog` + `blog_mgr` 的能力拆分：

### 5.1 路由挂载（`module/user/router/router.go`）

| 类型 | 方法/路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| 公开只读 | `GET /api/blog/memos` | 无（匿名） | 返回对外公开列表（受总开关+public 双重约束） |
| 登录管理 | `GET /api/memo/memos` | `Auth()` | 博主本人全部速记（含私有） |
| 登录管理 | `POST /api/memo/memos` | `Auth()` | 新建（默认 `public=false`） |
| 登录管理 | `PUT /api/memo/memos/:id` | `Auth()` | 更新 |
| 登录管理 | `DELETE /api/memo/memos/:id` | `Auth()` | 删除（软删除） |
| 登录管理 | `PATCH /api/memo/memos/:id/fields` | `Auth()` | 快捷改 `pinned` / `public` / 置顶 |
| 公开只读 | `GET /api/blog/memos/:id`（公开项才 404） | 无 | 单条公开详情（可选） |

### 5.2 后端分层

- `module/blog`：新增只读 API（`GET /api/blog/memos*`），复用 `service.published()` 的可见性思路，但 memo 是「public 开关 + 总开关」而非 publish_status。
- `module/blog_mgr`：新增登录管理 Handler（`api/memo_api.go`）+ service（`internal/service/memo_service.go`）。沿用现有泛化 Handler 模式（List/Create/Update/Delete/PATCH fields）但**不走 admin 白名单**的通用资源 CRUD，因为 memo 是登录即本人。

> 说明：为了复用 `module/user` 的 `Auth()` 中间件又保持「登录即博主」语义，管理类端点挂 `/api/memo/*`，只用 `Auth()`（不叠加 `RequireAdmin`，因为登录者即博主本人）。

### 5.3 前端

- 公开路由：`src/router/root-router/index.ts` 加 `/Blog/Memo`（`meta.public:true, blogModule:'memo'`），并注册到 `constants/blog-modules.ts` 的 `BLOG_MODULE_IDS` + `DEFAULT_BLOG_MODULES` + `isBlogModuleId`。
- 页面：`views/blog/MemoList.vue`（列表/新建），复用 `MarkdownPreview` + `MarkdownEditor`。
- 数据请求：公开读用 `public-blog-client`；登录写走 `myHttp`（自动带 X-Token / 401 跳登录）。

---

## 6. 可见性规则（核心逻辑汇总）

```
访问者身份 × 站点总开关 × 单条 public → 是否可见

访客（匿名）：
  总开关 OFF → 不可见（空态文案「仅博主可见」）
  总开关 ON  + public=true → 可见
  总开关 ON  + public=false → 不可见

博主本人（已登录）：
  任何开关 + 任何 record → 全部可见（私有项带「仅自己可见」角标）
```

服务端实现处：公开列表查询 = `WHERE public = true AND (总开关开启)`；管理列表 = 全部（带 public 字段）。总开关读取逻辑集中在 memo service 一处，避免漏判。

---

## 7. 里程碑 / 实施步骤

1. **后端模型**：`module/blog/data/model/model.go` 新增 `Memo` + 加进 `Tables()`；`module/user/global/db.go` 的 `TableMigrate` 经 `blogTables()` 自动覆盖（无需手工加）。
2. **站点总开关字段**：在 `blog_mgr/internal/service/site.go` 的 site `settings` 里登记 `memoPublicEnabled`（默认 false），并同步到 `module/blog/data/response` 的默认 site 结构与前端 `SiteManageInput` 类型。
3. **公开只读 API**：`module/blog` 新增 `GET /api/blog/memos`。
4. **登录管理 API**：`module/blog_mgr` 新增 memo 的 List/Create/Update/Delete/PATCH fields + `Auth()` 中间件。
5. **前端**：
   - `blog-modules.ts` 注册 memo 模块；
   - 公开路由 + `MemoList.vue`；
   - 「站点设置」新增「Memo 模块」开关（写回 `settings.memoPublicEnabled`）；
   - 复用 `MarkdownEditor`/`FloatingSaveButton`。
6. **API 变更登记**：在 `doc/api_docs/` 的 `blog/`、`blog-mgr/`（或新增）登记上述端点文件，聚合进 `openapi.yaml` 并重跑 bundle。
7. **可见性单测**：覆盖「总开关 × public × 身份」矩阵；文档更新。

---

## 8. 验证方式

- 后端：`go build ./...` + `go test ./...`（含可见性规则用例）；启动后 `GET /api/blog/memos` 匿名验证空/公开分界。
- 前端：`pnpm exec eslint <改动文件>` + `pnpm typecheck`；(可选) `pnpm run verify`。
- 浏览器端到端：匿名见公开态、登录后见私有态、总开关切换生效。

---

## 9. 已确认的实现范围

- **标签**：`Memo` 增加 `Tags JSON` 字段，列表提供按标签筛选，用于分类与搜索。
- **归档**：按创建时间归档（列表可分组浏览），已通过索引/排序满足。
- **搜索**：按标题/内容关键词搜索（服务端 LIKE 匹配）。
- **排序**：默认 `pinned DESC + updated_at DESC`（置顶优先），支持前端自定义排序字段与方向，服务端做字段白名单校验。
- **回收站/软删除恢复**：不用。删除走 GORM 软删除，但不提供前端恢复入口（保持软删除以便潜在数据保全，但 UI 不暴露）。

- **稳定性实测补充**：新增 `Memo.Images JSON` 字段（存 `/files/blog/...` 数组，AutoMigrate 自动加列）；速记形态定为**「文字 + 图片」（朋友圈式）**，前端新建/编辑不再使用 Markdown 编辑器。
- **后台管理页**：新增 `/BlogManage/memo`（`MemoManagePage.vue`），查看全部速记（含私有），支持编辑、删除、批量删除、公开/置顶快捷开关；新增 `POST /api/memo/memos/batch-delete` 端点。

## 10. 后续可扩展（暂不实现）

- 回收站恢复入口；批量标签管理；多用户体系。
