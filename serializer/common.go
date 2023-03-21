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

type ResponseOpt struct {
	f func(opt *Response)
}

func WithCode(code int) ResponseOpt {
	return ResponseOpt{func(opt *Response) {
		opt.Code = code
	}}
}

func WithData(data any) ResponseOpt {
	return ResponseOpt{func(opt *Response) {
		opt.Data = data
	}}
}

func WithMsg(msg string) ResponseOpt {
	return ResponseOpt{func(opt *Response) {
		opt.Message = msg
	}}
}

func WithErr(err error) ResponseOpt {
	return ResponseOpt{func(opt *Response) {
		opt.Error = err.Error()
	}}
}

func NewResp(ops ...ResponseOpt) Response {
	resp := &Response{}
	for _, do := range ops {
		do.f(resp)
	}
	return *resp
}

// ListData
// @Description: 列表数据
type ListData struct {
	Items any  `json:"items"`
	Total uint `json:"total"`
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

// BuildListResponse
// @Description: 返回列表数据resp
// @param item any
// @param total uint
// @return Response
func BuildListResponse(items any, total uint, c *gin.Context) Response {
	return Response{
		Code:    e.Success,
		Message: e.GetMessageByCode(e.Success, c),
		Data:    ListData{Items: items, Total: total},
	}
}
