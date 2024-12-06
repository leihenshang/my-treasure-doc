package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/middleware"

	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/data/request/note"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"
)

// NoteCreate 创建文档
func NoteCreate(c *gin.Context) {
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

	if d, ok := service.NoteCreate(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(d, c)
	}
}

// NoteDetail 文档详情
func NoteDetail(c *gin.Context) {
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
	if d, ok := service.NoteDetail(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(d, c)
	}

}

// NoteList 文档列表
func NoteList(c *gin.Context) {
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
	if d, ok := service.NoteList(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(d, c)
	}
}

// NoteUpdate 文档更新
func NoteUpdate(c *gin.Context) {
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
	if ok := service.NoteUpdate(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}

// NoteDelete 文档删除
func NoteDelete(c *gin.Context) {
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
	if ok := service.NoteDelete(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}
