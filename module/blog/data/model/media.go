package model

import "time"

// Media 上传媒体（/files/blog 下文件）的元数据目录。
// 文件本体仍以 sha256 命名存储在 files/blog/，这里只登记元信息与上传原始名，
// 用于列表检索、引用标记与（未来的）按原文件名/类型搜索。删除文件时同步删除该行。
type Media struct {
	// Name 存储文件名（<sha256>.<ext>），作为主键。
	Name string `gorm:"column:name;type:varchar(300);primaryKey"`
	// Size 文件字节数。
	Size int64 `gorm:"column:size;not null;default:0"`
	// Ext 存储后缀，含点，如 .png。
	Ext string `gorm:"column:ext;type:varchar(20);not null;default:''"`
	// Mime 内容嗅探得到的 MIME 类型。
	Mime string `gorm:"column:mime;type:varchar(80);not null;default:''"`
	// Sha256 内容摘要（不含后缀）。
	Sha256 string `gorm:"column:sha256;type:varchar(64);not null;default:''"`
	// OriginName 上传时的原始文件名，便于识图。
	OriginName string `gorm:"column:origin_name;type:varchar(255);not null;default:''"`
	// Width/Height 图片像素尺寸（尽力而为，非图片/不支持格式为 0）。
	Width  int `gorm:"column:width;not null;default:0"`
	Height int `gorm:"column:height;not null;default:0"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;autoUpdateTime"`
}

func (*Media) TableName() string { return "td_blog_media" }
