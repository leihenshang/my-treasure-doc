package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/api"
	"fastduck/treasure-doc/service/user/middleware"
)

func InitRoute(r *gin.Engine) {

	//静态资源
	r.Static("/static", "./static")

	//测试连通性
	base := r.Group("/")
	{
		base.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "pong!",
			})
		})
	}

	apiBase := r.Group("api")

	//user
	userRoute := apiBase.Group("user")
	//不需要验证登录参数
	{
		userRoute.POST("/reg", api.UserRegister)
		userRoute.POST("/login", api.UserLogin)
	}
	//需要验证登录参数
	userRoute.Use(middleware.Auth())
	{
		userRoute.POST("/logout", api.UserLogout)
		userRoute.POST("/update-profile", api.UserProfileUpdate)
	}

	//doc
	docRoute := apiBase.Group("doc").Use(middleware.Auth(), middleware.Cors())
	{
		docRoute.POST("/create", api.DocCreate)
		docRoute.POST("/detail", api.DocDetail)
		docRoute.GET("/list", api.DocList)
		docRoute.POST("/update", api.DocUpdate)
		docRoute.POST("/delete", api.DocDelete)
	}

	//doc group
	docGroupRoute := apiBase.Group("doc-group").Use(middleware.Auth())
	{
		docGroupRoute.POST("/create", api.DocGroupCreate)
		docGroupRoute.POST("/list", api.DocGroupList)
		docGroupRoute.POST("/update", api.DocGroupUpdate)
		docGroupRoute.POST("/delete", api.DocGroupDelete)
		docGroupRoute.GET("/tree", api.DocGroupTree)
	}

	// file upload
	fileGroupRoute := apiBase.Group("file").Use(middleware.Auth())
	{
		fileGroupRoute.POST("upload", api.FileUpload)
	}
}
