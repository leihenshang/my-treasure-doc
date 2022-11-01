package router

import (
	api "fastduck/treasure-doc/service/mall/api/doc"
	apiDoc "fastduck/treasure-doc/service/mall/api/doc"
	apiGoods "fastduck/treasure-doc/service/mall/api/mall/goods"
	apiOrder "fastduck/treasure-doc/service/mall/api/mall/order"
	apiPay "fastduck/treasure-doc/service/mall/api/mall/pay"
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

	mallGroupRouteNoAuth := r.Group("mall")
	{
		//-----商品-------
		//列表
		mallGroupRouteNoAuth.GET("/goods/list", apiGoods.List)
		//详情
		mallGroupRouteNoAuth.GET("/goods/detail", apiGoods.Detail)

	}

	mallGroupRoute := r.Group("mall").Use(auth.Auth())
	{
		//-----订单-----
		//创建
		mallGroupRoute.POST("/order/create", apiOrder.Create)
		//列表
		mallGroupRoute.GET("/order/list", apiOrder.List)
		//详情
		mallGroupRoute.GET("/order/detail", apiOrder.Detail)

		//-----支付-----
		//支付
		mallGroupRoute.POST("/pay/create", apiPay.Create)
	}

}
