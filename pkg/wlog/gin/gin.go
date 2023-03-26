package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	c "go_awd/pkg/wlog/common"
	"net/http"
	"time"
)

type OptGin struct {
	//example 'user/logout', this route will be ignored in log
	SkipRoute map[string]struct{}
	OptLogrus *c.OptLog
}

/*GetOpt
 * @msg get default parameters
 * @return: *OptGin
 */
func GetOpt() *OptGin {
	opt := c.InitOpt()
	opt.FileNamePrefix = "access.log"
	return &OptGin{
		OptLogrus: opt,
	}
}

/*New
 * @msg create logrus instance for gin middleware
 * @param opt
 * @return: *logrus.Logger  return logrus for configuration in advance
 * @return: gin.HandlerFunc
 * @return: error
 */
func New(opt *OptGin) (*logrus.Logger, gin.HandlerFunc, error) {

	if log, err := opt.OptLogrus.ConfigLogrus(); err != nil {
		return log, func(ctx *gin.Context) {}, errors.Cause(err)
	} else {
		return log,
			func(ctx *gin.Context) {
				if _, ok := opt.SkipRoute[ctx.Request.URL.Path]; ok {
					return
				}
				start := time.Now()
				path := ctx.Request.URL.Path
				raw := ctx.Request.URL.RawQuery
				ctx.Next()
				end := time.Now()
				latency := end.Sub(start) //记录请求处理时间
				clientIP := ctx.ClientIP()
				method := ctx.Request.Method
				statusCode := ctx.Writer.Status()
				//请求大小
				bodySize := ctx.Writer.Size()

				//记录url param
				if raw != "" {
					path = path + "?" + raw
				}
				//设置json字段内容
				entry := log.WithFields(logrus.Fields{
					"Code":    statusCode,
					"Latency": latency, // time to process
					"IP":      clientIP,
					"Method":  method,
					"Path":    path,
					"Size":    bodySize,
				})

				if len(ctx.Errors) > 0 {
					entry.Error(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
				} else {
					//base on response value to match log level
					//msg := fmt.Sprintf("%s - \"%s %s\" %d %d (%dms)", clientIP, method, path, statusCode, bodySize, latency)
					if statusCode >= http.StatusInternalServerError { //500 assign to error level
						entry.Error()
					} else if statusCode >= http.StatusBadRequest { //400 assign to warn level
						entry.Warn()
					} else {
						entry.Info()
					}
				}
			},
			nil
	}
}
