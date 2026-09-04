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
}

type Post struct {
	PostSummary
	Content string `json:"content"`
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
	Date      string   `json:"date"`
}

type PortfolioItem struct {
	PortfolioSummary
	Links   []PortfolioLink `json:"links"`
	Content string          `json:"content"`
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
	ID   string `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Path string `json:"path"`
}

type SiteMilestone struct {
	Date  string `json:"date"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type Site struct {
	Name       string          `json:"name"`
	Slogan     string          `json:"slogan"`
	Intro      string          `json:"intro"`
	TechStack  []string        `json:"techStack"`
	Modules    []SiteModule    `json:"modules"`
	Milestones []SiteMilestone `json:"milestones"`
}

type Stats struct {
	Posts   int64 `json:"posts"`
	Diaries int64 `json:"diaries"`
	Works   int64 `json:"works"`
}
