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
func DocGroupCreate(r doc.CreateDocGroupRequest, userId uint64) (dg *model.DocGroup, err error) {
	insertData := &model.DocGroup{
		Title:    r.Title,
		Icon:     r.Icon,
		PId:      uint64(r.PId),
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
func checkDocGroupTitleRepeat(title string, userId uint64) (dg *model.DocGroup, err error) {
	q := global.DB.Model(&model.DocGroup{}).Where("title = ? AND user_id = ?", title, userId)
	if err = q.First(&dg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}

// DocGroupList 文档分组列表
func DocGroupList(r request.PaginationWithSort, userId uint64) (res resp.ListResponse, err error) {
	var list []model.DocGroup
	q := global.DB.Model(&model.DocGroup{}).Where("user_id = ?", userId)
	q.Count(&r.Total)
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&list).Error
	res.List = list
	res.Pagination = r
	return
}

// DocGroupUpdate 文档分组更新
func DocGroupUpdate(r doc.UpdateDocGroupRequest, userId uint64) (err error) {
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
func DocGroupDelete(r doc.UpdateDocGroupRequest, userId uint64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.DB.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.DocGroup{}).Error; err != nil {
		errMsg := fmt.Sprintf("删除id 为 %d 的数据失败 %v ", r.Id, err)
		global.ZAPSUGAR.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}

func DocGroupTree(r doc.GroupTreeRequest, userId uint64) (docTree []*resp.DocTree, err error) {
	docTree = make([]*resp.DocTree, 0)
	var list model.DocGroups
	q := global.DB.Where("user_id = ?", userId).Where("p_id = ?", r.Pid)
	if err = q.Find(&list).Error; err != nil {
		global.ZAPSUGAR.Error(err)
		return nil, errors.New("查询分组信息失败")
	}

	var countList model.DocGroups
	countQ := global.DB.Debug().Select("p_id,count(*) as children_count").Where("p_id IN (?)", list.GetIds()).Group("p_id").Find(&countList)
	if err = countQ.Error; err != nil {
		global.ZAPSUGAR.Error(err)
		return nil, errors.New("查询分组统计信息失败")
	}

	countMap := countList.ToPidMap()
	for _, v := range list {
		if cnt, ok := countMap[v.Id]; ok {
			v.ChildrenCount = cnt.ChildrenCount
		}
		docTree = append(docTree, &resp.DocTree{
			DocGroup: v,
		})
	}
	return
}
