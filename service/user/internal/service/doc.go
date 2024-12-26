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

	if createDoc.GroupId != "" && createDoc.GroupId != global.RootGroup {
		errorMsg := fmt.Errorf("分组没有找到")
		groupList, err := getDocGroupByIds(userId, createDoc.GroupId)
		if err != nil {
			global.Log.Error(err)
			return nil, errorMsg
		} else if len(groupList) == 0 {
			return nil, errorMsg
		}
	} else {
		createDoc.GroupId = global.RootGroup
	}

	if err = global.Db.Create(createDoc).Error; err != nil {
		global.Log.Error(err)
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

	note := &model.Note{}
	if err = global.Db.Where("doc_id = ? AND user_id = ?", id, userId).First(&note).Error; err != nil {
		d.IsPin = 2
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Errorf("failed to query note,data:[%#v],error:[%v]", id, err)
		}
	} else {
		d.IsPin = 1
	}

	var group *model.DocGroup
	if err = global.Db.Where("user_id = ? AND id = ?", userId, d.GroupId).First(&group).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Error(id, err)
			return nil, err
		} else {
			return d, nil
		}
	}

	var parentGroups model.DocGroups
	if err = global.Db.Where("user_id = ? AND id IN (?)", userId, strings.Split(group.GroupPath, ",")).
		Order("created_at ASC").Find(&parentGroups).Error; err != nil {
		global.Log.Error(id, err)
	} else {
		d.GroupPath = parentGroups
	}

	return d, nil
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
	if r.Id == "" {
		return nil, nil
	}

	tx := global.Db.Begin()
	q := tx.Unscoped().Debug().Model(&model.Doc{}).Where("id = ? AND user_id = ?", r.Id, userId).
		Where("version = ?", r.Version)
	var dbDoc *model.Doc
	if err = q.First(&dbDoc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorDocIsEdited
		} else {
			errMsg := fmt.Errorf("id 为 %s 的数据没有找到", r.Id)
			global.Log.Error(err)
			tx.Rollback()
			return nil, errMsg
		}
	}

	if result := q.Updates(setDocUpdateData(r)); result.Error != nil {
		tx.Rollback()
		global.Log.Error(result.Error)
		return nil, fmt.Errorf("修改id 为 %s 的数据失败", r.Id)
	} else if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, ErrorDocIsEdited
	}

	if err = handleDocExtraData(tx, dbDoc, r); err != nil {
		return nil, err
	}

	tx.Commit()

	dbDoc.Version++
	return dbDoc.HiddenData(), nil
}

func setDocUpdateData(r doc.UpdateDocRequest) map[string]any {
	updateData := make(map[string]any)
	if r.Title != "" {
		updateData["Title"] = r.Title
	}
	if r.Content != "" {
		updateData["Content"] = r.Content
	}

	if r.GroupId != "" {
		updateData["GroupId"] = r.GroupId
	}

	if r.IsTop > 0 {
		updateData["IsTop"] = r.IsTop
	}

	if r.ReadOnly > 0 {
		updateData["ReadOnly"] = r.ReadOnly
	}

	updateData["version"] = gorm.Expr("version + ?", 1)
	return updateData
}

func handleDocExtraData(tx *gorm.DB, dbDoc *model.Doc, r doc.UpdateDocRequest) (err error) {
	if err = tx.Create(&model.DocHistory{
		BaseModel: model.BaseModel{},
		DocId:     dbDoc.Id,
		UserId:    dbDoc.UserId,
		Title:     dbDoc.Title,
		Content:   dbDoc.Content,
	}).Error; err != nil {
		tx.Rollback()
		global.Log.Error(dbDoc, err)
		return fmt.Errorf("保存id 为 %s 的历史数据失败", dbDoc.Id)
	}

	if r.IsPin == 1 {
		var dbNote *model.Note
		if err = tx.Where("user_id = ? AND doc_id = ? AND note_type = ?", dbDoc.UserId, dbDoc.Id, model.NoteTypeDoc).
			First(&dbNote).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&model.Note{
				BaseModel: model.BaseModel{},
				UserId:    dbDoc.UserId,
				DocId:     dbDoc.Id,
				NoteType:  model.NoteTypeDoc,
			}).Error; err != nil {
				tx.Rollback()
				global.Log.Error(err)
				return errors.New("保存文档笔记失败")
			}
		} else if err != nil {
			tx.Rollback()
			global.Log.Error(err)
			return errors.New("获取文档笔记失败")
		}
	}
	if r.IsPin == 2 {
		if err = tx.Unscoped().Where("doc_id = ? AND user_id = ? AND note_type = ?", dbDoc.Id, dbDoc.UserId, model.NoteTypeDoc).
			Delete(&model.Note{}).Error; err != nil {
			tx.Rollback()
			global.Log.Error(err)
			return errors.New("删除文档笔记失败")
		}
	}
	return nil
}

// Delete 文档删除
func (doc *DocService) Delete(id string, userId string) (err error) {
	if id == "" {
		return nil
	}

	q := global.Db.Where("id = ? AND user_id = ?", id, userId)
	if err = q.Delete(&model.Doc{}).Error; err != nil {
		global.Log.Error(err)
		return fmt.Errorf("删除id 为 %s 的数据失败 %v ", id, err)
	}

	return
}

// Recover 文档恢复
func (doc *DocService) Recover(id string, userId string) (err error) {
	if id == "" {
		return nil
	}

	result := global.Db.Unscoped().Model(&model.Doc{}).
		Where("id = ? AND user_id = ?", id, userId).Update("deleted_at", nil)
	if err = result.Error; err != nil {
		global.Log.Error(err)
		return fmt.Errorf("恢复id 为 %s 的数据失败 %v ", id, err)
	}

	return
}
