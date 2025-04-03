package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/internal/auth"

	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/doc"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
)

type DocGroupApi struct {
	DocGroupService *service.DocGroupService
}

func NewDocGroupApi() *DocGroupApi {
	return &DocGroupApi{DocGroupService: service.NewDocGroupService()}
}

// Create 创建文档分组
func (d *DocGroupApi) Create(c *gin.Context) {
	var req doc.CreateDocGroupRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	insertData := &model.DocGroup{
		Title:    req.Title,
		Icon:     req.Icon,
		PId:      req.PId,
		UserId:   u.Id,
		Priority: req.Priority,
	}
	if group, ok := d.DocGroupService.Create(insertData, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, group)
	}
}

// Detail 文档详情
func (d *DocGroupApi) Detail(c *gin.Context) {
	req := request.IDReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if resp, err := d.DocGroupService.List(doc.ListDocGroupRequest{
		Id: req.ID,
	}, u.Id); err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		if groups, ok := resp.List.(model.DocGroups); ok && len(groups) > 0 {
			response.OkWithData(c, groups[0])
			return
		}
		response.OkWithData(c, nil)
	}

}

// List 文档分组列表
func (d *DocGroupApi) List(c *gin.Context) {
	var req doc.ListDocGroupRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if group, ok := d.DocGroupService.List(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, group)
	}
}

// Update 文档分组更新
func (d *DocGroupApi) Update(c *gin.Context) {
	var req doc.UpdateDocGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	updateGroup := &model.DocGroup{
		BaseModel: model.BaseModel{
			Id: req.Id,
		},
		Title: req.Title,
		Icon:  req.Icon,
		PId:   req.PId,
	}
	if ok := d.DocGroupService.Update(updateGroup, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}

// Delete 文档分组删除
func (d *DocGroupApi) Delete(c *gin.Context) {
	var req request.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := d.DocGroupService.Delete(req.ID, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}

func (d *DocGroupApi) Tree(c *gin.Context) {
	var req doc.GroupTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if res, ok := d.DocGroupService.Tree(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, res)
	}
}
