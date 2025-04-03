package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/data/request"
)

type ListResponse struct {
	Pagination request.Pagination `json:"pagination"`
	List       interface{}        `json:"list"`
}

type Response struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Result(c *gin.Context, code ErrorCode, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(c, SUCCESS, map[string]interface{}{}, "操作成功")
}

func OkWithMessage(c *gin.Context, message string) {
	Result(c, SUCCESS, map[string]interface{}{}, message)
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "操作成功")
}

func OkWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, SUCCESS, data, message)
}

func Fail(c *gin.Context) {
	Result(c, ERROR, map[string]interface{}{}, "操作失败")
}

func FailWithMessage(c *gin.Context, message string, code ...ErrorCode) {
	if len(code) == 0 {
		Result(c, ERROR, map[string]interface{}{}, message)
		return
	}
	Result(c, code[0], map[string]interface{}{}, message)
}

func FailWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, ERROR, data, message)
}

type ErrorCode int
