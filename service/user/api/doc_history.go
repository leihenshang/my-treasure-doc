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

// DocHistoryDetail 文档详情
func DocHistoryDetail(c *gin.Context) {
	req := request.IDReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if d, ok := service.DocHistoryDetail(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}

}

// DocHistoryList 文档列表
func DocHistoryList(c *gin.Context) {
	var req doc.ListDocHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if d, ok := service.DocHistoryList(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

// DocHistoryRecover 文档更新
func DocHistoryRecover(c *gin.Context) {
	req := request.IDReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.DocHistoryRecover(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}
