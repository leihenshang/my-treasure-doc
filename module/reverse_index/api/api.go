package api

import (
	"fastduck/treasure-doc/module/reverse_index/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexReq struct {
	Content []string `json:"content"`
}

func Index(c *gin.Context) {
	var req IndexReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 校验content是否为空
	if len(req.Content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "content is empty",
		})
		return
	}
	err := service.Index(req.Content...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "index success",
	})
}
func Search(c *gin.Context) {
	keyword := c.Query("keyword")
	// 校验keyword是否为空
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "keyword is empty",
		})
		return
	}
	results, err := service.Search(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"results": service.List(),
	})
}
