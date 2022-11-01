package api

import (
	"fastduck/treasure-doc/service/mall/global"
	"fastduck/treasure-doc/service/mall/internal/service"
	"fastduck/treasure-doc/service/mall/middleware/auth"
	"fastduck/treasure-doc/service/mall/request"
	"fastduck/treasure-doc/service/mall/request/team"
	"fastduck/treasure-doc/service/mall/response"

	"github.com/gin-gonic/gin"
)

//TeamCreate 创建团队
func TeamCreate(c *gin.Context) {
	var req team.CreateOrUpdateTeamRequest
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

	if d, ok := service.TeamCreate(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

//TeamDetail 团队详情
func TeamDetail(c *gin.Context) {
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
	if d, ok := service.TeamDetail(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}

}

//TeamList 团队列表
func TeamList(c *gin.Context) {
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
	if d, ok := service.TeamList(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

//TeamUpdate 团队更新
func TeamUpdate(c *gin.Context) {
	var req team.CreateOrUpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.TeamUpdate(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}

//TeamDelete 团队删除
func TeamDelete(c *gin.Context) {
	var req team.CreateOrUpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if ok := service.TeamDelete(req, uint64(u.ID)); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.Ok(c)
	}
}
