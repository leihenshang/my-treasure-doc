package api

import (
	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"
	"fastduck/treasure-doc/service/user/middleware"
)

type UserManageApi struct {
	UserService *service.UserManageService
}

func NewUserManageApi() *UserManageApi {
	return &UserManageApi{UserService: service.NewUserManageService()}
}

func (u *UserManageApi) Create(c *gin.Context) {
	var req *model.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(c, global.ErrResp(err))
		return
	}

	user, err := middleware.GetUserInfoByCtx(c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	if createdUser, err := u.UserService.Create(user); err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.OkWithData(c, createdUser)
	}
}

func (u *UserManageApi) List(c *gin.Context) {
	response.OkWithData(c, u)
}

func (u *UserManageApi) Detail(c *gin.Context) {
	response.OkWithData(c, u)
}

func (u *UserManageApi) Delete(c *gin.Context) {
	response.OkWithData(c, u)
}

func (u *UserManageApi) Update(c *gin.Context) {
	response.OkWithData(c, u)
}

func (u *UserManageApi) ResetPwd(c *gin.Context) {
	response.OkWithData(c, u)
}
