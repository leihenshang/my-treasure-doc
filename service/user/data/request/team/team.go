package team

// CreateOrUpdateTeamRequest 更新文档
type CreateOrUpdateTeamRequest struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"max=250"` // 标题
}
