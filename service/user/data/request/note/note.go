package note

import (
	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
)

// CreateNoteRequest 创建文档
type CreateNoteRequest struct {
	Title    string         `json:"title" binding:"max=250"`     // 标题
	Content  string         `json:"content" binding:"required"`  // 文档内容
	NoteType model.NoteType `json:"noteType" binding:"required"` // 文档内容
	IsTop    int            `json:"isTop" binding:""`            // 是否置顶
	Color    string         `json:"color" binding:""`
	Icon     string         `json:"icon" binding:""`
	Priority int            `json:"priority" binding:""`
}

// UpdateNoteRequest 更新文档
type UpdateNoteRequest struct {
	Id       int            `json:"id" binding:"required"`
	Title    string         `json:"title" binding:"max=250"` // 标题
	Content  string         `json:"content" binding:""`      // 文档内容
	IsTop    int            `json:"isTop" binding:""`        // 是否置顶
	Color    string         `json:"color" binding:""`
	Icon     string         `json:"icon" binding:""`
	Priority int            `json:"priority" binding:""`
	NoteType model.NoteType `json:"noteType" binding:"required"` // 文档内容
}

type ListNoteRequest struct {
	TreeHole bool `json:"treeHole"`
	request.ListPagination
	request.ListSort
}

func (l ListNoteRequest) GetNoteTypeList() []model.NoteType {
	if l.TreeHole {
		return []model.NoteType{
			model.NoteTypeBookmark,
			model.NoteTypeTreeHole,
			model.NoteTypeTreeNote,
		}
	}
	return []model.NoteType{
		model.NoteTypeBookmark,
		model.NoteTypeTreeNote,
	}
}
