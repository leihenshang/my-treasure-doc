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

// DocGroupCreate 创建文档分组
func DocGroupCreate(c *gin.Context) {
	var req doc.CreateDocGroupRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if d, ok := service.DocGroupCreate(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

// DocGroupList 文档分组列表
func DocGroupList(c *gin.Context) {
	var req request.ListPagination
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if d, ok := service.DocGroupList(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

// DocGroupUpdate 文档分组更新
func DocGroupUpdate(c *gin.Context) {
	var req doc.UpdateDocGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.DocGroupUpdate(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}

// DocGroupDelete 文档分组删除
func DocGroupDelete(c *gin.Context) {
	var req doc.UpdateDocGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.DocGroupDelete(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}

// DocGroupTree 文档组树
func DocGroupTree(c *gin.Context) {
	var req doc.GroupTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if res, ok := service.DocGroupTree(req, u.Id); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(res, c)
	}
}
