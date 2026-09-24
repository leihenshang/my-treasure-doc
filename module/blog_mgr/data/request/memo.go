package request

// Memo 是前台速记本的管理入参：登录即博主本人，默认私有，可单独公开。
type Memo struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Images    []string `json:"images"`
	Pinned    bool     `json:"pinned"`
	Public    bool     `json:"public"`
	Tags      []string `json:"tags"`
	SortOrder int      `json:"sortOrder"`
	Version   int      `json:"version"`
}
