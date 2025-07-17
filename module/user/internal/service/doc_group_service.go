package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/doc"
	resp "fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"

	"gorm.io/gorm"
)

type DocGroupService struct{}

var docGroupService *DocGroupService

var docGroupOnce = sync.Once{}

func NewDocGroupService() *DocGroupService {
	docGroupOnce.Do(func() {
		docGroupService = &DocGroupService{}
	})
	return docGroupService
}

func GetDocGroupService() *DocGroupService {
	return docGroupService
}

// Create 创建文档分组
func (group *DocGroupService) Create(createGroup *model.DocGroup, userId string) (dg *model.DocGroup, err error) {
	if createGroup == nil {
		return nil, nil
	}

	if createGroup.PId == "" {
		createGroup.PId = global.RootGroup
	}

	if dbGroup, err := checkDocGroupTitleRepeat(createGroup.PId, createGroup.Title, userId); err != nil {
		global.Log.Error(createGroup, userId, err)
		return nil, errors.New("检查文档分组标题失败")
	} else if dbGroup != nil {
		return nil, errors.New("文档分组标题已存在")
	}

	tx := global.Db.Begin()
	if err = tx.Create(createGroup).Error; err != nil {
		global.Log.Error(createGroup, err)
		tx.Rollback()
		return nil, errors.New("创建文档分组失败")
	}

	groupPath, err := genGroupPath(createGroup, userId)
	if err != nil {
		global.Log.Error(createGroup, userId, err)
		tx.Rollback()
		return nil, errors.New("更新分组路径失败")
	}

	if err = tx.Model(&createGroup).Update("GroupPath", groupPath).Error; err != nil {
		global.Log.Error(createGroup, userId, err)
		tx.Rollback()
		return nil, errors.New("更新分组路径失败")
	}

	tx.Commit()
	return createGroup, nil
}

func genGroupPath(group *model.DocGroup, userId string) (string, error) {
	if group.PId == "" || group.PId == global.RootGroup {
		return fmt.Sprintf("%s,%s", global.RootGroup, group.Id), nil
	}
	var parentGroup *model.DocGroup
	if err := global.Db.Where("id = ? AND user_id = ?", group.PId, userId).First(&parentGroup).Error; err != nil {
		errorMsg := fmt.Errorf("查找父级分组失败")
		global.Log.Error(errorMsg, err)
		return "", errorMsg
	}
	paths := append(strings.Split(parentGroup.GroupPath, ","), group.Id)
	return strings.Join(paths, ","), nil
}

// checkDocGroupTitleRepeat 查询数据库检查文档分组标题是否重复
func checkDocGroupTitleRepeat(pid string, title string, userId string) (dg *model.DocGroup, err error) {
	q := global.Db.Model(&model.DocGroup{}).Where("title = ? AND user_id = ? AND p_id = ?", title, userId, pid)
	if err = q.First(&dg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}

// Detail 分组详情
func (group *DocGroupService) Detail(id string, userId string) (d *model.DocGroup, err error) {
	return d, global.Db.Unscoped().Where("id = ? AND user_id = ?", id, userId).First(&d).Error
}

// List 文档分组列表
func (group *DocGroupService) List(r doc.ListDocGroupRequest, userId string) (res resp.ListResponse, err error) {
	var list model.DocGroups
	q := global.Db.Model(&model.DocGroup{}).Where("user_id = ?", userId)
	if r.Id != "" {
		q = q.Where("id = ?", r.Id)
	}
	if r.PageSize > 0 {
		q.Count(&r.Total)
	}
	err = q.Limit(r.PageSize).Offset(r.Offset()).Find(&list).Error
	if err != nil {
		return res, err
	}

	var parentList model.DocGroups
	if err = global.Db.Model(&model.DocGroup{}).Where("user_id = ? AND id IN (?)", userId, list.GetPIds()).Find(&parentList).Error; err != nil {
		return res, err
	}

	parentMap := parentList.ToMap()
	parentMap["root"] = &model.DocGroup{
		BaseModel: model.BaseModel{
			Id: "root",
		},
		Title:     "root",
		GroupType: model.GroupTypeGroup,
		IsLeaf:    false,
	}
	for _, docGroup := range list {
		pathList := strings.Split(docGroup.GroupPath, ",")
		for _, path := range pathList {
			if p, ok := parentMap[path]; ok {
				docGroup.GroupPathList = append(docGroup.GroupPathList, p)
			}
		}
		docGroup.GroupPathList = append(docGroup.GroupPathList, &model.DocGroup{
			BaseModel: docGroup.BaseModel,
			UserId:    docGroup.UserId,
			Title:     docGroup.Title,
			Icon:      docGroup.Icon,
			PId:       docGroup.PId,
			Priority:  docGroup.Priority,
			GroupType: docGroup.GroupType,
			IsLeaf:    docGroup.IsLeaf,
			GroupPath: docGroup.GroupPath,
		})
	}

	res.List = list
	res.Pagination = r.Pagination
	return
}

// Update 文档分组更新
func (group *DocGroupService) Update(updateGroup *model.DocGroup, userId string) (err error) {
	if updateGroup.Id == "" {
		return nil
	}

	if updateGroup.PId == "" {
		updateGroup.PId = global.RootGroup
	}

	groupPath, err := genGroupPath(updateGroup, userId)
	if err != nil {
		global.Log.Error(updateGroup, userId, err)
		return err
	}

	q := global.Db.Model(&model.DocGroup{}).Where("id = ? AND user_id = ?", updateGroup.Id, userId)
	u := map[string]interface{}{"Title": updateGroup.Title, "PId": updateGroup.PId, "Icon": updateGroup.Icon,
		"GroupPath": groupPath}
	if err = q.Updates(u).Error; err != nil {
		errMsg := fmt.Sprintf("修改id 为 %s 的数据失败 %v ", updateGroup.Id, err)
		global.Log.Error(errMsg)
		return errors.New(errMsg)
	}

	return
}

// Delete 文档分组删除
func (group *DocGroupService) Delete(id string, userId string) (err error) {
	if id == "" {
		return nil
	}

	tx := global.Db.Begin()
	q := tx.Where("id = ? AND user_id = ?", id, userId)
	if err = q.Delete(&model.DocGroup{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除id 为 %s 的数据失败 %v ", id, err)
	}
	if err = tx.Where("group_id = ? AND user_id = ?", id, userId).Delete(&model.Doc{}).Error; err != nil {
		tx.Rollback()
		errMsg := fmt.Sprintf("删除id 为 %s 的分组下的文档失败 %v ", id, err)
		return errors.New(errMsg)
	}
	tx.Commit()
	return
}

func (group *DocGroupService) Tree(r doc.GroupTreeRequest, userId string) (docTree resp.DocTrees, err error) {
	var list model.DocGroups
	if err = global.Db.Where("user_id = ?", userId).Where("p_id = ?", r.Pid).Order("created_at ASC").Find(&list).Error; err != nil {
		global.Log.Error(err)
		return nil, errors.New("查询分组信息失败")
	}

	var children model.DocGroups
	if err = global.Db.Where("user_id = ?", userId).Where("p_id IN (?)", list.GetIds()).Find(&children).Error; err != nil {
		global.Log.Error(err)
		return nil, errors.New("查询下级分组信息失败")
	}

	childrenDoc, err := getDocByGroupIds(userId, list.GetIds()...)
	if err != nil {
		return nil, err
	}

	childrenPidMap := children.ToPidMap()
	docMapWithGroupId := childrenDoc.ToGroupIdMap()
	excludeMap := request.GetUniqueMapFromDotStr(r.ExcludeIds)
	for _, v := range list {
		if _, ok := excludeMap[v.Id]; ok {
			continue
		}

		if _, ok := excludeMap[v.PId]; ok {
			excludeMap[v.Id] = v.Id
			continue
		}

		if _, ok := childrenPidMap[v.Id]; !ok {
			v.IsLeaf = true
		}

		if r.WithDoc {
			if _, ok := docMapWithGroupId[v.Id]; ok {
				v.IsLeaf = false
			}
		}

		v.GroupType = model.GroupTypeGroup
		docTree = append(docTree, &resp.DocTree{
			DocGroup: v,
		})
	}

	if !r.WithDoc {
		return
	}

	docs, err := getDocByGroupIds(userId, r.Pid)
	if err != nil {
		return nil, err
	}

	for _, d := range docs {
		docTree = append(docTree, &resp.DocTree{
			DocGroup: &model.DocGroup{
				BaseModel: d.BaseModel,
				Title:     d.Title,
				Priority:  d.Priority,
				GroupType: model.GroupTypeDoc,
				IsLeaf:    true,
				PId:       r.Pid,
			},
		})
	}
	return
}

func getDocByGroupIds(userId string, groupId ...string) (res model.Docs, err error) {
	err = global.Db.Where("user_id = ?", userId).Where("group_id IN (?)", groupId).Find(&res).Error
	return
}

func getDocGroupByIds(userId string, groupId ...string) (res model.DocGroups, err error) {
	err = global.Db.Where("user_id = ?", userId).Where("id IN (?)", groupId).Find(&res).Error
	return
}
