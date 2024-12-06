package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/middleware"

	"fastduck/treasure-doc/service/user/data/request/user"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"
)

type UserApi struct {
	UserService *service.UserService
}

func NewUserApi() *UserApi {
	return &UserApi{UserService: service.NewUserService()}
}

// UserRegister 用户注册
func (u *UserApi) UserRegister(c *gin.Context) {
	var reg user.UserRegisterRequest
	err := c.ShouldBindJSON(&reg)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if u, ok := u.UserService.UserRegister(reg); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(u, c)
	}
}

// UserLogin 用户登录，账号字段支持填入账号和邮箱，因为都是唯一的
func (u *UserApi) UserLogin(c *gin.Context) {
	var login user.UserLoginRequest
	err := c.ShouldBindJSON(&login)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	if u, ok := u.UserService.UserLogin(login, c.ClientIP()); ok != nil {
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(u, c)
	}
}

// UserLogout 用户退出登陆
func (u *UserApi) UserLogout(c *gin.Context) {
	loginUser, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if err := u.UserService.UserLogout(loginUser.Id, loginUser.Token); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	response.Ok(c)
}

// UserProfileUpdate 更新用户个人资料
func (u *UserApi) UserProfileUpdate(c *gin.Context) {
	var profile user.UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&profile); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	loginUser, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if _, err := u.UserService.UserProfileUpdate(profile, loginUser.Id); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	} else {
		response.Ok(c)
	}

}
