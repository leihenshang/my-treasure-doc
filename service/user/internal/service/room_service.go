package service

import (
	"errors"
	"fmt"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request/room"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
)

type RoomService struct{}

func NewRoomService() *RoomService {
	return &RoomService{}
}

func (s *RoomService) Create(req room.CreateRoomRequest, userId string) (res *model.Room, err error) {
	if req.Name == "" {
		return nil, errors.New("房间名称不能为空")
	}

	roomObj := &model.Room{
		Name:   req.Name,
		UserId: fmt.Sprintf("%d", userId), // 将uint转换为string
		Status: model.RoomStatusNormal,
	}

	if err = global.Db.Create(roomObj).Error; err != nil {
		global.Log.Errorf("RoomService.Create error:%v", err)
		return nil, errors.New("创建房间失败")
	}
	return roomObj, nil
}

func (s *RoomService) Detail(id string, userId string) (res *model.Room, err error) {
	if id == "" {
		return nil, errors.New("房间ID不能为空")
	}

	var roomObj model.Room
	err = global.Db.Where("id = ? AND user_id = ?", id, fmt.Sprintf("%d", userId)).First(&roomObj).Error
	if err != nil {
		global.Log.Errorf("RoomService.Detail error:%v", err)
		return nil, errors.New("获取房间详情失败")
	}

	return &roomObj, nil
}

func (s *RoomService) List(req room.ListRoomRequest, userId string) (res response.ListResponse, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if userId == "" {
		return res, errors.New("用户ID不能为空")
	}

	var rooms []*model.Room
	var total int64

	db := global.Db.Model(&model.Room{}).Where("user_id = ?", userId)
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	err = db.Count(&total).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&rooms).Error
	if err != nil {
		global.Log.Errorf("RoomService.List error:%v", err)
		return res, errors.New("获取房间列表失败")
	}

	res.List = rooms
	res.Pagination = req.Pagination
	return
}

func (s *RoomService) Update(req room.UpdateRoomRequest, userId string) (res *model.Room, err error) {
	if req.Id == "" {
		return nil, errors.New("房间ID不能为空")
	}

	roomObj := &model.Room{
		BaseModel: model.BaseModel{Id: req.Id},
		Name:      req.Name,
		Status:    model.RoomStatus(req.Status),
	}

	err = global.Db.Model(roomObj).
		Where("id = ? AND user_id = ?", req.Id, fmt.Sprintf("%d", userId)).
		Updates(roomObj).Error
	if err != nil {
		global.Log.Errorf("RoomService.Update error:%v", err)
		return nil, errors.New("更新房间失败")
	}

	return s.Detail(req.Id, userId)
}

func (s *RoomService) Delete(id string, userId string) (err error) {
	if id == "" {
		return errors.New("房间ID不能为空")
	}

	err = global.Db.Where("id = ? AND user_id = ?", id, fmt.Sprintf("%d", userId)).
		Delete(&model.Room{}).Error
	if err != nil {
		global.Log.Errorf("RoomService.Delete error:%v", err)
		return errors.New("删除房间失败")
	}
	return nil
}
