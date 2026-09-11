package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
)

var mockUser = &model.User{
	BaseModel: model.BaseModel{
		Id: "9999999999",
	},
	Nickname:   "mockUser9999999999",
	Account:    "mockUser9999999999",
	Email:      "9999999999",
	Password:   "9999999999",
	UserType:   100,
	UserStatus: 1,
	Mobile:     "9999999999",
	Avatar:     "",
	Bio:        "mockUser9999999999",
	Token:      "mockUser9999999999",
}

// Auth 身份验证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.GetHeader("X-Token")
		result := &response.Response{Code: response.ERROR}
		cfg := global.GetConf()

		if cfg != nil && cfg.App.IsDev() && cfg.Debug.EnableMockLogin {
			// mock 登录仅在开发模式显式开启时生效；生产（release）模式下永远走真实鉴权，
			// 不会因误配而暴露管理员后台。
			requestUser := *mockUser
			if cfg.Debug.MockUserId != "" {
				requestUser.Id = cfg.Debug.MockUserId
			}
			c.Set(global.UserInfoKey, &requestUser)
		} else {
			if authKey == "" {
				result.Msg = "请先登录"
				c.AbortWithStatusJSON(http.StatusUnauthorized, result)
				return
			}

			u, err := service.GetUserByToken(authKey)
			if err != nil {
				global.Log.Error(err)
				result.Msg = "认证失败，请重新登录"
				c.AbortWithStatusJSON(http.StatusUnauthorized, result)
				return
			}
			c.Set(global.UserInfoKey, u)
		}

		c.Next()
	}

}
