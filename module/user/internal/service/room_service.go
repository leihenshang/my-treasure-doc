package service

import (
	"errors"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/data/request/room"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
)

type RoomService struct{}

var roomSrv *RoomService

func NewRoomService() *RoomService {
	roomSrv = &RoomService{}
	return roomSrv
}

func GetRoomService() *RoomService {
	return roomSrv
}

func (s *RoomService) Create(req room.CreateRoomRequest, userId string) (res *model.Room, err error) {
	if req.Name == "" {
		return nil, errors.New("房间名称不能为空")
	}

	// 检查用户是否已有同名房间
	var count int64
	if err = global.Db.Model(&model.Room{}).
		Where("user_id = ? AND name = ?", userId, req.Name).
		Count(&count).Error; err != nil {
		global.Log.Errorf("RoomService.Create check duplicate error:%v", err)
		return nil, errors.New("检查房间名称失败")
	}
	if count > 0 {
		return nil, errors.New("房间名称已存在")
	}

	roomObj := &model.Room{
		Name:   req.Name,
		UserId: userId,
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
	err = global.Db.Where("id = ? AND user_id = ?", id, userId).First(&roomObj).Error
	if err != nil {
		global.Log.Errorf("RoomService.Detail error:%v", err)
		return nil, errors.New("获取房间详情失败")
	}

	return &roomObj, nil
}

func (s *RoomService) GetDefaultRoom(userId string) (res *model.Room, err error) {
	var room model.Room
	err = global.Db.Where("user_id = ? AND is_default = ?", userId, global.RoomIsDefault).First(&room).Error
	if err != nil {
		global.Log.Errorf("RoomService.GetDefaultRoom error:%v", err)
		return nil, errors.New("获取默认房间详情失败")
	}

	return &room, nil
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

	db := global.Db.Model(&model.Room{}).Where("user_id = ?", userId)
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if req.PageSize > 0 {
		db.Count(&req.Total)
	}

	err = db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&rooms).Error
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
		Where("id = ? AND user_id = ?", req.Id, userId).
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

	// 检查是否是默认空间
	var roomObj model.Room
	if err = global.Db.Where("id = ? AND user_id = ?", id, userId).First(&roomObj).Error; err != nil {
		global.Log.Errorf("RoomService.Delete query error:%v", err)
		return errors.New("查询房间信息失败")
	}
	if roomObj.IsDefault == 1 {
		return errors.New("默认空间不能删除")
	}

	err = global.Db.Where("id = ? AND user_id = ?", id, userId).
		Delete(&model.Room{}).Error
	if err != nil {
		global.Log.Errorf("RoomService.Delete error:%v", err)
		return errors.New("删除房间失败")
	}
	return nil
}
