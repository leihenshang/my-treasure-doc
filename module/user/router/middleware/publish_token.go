package middleware

import (
	"net/http"
	"strings"

	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

// PublishTokenHeader 内容发布接口的机器令牌请求头。
const PublishTokenHeader = "X-Publish-Token"

// RequirePublishToken 校验内容发布接口的机器令牌（系统设置 publish.apiToken）中间件。
// 仅用于 /api/publish/*。令牌未配置或校验失败一律 401；发布是写操作，务必确认令牌强度与来源。
func RequirePublishToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := global.EffectivePublishApiToken()
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40100, "msg": "未配置发布令牌", "data": nil})
			return
		}
		// 恒定时间比较，避免简单时序侧信道；EqualFold 在意外长度时提前返回，可接受
		if !strings.EqualFold(c.GetHeader(PublishTokenHeader), token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40100, "msg": "发布令牌无效", "data": nil})
			return
		}
		c.Next()
	}
}
