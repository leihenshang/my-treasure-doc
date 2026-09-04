package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"fastduck/treasure-doc/module/user/global/gid"

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

type JSON []byte

func NewJSON(value interface{}) JSON {
	data, _ := json.Marshal(value)
	return data
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return []byte("[]"), nil
	}
	return []byte(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("[]")
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSON", value)
	}
	*j = append((*j)[:0], bytes...)
	return nil
}

type BaseModel struct {
	ID        string         `gorm:"column:id;type:varchar(100);primaryKey"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index"`
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
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamptz;not null;index"`
	Pinned        bool      `gorm:"column:pinned;not null;default:false;index:idx_blog_post_public,priority:2"`
	Version       int       `gorm:"column:version;not null;default:1"`
}

func (*Post) TableName() string { return "td_blog_post" }

type PostTag struct {
	PostID    string    `gorm:"column:post_id;type:varchar(100);primaryKey"`
	TagID     string    `gorm:"column:tag_id;type:varchar(100);primaryKey;index"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null"`
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
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamptz;not null;index"`
	Pinned        bool      `gorm:"column:pinned;not null;default:false;index:idx_blog_diary_public,priority:2"`
	Version       int       `gorm:"column:version;not null;default:1"`
}

func (*Diary) TableName() string { return "td_blog_diary" }

type DiaryTag struct {
	DiaryID   string    `gorm:"column:diary_id;type:varchar(100);primaryKey"`
	TagID     string    `gorm:"column:tag_id;type:varchar(100);primaryKey;index"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null"`
}

func (*DiaryTag) TableName() string { return "td_blog_diary_tag" }

type PortfolioItem struct {
	BaseModel
	Slug          string    `gorm:"column:slug;type:varchar(128);not null;uniqueIndex"`
	Title         string    `gorm:"column:title;type:varchar(200);not null"`
	Summary       string    `gorm:"column:summary;type:text;not null"`
	CategoryID    string    `gorm:"column:category_id;type:varchar(128);not null;index"`
	Cover         string    `gorm:"column:cover;type:varchar(500);not null"`
	TechStack     JSON      `gorm:"column:tech_stack;type:jsonb;not null;default:'[]'"`
	Links         JSON      `gorm:"column:links;type:jsonb;not null;default:'[]'"`
	Content       string    `gorm:"column:content;type:text;not null"`
	PublishStatus string    `gorm:"column:publish_status;type:varchar(16);not null;default:'draft';index"`
	PublishedOn   time.Time `gorm:"column:published_on;type:date;not null;index"`
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamptz;not null;index"`
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
	PublishedAt       time.Time `gorm:"column:published_at;type:timestamptz;not null;index"`
	SortOrder         int       `gorm:"column:sort_order;not null;default:0;index"`
	Version           int       `gorm:"column:version;not null;default:1"`
}

func (*Tool) TableName() string { return "td_blog_tool" }

type Bookmark struct {
	BaseModel
	PublicID      string    `gorm:"column:public_id;type:varchar(128);not null;uniqueIndex"`
	Title         string    `gorm:"column:title;type:varchar(200);not null"`
	URL           string    `gorm:"column:url;type:varchar(1000);not null"`
	Description   string    `gorm:"column:description;type:text;not null"`
	CategoryID    string    `gorm:"column:category_id;type:varchar(128);not null;index"`
	Icon          string    `gorm:"column:icon;type:varchar(500);not null"`
	PublishStatus string    `gorm:"column:publish_status;type:varchar(16);not null;default:'draft';index"`
	PublishedAt   time.Time `gorm:"column:published_at;type:timestamptz;not null;index"`
	SortOrder     int       `gorm:"column:sort_order;not null;default:0;index"`
	Version       int       `gorm:"column:version;not null;default:1"`
}

func (*Bookmark) TableName() string { return "td_blog_bookmark" }

type BookmarkTag struct {
	BookmarkID string    `gorm:"column:bookmark_id;type:varchar(100);primaryKey"`
	TagID      string    `gorm:"column:tag_id;type:varchar(100);primaryKey;index"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamptz;not null"`
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
	Links      JSON   `gorm:"column:links;type:jsonb;not null;default:'[]'"`
	Skills     JSON   `gorm:"column:skills;type:jsonb;not null;default:'[]'"`
}

func (*Profile) TableName() string { return "td_blog_profile" }

type Site struct {
	BaseModel
	SiteKey    string `gorm:"column:site_key;type:varchar(50);not null;uniqueIndex"`
	Name       string `gorm:"column:name;type:varchar(100);not null"`
	Slogan     string `gorm:"column:slogan;type:varchar(200);not null"`
	Intro      string `gorm:"column:intro;type:text;not null"`
	TechStack  JSON   `gorm:"column:tech_stack;type:jsonb;not null;default:'[]'"`
	Modules    JSON   `gorm:"column:modules;type:jsonb;not null;default:'[]'"`
	Milestones JSON   `gorm:"column:milestones;type:jsonb;not null;default:'[]'"`
}

func (*Site) TableName() string { return "td_blog_site" }

func Tables() []interface{} {
	return []interface{}{
		&Category{}, &Tag{}, &Post{}, &PostTag{}, &Diary{}, &DiaryTag{},
		&PortfolioItem{}, &Tool{}, &Bookmark{}, &BookmarkTag{}, &Profile{}, &Site{},
	}
}
