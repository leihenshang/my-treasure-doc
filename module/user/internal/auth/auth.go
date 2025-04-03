package auth

import (
	"errors"
	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

// GetUserInfoByCtx 从上下文获取用户信息
func GetUserInfoByCtx(c *gin.Context) (u *model.User, err error) {
	if v, exists := c.Get(global.UserInfoKey); !exists {
		return nil, errors.New("从上下文中获取用户信息失败")
	} else {
		if u, ok := v.(*model.User); ok {
			return u, nil
		}
	}

	return nil, errors.New("从上下文解析用户信息失败")
}
