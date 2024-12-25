package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/middleware"

	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"
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

	if newDoc, ok := d.DocService.DocCreate(createDoc, u.Id); ok != nil {
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
	if docObj, ok := d.DocService.DocDetail(req, u.Id); ok != nil {
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
	if list, ok := d.DocService.DocList(req, u.Id); ok != nil {
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
	if newDoc, err := d.DocService.DocUpdate(req, u.Id); err != nil {
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
	var req doc.DeleteDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := d.DocService.DocDelete(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}

// Tree 文档树
func (d *DocApi) Tree(c *gin.Context) {
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

	if res, ok := d.DocService.DocTree(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, res)
	}
}
