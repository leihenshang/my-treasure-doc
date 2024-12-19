package service

import (
	"errors"
	"fmt"
	"sync"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
)

type DocHistoryService struct{}

var docHistoryService *DocHistoryService

var docHistoryOnce = sync.Once{}

func NewDocHistoryService() *DocHistoryService {
	docHistoryOnce.Do(func() {
		docHistoryService = &DocHistoryService{}
	})
	return docHistoryService
}

func (dh *DocHistoryService) DocHistoryRecover(r request.IDReq, userId int64) (err error) {
	history, err := docHistoryService.DocHistoryDetail(r, userId)
	if err != nil {
		return err
	} else if history == nil {
		return errors.New("文档历史没有找到")
	}

	dbDoc, err := docService.DocDetail(request.IDReq{ID: history.DocId}, userId)
	if err != nil {
		return err
	} else if dbDoc == nil {
		return errors.New("文档没有找到")
	}

	tx := global.Db.Begin()
	if err = tx.Create(&model.DocHistory{
		BasicModel: model.BasicModel{},
		DocId:      dbDoc.Id,
		UserId:     dbDoc.UserId,
		Title:      dbDoc.Title,
		Content:    dbDoc.Content,
	}).Error; err != nil {
		tx.Rollback()
		errMsg := fmt.Errorf("保存id 为 %d 的历史数据失败 %v ", r.ID, err)
		global.Log.Error(errMsg)
		return errors.New("操作失败")
	}

	if err = tx.Unscoped().Model(&model.Doc{}).Where("id = ? AND user_id = ?", history.DocId, userId).Updates(map[string]any{"Content": history.Content}).Error; err != nil {
		errMsg := fmt.Errorf("修改id 为 %d 的数据失败 %v ", r.ID, err)
		global.Log.Error(errMsg)
		tx.Rollback()
		return errors.New("操作失败")
	}
	tx.Commit()

	return nil
}

// DocHistoryDetail 文档历史详情
func (dh *DocHistoryService) DocHistoryDetail(r request.IDReq, userId int64) (d *model.DocHistory, err error) {
	q := global.Db.Unscoped().Model(&model.DocHistory{}).Where("id = ? AND user_id = ?", r.ID, userId)
	err = q.First(&d).Error
	return
}

// DocHistoryList 文档列表
func (dh *DocHistoryService) DocHistoryList(r doc.ListDocHistoryRequest, userId int64) (res response.ListResponse, err error) {
	q := global.Db.Model(&model.DocHistory{}).Where("user_id = ?", userId)
	global.Log.Infof(`requet:%+v`, r)

	if r.DocId > 0 {
		q = q.Where("doc_id = ?", r.DocId)
	}

	if r.Pagination.PageSize > 0 {
		q.Count(&r.Total)
	}

	if sortStr, err := r.Sort.Sort(map[string]string{"createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	}

	var list model.DocHistories
	if r.Pagination.Page > 0 && r.Pagination.PageSize > 0 {
		q = q.Limit(r.PageSize).Offset(r.Offset())
	}

	err = q.Find(&list).Error
	res.List = list
	res.Pagination = r.Pagination
	return
}
