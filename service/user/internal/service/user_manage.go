package service

import (
	"errors"
	"fmt"
	"sync"

	"fastduck/treasure-doc/service/user/data/model"
	userReq "fastduck/treasure-doc/service/user/data/request/user"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
)

type UserManageService struct{}

var userManageService *UserManageService

var userManageOnce = sync.Once{}

func NewUserManageService() *UserManageService {
	userManageOnce.Do(func() {
		userManageService = &UserManageService{}
	})
	return userManageService
}

func (u *UserManageService) List(r userReq.ListUserManageRequest) (res response.ListResponse, err error) {
	q := global.Db.Model(&model.User{})
	if r.Keyword != "" {
		likeStr := fmt.Sprintf(`%%%s%%`, r.Keyword)
		q = q.Where("account LIKE ? OR email LIKE ?", likeStr, likeStr)
	}

	if r.Id > 0 {
		q = q.Where("id = ?", r.Id)
	}

	if r.Pagination.PageSize > 0 {
		q.Count(&r.Total)
	}

	if sortStr, err := r.Sort.Sort(map[string]string{"createdAt": "created_at", "id": "id"}); err == nil {
		q = q.Order(sortStr)
	}

	var list model.Users
	if r.Pagination.Page > 0 && r.Pagination.PageSize > 0 {
		q = q.Limit(r.PageSize).Offset(r.Offset())
	}

	err = q.Find(&list).Error
	res.List = list
	res.Pagination = r.Pagination
	return
}

func (u *UserManageService) Create(user *model.User) (createdUser *model.User, err error) {
	regRequest := &userReq.UserRegisterRequest{
		Password:   user.Password,
		RePassword: user.Password,
		Account:    user.Account,
		Email:      user.Email,
	}

	if createdUser, err = userService.UserRegister(regRequest); err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (u *UserManageService) Delete(userId int64) error {
	if userId <= 0 {
		errMsg := fmt.Sprintf("id 为 %d 的数据没有找到", userId)
		global.Log.Error(errMsg)
		return errors.New(errMsg)
	}

	q := global.Db.Where("id = ? AND user_type = ?", userId, model.UserTypeUser)
	if err := q.Delete(&model.Doc{}).Error; err != nil {
		errMsg := fmt.Sprintf("删除id 为 %d 的数据失败 %v ", userId, err)
		global.Log.Error(errMsg)
		return errors.New("操作失败")
	}

	return nil
}

func (u *UserManageService) Update(user *model.User) (updatedUser *model.User, err error) {
	return user, global.Db.Select("NickName", "UserStatus", "Mobile", "avatar", "Bio").Save(user).Error
}

func (u *UserManageService) Detail(userId int64) (user *model.User, err error) {
	if res, err := u.List(userReq.ListUserManageRequest{Id: userId}); err != nil {
		return nil, err
	} else if userList, ok := res.List.(model.Users); !ok {
		return nil, errors.New("用户没有找打")
	} else {
		return userList[0], nil
	}
}

func (u *UserManageService) RestPwd(req *userReq.RestPwdRequest) error {
	var user *model.User
	if res, err := u.List(userReq.ListUserManageRequest{Account: req.Account}); err != nil {
		return err
	} else if userList, ok := res.List.(model.Users); !ok {
		return errors.New("用户没有找打")
	} else {
		user = userList[0]
	}

	return ResetPwd(user.Account, req.Password, req.RePassword)
}
