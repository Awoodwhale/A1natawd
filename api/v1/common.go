package v1

import (
	"github.com/gin-gonic/gin"
	"go_awd/pkg/e"
	"go_awd/serializer"
)

// ErrorResponse
// @Description: 获取错误信息resp
// @param err error
// @param service any
// @return serializer.Response
func ErrorResponse(err error, service any, c *gin.Context) serializer.Response {
	return serializer.Response{
		Code:    e.Error,
		Message: e.GetMessageByCode(e.Error, c),
		Error:   e.HandleBindingError(err, service, c),
	}
}
