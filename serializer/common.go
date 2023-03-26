package serializer

import (
	"github.com/gin-gonic/gin"
	"go_awd/pkg/e"
)

// Response
// @Description: 响应
type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ListData
// @Description: 列表数据
type ListData struct {
	Items any `json:"items"`
	Total int `json:"total"`
}

func RespSuccess(code int, data any, c *gin.Context) Response {
	return Response{
		Code:    code,
		Message: e.GetMessageByCode(code, c),
		Data:    data,
	}
}

func RespErr(code int, err error, c *gin.Context) Response {
	return Response{
		Code:    code,
		Message: e.GetMessageByCode(code, c),
		Error:   err.Error(),
	}
}

func RespCode(code int, c *gin.Context) Response {
	return Response{
		Code:    code,
		Message: e.GetMessageByCode(code, c),
	}
}

func RespList(code int, items any, total int, c *gin.Context) Response {
	return Response{
		Code:    code,
		Message: e.GetMessageByCode(code, c),
		Data:    ListData{Items: items, Total: total},
	}
}
