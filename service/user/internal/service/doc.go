package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
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

// DocCreate 创建文档
func (doc *DocService) DocCreate(r doc.CreateDocRequest, userId int64) (d *model.Doc, err error) {
	insertData := &model.Doc{
		UserId:  userId,
		Title:   r.Title,
		Content: r.Content,
		GroupId: r.GroupId,
		IsTop:   r.IsTop,
	}

	//if existed, checkErr := checkDocTitleIsDuplicates(insertData.Title, userId); checkErr != nil {
	//	global.Log.Error(r, userId, "检查文档标题失败")
	//	return nil, errors.New("检查文档标题失败")
	//} else {
	//	if existed != nil {
	//		return nil, errors.New("文档标题已存在")
	//	}
	//}

	if insertData.GroupId > 0 {
		groupList, err := getDocGroupByIds(insertData.GroupId)
		if err != nil {
			return nil, err
		} else if len(groupList) == 0 {
			return nil, errors.New("分组没有找到")
		}
	}

	if err = global.Db.Create(insertData).Error; err != nil {
		global.Log.Error(r, err)
		return nil, errors.New("创建文档失败")
	}

	return insertData, nil
}

// checkDocTitleIsDuplicates 检查文档标题是否重复
func checkDocTitleIsDuplicates(title string, userId int64) (doc *model.Doc, err error) {
	q := global.Db.Model(&model.Doc{}).Where("title = ? AND user_id = ?", title, userId)
	if err = q.First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}

// DocDetail 文档详情
func (doc *DocService) DocDetail(r request.IDReq, userId int64) (d *model.Doc, err error) {
	err = global.Db.Unscoped().Where("id = ? AND user_id = ?", r.ID, userId).First(&d).Error
	if err != nil {
		return
	}

	d.IsPin = 1
	note := &model.Note{}
	if err := global.Db.Where("doc_id = ? AND user_id = ?", r.ID, userId).First(&note).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			d.IsPin = 2
		}
	}

	var group *model.DocGroup
	if err := global.Db.Where("user_id = ? AND id = ?", userId, d.GroupId).First(&group).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Error(r, err)
			return nil, err
		}
	} else {
		var parentGroups model.DocGroups
		if err := global.Db.Where("user_id = ? AND id IN (?)", userId, strings.Split(group.GroupPath, ",")).Order("created_at ASC").Find(&parentGroups).Error; err != nil {
			global.Log.Error(r, err)
			return nil, err
		} else {
			d.GroupPath = parentGroups
		}
	}

	return
}

// DocList 文档列表
func (doc *DocService) DocList(r doc.ListDocRequest, userId int64) (res response.ListResponse, err error) {
	q := global.Db.Model(&model.Doc{}).Where("user_id = ?", userId)
	if r.RecycleBin == 1 {
		q = q.Unscoped().Where("deleted_at is not null")
	}
	if r.GroupId > 0 {
		q = q.Where("group_id = ?", r.GroupId)
	}

	if r.Pid > 0 {
		q = q.Where("pid = ?", r.Pid)
	}

	if r.IsTop > 0 {
		q = q.Where("is_top = ?", r.IsTop)
	}

	if r.Keyword != "" {
		likeStr := fmt.Sprintf(`%%%s%%`, r.Keyword)
		q = q.Where("title LIKE ? OR content LIKE ?", likeStr, likeStr)
	}

	if r.ListPagination.PageSize > 0 {
		q.Count(&r.Total)
	}

	if sortStr, err := r.ListSort.Sort(map[string]string{"createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	}

	var list []*model.Doc

	if r.ListPagination.Page > 0 && r.ListPagination.PageSize > 0 {
		q = q.Limit(r.PageSize).Offset(r.Offset())
	}

	err = q.Find(&list).Error
	res.List = list
	res.Pagination = r.ListPagination
	return
}

func fillGroupPath(docs model.Docs) model.Docs {
	docs.GetGroupIds(true)
	return docs
}

// DocUpdate 文档更新
func (doc *DocService) DocUpdate(r doc.UpdateDocRequest, userId int64) (err error) {
	errMsg := fmt.Errorf("id 为 %d 的数据没有找到", r.Id)
	if r.Id <= 0 {
		global.Log.Error(errMsg)
		return errMsg
	}

	tx := global.Db.Begin()
	q := tx.Unscoped().Model(&model.Doc{}).Where("id = ? AND user_id = ?", r.Id, userId)
	var oldDoc *model.Doc
	if err = q.First(&oldDoc).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Error(err)
			tx.Rollback()
			return errMsg
		} else {
			global.Log.Error(errMsg)
			tx.Rollback()
			return errMsg
		}
	}

	u := map[string]interface{}{}
	if r.Title != "" {
		u["Title"] = r.Title
	}
	if r.Content != "" {
		u["Content"] = r.Content
	}
	if r.Pid > 0 {
		u["Pid"] = r.Pid
	}

	if r.GroupId > 0 {
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
		BasicModel: model.BasicModel{},
		DocId:      oldDoc.Id,
		UserId:     oldDoc.UserId,
		Title:      oldDoc.Title,
		Content:    oldDoc.Content,
	}).Error; err != nil {
		tx.Rollback()
		errMsg = fmt.Errorf("保存id 为 %d 的历史数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return errors.New("操作失败")
	}

	if r.IsPin == 1 {
		var dbNote *model.Note
		if err := tx.Where("user_id = ? AND doc_id = ? AND note_type = ?", userId, r.Id, model.NoteTypeDoc).First(&dbNote).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&model.Note{
				BasicModel: model.BasicModel{},
				UserId:     userId,
				DocId:      r.Id,
				NoteType:   model.NoteTypeDoc,
			}).Error; err != nil {
				tx.Rollback()
				errMsg = fmt.Errorf("保存id 为 %d 的笔记失败 %v ", r.Id, err)
				global.Log.Error(errMsg)
				return errors.New("操作失败")
			}
		} else if err != nil {
			tx.Rollback()
			global.Log.Error(err)
			return errors.New("操作失败")
		}
	} else if r.IsPin == 2 {
		if err := tx.Unscoped().Where("doc_id = ? AND user_id = ? AND note_type = ?", r.Id, userId, model.NoteTypeDoc).Delete(&model.Note{}).Error; err != nil {
			tx.Rollback()
			errMsg := fmt.Sprintf("删除id 为 %d 的笔记数据失败 %v ", r.Id, err)
			global.Log.Error(errMsg)
			return errors.New("操作失败")
		}
	}

	if err = q.Updates(u).Error; err != nil {
		errMsg = fmt.Errorf("修改id 为 %d 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		tx.Rollback()
		return errors.New("操作失败")
	}

	tx.Commit()
	return
}

// DocDelete 文档删除
func (doc *DocService) DocDelete(r doc.UpdateDocRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.Log.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.Db.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.Doc{}).Error; err != nil {
		errMsg := fmt.Sprintf("删除id 为 %d 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}

func (doc *DocService) DocTree(r doc.ListDocRequest, userId int64) (res model.Docs, err error) {
	q := global.Db.Model(&model.Doc{}).Select("id,pid,title").Where("user_id = ?", userId).Where("pid = ?", r.Pid)
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&res).Error
	return
}
