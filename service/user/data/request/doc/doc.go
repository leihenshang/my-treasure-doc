package doc

import (
	"fastduck/treasure-doc/service/user/data/request"
)

// CreateDocRequest 创建文档
type CreateDocRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=250"`    // 标题
	Content string `json:"content" binding:"required,min=1,max=2500"` // 文档内容
	Pid     uint64 `json:"pid" binding:""`
	GroupId int    `json:"groupId" binding:""`
	IsTop   int    `json:"isTop" binding:""` // 是否置顶
}

// UpdateDocRequest 更新文档
type UpdateDocRequest struct {
	Id      int    `json:"id" binding:"required"`
	Pid     uint64 `json:"pid" binding:""`
	GroupId int    `json:"groupId" binding:""`
	Title   string `json:"title" binding:"max=250"`    // 标题
	Content string `json:"content" binding:"max=2500"` // 文档内容
	IsTop   int    `json:"isTop" binding:""`           // 是否置顶
}

type ListDocRequest struct {
	GroupId int `json:"groupId" binding:""`
	Pid     int `json:"pid" form:"pid" binding:""`
	request.PaginationWithSort
}
