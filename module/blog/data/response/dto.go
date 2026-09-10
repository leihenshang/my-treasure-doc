package response

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PostSummary struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Author   string   `json:"author"`
	Date     string   `json:"date"`
	Pinned   bool     `json:"pinned"`
	Views    int64    `json:"views"`
}

type Post struct {
	PostSummary
	Content string `json:"content"`
	// Prev/Next 为按发布时间相邻的文章（Prev 更早、Next 更新）
	Prev *PostSummary `json:"prev,omitempty"`
	Next *PostSummary `json:"next,omitempty"`
	// Related 为同分类或同标签的推荐文章
	Related []PostSummary `json:"related"`
}

// ArchiveGroup 是归档页按月份聚合的文章。
type ArchiveGroup struct {
	Month string        `json:"month"`
	Posts []PostSummary `json:"posts"`
}

type DiarySummary struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
	Date    string   `json:"date"`
	Mood    string   `json:"mood"`
	Weather string   `json:"weather"`
	Pinned  bool     `json:"pinned"`
	Views   int64    `json:"views"`
}

type Diary struct {
	DiarySummary
	Content string `json:"content"`
}

type PortfolioLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type PortfolioSummary struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Category  string   `json:"category"`
	Cover     string   `json:"cover"`
	TechStack []string `json:"techStack"`
	Status    string   `json:"status"`
	Date      string   `json:"date"`
	Views     int64    `json:"views"`
}

type PortfolioItem struct {
	PortfolioSummary
	Links   []PortfolioLink `json:"links"`
	// Gallery 为项目截图/效果图，DemoURL 与 RepoURL 为在线演示与仓库地址
	Gallery []string `json:"gallery"`
	DemoURL string   `json:"demoUrl"`
	RepoURL string   `json:"repoUrl"`
	// Role 为本人角色，Metrics 为成果指标（如「10k+ 用户」）
	Role    string   `json:"role"`
	Metrics []string `json:"metrics"`
	Content string   `json:"content"`
}

type Tool struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	URL     string `json:"url,omitempty"`
	Cover   string `json:"cover,omitempty"`
	Status  string `json:"status,omitempty"`
	Content string `json:"content,omitempty"`
}

type Bookmark struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	URL      string   `json:"url"`
	Desc     string   `json:"desc"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Icon     string   `json:"icon"`
}

type ProfileLink struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
	URL   string `json:"url,omitempty"`
	Icon  string `json:"icon"`
}

type ProfileSkill struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Group string `json:"group"`
}

type Profile struct {
	Name     string         `json:"name"`
	Avatar   string         `json:"avatar"`
	Role     string         `json:"role"`
	Location string         `json:"location"`
	Motto    string         `json:"motto"`
	Bio      string         `json:"bio"`
	Links    []ProfileLink  `json:"links"`
	Skills   []ProfileSkill `json:"skills"`
}

type SiteModule struct {
	ID      string `json:"id"`
	Icon    string `json:"icon"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Path    string `json:"path"`
	Marker  string `json:"marker"`
	Visible bool   `json:"visible"`
}

type SiteMilestone struct {
	Date  string `json:"date"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type SiteHomeAI struct {
	Eyebrow     string `json:"eyebrow"`
	Title       string `json:"title"`
	Description string `json:"description"`
	LinkText    string `json:"linkText"`
	LinkURL     string `json:"linkUrl"`
	ImageURL    string `json:"imageUrl"`
}

// SiteHomeTerminal 描述主页右侧终端卡片（无 AI 模块配图时展示）。
type SiteHomeTerminal struct {
	// Title 是终端窗口标题，例如 creative-workspace
	Title string `json:"title"`
	// Command 是命令行内容，例如 imagine "一个值得被创造的世界"
	Command string `json:"command"`
	// Lines 是命令行下方的输出行（每条自动加 ✓ 前缀）
	Lines []string `json:"lines"`
}

type SiteHome struct {
	Title             string           `json:"title"`
	Subtitle          string           `json:"subtitle"`
	AI                SiteHomeAI       `json:"ai"`
	Terminal          SiteHomeTerminal `json:"terminal"`
	PortfolioImageURL string           `json:"portfolioImageUrl"`
	BookmarkImageURL  string           `json:"bookmarkImageUrl"`
}

type SiteFooter struct {
	Text         string `json:"text"`
	LinkText     string `json:"linkText"`
	LinkURL      string `json:"linkUrl"`
	ICPNumber    string `json:"icpNumber"`
	ICPURL       string `json:"icpUrl"`
	PoliceNumber string `json:"policeNumber"`
	PoliceURL    string `json:"policeUrl"`
}

type SiteBanner struct {
	Enabled         bool   `json:"enabled"`
	Text            string `json:"text"`
	BackgroundColor string `json:"backgroundColor"`
	TextColor       string `json:"textColor"`
}

type Site struct {
	Name       string          `json:"name"`
	Slogan     string          `json:"slogan"`
	Intro      string          `json:"intro"`
	TechStack  []string        `json:"techStack"`
	Modules    []SiteModule    `json:"modules"`
	Milestones []SiteMilestone `json:"milestones"`
	Home            SiteHome        `json:"home"`
	Footer          SiteFooter      `json:"footer"`
	Banner          SiteBanner      `json:"banner"`
	MaintenanceMode bool            `json:"maintenanceMode"`
}

type Stats struct {
	Posts   int64 `json:"posts"`
	Diaries int64 `json:"diaries"`
	Works   int64 `json:"works"`
}
