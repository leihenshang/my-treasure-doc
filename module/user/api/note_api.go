package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/router/middleware"

	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/note"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
)

type NoteApi struct {
	NoteService *service.NoteService
}

func NewNoteApi() *NoteApi {
	return &NoteApi{NoteService: service.NewNoteService()}
}

// NoteCreate 创建文档
func (n *NoteApi) NoteCreate(c *gin.Context) {
	var req note.CreateNoteRequest
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

	if d, ok := n.NoteService.NoteCreate(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}
}

// NoteDetail 文档详情
func (n *NoteApi) NoteDetail(c *gin.Context) {
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
	if d, ok := n.NoteService.NoteDetail(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}

}

// NoteList 文档列表
func (n *NoteApi) NoteList(c *gin.Context) {
	var req note.ListNoteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if d, ok := n.NoteService.NoteList(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}
}

// NoteUpdate 文档更新
func (n *NoteApi) NoteUpdate(c *gin.Context) {
	var req note.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if d, ok := n.NoteService.Update(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, d)
	}
}

// NoteDelete 文档删除
func (n *NoteApi) NoteDelete(c *gin.Context) {
	var req note.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if ok := n.NoteService.NoteDelete(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}
