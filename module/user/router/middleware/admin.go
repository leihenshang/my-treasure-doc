package middleware

import (
	"net/http"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(global.UserInfoKey)
		user, ok := value.(*model.User)
		if !exists || !ok || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40100, "msg": "未登录或登录已失效", "data": nil})
			return
		}

		if user.UserType != model.UserTypeAdmin && user.UserType != model.UserTypeRoot {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 40300, "msg": "无管理权限", "data": nil})
			return
		}

		c.Next()
	}
}
