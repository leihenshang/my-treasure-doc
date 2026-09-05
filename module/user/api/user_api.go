package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/request/user"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/auth"
	"fastduck/treasure-doc/module/user/internal/service"
)

type UserApi struct {
	UserService *service.UserService
}

func NewUserApi() *UserApi {
	return &UserApi{UserService: service.NewUserService()}
}

// UserCaptcha 生成图形验证码，登录前匿名获取
func (u *UserApi) UserCaptcha(c *gin.Context) {
	captcha, err := service.GenCaptcha(c)
	if err != nil {
		global.Log.Errorf("failed to generate captcha:%v", err)
		response.FailWithMessage(c, "验证码生成失败")
		return
	}

	response.OkWithData(c, captcha)
}

// UserLogin 用户登录，账号字段支持填入账号和邮箱，因为都是唯一的
func (u *UserApi) UserLogin(c *gin.Context) {
	var login user.LoginRequest
	err := c.ShouldBindJSON(&login)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	if u, ok := u.UserService.UserLogin(login, c.ClientIP()); ok != nil {
		if errors.Is(ok, service.ErrCaptchaInvalid) {
			response.FailWithMessage(c, ok.Error(), response.CaptchaInvalid)
			return
		}
		response.FailWithMessage(c, ok.Error())
	} else {
		response.OkWithData(c, u)
	}
}

// UserLogout 用户退出登陆
func (u *UserApi) UserLogout(c *gin.Context) {
	loginUser, err := auth.GetUserInfoByCtx(c)
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
