package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/middleware"

	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/doc"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
)

type DocApi struct {
	DocService *service.DocService
}

func NewDocApi() *DocApi {
	return &DocApi{DocService: service.NewDocService()}
}

// Create 创建文档
func (d *DocApi) Create(c *gin.Context) {
	var req doc.CreateDocRequest
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

	createDoc := &model.Doc{
		UserId:  u.Id,
		Title:   req.Title,
		Content: req.Content,
		GroupId: req.GroupId,
		IsTop:   req.IsTop,
	}

	if newDoc, ok := d.DocService.Create(createDoc, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, newDoc)
	}
}

// Detail 文档详情
func (d *DocApi) Detail(c *gin.Context) {
	req := request.IDReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if docObj, ok := d.DocService.Detail(req.ID, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, docObj)
	}

}

// List 文档列表
func (d *DocApi) List(c *gin.Context) {
	var req doc.ListDocRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if list, ok := d.DocService.List(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, list)
	}
}

// Update 文档更新
func (d *DocApi) Update(c *gin.Context) {
	var req doc.UpdateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if newDoc, err := d.DocService.Update(req, u.Id); err != nil {
		if errors.Is(err, service.ErrorDocIsEdited) {
			response.FailWithMessage(c, err.Error(), response.DocIsEdited)
			return
		}
		response.FailWithMessage(c, err.Error())
	} else {
		response.OkWithData(c, newDoc)
	}
}

// Delete 文档删除
func (d *DocApi) Delete(c *gin.Context) {
	var req request.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := d.DocService.Delete(req.ID, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}

func (d *DocApi) Recover(c *gin.Context) {
	var req request.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := d.DocService.Recover(req.ID, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}
