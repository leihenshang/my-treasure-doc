package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess          = 0
	CodeInvalidQuery     = 40001
	CodeUnsupportedSort  = 40002
	CodePostNotFound     = 40401
	CodeDiaryNotFound    = 40402
	CodePortfolioMissing = 40403
	CodeToolNotFound     = 40404
	CodeInternal         = 50000
)

type Envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

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
	c.JSON(http.StatusOK, Envelope{Code: CodeSuccess, Msg: "", Data: data})
}

func Error(c *gin.Context, status, code int, message string) {
	c.JSON(status, Envelope{Code: code, Msg: message, Data: nil})
}
