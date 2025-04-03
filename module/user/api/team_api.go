package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/team"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/auth"
	"fastduck/treasure-doc/module/user/internal/service"
)

type TeamApi struct {
	TeamService *service.TeamService
}

func NewTeamApi() *TeamApi {
	return &TeamApi{TeamService: service.NewTeamService()}
}

// TeamCreate 创建团队
func (t *TeamApi) TeamCreate(c *gin.Context) {
	var req team.CreateOrUpdateTeamRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if d, ok := t.TeamService.TeamCreate(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}
}

// TeamDetail 团队详情
func (t *TeamApi) TeamDetail(c *gin.Context) {
	req := request.IDReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if d, ok := t.TeamService.TeamDetail(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}

}

// TeamList 团队列表
func (t *TeamApi) TeamList(c *gin.Context) {
	var req request.Pagination
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if d, ok := t.TeamService.TeamList(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}
}

// TeamDelete 团队删除
func (t *TeamApi) TeamDelete(c *gin.Context) {
	var req team.CreateOrUpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := t.TeamService.TeamDelete(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}
