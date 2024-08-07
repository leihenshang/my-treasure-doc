package api

import (
	"fastduck/treasure-doc/service/user/middleware"
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/data/request/user"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var reg user.UserRegisterRequest
	err := c.ShouldBindJSON(&reg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if u, ok := service.UserRegister(reg); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(u, c)
	}
}

// UserLogin 用户登录，账号字段支持填入账号和邮箱，因为都是唯一的
func UserLogin(c *gin.Context) {
	var login user.UserLoginRequest
	err := c.ShouldBindJSON(&login)
	if err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	if u, ok := service.UserLogin(login, c.ClientIP()); ok != nil {
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(u, c)
	}
}

// UserLogout 用户退出登陆
func UserLogout(c *gin.Context) {
	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.UserLogout(u.Id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.Ok(c)
}

// UserProfileUpdate 更新用户个人资料
func UserProfileUpdate(c *gin.Context) {
	var profile user.UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&profile); err != nil {
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := service.UserProfileUpdate(profile, u.Id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.Ok(c)
	}

}
