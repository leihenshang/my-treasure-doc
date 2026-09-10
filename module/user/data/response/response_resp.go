package response

import (
	"net/http"

	commonresponse "fastduck/treasure-doc/module/common/response"

	"github.com/gin-gonic/gin"
)

// Response 与统一响应体结构一致，保留别名便于按包引用。
type Response = commonresponse.Envelope

// ErrorCode 为业务码，取值见 code_error_resp.go。
type ErrorCode = int

// Result 始终返回 HTTP 200，业务结果由 code 表达。
func Result(c *gin.Context, code ErrorCode, data interface{}, msg string) {
	commonresponse.JSON(c, http.StatusOK, code, msg, data)
}

func Ok(c *gin.Context) {
	Result(c, SUCCESS, map[string]interface{}{}, "操作成功")
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "操作成功")
}

// FailWithMessage 失败响应；不传 code 时使用通用失败码。
func FailWithMessage(c *gin.Context, message string, code ...ErrorCode) {
	if len(code) == 0 {
		Result(c, ERROR, map[string]interface{}{}, message)
		return
	}
	Result(c, code[0], map[string]interface{}{}, message)
}
