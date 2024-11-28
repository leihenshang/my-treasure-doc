package service

import (
	"errors"
	"fmt"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
)

func DocHistoryRecover(r request.IDReq, userId int64) (err error) {
	history, err := DocHistoryDetail(r, userId)
	if err != nil {
		return err
	} else if history == nil {
		return errors.New("文档历史没有找到")
	}

	dbDoc, err := DocDetail(request.IDReq{ID: history.DocId}, userId)
	if err != nil {
		return err
	} else if dbDoc == nil {
		return errors.New("文档没有找到")
	}

	tx := global.DB.Begin()
	if err = tx.Create(&model.DocHistory{
		BasicModel: model.BasicModel{},
		DocId:      dbDoc.Id,
		UserId:     dbDoc.UserId,
		Title:      dbDoc.Title,
		Content:    dbDoc.Content,
	}).Error; err != nil {
		tx.Rollback()
		errMsg := fmt.Errorf("保存id 为 %d 的历史数据失败 %v ", r.ID, err)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New("操作失败")
	}

	if err = tx.Unscoped().Model(&model.Doc{}).Where("id = ? AND user_id = ?", history.DocId, userId).Updates(map[string]any{"Content": history.Content}).Error; err != nil {
		errMsg := fmt.Errorf("修改id 为 %d 的数据失败 %v ", r.ID, err)
		global.ZAPSUGAR.Error(errMsg)
		tx.Rollback()
		return errors.New("操作失败")
	}
	tx.Commit()

	return nil
}

// DocHistoryDetail 文档历史详情
func DocHistoryDetail(r request.IDReq, userId int64) (d *model.DocHistory, err error) {
	q := global.DB.Unscoped().Model(&model.DocHistory{}).Where("id = ? AND user_id = ?", r.ID, userId)
	err = q.First(&d).Error
	return
}

// DocHistoryList 文档列表
func DocHistoryList(r doc.ListDocHistoryRequest, userId int64) (res response.ListResponse, err error) {
	q := global.DB.Model(&model.DocHistory{}).Where("user_id = ?", userId)
	global.ZAPSUGAR.Infof(`requet:%+v`, r)

	if r.DocId > 0 {
		q = q.Where("doc_id = ?", r.DocId)
	}

	if r.ListPagination.PageSize > 0 {
		q.Count(&r.Total)
	}

	if sortStr, err := r.ListSort.Sort(map[string]string{"createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	}

	var list model.DocHistories
	if r.ListPagination.Page > 0 && r.ListPagination.PageSize > 0 {
		q = q.Limit(r.PageSize).Offset(r.Offset())
	}

	err = q.Find(&list).Error
	res.List = list
	res.Pagination = r.ListPagination
	return
}
