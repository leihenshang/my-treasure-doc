package service

import (
	"sync"

	"fastduck/treasure-doc/service/user/data/model"
	userReq "fastduck/treasure-doc/service/user/data/request/user"
)

type UserManageService struct{}

var userManageService *UserManageService

var userManageOnce = sync.Once{}

func NewUserManageService() *UserManageService {
	userOnce.Do(func() {
		userManageService = &UserManageService{}
	})
	return userManageService
}

func (u *UserManageService) List() {

}

func (u *UserManageService) Create(user *model.User) (createdUser *model.User, err error) {
	regRequest := userReq.UserRegisterRequest{
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

func (u *UserManageService) Delete() {

}

func (u *UserManageService) Update() {

}

func (u *UserManageService) Detail() {

}

func (u *UserManageService) RestPwd() {

}
