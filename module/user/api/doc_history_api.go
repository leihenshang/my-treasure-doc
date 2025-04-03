package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/internal/auth"

	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/doc"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
)

type DocHistoryApi struct {
	DocHistoryService *service.DocHistoryService
}

func NewDocHistoryApi() *DocHistoryApi {
	return &DocHistoryApi{DocHistoryService: service.NewDocHistoryService()}
}

// Detail 文档详情
func (h *DocHistoryApi) Detail(c *gin.Context) {
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
	if d, ok := h.DocHistoryService.Detail(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}

}

// List 文档列表
func (h *DocHistoryApi) List(c *gin.Context) {
	var req doc.ListDocHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if d, ok := h.DocHistoryService.List(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}
}

// Recover 恢复
func (h *DocHistoryApi) Recover(c *gin.Context) {
	req := request.IDReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := h.DocHistoryService.Recover(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}
