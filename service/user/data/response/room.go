package response

import "fastduck/treasure-doc/service/user/data/model"

// RoomDetail 房间详情响应
type RoomDetail struct {
	*model.Room
}

// RoomList 房间列表响应
type RoomList struct {
	List       []*model.Room `json:"list"`
	Pagination struct {
		Total    int64 `json:"total"`
		Page     int   `json:"page"`
		PageSize int   `json:"pageSize"`
	} `json:"pagination"`
}
