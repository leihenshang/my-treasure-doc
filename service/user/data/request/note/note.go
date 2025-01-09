package note

import (
	"strings"

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
	Id       string         `json:"id" binding:"required"`
	Title    string         `json:"title" binding:"max=250"` // 标题
	Content  string         `json:"content" binding:""`      // 文档内容
	IsTop    int            `json:"isTop" binding:""`        // 是否置顶
	Color    string         `json:"color" binding:""`
	Icon     string         `json:"icon" binding:""`
	Priority int            `json:"priority" binding:""`
	NoteType model.NoteType `json:"noteType" binding:""` // 文档内容
}

type ListNoteRequest struct {
	NoteTypes ReqNoteTypes `json:"noteTypes" form:"noteTypes"`
	request.Pagination
	request.Sort
}

type ReqNoteTypes string

func (r ReqNoteTypes) GetNoteTypeList() model.NoteTypes {
	str := string(r)
	res := strings.Split(strings.TrimSpace(str), ",")
	if len(res) == 1 && res[0] == "" {
		return model.NoteTypes{
			model.NoteTypeBookmark,
			model.NoteTypeTreeNote,
			model.NoteTypeDoc,
		}
	}

	return res
}
