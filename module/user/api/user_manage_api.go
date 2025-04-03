package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/model"
	userReq "fastduck/treasure-doc/module/user/data/request/user"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/auth"
	"fastduck/treasure-doc/module/user/internal/service"
)

type UserManageApi struct {
	UserManageService *service.UserManageService
}

func NewUserManageApi() *UserManageApi {
	return &UserManageApi{UserManageService: service.NewUserManageService()}
}

func (u *UserManageApi) Create(c *gin.Context) {
	var req *userReq.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	_, err = auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if createdUser, err := u.UserManageService.Create(req); err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.OkWithData(c, createdUser)
	}
}

func (u *UserManageApi) List(c *gin.Context) {
	var req userReq.ListUserManageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	_, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if res, err := u.UserManageService.List(req); err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.OkWithData(c, res)
	}
}

func (u *UserManageApi) Detail(c *gin.Context) {
	var req userReq.ListUserManageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	_, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if res, err := u.UserManageService.Detail(req.Id); err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.OkWithData(c, res)
	}
}

func (u *UserManageApi) Delete(c *gin.Context) {
	var req *model.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	_, err = auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if err = u.UserManageService.Delete(req.Id); err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.Ok(c)
	}
}

func (u *UserManageApi) Update(c *gin.Context) {
	var req *model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}
	_, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if updatedUser, err := u.UserManageService.Update(req); err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.OkWithData(c, updatedUser)
	}
}

func (u *UserManageApi) ResetPwd(c *gin.Context) {
	var req *userReq.RestPwdRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	_, err = auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if err = u.UserManageService.RestPwd(req); err != nil {
		response.FailWithMessage(c, err.Error())
	}
	response.Ok(c)
}
