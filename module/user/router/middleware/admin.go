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

		// 默认管理员首次启动会被标记强制改密：改密完成前不允许调用管理接口。
		// 修改密码接口挂在 Auth 组（/api/user/change-pwd），不经过这里，因此不会把自己锁死。
		if user.RequirePwdReset {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 40301, "msg": "请先修改默认密码", "data": nil})
			return
		}

		c.Next()
	}
}
