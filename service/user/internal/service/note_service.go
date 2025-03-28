package service

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/note"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
)

type NoteService struct{}

var noteService *NoteService

var noteOnce = sync.Once{}

func NewNoteService() *NoteService {
	noteOnce.Do(func() {
		noteService = &NoteService{}
	})
	return noteService
}

// NoteCreate 创建笔记
func (n *NoteService) NoteCreate(r note.CreateNoteRequest, userId string) (d *model.Note, err error) {
	insertData := &model.Note{
		UserId:   userId,
		Content:  r.Content,
		IsTop:    r.IsTop,
		NoteType: r.NoteType,
		Title:    r.Title,
		Priority: r.Priority,
		Color:    r.Color,
		Icon:     r.Icon,
	}

	if err = global.Db.Create(insertData).Error; err != nil {
		global.Log.Error(r, err)
		return nil, errors.New("创建笔记失败")
	}

	return insertData, nil
}

// NoteDetail 笔记详情
func (n *NoteService) NoteDetail(r request.IDReq, userId string) (d *model.Note, err error) {
	q := global.Db.Model(&model.Note{}).Where("id = ? AND user_id = ?", r.ID, userId)
	err = q.First(&d).Error
	if err != nil {
		return
	}

	doc := &model.Doc{}
	if d.NoteType == model.NoteTypeDoc {
		if err = global.Db.Where("id = ? AND user_id = ?", d.DocId, userId).First(&doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("文档不存在")
			} else {
				global.Log.Errorf("failed to query doc:[%v],error:[%v]", d.DocId, err)
				return d, err
			}
		} else {
			d.Title = doc.Title
			d.Content = doc.Content
		}
	}
	return
}

// NoteList 笔记列表
func (n *NoteService) NoteList(r note.ListNoteRequest, userId string) (res response.ListResponse, err error) {
	q := global.Db.Model(&model.Note{}).Where("user_id = ?", userId).Where("note_type IN ?", r.NoteTypes.GetNoteTypeList())
	q.Count(&r.Total)
	r.Sort.OrderBy = "isTop_desc,priority_desc,createdAt_desc,id_asc"
	if sortStr, err := r.Sort.Sort(map[string]string{"isTop": "is_top", "priority": "priority", "createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	} else {
		global.Log.Error(r, err)
	}

	var list []*model.Note
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&list).Error
	if err != nil {
		global.Log.Error(r, err)
		return res, err
	}

	if err = FillDoc(list); err != nil {
		global.Log.Error(r, err)
		return res, err
	}

	res.List = list
	res.Pagination = r.Pagination
	return
}

func FillDoc(notes model.Notes) error {
	if len(notes) == 0 {
		return nil
	}
	var docs model.Docs
	if err := global.Db.Where("id IN (?) AND user_id = ?", notes.GetDocIds(), notes[0].UserId).Find(&docs).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if len(docs) == 0 {
		return nil
	}
	docMap := docs.ToMap()
	for _, n := range notes {
		if d, ok := docMap[n.DocId]; ok {
			n.Title = d.Title
			//n.Content = d.Content
		}
	}

	return nil
}

// Update 笔记更新
func (n *NoteService) Update(r note.UpdateNoteRequest, userId string) (d *model.Note, err error) {
	if r.Id == "" {
		return nil, nil
	}

	if err = global.Db.Where("id = ? AND user_id = ?", r.Id, userId).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("没有找到笔记")
		}
		global.Log.Error(r, err)
		return nil, errors.New("查询笔记失败")
	}

	u := map[string]interface{}{}
	if r.Title != "" {
		u["Title"] = r.Title
	}

	if r.NoteType != "" {
		u["NoteType"] = r.NoteType
	}

	u["Content"] = r.Content
	u["IsTop"] = r.IsTop
	u["Color"] = r.Color
	u["Icon"] = r.Icon
	u["Priority"] = r.Priority

	if err = global.Db.Model(&d).Updates(u).Error; err != nil {
		errMsg := fmt.Sprintf("修改id 为 %s 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return nil, errors.New("操作失败")
	}

	return
}

// NoteDelete 笔记删除
func (n *NoteService) NoteDelete(r note.UpdateNoteRequest, userId string) (err error) {
	if r.Id == "" {
		return nil
	}

	q := global.Db.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.Note{}).Error; err != nil {
		errMsg := fmt.Sprintf("删除id 为 %s 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}
