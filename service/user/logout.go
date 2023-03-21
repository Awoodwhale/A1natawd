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

func (u *EmptyService) Logout(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	// 删除token认证
	DeleteAccessToken(claims.ID)
	if err := DeleteRefreshToken(claims.ID, dao.NewTokenDao(c)); err != nil {
		service.Debugln("delete refresh_token from mysql error,", err.Error())
	}
	return serializer.RespSuccess(e.SuccessWithLogout, nil, c)
}

func DeleteRefreshToken(userID int64, tokenDao *dao.TokenDao) error {
	return tokenDao.DeleteRefreshTokenByUserID(userID)
}

func DeleteAccessToken(userID int64) {
	if tokenKey, err := cache.RedisClient.Get(cache.StoreAccessTokenKeyKey(userID)).Result(); err == nil {
		cache.RedisClient.Del(cache.AccessTokenKey(tokenKey), cache.StoreAccessTokenKeyKey(userID))
	}
}
