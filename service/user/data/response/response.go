package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/data/request"
)

type ListResponse struct {
	Pagination request.ListPagination `json:"pagination"`
	List       interface{}            `json:"list"`
}

type Response struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Result(code ErrorCode, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(c *gin.Context, message string, code ...ErrorCode) {
	if len(code) == 0 {
		Result(ERROR, map[string]interface{}{}, message, c)
		return
	}
	Result(code[0], map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

type ErrorCode int
