package api

import (
	"fastduck/treasure-doc/service/mall/data/request"
	"fastduck/treasure-doc/service/mall/data/request/doc"
	"fastduck/treasure-doc/service/mall/data/response"
	"fastduck/treasure-doc/service/mall/global"
	"fastduck/treasure-doc/service/mall/internal/service"
	"fastduck/treasure-doc/service/mall/middleware/auth"

	"github.com/gin-gonic/gin"
)

//DocCreate 创建文档
func DocCreate(c *gin.Context) {
	var req doc.CreateDocRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if d, ok := service.DocCreate(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

//DocDetail 文档详情
func DocDetail(c *gin.Context) {
	req := request.IdRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if d, ok := service.DocDetail(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}

}

//DocList 文档列表
func DocList(c *gin.Context) {
	var req request.ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if d, ok := service.DocList(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

//DocUpdate 文档更新
func DocUpdate(c *gin.Context) {
	var req doc.UpdateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.DocUpdate(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}

//DocDelete 文档删除
func DocDelete(c *gin.Context) {
	var req doc.UpdateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.DocDelete(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}
