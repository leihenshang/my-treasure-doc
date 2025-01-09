package doc

import "fastduck/treasure-doc/service/user/data/request"

// CreateDocGroupRequest 创建文档分组
type CreateDocGroupRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=250"`
	PId      string `json:"pId" binding:"required,alphanum"`
	Icon     string `json:"icon" binding:"max=250"`
	Priority int    `json:"priority" binding:""`
}

// UpdateDocGroupRequest 更新文档分组
type UpdateDocGroupRequest struct {
	Id    string `json:"id" binding:"required,alphanum"`
	Title string `json:"title" binding:"required,max=250,min=1"`
	PId   string `json:"pId" binding:"alphanum"`
	Icon  string `json:"icon" binding:"max=250"`
}

// GroupTreeRequest 文档分组树
type GroupTreeRequest struct {
	Pid     string `json:"pid" form:"pid" binding:"required,alphanum"`
	WithDoc bool   `json:"withDoc" form:"withDoc"`
}

type ListDocGroupRequest struct {
	Id string `json:"id" form:"id" binding:""`
	request.Pagination
	request.Sort
}
