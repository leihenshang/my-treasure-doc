package router

import (
	"net/http"

	"fastduck/treasure-doc/service/user/config"
	"fastduck/treasure-doc/service/user/middleware"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/api"
)

func InitRoute(r *gin.Engine) {
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

	//user
	userRoute := apiBase.Group("user")
	{
		userRoute.POST("/reg", api.UserRegister)
		userRoute.POST("/login", api.UserLogin)
	}

	userRoute.Use(middleware.Auth())
	{
		userRoute.POST("/logout", api.UserLogout)
		userRoute.POST("/update-profile", api.UserProfileUpdate)
	}

	//doc
	docRoute := apiBase.Group("doc").Use(middleware.Auth(), middleware.Cors())
	{
		docRoute.POST("/create", api.DocCreate)
		docRoute.GET("/detail", api.DocDetail)
		docRoute.GET("/list", api.DocList)
		docRoute.POST("/update", api.DocUpdate)
		docRoute.POST("/delete", api.DocDelete)
		docRoute.GET("/tree", api.DocTree)
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
