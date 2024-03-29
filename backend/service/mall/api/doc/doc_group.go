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

//DocGroupCreate 创建文档分组
func DocGroupCreate(c *gin.Context) {
	var req doc.CreateDocGroupRequest
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

	if d, ok := service.DocGroupCreate(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

//DocGroupList 文档分组列表
func DocGroupList(c *gin.Context) {
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
	if d, ok := service.DocGroupList(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

//DocGroupUpdate 文档分组更新
func DocGroupUpdate(c *gin.Context) {
	var req doc.UpdateDocGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.DocGroupUpdate(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}

//DocGroupDelete 文档分组删除
func DocGroupDelete(c *gin.Context) {
	var req doc.UpdateDocGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.DocGroupDelete(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}
