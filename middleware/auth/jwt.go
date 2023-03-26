package auth

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/pkg/e"
	"go_awd/pkg/util"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
	"net/http"
	"time"
)

// JWT
// @Description: jwt中间件
// @return gin.HandlerFunc
func JWT() gin.HandlerFunc {
	// token验证中间件
	return func(c *gin.Context) {
		var claims *wjwt.Claims
		nowTime := time.Now()
		tokenKey := c.GetHeader("Authorization")
		// 查看accessToken是否过期
		accessToken, err := cache.RedisClient.Get(cache.AccessTokenKey(tokenKey)).Result()
		if err != nil { // accessToken过期了
			// 首先去MySQL中查询refreshToken
			tokenDao := dao.NewTokenDao(c)
			refreshToken, err := tokenDao.GetByTokenKey(tokenKey)
			if err != nil { // 不存在refreshToken
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
			// 找到了refreshToken，验证其有效期
			if refreshClaims, err := wjwt.ParseJWT(refreshToken.RefreshToken); err != nil || nowTime.Unix() > refreshClaims.ExpiresAt {
				// refreshToken也过期了，需要用户重新登录
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
			// 如果refreshToken没有过期，那么生成新的token
			user, err := dao.NewUserDaoByDB(tokenDao.DB).GetUserByID(refreshToken.UserID)
			if err != nil { // 找不到MySQL中的refreshToken
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
			// 如果登录用户当前存在tokenKey在redis中，先删除改tokenKey
			if storedTokenKey, err := cache.RedisClient.Get(cache.StoreAccessTokenKeyKey(user.ID)).Result(); err == nil { // 如果当前用户存在tokenKey，那么删除
				cache.RedisClient.Del(cache.AccessTokenKey(storedTokenKey))
			}
			newAccessToken, err := wjwt.GenAccessToken(time.Now(), user.ID, user.Username, user.Role)
			if err != nil { // 生成new access token失败
				service.Debugln("create new access token error, ", err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
			newTokenKey := util.Sha1(newAccessToken) // 新的accessToken与tokenKey存入redis
			if err := cache.RedisClient.Set(cache.AccessTokenKey(newTokenKey), newAccessToken, wjwt.AccessExpireTime).Err(); err != nil {
				service.Errorln("update newTokenKey to redis error,", err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
			// refreshToken的tokenKey更新
			refreshToken.TokenKey = newTokenKey
			if err := tokenDao.CreateRefreshToken(refreshToken); err != nil {
				service.Errorln("update refreshToken tokenKey error,", err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
			// storeAccessTokenKeyKey存入redis
			cache.RedisClient.Set(cache.StoreAccessTokenKeyKey(user.ID), newTokenKey, wjwt.AccessExpireTime)
			// then
			claims, _ = wjwt.ParseJWT(newAccessToken)
			// 如果更新了accessToken，就返回前端响应头中一个Authorization
			c.Header("Authorization", newTokenKey)
		} else {
			// accessToken没有过期
			claims, err = wjwt.ParseJWT(accessToken)
			if err != nil || nowTime.Unix() > claims.ExpiresAt { // 解析失败,过期了也会解析失败
				c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.RespCode(e.InvalidWithAuth, c))
				return
			}
		}
		// 解析成功，将claims的值进行传递
		c.Set("claims", claims)
		c.Next()
	}
}
