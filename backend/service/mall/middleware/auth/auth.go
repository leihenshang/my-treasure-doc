package auth

import (
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/global"
	"fastduck/treasure-doc/service/mall/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserInfoKey = "userinfo"
)

//Auth 身份验证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.GetHeader("X-Token")
		if authKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "身份认证参数错误")
			return
		}

		u, err := service.GetUserByToken(authKey)
		if err != nil {
			global.ZapSugar.Errorf("[auth|service.GetUserByToken] an error occurred.err:%+v, authKey:%+v ", err, authKey)
			c.AbortWithStatusJSON(http.StatusUnauthorized, "查询用户信息失败")
			return
		}
		if u == nil {
			global.ZapSugar.Errorf("[auth|service.GetUserByToken] an error occurred.err:%+v, authKey:%+v ", "没有找到用户信息")
			c.AbortWithStatusJSON(http.StatusUnauthorized, "没有找到用户信息")
			return
		}

		//把用户信息写入 context
		c.Set("userinfo", u)

		c.Next()
	}

}

func userAccessParamsCheck(c *gin.Context) {

}

//GetUserInfoByCtx 从上下文获取用户信息
func GetUserInfoByCtx(c *gin.Context) (u *model.User, err error) {
	if v, exists := c.Get(UserInfoKey); !exists {
		return nil, errors.New("从上下文中获取用户信息键失败")
	} else {
		if u, ok := v.(*model.User); ok {
			return u, nil
		}
	}

	return nil, errors.New("从用户键值解析用户信息失败")
}
