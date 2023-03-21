package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	"go_awd/pkg/util"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
)

// Login
// @Description: 用户登录
// @receiver u *LoginService
// @param c *gin.Context
// @return serializer.Response
func (u *LoginService) Login(c *gin.Context) serializer.Response {
	// 先验证验证码
	if !ImgCaptchaContainer.Store.Verify(u.CaptchaKey, u.CaptchaValue, true) {
		service.Debugln("captcha verify failed")
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	// 再去对比密码
	userDao := dao.NewUserDao(c)
	user, exist := userDao.ExistOrNotByUserNameOrEmail(u.Username, u.Username)
	if !exist || !user.CheckPassword(u.Password) { // 不存在对应的用户或密码错误
		service.Debugln("user not found or pwd error")
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	// 如果登录用户当前存在tokenKey在redis中，先删除改tokenKey
	if storedTokenKey, err := cache.RedisClient.Get(cache.StoreAccessTokenKeyKey(user.ID)).Result(); err == nil { // 如果当前用户存在tokenKey，那么删除
		cache.RedisClient.Del(cache.AccessTokenKey(storedTokenKey))
	}
	// 生成double token
	accessToken, refreshToken, err := wjwt.GenAccessAndRefreshJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return serializer.RespCode(e.InvalidWithGenJwt, c)
	}
	// accessToken存入redis
	tokenKey := util.Sha1(accessToken)
	if err := cache.RedisClient.Set(cache.AccessTokenKey(tokenKey), accessToken, wjwt.AccessExpireTime).Err(); err != nil {
		service.Errorln("accessToken to redis error,", err.Error())
		return serializer.RespCode(e.InvalidWithGenJwt, c)
	}
	// storeAccessTokenKeyKey存入redis
	cache.RedisClient.Set(cache.StoreAccessTokenKeyKey(user.ID), tokenKey, wjwt.AccessExpireTime)
	// refreshToken存入MySQL
	if err := dao.NewTokenDaoByDB(userDao.DB).CreateRefreshToken(
		&model.RefreshToken{
			UserID:       user.ID,
			TokenKey:     tokenKey,
			RefreshToken: refreshToken,
		}); err != nil {
		service.Errorln("refreshToken to mysql error,", err.Error())
		return serializer.RespCode(e.InvalidWithGenJwt, c)
	}
	// 返回给用户tokenKey
	return serializer.RespSuccess(e.SuccessWithLogin, gin.H{"token": tokenKey}, c)
}
