package response

import (
	"net/http"

	commonresponse "fastduck/treasure-doc/module/common/response"

	"github.com/gin-gonic/gin"
)

// 公开只读接口的业务码
const (
	CodeSuccess          = commonresponse.CodeSuccess
	CodeInvalidQuery     = 40001
	CodeUnsupportedSort  = 40002
	CodePostNotFound     = 40401
	CodeDiaryNotFound    = 40402
	CodePortfolioMissing = 40403
	CodeToolNotFound     = 40404
	CodeInternal         = 50000
)

type Envelope = commonresponse.Envelope

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Total    int64  `json:"total"`
	OrderBy  string `json:"orderBy"`
}

type Page struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

func OK(c *gin.Context, data interface{}) {
	commonresponse.JSON(c, http.StatusOK, CodeSuccess, "", data)
}

func Error(c *gin.Context, status, code int, message string) {
	commonresponse.Error(c, status, code, message)
}
