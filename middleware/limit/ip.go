package ip

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/pkg/e"
	"go_awd/serializer"
	"net/http"
	"time"
)

// LimitMiddleware
// @Description: 限制ip在一定时间内的访问接口的次数
// @param limit int64
// @param duration time.Duration
// @return gin.HandlerFunc
func LimitMiddleware(limit int64, duration time.Duration, operation string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := cache.ClientOperationIPKey(operation, c.ClientIP())
		count, err := cache.RedisClient.Incr(key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, serializer.RespErr(e.ErrorWithRedis, err, c))
			return
		}
		if count == 1 {
			if err := cache.RedisClient.Expire(key, duration).Err(); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, serializer.RespErr(e.ErrorWithRedis, err, c))
				return
			}
		}
		if count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, serializer.RespCode(e.InvalidTooManyRequest, c))
			return
		}
		c.Next()
	}
}
