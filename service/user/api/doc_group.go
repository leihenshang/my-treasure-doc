package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/middleware"

	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"
)

type DocGroupApi struct {
	DocGroupService *service.DocGroupService
}

func NewDocGroupApi() *DocGroupApi {
	return &DocGroupApi{DocGroupService: service.NewDocGroupService()}
}

// DocGroupCreate 创建文档分组
func (d *DocGroupApi) DocGroupCreate(c *gin.Context) {
	var req doc.CreateDocGroupRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
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
	if group, ok := d.DocGroupService.DocGroupCreate(insertData, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, group)
	}
}

// DocGroupList 文档分组列表
func (d *DocGroupApi) DocGroupList(c *gin.Context) {
	var req request.Pagination
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if group, ok := d.DocGroupService.DocGroupList(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, group)
	}
}

// DocGroupUpdate 文档分组更新
func (d *DocGroupApi) DocGroupUpdate(c *gin.Context) {
	var req doc.UpdateDocGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
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
	if ok := d.DocGroupService.DocGroupUpdate(updateGroup, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}

// DocGroupDelete 文档分组删除
func (d *DocGroupApi) DocGroupDelete(c *gin.Context) {
	var req doc.UpdateDocGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := d.DocGroupService.DocGroupDelete(req.Id, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}

// DocGroupTree 文档组树
func (d *DocGroupApi) DocGroupTree(c *gin.Context) {
	var req doc.GroupTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if res, ok := d.DocGroupService.DocGroupTree(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, res)
	}
}
