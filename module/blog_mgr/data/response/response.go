package response

import (
	"net/http"

	commonresponse "fastduck/treasure-doc/module/common/response"

	"github.com/gin-gonic/gin"
)

type Envelope = commonresponse.Envelope

// OK / Created 返回的数据来自数据库模型（无 json tag），因此先做 key 归一化。
func OK(c *gin.Context, data interface{}) {
	commonresponse.JSON(c, http.StatusOK, commonresponse.CodeSuccess, "", commonresponse.Normalize(data))
}

func Created(c *gin.Context, data interface{}) {
	commonresponse.JSON(c, http.StatusCreated, commonresponse.CodeSuccess, "", commonresponse.Normalize(data))
}

func Error(c *gin.Context, status, code int, message string) {
	commonresponse.Error(c, status, code, message)
}
