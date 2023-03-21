package auth

import (
	"github.com/gin-gonic/gin"
	"go_awd/dao"
	"go_awd/pkg/e"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"net/http"
)

func Role(role string) gin.HandlerFunc {
	// role权限校验中间件
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
			return
		}
		if role == "admin" && claims.(*wjwt.Claims).Role == role { // 判断是否是管理员
			c.Next()
		} else if role == "leader" { // 判断是否是队长
			userDao := dao.NewUserDao(c)
			user, err := userDao.GetUserByID(claims.(*wjwt.Claims).ID)
			if err != nil || !user.IsTeamLeader {
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
			return
		}
	}
}
