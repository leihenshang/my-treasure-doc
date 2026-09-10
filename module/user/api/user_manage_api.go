package api

import (
	"fastduck/treasure-doc/module/user/data/request/user"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/auth"

	"github.com/gin-gonic/gin"
)

// UserChangePwd 当前登录用户修改自己的密码，需要校验原密码
func (u *UserApi) UserChangePwd(c *gin.Context) {
	var req user.ChangePwdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	loginUser, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if err := u.UserService.ChangePassword(loginUser.Id, loginUser.Token, req.OldPassword, req.Password, req.RePassword); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	response.Ok(c)
}
