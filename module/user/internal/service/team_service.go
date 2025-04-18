package service

import (
	"errors"
	"fmt"
	"sync"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/team"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"

	"gorm.io/gorm"
)

type TeamService struct{}

var teamService *TeamService

var teamOnce = sync.Once{}

func NewTeamService() *TeamService {
	teamOnce.Do(func() {
		teamService = &TeamService{}
	})
	return teamService
}

// TeamCreate 创建文档
func (t *TeamService) TeamCreate(r team.CreateOrUpdateTeamRequest, userId string) (d *model.Team, err error) {
	insertData := &model.Team{}

	if existed, checkErr := checkTeamTitleRepeat(insertData.Name, userId); checkErr != nil {
		global.Log.Error(r, userId, "检查文档标题失败")
		return nil, errors.New("检查文档标题失败")
	} else {
		if existed != nil {
			return nil, errors.New("文档标题已存在")
		}
	}

	if err = global.Db.Create(insertData).Error; err != nil {
		global.Log.Error(r, err)
		return nil, errors.New("创建文档失败")
	}

	return
}

// checkTeamTitleRepeat 查询数据库检查文档标题是否重复
func checkTeamTitleRepeat(title string, userId string) (team *model.Team, err error) {
	q := global.Db.Model(&model.Team{}).Where("title = ? AND user_id = ?", title, userId)
	if err = q.First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}

// TeamDetail 文档详情
func (t *TeamService) TeamDetail(r request.IDReq, userId string) (d *model.Team, err error) {
	q := global.Db.Model(&model.Team{}).Where("id = ? AND user_id = ?", r.ID, userId)
	err = q.First(&d).Error
	return
}

// TeamList 文档列表
func (t *TeamService) TeamList(r request.Pagination, userId string) (res response.ListResponse, err error) {
	offset := (r.Page - 1) * r.PageSize
	if offset < 0 {
		offset = 1
	}

	var list []model.Team
	q := global.Db.Model(&model.Team{}).Where("user_id = ?", userId)
	q.Count(&res.Pagination.Total)
	err = q.Limit(r.PageSize).Offset(offset).Find(&list).Error
	res.List = list
	res.Pagination = r
	return
}

// TeamUpdate 文档更新
func (t *TeamService) TeamUpdate(r team.CreateOrUpdateTeamRequest, userId string) (err error) {
	if r.Id == "" {
		return nil
	}

	q := global.Db.Model(&model.Team{}).Where("id = ? AND user_id = ?", r.Id, userId)
	u := map[string]interface{}{"Name": r.Name}
	if err = q.Updates(u).Error; err != nil {
		errMsg := fmt.Sprintf("修改id 为 %s 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}

// TeamDelete 文档删除
func (t *TeamService) TeamDelete(r team.CreateOrUpdateTeamRequest, userId string) (err error) {
	if r.Id == "" {
		return nil
	}

	q := global.Db.Where("id = ? AND user_id = ?", r.Id, userId)
	if err = q.Delete(&model.Team{}).Error; err != nil {
		errMsg := fmt.Sprintf("删除id 为 %s 的数据失败 %v ", r.Id, err)
		global.Log.Error(errMsg)
		return errors.New("操作失败")
	}

	return
}
