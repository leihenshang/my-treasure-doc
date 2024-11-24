package doc

import (
	"fastduck/treasure-doc/service/user/data/request"
)

// CreateDocRequest 创建文档
type CreateDocRequest struct {
	Title   string `json:"title" binding:"max=250"` // 标题
	Content string `json:"content" binding:""`      // 文档内容
	Pid     int64  `json:"pid" binding:""`
	GroupId int64  `json:"groupId" binding:""`
	IsTop   int    `json:"isTop" binding:""`
}

// UpdateDocRequest 更新文档
type UpdateDocRequest struct {
	Id      int    `json:"id" binding:"required"`
	Pid     int64  `json:"pid" binding:""`
	GroupId int    `json:"groupId" binding:""`
	Title   string `json:"title" binding:"max=250"` // 标题
	Content string `json:"content" binding:""`      // 文档内容
	IsTop   int    `json:"isTop" binding:""`        // 是否置顶
}

type ListDocRequest struct {
	GroupId    int `json:"groupId" form:"groupId" binding:""`
	Pid        int `json:"pid" form:"pid" binding:""`
	IsTop      int `json:"isTop" form:"isTop" binding:""`
	RecycleBin int `json:"recycleBin" form:"recycleBin" binding:""`
	request.ListPagination
	request.ListSort
}
