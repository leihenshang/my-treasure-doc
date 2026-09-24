package model

import (
	"encoding/json"
	"time"

	"fastduck/treasure-doc/module/user/global/gid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"

	CategoryPost      = "post"
	CategoryPortfolio = "portfolio"
	CategoryBookmark  = "bookmark"
)

// JSON 跨方言 JSON 字段类型：PostgreSQL 下映射为 jsonb，SQLite 下映射为 json/text。
// 复用 GORM datatypes.JSON 的序列化能力，避免自定义类型在 SQLite（无 jsonb）下建出 NUMERIC 亲和列。
type JSON = datatypes.JSON

func NewJSON(value interface{}) JSON {
	data, _ := json.Marshal(value)
	return data
}

type BaseModel struct {
	ID        string         `gorm:"column:id;type:varchar(100);primaryKey"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index"`
}

func (m *BaseModel) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		m.ID = gid.GenId()
	}
	return nil
}

type Category struct {
	BaseModel
	Scope     string `gorm:"column:scope;type:varchar(20);not null;uniqueIndex:uq_blog_category_scope_slug"`
	Slug      string `gorm:"column:slug;type:varchar(128);not null;uniqueIndex:uq_blog_category_scope_slug"`
	Name      string `gorm:"column:name;type:varchar(100);not null"`
	SortOrder int    `gorm:"column:sort_order;not null;default:0;index"`
	Enabled   bool   `gorm:"column:enabled;not null;default:true"`
}

func (*Category) TableName() string { return "td_blog_category" }

type Tag struct {
	BaseModel
	Name           string `gorm:"column:name;type:varchar(100);not null"`
	NormalizedName string `gorm:"column:normalized_name;type:varchar(100);not null;uniqueIndex"`
}

func (*Tag) TableName() string { return "td_blog_tag" }

type Post struct {
	BaseModel
	Slug          string    `gorm:"column:slug;type:varchar(128);not null;uniqueIndex"`
	Title         string    `gorm:"column:title;type:varchar(200);not null"`
	Summary       string    `gorm:"column:summary;type:text;not null"`
	CategoryID    string    `gorm:"column:category_id;type:varchar(128);not null;index"`
	Author        string    `gorm:"column:author;type:varchar(100);not null"`
	Content       string    `gorm:"column:content;type:text;not null"`
	PublishStatus string    `gorm:"column:publish_status;type:varchar(16);not null;default:'draft';index:idx_blog_post_public,priority:1"`
	PublishedOn   time.Time `gorm:"column:published_on;type:date;not null;index:idx_blog_post_public,priority:3"`
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamp;not null;index"`
	Pinned        bool      `gorm:"column:pinned;not null;default:false;index:idx_blog_post_public,priority:2"`
	ViewCount     int64     `gorm:"column:view_count;not null;default:0"`
	Version       int       `gorm:"column:version;not null;default:1"`
	// 发布来源（default/siyuan…可扩展）：source_id 在该来源内的稳定标识，用于覆盖更新精确定位
	PublishSource string `gorm:"column:publish_source;type:varchar(32);not null;default:'default';index:idx_blog_post_source"`
	SourceID      string `gorm:"column:source_id;type:varchar(128);not null;default:'';index:idx_blog_post_source"`
}

func (*Post) TableName() string { return "td_blog_post" }

type PostTag struct {
	PostID    string    `gorm:"column:post_id;type:varchar(100);primaryKey"`
	TagID     string    `gorm:"column:tag_id;type:varchar(100);primaryKey;index"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
}

func (*PostTag) TableName() string { return "td_blog_post_tag" }

type Diary struct {
	BaseModel
	PublicID      string    `gorm:"column:public_id;type:varchar(128);not null;uniqueIndex"`
	Title         string    `gorm:"column:title;type:varchar(200);not null"`
	Summary       string    `gorm:"column:summary;type:text;not null"`
	Content       string    `gorm:"column:content;type:text;not null"`
	Mood          string    `gorm:"column:mood;type:varchar(50);not null"`
	Weather       string    `gorm:"column:weather;type:varchar(50);not null"`
	PublishStatus string    `gorm:"column:publish_status;type:varchar(16);not null;default:'draft';index:idx_blog_diary_public,priority:1"`
	PublishedOn   time.Time `gorm:"column:published_on;type:date;not null;index:idx_blog_diary_public,priority:3"`
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamp;not null;index"`
	Pinned        bool      `gorm:"column:pinned;not null;default:false;index:idx_blog_diary_public,priority:2"`
	ViewCount     int64     `gorm:"column:view_count;not null;default:0"`
	Version       int       `gorm:"column:version;not null;default:1"`
	// 发布来源（default/siyuan…可扩展）：source_id 在该来源内的稳定标识，用于覆盖更新精确定位
	PublishSource string `gorm:"column:publish_source;type:varchar(32);not null;default:'default';index:idx_blog_diary_source"`
	SourceID      string `gorm:"column:source_id;type:varchar(128);not null;default:'';index:idx_blog_diary_source"`
}

func (*Diary) TableName() string { return "td_blog_diary" }

type DiaryTag struct {
	DiaryID   string    `gorm:"column:diary_id;type:varchar(100);primaryKey"`
	TagID     string    `gorm:"column:tag_id;type:varchar(100);primaryKey;index"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
}

func (*DiaryTag) TableName() string { return "td_blog_diary_tag" }

type PortfolioItem struct {
	BaseModel
	Slug          string    `gorm:"column:slug;type:varchar(128);not null;uniqueIndex"`
	Title         string    `gorm:"column:title;type:varchar(200);not null"`
	Summary       string    `gorm:"column:summary;type:text;not null"`
	CategoryID    string    `gorm:"column:category_id;type:varchar(128);not null;index"`
	Cover         string    `gorm:"column:cover;type:varchar(500);not null"`
	TechStack     JSON      `gorm:"column:tech_stack;not null;default:'[]'"`
	Links         JSON      `gorm:"column:links;not null;default:'[]'"`
	Gallery       JSON      `gorm:"column:gallery;not null;default:'[]'"`
	Metrics       JSON      `gorm:"column:metrics;not null;default:'[]'"`
	DemoURL       string    `gorm:"column:demo_url;type:varchar(1000);not null;default:''"`
	RepoURL       string    `gorm:"column:repo_url;type:varchar(1000);not null;default:''"`
	Status        string    `gorm:"column:status;type:varchar(50);not null;default:''"`
	Role          string    `gorm:"column:role;type:varchar(100);not null;default:''"`
	Content       string    `gorm:"column:content;type:text;not null"`
	PublishStatus string    `gorm:"column:publish_status;type:varchar(16);not null;default:'draft';index"`
	PublishedOn   time.Time `gorm:"column:published_on;type:date;not null;index"`
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamp;not null;index"`
	ViewCount     int64     `gorm:"column:view_count;not null;default:0"`
	Version       int       `gorm:"column:version;not null;default:1"`
}

func (*PortfolioItem) TableName() string { return "td_blog_portfolio_item" }

type Tool struct {
	BaseModel
	Slug              string    `gorm:"column:slug;type:varchar(128);not null;uniqueIndex"`
	Kind              string    `gorm:"column:kind;type:varchar(10);not null;index"`
	Name              string    `gorm:"column:name;type:varchar(100);not null"`
	Description       string    `gorm:"column:description;type:text;not null"`
	URL               string    `gorm:"column:url;type:varchar(1000);not null;default:''"`
	Cover             string    `gorm:"column:cover;type:varchar(500);not null;default:''"`
	DevelopmentStatus string    `gorm:"column:development_status;type:varchar(30);not null;default:''"`
	Content           string    `gorm:"column:content;type:text;not null;default:''"`
	PublishStatus     string    `gorm:"column:publish_status;type:varchar(16);not null;default:'draft';index"`
	PublishedAt       time.Time `gorm:"column:published_at;type:timestamp;not null;index"`
	SortOrder         int       `gorm:"column:sort_order;not null;default:0;index"`
	Version           int       `gorm:"column:version;not null;default:1"`
}

func (*Tool) TableName() string { return "td_blog_tool" }

type Bookmark struct {
	BaseModel
	Title         string    `gorm:"column:title;type:varchar(200);not null"`
	URL           string    `gorm:"column:url;type:varchar(1000);not null"`
	Description   string    `gorm:"column:description;type:text;not null"`
	CategoryID    string    `gorm:"column:category_id;type:varchar(128);not null;index"`
	Icon          string    `gorm:"column:icon;type:varchar(500);not null"`
	PublishStatus string    `gorm:"column:publish_status;type:varchar(16);not null;default:'draft';index"`
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamp;not null;index"`
	SortOrder     int       `gorm:"column:sort_order;not null;default:0;index"`
	OpenInNewTab  bool      `gorm:"column:open_in_new_tab;not null;default:true"`
	Version       int       `gorm:"column:version;not null;default:1"`
}

func (*Bookmark) TableName() string { return "td_blog_bookmark" }

type BookmarkTag struct {
	BookmarkID string    `gorm:"column:bookmark_id;type:varchar(100);primaryKey"`
	TagID      string    `gorm:"column:tag_id;type:varchar(100);primaryKey;index"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;not null"`
}

func (*BookmarkTag) TableName() string { return "td_blog_bookmark_tag" }

type Profile struct {
	BaseModel
	ProfileKey string `gorm:"column:profile_key;type:varchar(50);not null;uniqueIndex"`
	Name       string `gorm:"column:name;type:varchar(100);not null"`
	Avatar     string `gorm:"column:avatar;type:varchar(500);not null"`
	Role       string `gorm:"column:role;type:varchar(100);not null"`
	Location   string `gorm:"column:location;type:varchar(100);not null"`
	Motto      string `gorm:"column:motto;type:varchar(200);not null"`
	Bio        string `gorm:"column:bio;type:text;not null"`
	Links      JSON   `gorm:"column:links;not null;default:'[]'"`
	Skills     JSON   `gorm:"column:skills;not null;default:'[]'"`
}

func (*Profile) TableName() string { return "td_blog_profile" }

type Site struct {
	BaseModel
	SiteKey         string `gorm:"column:site_key;type:varchar(50);not null;uniqueIndex"`
	Name            string `gorm:"column:name;type:varchar(100);not null"`
	Slogan          string `gorm:"column:slogan;type:varchar(200);not null"`
	Intro           string `gorm:"column:intro;type:text;not null"`
	TechStack       JSON   `gorm:"column:tech_stack;not null;default:'[]'"`
	Modules         JSON   `gorm:"column:modules;not null;default:'[]'"`
	Milestones      JSON   `gorm:"column:milestones;not null;default:'[]'"`
	Home            JSON   `gorm:"column:home;not null;default:'{}'"`
	Footer          JSON   `gorm:"column:footer;not null;default:'{}'"`
	Banner          JSON   `gorm:"column:banner;not null;default:'{}'"`
	MaintenanceMode bool   `gorm:"column:maintenance_mode;not null;default:false"`
	// MemoPublicEnabled 站点总开关：为 true 时公开 memo（public=true）才对访客可见。
	MemoPublicEnabled bool `gorm:"column:memo_public_enabled;not null;default:false"`
}

func (*Site) TableName() string { return "td_blog_site" }

// Memo 前台速记本：博主登录后记录碎片内容，默认私有，可单独公开。
// public=false 仅博主可见；public=true 在「站点总开关 memoPublicEnabled」开启时对访客公开。
type Memo struct {
	BaseModel
	Title     string    `gorm:"column:title;type:varchar(200);not null;default:''"`
	Content   string    `gorm:"column:content;type:text;not null"`
	Images    JSON      `gorm:"column:images;not null;default:'[]'"`
	Weather   string    `gorm:"column:weather;type:varchar(10);not null;default:''"`
	Mood      string    `gorm:"column:mood;type:varchar(10);not null;default:''"`
	Pinned    bool      `gorm:"column:pinned;not null;default:false;index"`
	Public    bool      `gorm:"column:public;not null;default:false;index"`
	PublicAt  time.Time `gorm:"column:public_at;type:timestamp;index"`
	Tags      JSON      `gorm:"column:tags;not null;default:'[]'"`
	SortOrder int       `gorm:"column:sort_order;not null;default:0;index"`
	Version   int       `gorm:"column:version;not null;default:1"`
}

func (*Memo) TableName() string { return "td_blog_memo" }

// MemoTag 速记标签关联（复用 td_blog_tag 的标签字典，与文章/日记一致）。
type MemoTag struct {
	MemoID    string    `gorm:"column:memo_id;type:varchar(100);primaryKey"`
	TagID     string    `gorm:"column:tag_id;type:varchar(100);primaryKey;index"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
}

func (*MemoTag) TableName() string { return "td_blog_memo_tag" }

// EditHistory 文章/日记的编辑历史：每次保存记录一条完整快照，支持查看与恢复。
// 每个 ref（resource+ref_id）保留最近 maxEditHistory 条。
type EditHistory struct {
	BaseModel
	Resource string `gorm:"column:resource;type:varchar(20);not null;index:idx_edit_history_ref"`
	RefID    string `gorm:"column:ref_id;type:varchar(100);not null;index:idx_edit_history_ref"`
	Seq      int    `gorm:"column:seq;not null"`
	Version  int    `gorm:"column:version;not null;default:0"`
	Summary  string `gorm:"column:summary;type:varchar(200);not null;default:''"`
	// Snapshot 完整实体快照（含 categoryID/tagIds 等，与详情返回一致）。
	Snapshot JSON `gorm:"column:snapshot;not null"`
}

func (*EditHistory) TableName() string { return "td_edit_history" }

func Tables() []interface{} {
	return []interface{}{
		&Category{}, &Tag{}, &Post{}, &PostTag{}, &Diary{}, &DiaryTag{},
		&PortfolioItem{}, &Tool{}, &Bookmark{}, &BookmarkTag{}, &Profile{}, &Site{},
		&Memo{}, &MemoTag{}, &Media{}, &VisitorLog{}, &EditHistory{},
	}
}
