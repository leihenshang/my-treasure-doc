package middleware

import (
	"net/http"

	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

// Cors 按白名单反射跨域来源：仅对明确允许的站点返回 Access-Control-Allow-Origin，
// 不使用通配符 "*"，且只在命中具体来源时才携带凭据头，避免任意站点读取后台接口响应。
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := global.GetConf()
		origin := c.GetHeader("Origin")
		if origin != "" && cfg != nil && matchAllowedOrigin(cfg.App.CorsAllowOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-Token")
		c.Header("Access-Control-Max-Age", "3600")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// matchAllowedOrigin 精确匹配白名单来源（区分大小写）。
func matchAllowedOrigin(allowed []string, origin string) bool {
	for _, o := range allowed {
		if o != "" && o == origin {
			return true
		}
	}
	return false
}
