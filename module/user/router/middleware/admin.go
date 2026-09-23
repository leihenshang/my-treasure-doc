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

		// 强制改密改为「前端提醒、可单次忽略」：中间件不再拦截管理接口，
		// 未改密也能继续使用后台；每次登录响应都会带 requirePwdReset，由前端提醒，
		// 重新登录后若仍未改密会再次提醒。
		c.Next()
	}
}
