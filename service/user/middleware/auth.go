package middleware

import (
	"errors"
	"net/http"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/response"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	UserInfoKey = "userinfo"
)

// Auth 身份验证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.GetHeader("X-Token")
		result := &response.Response{Code: response.ERROR}
		if authKey == "" {
			result.Msg = "参数错误"
			c.AbortWithStatusJSON(http.StatusOK, result)
			return
		}

		u, err := service.GetUserByToken(authKey)
		if err != nil {
			global.ZAPSUGAR.Error(err)
			result.Msg = "查询用户信息失败"
			c.AbortWithStatusJSON(http.StatusOK, result)
			return
		}
		if u == nil {
			global.ZAPSUGAR.Error("获取用户信息失败")
			result.Msg = "获取用户信息失败"
			c.AbortWithStatusJSON(http.StatusOK, result)
			return
		}

		//把用户信息写入 context
		c.Set("userinfo", u)

		c.Next()
	}

}

func userAccessParamsCheck(c *gin.Context) {

}

// GetUserInfoByCtx 从上下文获取用户信息
func GetUserInfoByCtx(c *gin.Context) (u *model.User, err error) {
	if v, exists := c.Get(UserInfoKey); !exists {
		return nil, errors.New("从上下文中获取用户信息失败")
	} else {
		if u, ok := v.(*model.User); ok {
			return u, nil
		}
	}

	return nil, errors.New("从上下文解析用户信息失败")
}
