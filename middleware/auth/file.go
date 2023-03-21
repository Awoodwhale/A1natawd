package auth

import (
	"github.com/gin-gonic/gin"
	"go_awd/conf"
	"go_awd/pkg/e"
	"go_awd/serializer"
	"net/http"
	"path/filepath"
	"strings"
)

func FileTypeMiddleware() gin.HandlerFunc {
	// token验证中间件
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err == nil { // 存在文件
			extension := filepath.Ext(file.Filename)
			// 验证扩展名
			for _, ext := range conf.ImageTypes {
				if strings.HasSuffix(extension, ext) {
					// 允许上传
					c.Next()
					return
				} else {
					c.AbortWithStatusJSON(http.StatusBadRequest, serializer.RespCode(e.InvalidWithFileType, c))
					return
				}
			}
		}
		// 不存在文件那就直接通过
		c.Next()
	}
}
