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

	{
		userMgrApi := api.NewUserManageApi()
		userMgrRoute := apiBase.Group("user-manage").Use(middleware.Auth(), middleware.Cors())
		userMgrRoute.POST("/create", userMgrApi.Create)
		userMgrRoute.GET("/detail", userMgrApi.Detail)
		userMgrRoute.GET("/list", userMgrApi.List)
		userMgrRoute.POST("/update", userMgrApi.Update)
		userMgrRoute.POST("/delete", userMgrApi.Delete)
		userMgrRoute.POST("/reset-pwd", userMgrApi.ResetPwd)
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
		docHistoryApi := api.NewDocHistoryApi()
		docHistoryRoute.POST("/recover", docHistoryApi.DocHistoryRecover)
		docHistoryRoute.GET("/detail", docHistoryApi.DocHistoryDetail)
		docHistoryRoute.GET("/list", docHistoryApi.DocHistoryList)
	}

	//note
	noteRoute := apiBase.Group("note").Use(middleware.Auth(), middleware.Cors())
	{
		noteApi := api.NewNoteApi()
		noteRoute.POST("/create", noteApi.NoteCreate)
		noteRoute.GET("/detail", noteApi.NoteDetail)
		noteRoute.GET("/list", noteApi.NoteList)
		noteRoute.POST("/update", noteApi.NoteUpdate)
		noteRoute.POST("/delete", noteApi.NoteDelete)
	}

	//doc group
	docGroupRoute := apiBase.Group("doc-group").Use(middleware.Auth())
	{
		docGroupApi := api.NewDocGroupApi()
		docGroupRoute.POST("/create", docGroupApi.DocGroupCreate)
		docGroupRoute.POST("/list", docGroupApi.DocGroupList)
		docGroupRoute.POST("/update", docGroupApi.DocGroupUpdate)
		docGroupRoute.POST("/delete", docGroupApi.DocGroupDelete)
		docGroupRoute.GET("/tree", docGroupApi.DocGroupTree)
	}

	// file upload
	fileGroupRoute := apiBase.Group("file").Use(middleware.Auth())
	{
		fileApi := api.NewFileApi()
		fileGroupRoute.POST("upload", fileApi.FileUpload)
	}
}
