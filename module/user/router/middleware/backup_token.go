package middleware

import (
	"net/http"
	"strings"

	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

// BackupTokenHeader NAS 备份下载的机器令牌请求头。
const BackupTokenHeader = "X-Backup-Token"

// BackupToken 校验 NAS 机器令牌（后台配置 backup.apiToken，回退到配置文件 [backup].apiToken）的中间件，
// 仅用于 /api/backup/export。令牌未配置或校验失败一律 401；不参与登录/验证码链，适合定时任务免登录拉取。
func BackupToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := global.EffectiveBackupApiToken()
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40100, "msg": "未配置 NAS 备份令牌", "data": nil})
			return
		}
		if !strings.EqualFold(c.GetHeader(BackupTokenHeader), token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40100, "msg": "NAS 备份令牌无效", "data": nil})
			return
		}
		c.Next()
	}
}
