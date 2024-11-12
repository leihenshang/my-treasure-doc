package service

import (
	"errors"
	"fmt"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"

	"gorm.io/gorm"
)

// DocCreate 创建文档
func DocCreate(r doc.CreateDocRequest, userId int64) (d *model.Doc, err error) {
	insertData := &model.Doc{
		UserId:  userId,
		Title:   r.Title,
		Content: r.Content,
		Pid:     r.Pid,
		GroupId: r.GroupId,
		IsTop:   r.IsTop,
	}

	//if existed, checkErr := checkDocTitleIsDuplicates(insertData.Title, userId); checkErr != nil {
	//	global.ZAPSUGAR.Error(r, userId, "检查文档标题失败")
	//	return nil, errors.New("检查文档标题失败")
	//} else {
	//	if existed != nil {
	//		return nil, errors.New("文档标题已存在")
	//	}
	//}

	if err = global.DB.Create(insertData).Error; err != nil {
		global.ZAPSUGAR.Error(r, err)
		return nil, errors.New("创建文档失败")
	}

	return insertData, nil
}

// checkDocTitleIsDuplicates 检查文档标题是否重复
func checkDocTitleIsDuplicates(title string, userId int64) (doc *model.Doc, err error) {
	q := global.DB.Model(&model.Doc{}).Where("title = ? AND user_id = ?", title, userId)
	if err = q.First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}

// DocDetail 文档详情
func DocDetail(r request.IDReq, userId int64) (d *model.Doc, err error) {
	q := global.DB.Model(&model.Doc{}).Where("id = ? AND user_id = ?", r.ID, userId)
	err = q.First(&d).Error
	return
}

// DocList 文档列表
func DocList(r doc.ListDocRequest, userId int64) (res response.ListResponse, err error) {
	q := global.DB.Model(&model.Doc{}).Where("user_id = ?", userId)

	if r.GroupId > 0 {
		q = q.Where("group_id = ?", r.GroupId)
	}

	if r.Pid > 0 {
		q = q.Where("pid = ?", r.Pid)
	}

	q.Count(&r.Total)
	if sortStr, err := r.ListSort.Sort(map[string]string{"createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	}

	var list []*model.Doc
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&list).Error
	res.List = list
	res.Pagination = r.ListPagination
	return
}

// DocUpdate 文档更新
func DocUpdate(r doc.UpdateDocRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.DB.Model(&model.Doc{}).Where("id = ? AND user_id = ?", r.Id, userId)
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

	if err = q.Updates(u).Error; err != nil {
		errMsg := fmt.Sprintf("修改id 为 %d 的数据失败 %v ", r.Id, err)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}

// DocDelete 文档删除
func DocDelete(r doc.UpdateDocRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.DB.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.Doc{}).Error; err != nil {
		errMsg := fmt.Sprintf("删除id 为 %d 的数据失败 %v ", r.Id, err)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}

func DocTree(r doc.ListDocRequest, userId int64) (res model.Docs, err error) {
	q := global.DB.Model(&model.Doc{}).Select("id,pid,title").Where("user_id = ?", userId).Where("pid = ?", r.Pid)
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&res).Error
	return
}
