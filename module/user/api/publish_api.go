package api

import (
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

// PublishTokenRequest 更新内容发布令牌的请求体。
type PublishTokenRequest struct {
	// Token 新发布令牌；空字符串表示清除（禁用发布接口）。
	Token string `json:"token"`
}

// PublishApi 管理内容发布接口的机器令牌（仅后台管理员可操作）。
type PublishApi struct{}

func NewPublishApi() *PublishApi { return &PublishApi{} }

// GetToken 返回当前生效的发布令牌。
func (p *PublishApi) GetToken(c *gin.Context) {
	response.OkWithData(c, gin.H{"token": global.EffectivePublishApiToken()})
}

// PutToken 更新发布令牌（写入 DB，即时生效）。
func (p *PublishApi) PutToken(c *gin.Context) {
	var req PublishTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, "参数错误")
		return
	}
	if err := global.SetSystemSetting(global.SystemSettingPublishTokenKey, req.Token); err != nil {
		response.FailWithMessage(c, "保存失败")
		return
	}
	response.OkWithData(c, gin.H{"token": req.Token})
}
