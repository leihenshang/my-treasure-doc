// Package response 提供后端统一的 HTTP 响应封装。
//
// 所有接口都返回 {code, msg, data}，其中 code=0 表示成功；业务错误码由各模块自行定义。
package response

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

const (
	// CodeSuccess 成功业务码
	CodeSuccess = 0
	// CodeError 通用失败业务码
	CodeError = 1
)

// Envelope 是所有接口共用的响应体。
type Envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// JSON 按给定 HTTP 状态码与业务码写出响应，data 原样序列化。
func JSON(c *gin.Context, status, code int, msg string, data interface{}) {
	c.JSON(status, Envelope{Code: code, Msg: msg, Data: data})
}

// OK 成功响应：HTTP 200 + code 0。
func OK(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, CodeSuccess, "", data)
}

// Created 创建成功响应：HTTP 201 + code 0。
func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, CodeSuccess, "", data)
}

// Error 失败响应：HTTP 状态码与业务码分离，data 固定为 null。
func Error(c *gin.Context, status, code int, msg string) {
	JSON(c, status, code, msg, nil)
}

// Normalize 把结构体字段名转为小驼峰（ID → id、PublicID → publicID），
// 供没有 json tag 的模型直接返回给前端；已带小驼峰 tag 的字段保持不变。
func Normalize(value interface{}) interface{} {
	data, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var decoded interface{}
	if json.Unmarshal(data, &decoded) != nil {
		return value
	}
	return lowerKeys(decoded)
}

func lowerKeys(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			result[lowerCamel(key)] = lowerKeys(item)
		}
		return result
	case []interface{}:
		for index := range typed {
			typed[index] = lowerKeys(typed[index])
		}
	}
	return value
}

// lowerCamel 只把开头连续的缩写或首字母小写：ID → id、PublicID → publicID、URLValue → urlValue。
func lowerCamel(value string) string {
	if value == "" {
		return value
	}
	runes := []rune(value)
	end := 1
	for end < len(runes) && unicode.IsUpper(runes[end]) {
		if end+1 < len(runes) && unicode.IsLower(runes[end+1]) {
			break
		}
		end++
	}
	return strings.ToLower(string(runes[:end])) + string(runes[end:])
}
