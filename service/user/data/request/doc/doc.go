package doc

import (
	"fastduck/treasure-doc/service/user/data/request"
)

// CreateDocRequest 创建文档
type CreateDocRequest struct {
	Title   string `json:"title" binding:"max=250"` // 标题
	Content string `json:"content" binding:""`      // 文档内容
	GroupId string `json:"groupId" binding:""`
	IsTop   int8   `json:"isTop" binding:""`
}

// UpdateDocRequest 更新文档
type UpdateDocRequest struct {
	Id        string `json:"id" binding:"required"`
	GroupId   string `json:"groupId" binding:""`
	Title     string `json:"title" binding:"max=250"` // 标题
	Content   string `json:"content" binding:""`      // 文档内容
	IsTop     int8   `json:"isTop" binding:""`        // 是否置顶
	IsRecover bool   `json:"isRecover" binding:""`
	IsPin     int8   `json:"isPin" binding:""`
	ReadOnly  int8   `json:"readOnly" binding:""`
	Version   *int   `json:"version" binding:"required"`
}

// DeleteDocRequest 更新文档
type DeleteDocRequest struct {
	Id string `json:"id" binding:"required"`
}

type ListDocRequest struct {
	GroupId    string `json:"groupId" form:"groupId" binding:""`
	IsTop      int    `json:"isTop" form:"isTop" binding:""`
	RecycleBin int    `json:"recycleBin" form:"recycleBin" binding:""`
	Keyword    string `json:"keyword" form:"keyword" binding:""`
	request.Pagination
	request.Sort
}
