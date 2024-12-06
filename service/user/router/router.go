package router

import (
	"net/http"

	"fastduck/treasure-doc/service/user/config"
	"fastduck/treasure-doc/service/user/middleware"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/api"
)

func InitRouter(r *gin.Engine) {
	r.Static("/web", config.WebPath)
	r.Static("/files", config.FilesPath)

	r.Any("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong!",
		})
	})

	apiBase := r.Group("api")
	{
		userApi := api.NewUserApi()
		userRoute := apiBase.Group("user").Use(middleware.Cors())
		userRoute.POST("/reg", userApi.UserRegister)
		userRoute.POST("/login", userApi.UserLogin)
		userRoute.Use(middleware.Auth(), middleware.Cors())
		userRoute.POST("/logout", userApi.UserLogout)
		userRoute.POST("/update-profile", userApi.UserProfileUpdate)
	}

	//doc
	docRoute := apiBase.Group("doc").Use(middleware.Auth(), middleware.Cors())
	{
		docApi := api.NewDocApi()
		docRoute.POST("/create", docApi.DocCreate)
		docRoute.GET("/detail", docApi.DocDetail)
		docRoute.GET("/list", docApi.DocList)
		docRoute.POST("/update", docApi.DocUpdate)
		docRoute.POST("/delete", docApi.DocDelete)
		docRoute.GET("/tree", docApi.DocTree)
	}

	//doc-history
	docHistoryRoute := apiBase.Group("doc-history").Use(middleware.Auth(), middleware.Cors())
	{
		docHistoryRoute.POST("/recover", api.DocHistoryRecover)
		docHistoryRoute.GET("/detail", api.DocHistoryDetail)
		docHistoryRoute.GET("/list", api.DocHistoryList)
	}

	//note
	noteRoute := apiBase.Group("note").Use(middleware.Auth(), middleware.Cors())
	{
		noteRoute.POST("/create", api.NoteCreate)
		noteRoute.GET("/detail", api.NoteDetail)
		noteRoute.GET("/list", api.NoteList)
		noteRoute.POST("/update", api.NoteUpdate)
		noteRoute.POST("/delete", api.NoteDelete)
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
