package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/request"
	"fastduck/treasure-doc/module/user/data/request/room"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
	"fastduck/treasure-doc/module/user/middleware"
)

type RoomApi struct {
	RoomService *service.RoomService
}

func NewRoomApi() *RoomApi {
	return &RoomApi{RoomService: service.NewRoomService()}
}

// Create 创建房间
func (r *RoomApi) Create(c *gin.Context) {
	var req room.CreateRoomRequest
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

	if newRoom, ok := r.RoomService.Create(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, newRoom)
	}
}

// Detail 房间详情
func (r *RoomApi) Detail(c *gin.Context) {
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

	if roomObj, ok := r.RoomService.Detail(req.ID, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, roomObj)
	}
}

// List 房间列表
func (r *RoomApi) List(c *gin.Context) {
	var req room.ListRoomRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if list, ok := r.RoomService.List(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, list)
	}
}

// Update 更新房间
func (r *RoomApi) Update(c *gin.Context) {
	var req room.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if updatedRoom, ok := r.RoomService.Update(req, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, updatedRoom)
	}
}

// Delete 删除房间
func (r *RoomApi) Delete(c *gin.Context) {
	var req request.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if ok := r.RoomService.Delete(req.ID, u.Id); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.Ok(c)
	}
}
