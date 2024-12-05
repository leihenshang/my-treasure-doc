package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	resp "fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"

	"gorm.io/gorm"
)

// DocGroupCreate 创建文档分组
func DocGroupCreate(r doc.CreateDocGroupRequest, userId int64) (dg *model.DocGroup, err error) {
	if existed, err := checkDocGroupTitleRepeat(r.PId, r.Title, userId); err != nil {
		global.Log.Error(r, userId, err)
		return nil, errors.New("检查文档分组标题失败")
	} else if existed != nil {
		return nil, errors.New("文档分组标题已存在")
	}

	parentGroup := &model.DocGroup{
		BasicModel: model.BasicModel{
			Id: r.PId,
		},
	}
	if parentGroup.PId > 0 {
		if err = global.Db.Where("id = ? AND user_id = ?", r.PId, userId).First(&parentGroup).Error; err != nil {
			errorMsg := fmt.Errorf("查找父级分组失败")
			global.Log.Error(errorMsg, err)
			return nil, errorMsg
		}
	} else {
		parentGroup.GroupPath = strconv.FormatInt(r.PId, 10)
	}

	insertData := &model.DocGroup{
		Title:    r.Title,
		Icon:     r.Icon,
		PId:      r.PId,
		UserId:   userId,
		Priority: r.Priority,
	}

	tx := global.Db.Begin()

	if err = tx.Create(insertData).Error; err != nil {
		global.Log.Error(r, err)
		tx.Rollback()
		return nil, errors.New("创建文档分组失败")
	}

	groupPath, err := genGroupPath(insertData.Id, r.PId, userId)
	if err != nil {
		global.Log.Error(r, userId, err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&insertData).Update("GroupPath", groupPath).Error; err != nil {
		global.Log.Error(r, userId, err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return insertData, nil
}

func genGroupPath(id, pid int64, userId int64) (string, error) {
	parentGroup := &model.DocGroup{
		BasicModel: model.BasicModel{
			Id: pid,
		},
	}
	if pid > 0 {
		if err := global.Db.Where("id = ? AND user_id = ?", pid, userId).First(&parentGroup).Error; err != nil {
			errorMsg := fmt.Errorf("查找父级分组失败")
			global.Log.Error(errorMsg, err)
			return "", errorMsg
		}
	} else {
		parentGroup.GroupPath = strconv.FormatInt(pid, 10)
		return fmt.Sprintf("%d,%d", 0, id), nil
	}

	paths := append(strings.Split(parentGroup.GroupPath, ","), fmt.Sprintf(`%d`, id))
	return strings.Join(paths, ","), nil
}

// checkDocGroupTitleRepeat 查询数据库检查文档分组标题是否重复
func checkDocGroupTitleRepeat(pid int64, title string, userId int64) (dg *model.DocGroup, err error) {
	q := global.Db.Model(&model.DocGroup{}).Where("title = ? AND user_id = ? AND p_id = ?", title, userId, pid)
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
	q := global.Db.Model(&model.DocGroup{}).Where("user_id = ?", userId)
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

	groupPath, err := genGroupPath(r.Id, r.PId, userId)
	if err != nil {
		global.Log.Error(r, userId, err)
		return err
	}

	q := global.Db.Model(&model.DocGroup{}).Where("id = ? AND user_id = ?", r.Id, userId)
	u := map[string]interface{}{"Title": r.Title, "PId": r.PId, "Icon": r.Icon, "GroupPath": groupPath}
	if err = q.Updates(u).Error; err != nil {
		errMsg := fmt.Sprintf("修改id 为 %d 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return errors.New(errMsg)
	}

	return
}

// DocGroupDelete 文档分组删除
func DocGroupDelete(r doc.UpdateDocGroupRequest, userId int64) (err error) {
	if r.Id <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", r.Id)
		global.Log.Error(errMsg)
		return errors.New(errMsg)
	}

	tx := global.Db.Begin()
	q := tx.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.DocGroup{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除id 为 %d 的数据失败 %v ", r.Id, err)
	}
	if err = tx.Where("group_id = ? AND user_id = ?", r.Id, userId).Delete(&model.Doc{}).Error; err != nil {
		// if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 	tx.Rollback()
		// 	errMsg := fmt.Sprintf("删除id 为 %d 的分组下的文档失败 %v ", r.Id, err)
		// 	return errors.New(errMsg)
		// }
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
	if err = global.Db.Where("user_id = ?", userId).Where("p_id = ?", r.Pid).Order("created_at ASC").Find(&list).Error; err != nil {
		global.Log.Error(err)
		return nil, errors.New("查询分组信息失败")
	}

	var children model.DocGroups
	if err = global.Db.Where("user_id = ?", userId).Where("p_id IN (?)", list.GetIds()).Find(&children).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Error(err)
			return nil, errors.New("查询分组信息失败")
		}
	}

	childrenDoc, err := getDocByGroupIds(list.GetIds()...)
	if err != nil {
		return nil, err
	}

	childrenPidMap := children.ToPidMap()
	docGroupIdMap := childrenDoc.ToGroupIdMap()
	for _, v := range list {
		if _, ok := childrenPidMap[v.Id]; ok {
			v.IsLeaf = false
		} else {
			v.IsLeaf = true
		}
		if _, ok := docGroupIdMap[v.Id]; ok {
			v.IsLeaf = false
		}
		v.GroupType = model.GroupTypeGroup
		docTree = append(docTree, &resp.DocTree{
			DocGroup: v,
		})
	}

	if !r.WithDoc {
		return
	}

	currentDocs, err := getDocByGroupIds(r.Pid)
	if err != nil {
		return nil, err
	}

	for _, d := range currentDocs {
		docTree = append(docTree, &resp.DocTree{
			DocGroup: &model.DocGroup{
				BasicModel: model.BasicModel{
					Id: d.Id,
				},
				Title:     d.Title,
				Priority:  d.Priority,
				GroupType: model.GroupTypeDoc,
				IsLeaf:    true,
			},
		})
	}

	return
}

func getDocByGroupIds(groupId ...int64) (res model.Docs, err error) {
	err = global.Db.Where("group_id IN (?)", groupId).Find(&res).Error
	return
}

func getDocGroupByIds(groupId ...int64) (res model.DocGroups, err error) {
	err = global.Db.Where("id IN (?)", groupId).Find(&res).Error
	return
}
