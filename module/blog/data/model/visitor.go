package model

// VisitorCategory 是访客访问内容的分类标识，由 Path 解析而来，用于后台统计区分内容类型。
type VisitorCategory string

const (
	VisitorCategoryPost        VisitorCategory = "post"
	VisitorCategoryDiary       VisitorCategory = "diary"
	VisitorCategoryPortfolio   VisitorCategory = "portfolio"
	VisitorCategoryTool        VisitorCategory = "tool"
	VisitorCategoryBookmark    VisitorCategory = "bookmark"
	VisitorCategorySite        VisitorCategory = "site"
	VisitorCategoryUnresolved  VisitorCategory = "other"
)

// VisitorLog 记录公开博客的一次访客访问（来源 IP + 目标路径 + 内容分类），
// 供后台概览统计时间范围内的访问 IP 与次数。逐请求插入，
// 个人博客流量下开销可忽略；记录失败不阻断公开接口响应。
type VisitorLog struct {
	BaseModel
	IP        string          `gorm:"column:ip;type:varchar(64);not null;index"`
	Path      string          `gorm:"column:path;type:varchar(500);not null"`
	Category  VisitorCategory `gorm:"column:category;type:varchar(32);not null;default:'other';index"`
	UserAgent string          `gorm:"column:user_agent;type:varchar(500);not null;default:''"`
}

func (*VisitorLog) TableName() string { return "td_blog_visitor_log" }
