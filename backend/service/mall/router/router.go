package router

import (
	api "fastduck/treasure-doc/service/mall/api/doc"
	apiDoc "fastduck/treasure-doc/service/mall/api/doc"
	apiUser "fastduck/treasure-doc/service/mall/api/user"
	"fastduck/treasure-doc/service/mall/middleware/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {
	//测试连通性
	base := r.Group("/")
	{
		base.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "pong!",
			})
		})
	}

	//user
	userRoute := r.Group("user")
	//不需要验证登录参数
	{
		userRoute.POST("/reg", apiUser.UserRegister)
		userRoute.POST("/login", apiUser.UserLogin)
	}
	//需要验证登录参数
	userRoute.Use(auth.Auth())
	{
		userRoute.POST("/logout", apiUser.UserLogout)
		userRoute.POST("/updateProfile", apiUser.UserProfileUpdate)
	}

	//doc
	docRoute := r.Group("doc").Use(auth.Auth())
	{
		docRoute.POST("/create", api.DocCreate)
		docRoute.POST("/detail", api.DocDetail)
		docRoute.POST("/list", api.DocList)
		docRoute.POST("/update", api.DocUpdate)
		docRoute.POST("/delete", api.DocDelete)
	}

	//doc group
	docGroupRoute := r.Group("doc-group").Use(auth.Auth())
	{
		docGroupRoute.POST("/create", apiDoc.DocGroupCreate)
		docGroupRoute.POST("/list", apiDoc.DocGroupList)
		docGroupRoute.POST("/update", apiDoc.DocGroupUpdate)
		docGroupRoute.POST("/delete", apiDoc.DocGroupDelete)
	}

	mallGroupRoute := r.Group("mall").Use(auth.Auth())
	{
		//-----商品-------
		//列表
		mallGroupRoute.GET("/goods/list", nil)
		//详情
		mallGroupRoute.GET("/goods/detail", nil)

		//-----订单-----
		//创建
		mallGroupRoute.POST("/order/create", nil)
		//列表
		mallGroupRoute.GET("/order/list", nil)
		//详情
		mallGroupRoute.GET("/order/detail", nil)

		//-----支付-----
		//支付
		mallGroupRoute.POST("/pay/create", nil)
	}

}
