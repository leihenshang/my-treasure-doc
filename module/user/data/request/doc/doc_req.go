package doc

import (
	"fastduck/treasure-doc/module/user/data/request"
)

// CreateDocRequest 创建文档
type CreateDocRequest struct {
	Title   string `json:"title" binding:"required,max=250,min=2"`
	Content string `json:"content" binding:""`
	GroupId string `json:"groupId" binding:"required,alphanum"`
	IsTop   int8   `json:"isTop" binding:""`
}

// UpdateDocRequest 更新文档
type UpdateDocRequest struct {
	Id       string `json:"id" binding:"required,alphanum"`
	GroupId  string `json:"groupId" binding:""`
	Title    string `json:"title" binding:"max=260"`
	Content  string `json:"content" binding:""`
	IsTop    int8   `json:"isTop" binding:""`
	IsPin    int8   `json:"isPin" binding:""`
	ReadOnly int8   `json:"readOnly" binding:""`
	Version  int    `json:"version" binding:"number"`
}

type ListDocRequest struct {
	GroupId    string `json:"groupId" form:"groupId" binding:""`
	IsTop      int    `json:"isTop" form:"isTop" binding:""`
	RecycleBin int    `json:"recycleBin" form:"recycleBin" binding:""`
	Keyword    string `json:"keyword" form:"keyword" binding:""`
	request.Pagination
	request.Sort
}
