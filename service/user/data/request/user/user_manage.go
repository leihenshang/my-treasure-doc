package user

import "fastduck/treasure-doc/service/user/data/request"

type RestPwdRequest struct {
	Password   string `json:"password" form:"password" binding:"required,min=8,max=16"`
	RePassword string `json:"rePassword" form:"rePassword" binding:"required,min=8,max=16"`
	Account    string `json:"account" form:"account" binding:"required"`
}

type ListUserManageRequest struct {
	Id      int64  `json:"id" form:"id" param:"id"`
	Keyword string `json:"keyword" form:"keyword" binding:""`
	Account string `json:"account" form:"account" binding:""`
	request.Pagination
	request.Sort
}
