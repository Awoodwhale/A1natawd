package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/pkg/e"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
)

// Logout
// @Description: 退出登录
// @receiver u *EmptyService
// @param c *gin.Context
// @return serializer.Response
func (u *EmptyService) Logout(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	// 删除token认证
	DeleteAccessToken(claims.ID)
	if err := DeleteRefreshToken(claims.ID, dao.NewTokenDao(c)); err != nil {
		service.Debugln("delete refresh_token from mysql error,", err.Error())
	}
	return serializer.RespSuccess(e.SuccessWithLogout, nil, c)
}

// DeleteRefreshToken
// @Description: 从MySQL删除refresh token
// @param userID int64
// @param tokenDao *dao.TokenDao
// @return error
func DeleteRefreshToken(userID int64, tokenDao *dao.TokenDao) error {
	return tokenDao.DeleteRefreshTokenByUserID(userID)
}

// DeleteAccessToken
// @Description: 从redis删除access token
// @param userID int64
func DeleteAccessToken(userID int64) {
	if tokenKey, err := cache.RedisClient.Get(cache.StoreAccessTokenKeyKey(userID)).Result(); err == nil {
		cache.RedisClient.Del(cache.AccessTokenKey(tokenKey), cache.StoreAccessTokenKeyKey(userID))
	}
}
