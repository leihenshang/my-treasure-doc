package room

import "fastduck/treasure-doc/service/user/data/request"

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=250"`          // 房间名称
	Status string `json:"status" binding:"omitempty,oneof=normal closed"` // 房间状态
}

// UpdateRoomRequest 更新房间请求
type UpdateRoomRequest struct {
	Id     string `json:"id" binding:"required"`                          // 房间ID
	Name   string `json:"name" binding:"required,min=1,max=250"`          // 房间名称
	Status string `json:"status" binding:"omitempty,oneof=normal closed"` // 房间状态
}

// ListRoomRequest 房间列表请求
type ListRoomRequest struct {
	Name   string `json:"name" form:"name" binding:"omitempty,max=250"`                 // 房间名称模糊查询
	Status string `json:"status" form:"status" binding:"omitempty,oneof=normal closed"` // 房间状态
	request.Pagination
	request.Sort
}
