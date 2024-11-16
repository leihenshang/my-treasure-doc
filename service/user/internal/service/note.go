package service

import (
	"errors"
	"fmt"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/note"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
)

// NoteCreate 创建笔记
func NoteCreate(r note.CreateNoteRequest, userId int64) (d *model.Note, err error) {
	insertData := &model.Note{
		UserId:   userId,
		Content:  r.Content,
		IsTop:    r.IsTop,
		NoteType: r.NoteType,
		Title:    r.Title,
	}

	if err = global.DB.Create(insertData).Error; err != nil {
		global.ZAPSUGAR.Error(r, err)
		return nil, errors.New("创建笔记失败")
	}

	return insertData, nil
}

// NoteDetail 笔记详情
func NoteDetail(r request.IDReq, userId int64) (d *model.Note, err error) {
	q := global.DB.Model(&model.Note{}).Where("id = ? AND user_id = ?", r.ID, userId)
	err = q.First(&d).Error
	return
}

// NoteList 笔记列表
func NoteList(r note.ListNoteRequest, userId int64) (res response.ListResponse, err error) {
	q := global.DB.Model(&model.Note{}).Where("user_id = ?", userId)
	q.Count(&r.Total)
	if sortStr, err := r.ListSort.Sort(map[string]string{"createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	}

	var list []*model.Note
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&list).Error
	res.List = list
	res.Pagination = r.ListPagination
	return
}

// NoteUpdate 笔记更新
func NoteUpdate(r note.UpdateNoteRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.DB.Model(&model.Note{}).Where("id = ? AND user_id = ?", r.Id, userId)
	u := map[string]interface{}{}
	if r.Content != "" {
		u["Content"] = r.Content
	}
	if r.Title != "" {
		u["Title"] = r.Title
	}

	if err = q.Updates(u).Error; err != nil {
		errMsg := fmt.Sprintf("修改id 为 %d 的数据失败 %v ", r.Id, err)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}

// NoteDelete 笔记删除
func NoteDelete(r note.UpdateNoteRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.DB.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.Note{}).Error; err != nil {
		errMsg := fmt.Sprintf("删除id 为 %d 的数据失败 %v ", r.Id, err)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}
