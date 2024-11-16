package note

import (
	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
)

// CreateNoteRequest 创建文档
type CreateNoteRequest struct {
	Title    string         `json:"title" binding:"min=1,max=250"` // 标题
	Content  string         `json:"content" binding:"required"`    // 文档内容
	NoteType model.NoteType `json:"noteType" binding:"required"`   // 文档内容
	IsTop    int            `json:"isTop" binding:""`              // 是否置顶
}

// UpdateNoteRequest 更新文档
type UpdateNoteRequest struct {
	Id      int    `json:"id" binding:"required"`
	Title   string `json:"title" binding:"min=1,max=250"` // 标题
	Content string `json:"content" binding:""`            // 文档内容
	IsTop   int    `json:"isTop" binding:""`              // 是否置顶
}

type ListNoteRequest struct {
	request.ListPagination
	request.ListSort
}
