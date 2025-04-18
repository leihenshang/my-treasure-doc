package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/config"
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

		if config.GetConfig().Debug.EnableMockLogin {
			if config.GetConfig().Debug.MockUserId != "" {
				mockUser.Id = config.GetConfig().Debug.MockUserId
			}
			c.Set("userinfo", mockUser)
		} else {
			if authKey == "" {
				result.Msg = "参数错误"
				c.AbortWithStatusJSON(http.StatusUnauthorized, result)
				return
			}

			u, err := service.GetUserByToken(authKey)
			if err != nil {
				global.Log.Error(err)
				result.Msg = err.Error()
				c.AbortWithStatusJSON(http.StatusUnauthorized, result)
				return
			}
			c.Set("userinfo", u)
		}

		c.Next()
	}

}
