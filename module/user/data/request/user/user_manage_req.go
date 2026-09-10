package user

import "fastduck/treasure-doc/module/user/data/request"

// ChangePwdRequest 当前登录用户修改自己的密码，需要校验原密码
type ChangePwdRequest struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" binding:"required"`
	Password    string `json:"password" form:"password" binding:"required,min=8,max=16"`
	RePassword  string `json:"rePassword" form:"rePassword" binding:"required,min=8,max=16"`
}

// RestPwdRequest 管理员重置指定账号的密码，无需原密码
type RestPwdRequest struct {
	Password   string `json:"password" form:"password" binding:"required,min=8,max=16"`
	RePassword string `json:"rePassword" form:"rePassword" binding:"required,min=8,max=16"`
	Account    string `json:"account" form:"account" binding:"required"`
}

type ListUserManageRequest struct {
	Id      string `json:"id" form:"id" param:"id"`
	Keyword string `json:"keyword" form:"keyword" binding:""`
	Account string `json:"account" form:"account" binding:""`
	request.Pagination
	request.Sort
}
