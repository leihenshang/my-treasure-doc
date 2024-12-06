package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/middleware"

	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/doc"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"
)

type DocHistoryApi struct {
	DocHistoryService *service.DocHistoryService
}

func NewDocHistoryApi() *DocHistoryApi {
	return &DocHistoryApi{DocHistoryService: service.NewDocHistoryService()}
}

// DocHistoryDetail 文档详情
func (h *DocHistoryApi) DocHistoryDetail(c *gin.Context) {
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
	if d, ok := h.DocHistoryService.DocHistoryDetail(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}

}

// DocHistoryList 文档列表
func (h *DocHistoryApi) DocHistoryList(c *gin.Context) {
	var req doc.ListDocHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if d, ok := h.DocHistoryService.DocHistoryList(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}
}

// DocHistoryRecover 文档更新
func (h *DocHistoryApi) DocHistoryRecover(c *gin.Context) {
	req := request.IDReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := h.DocHistoryService.DocHistoryRecover(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}
