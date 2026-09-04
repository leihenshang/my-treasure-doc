package response

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Code: 0, Msg: "", Data: normalize(data)})
}
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Code: 0, Msg: "", Data: normalize(data)})
}
func Error(c *gin.Context, status, code int, message string) {
	c.JSON(status, Envelope{Code: code, Msg: message, Data: nil})
}

func normalize(value interface{}) interface{} {
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
