package service

import (
	"errors"
	"fmt"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	resp "fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"

	"gorm.io/gorm"
)

// DocGroupCreate 创建文档分组
func DocGroupCreate(r doc.CreateDocGroupRequest, userId int64) (dg *model.DocGroup, err error) {
	insertData := &model.DocGroup{
		Title:    r.Title,
		Icon:     r.Icon,
		PId:      r.PId,
		UserId:   userId,
		Priority: r.Priority,
	}

	if existed, err := checkDocGroupTitleRepeat(insertData.Title, userId); err != nil {
		global.ZAPSUGAR.Error(r, userId, err)
		return nil, errors.New("检查文档分组标题失败")
	} else if existed != nil {
		return nil, errors.New("文档分组标题已存在")
	}

	if err = global.DB.Create(insertData).Error; err != nil {
		global.ZAPSUGAR.Error(r, err)
		return nil, errors.New("创建文档分组失败")
	}

	return insertData, nil
}

// checkDocGroupTitleRepeat 查询数据库检查文档分组标题是否重复
func checkDocGroupTitleRepeat(title string, userId int64) (dg *model.DocGroup, err error) {
	q := global.DB.Model(&model.DocGroup{}).Where("title = ? AND user_id = ?", title, userId)
	if err = q.First(&dg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}

// DocGroupList 文档分组列表
func DocGroupList(r request.ListPagination, userId int64) (res resp.ListResponse, err error) {
	var list []model.DocGroup
	q := global.DB.Model(&model.DocGroup{}).Where("user_id = ?", userId)
	q.Count(&r.Total)
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&list).Error
	res.List = list
	res.Pagination = r
	return
}

// DocGroupUpdate 文档分组更新
func DocGroupUpdate(r doc.UpdateDocGroupRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		return errors.New(errMsg)
	}

	q := global.DB.Model(&model.DocGroup{}).Where("id = ? AND user_id = ?", r.Id, userId)
	u := map[string]interface{}{"Title": r.Title, "PId": r.PId, "Icon": r.Icon}
	if err = q.Updates(u).Error; err != nil {
		errMsg := fmt.Sprintf("修改id 为 %d 的数据失败 %v ", r.Id, err)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New(errMsg)
	}

	return
}

// DocGroupDelete 文档分组删除u
func DocGroupDelete(r doc.UpdateDocGroupRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New(errMsg)
	}

	tx := global.DB.Begin()
	q := tx.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.DocGroup{}).Error; err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("删除id 为 %d 的数据失败 %v ", r.Id, err))
	}
	if err = q.Where("group_id = ? AND user_id = ?", r.Id, userId).Delete(&model.Doc{}).Error; err != nil {
		tx.Rollback()
		errMsg := fmt.Sprintf("删除id 为 %d 的分组下的文档失败 %v ", r.Id, err)
		return errors.New(errMsg)
	}

	tx.Commit()

	return
}

func DocGroupTree(r doc.GroupTreeRequest, userId int64) (docTree resp.DocTrees, err error) {
	docTree = make([]*resp.DocTree, 0)
	var list model.DocGroups
	q := global.DB.Where("user_id = ?", userId).Where("p_id = ?", r.Pid)
	if err = q.Find(&list).Error; err != nil {
		global.ZAPSUGAR.Error(err)
		return nil, errors.New("查询分组信息失败")
	}
	for _, v := range list {
		docTree = append(docTree, &resp.DocTree{
			DocGroup: v,
		})
	}
	if !r.WithChildren {
		return
	}

	groupDocs, err := fillGroupDoc(r.Pid)
	if err != nil {
		global.ZAPSUGAR.Error(err)
		return nil, err
	}
	docTree = append(docTree, groupDocs...)
	return
}

func fillGroupDoc(pid int64) (docTree resp.DocTrees, err error) {
	docs, err := getDocByGroupIds(pid)
	if err != nil {
		return nil, err
	}
	for _, d := range docs {
		docTree = append(docTree, &resp.DocTree{
			DocGroup: &model.DocGroup{
				BasicModel: model.BasicModel{
					Id: d.Id,
				},
				Title:     d.Title,
				Priority:  d.Priority,
				GroupType: model.GroupTypeDoc,
			},
		})
	}
	return
}

func getDocByGroupIds(groupId ...int64) (res model.Docs, err error) {
	err = global.DB.Where("group_id IN (?)", groupId).Find(&res).Error
	return
}

func getDocGroupByIds(groupId ...int64) (res model.DocGroups, err error) {
	err = global.DB.Where("id IN (?)", groupId).Find(&res).Error
	return
}
