package team

import "fastduck/treasure-doc/service/user/gid"

// CreateOrUpdateTeamRequest 更新文档
type CreateOrUpdateTeamRequest struct {
	Id   gid.Gid `json:"id" binding:"required"`
	Name string  `json:"name" binding:"max=250"` // 标题
}
