package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request/doc"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"

	"gorm.io/gorm"
)

type DocService struct{}

var docService *DocService

var docOnce = sync.Once{}

func NewDocService() *DocService {
	docOnce.Do(func() {
		docService = &DocService{}
	})
	return docService
}

// Create 创建文档
func (doc *DocService) Create(createDoc *model.Doc, userId string) (d *model.Doc, err error) {
	if createDoc == nil {
		return nil, nil
	}

	if createDoc.GroupId != "" {
		errorMsg := fmt.Errorf("分组没有找到")
		groupList, err := getDocGroupByIds(userId, createDoc.GroupId)
		if err != nil {
			global.Log.Errorf("failed to query user group,userId:[%s], groupId:[%s],error:[%v] ",
				userId, createDoc.GroupId, err)
			return nil, errorMsg
		} else if len(groupList) == 0 {
			return nil, errorMsg
		}
	} else {
		createDoc.GroupId = global.RootGroup
	}

	if err = global.Db.Create(createDoc).Error; err != nil {
		global.Log.Errorf("failed to create doc,data:[%#v],error:%v", createDoc, err)
		return nil, errors.New("创建文档失败")
	}

	return createDoc, nil
}

// checkDocTitleIsDuplicates 检查文档标题是否重复
func checkDocTitleIsDuplicates(title string, userId string) (doc *model.Doc, err error) {
	q := global.Db.Model(&model.Doc{}).Where("title = ? AND user_id = ?", title, userId)
	if err = q.First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}

// Detail 文档详情
func (doc *DocService) Detail(id string, userId string) (d *model.Doc, err error) {
	err = global.Db.Unscoped().Where("id = ? AND user_id = ?", id, userId).First(&d).Error
	if err != nil {
		return
	}

	d.IsPin = 2
	note := &model.Note{}
	if err = global.Db.Where("doc_id = ? AND user_id = ?", id, userId).First(&note).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			d.IsPin = 1
		} else {
			global.Log.Errorf("failed to query note,data:[%#v],error:[%v]", id, err)
		}
	}

	var group *model.DocGroup
	if err = global.Db.Where("user_id = ? AND id = ?", userId, d.GroupId).First(&group).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Error(id, err)
		}
	} else {
		var parentGroups model.DocGroups
		if err = global.Db.Where("user_id = ? AND id IN (?)", userId, strings.Split(group.GroupPath, ",")).
			Order("created_at ASC").Find(&parentGroups).Error; err != nil {
			global.Log.Error(id, err)
		} else {
			d.GroupPath = parentGroups
		}
	}

	return
}

// List 文档列表
func (doc *DocService) List(r doc.ListDocRequest, userId string) (res response.ListResponse, err error) {
	q := global.Db.Model(&model.Doc{}).Where("user_id = ?", userId)
	if r.RecycleBin == 1 {
		q = q.Unscoped().Where("deleted_at is not null")
	}

	if r.GroupId != "" {
		q = q.Where("group_id = ?", r.GroupId)
	}

	if r.IsTop > 0 {
		q = q.Where("is_top = ?", r.IsTop)
	}

	if r.Keyword != "" {
		likeStr := fmt.Sprintf(`%%%s%%`, r.Keyword)
		q = q.Where("title LIKE ? OR content LIKE ?", likeStr, likeStr)
	}

	if r.Pagination.PageSize > 0 {
		q.Count(&r.Total)
	}

	if sortStr, err := r.Sort.Sort(map[string]string{"createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	}

	var list []*model.Doc
	if r.Pagination.Page > 0 && r.Pagination.PageSize > 0 {
		q = q.Limit(r.PageSize).Offset(r.Offset())
	}

	err = q.Find(&list).Error
	res.List = list
	res.Pagination = r.Pagination
	return
}

var ErrorDocIsEdited = errors.New("数据已在其他位置更新,请刷新后再试~")

// Update 文档更新
func (doc *DocService) Update(r doc.UpdateDocRequest, userId string) (newDoc *model.Doc, err error) {
	errMsg := fmt.Errorf("id 为 %s 的数据没有找到", r.Id)
	if r.Id == "" {
		global.Log.Error(errMsg)
		return nil, errMsg
	}

	tx := global.Db.Begin()
	q := tx.Unscoped().Model(&model.Doc{}).Where("id = ? AND user_id = ?", r.Id, userId)
	var oldDoc *model.Doc
	if err = q.First(&oldDoc).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Error(err)
			tx.Rollback()
			return nil, errMsg
		} else {
			global.Log.Error(errMsg)
			tx.Rollback()
			return nil, errMsg
		}
	}
	if oldDoc.Version != *r.Version {
		return nil, ErrorDocIsEdited
	}

	u := map[string]interface{}{}
	if r.Title != "" {
		u["Title"] = r.Title
	}
	if r.Content != "" {
		u["Content"] = r.Content
	}

	if r.GroupId != "" {
		u["GroupId"] = r.GroupId
	}

	if r.IsTop > 0 {
		u["IsTop"] = r.IsTop
	}

	if r.IsRecover {
		u["DeletedAt"] = nil
		u["GroupId"] = 0
	}

	if r.ReadOnly > 0 {
		u["ReadOnly"] = r.ReadOnly
	}

	if err = tx.Create(&model.DocHistory{
		BaseModel: model.BaseModel{},
		DocId:     oldDoc.Id,
		UserId:    oldDoc.UserId,
		Title:     oldDoc.Title,
		Content:   oldDoc.Content,
	}).Error; err != nil {
		tx.Rollback()
		errMsg = fmt.Errorf("保存id 为 %s 的历史数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return nil, errors.New("操作失败")
	}

	if r.IsPin == 1 {
		var dbNote *model.Note
		if err := tx.Where("user_id = ? AND doc_id = ? AND note_type = ?", userId, r.Id, model.NoteTypeDoc).First(&dbNote).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&model.Note{
				BaseModel: model.BaseModel{},
				UserId:    userId,
				DocId:     r.Id,
				NoteType:  model.NoteTypeDoc,
			}).Error; err != nil {
				tx.Rollback()
				errMsg = fmt.Errorf("保存id 为 %s 的笔记失败 %v ", r.Id, err)
				global.Log.Error(errMsg)
				return nil, errors.New("操作失败")
			}
		} else if err != nil {
			tx.Rollback()
			global.Log.Error(err)
			return nil, errors.New("操作失败")
		}
	} else if r.IsPin == 2 {
		if err := tx.Unscoped().Where("doc_id = ? AND user_id = ? AND note_type = ?", r.Id, userId, model.NoteTypeDoc).Delete(&model.Note{}).Error; err != nil {
			tx.Rollback()
			errMsg := fmt.Sprintf("删除id 为 %s 的笔记数据失败 %v ", r.Id, err)
			global.Log.Error(errMsg)
			return nil, errors.New("操作失败")
		}
	}
	oldDoc.Version++
	u["version"] = oldDoc.Version
	if err = q.Updates(u).Error; err != nil {
		errMsg = fmt.Errorf("修改id 为 %s 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		tx.Rollback()
		return nil, errors.New("操作失败")
	}

	tx.Commit()
	return oldDoc.HiddenUnnecessary(), nil
}

// Delete 文档删除
func (doc *DocService) Delete(id string, userId string) (err error) {
	if id == "" {
		errMsg := fmt.Sprintf("id 为 %s 的数据没有找到", id)
		global.Log.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.Db.Where("id = ? AND user_id = ?", id, userId)
	if err = q.Delete(&model.Doc{}).Error; err != nil {
		errMsg := fmt.Errorf("删除id 为 %s 的数据失败 %v ", id, err)
		global.Log.Error(errMsg)
		return errMsg
	}

	return
}
