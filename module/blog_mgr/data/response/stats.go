package response

// StatusCounts 按发布状态统计的资源数量。
type StatusCounts struct {
	Total     int64 `json:"total"`
	Draft     int64 `json:"draft"`
	Published int64 `json:"published"`
	Archived  int64 `json:"archived"`
}

// Stats 是后台仪表盘的总览数据。
type Stats struct {
	Posts      StatusCounts `json:"posts"`
	Diaries    StatusCounts `json:"diaries"`
	Portfolio  StatusCounts `json:"portfolioItems"`
	Tools      StatusCounts `json:"tools"`
	Bookmarks  StatusCounts `json:"bookmarks"`
	Categories int64        `json:"categories"`
	Tags       int64        `json:"tags"`
	TotalViews int64        `json:"totalViews"`
}
