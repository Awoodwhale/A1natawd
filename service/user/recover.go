package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/pkg/e"
	"go_awd/pkg/util"
	"go_awd/serializer"
	"go_awd/service"
)

// RecoverPwd
// @Description: 恢复密码
// @receiver u *RecoverPwdService
// @param c *gin.Context
// @return serializer.Response
func (u *RecoverPwdService) RecoverPwd(c *gin.Context) serializer.Response {
	// 查询redis判断邮箱验证码是否相同
	emailKey := cache.EmailCaptchaKey("recover", u.Email)
	if captchaInRedis, err := cache.RedisClient.Get(emailKey).Result(); err != nil {
		return serializer.RespCode(e.InvalidWithCaptcha, c) // 获取不到值
	} else {
		if captchaInRedis != u.EmailCaptchaValue {
			service.Infoln("user RecoverPwd captcha not equal", u.Email)
			return serializer.RespCode(e.InvalidWithCaptcha, c) // 两次的验证码不相等
		}
	}

	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByUsernameAndEmail(u.Username, u.Email)
	if err != nil { // 邮箱和用户名不匹配
		return serializer.RespCode(e.InvalidWithAuth, c)
	}

	// 重制密码前，先把token全删了
	DeleteAccessToken(user.ID)
	_ = DeleteRefreshToken(user.ID, dao.NewTokenDaoByDB(userDao.DB))

	// 重制密码，随机16位的新密码
	randPwd := util.RandStringBytes(16)

	if err := user.SetPassword(randPwd); err != nil {
		return serializer.RespCode(e.ErrorWithEncryption, c)
	}

	if err := userDao.UpdateByID(user.ID, user); err != nil {
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}

	go func() { // 开启协程去发送邮件
		if err := service.SendUserRecoveredPwd(u.Email, randPwd); err != nil {
			service.Errorln("send recover email error,", err.Error())
		}
	}()

	// 从redis删除邮箱验证码
	if err := cache.RedisClient.Del(emailKey).Err(); err != nil {
		service.Errorln("user RecoverPwd redis del,", err.Error())
	}

	return serializer.RespSuccess(e.SuccessWithRecoverPwd, nil, c)
}

// ResetPwdByID
// @Description: 管理员重制用户密码
// @receiver u *EmptyService
// @param c *gin.Context
// @param id int64
// @return serializer.Response
func (u *EmptyService) ResetPwdByID(c *gin.Context, id int64) serializer.Response {
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWithNotExistUser, c)
	}

	// 重制密码前，先把token全删了
	DeleteAccessToken(user.ID)
	_ = DeleteRefreshToken(user.ID, dao.NewTokenDaoByDB(userDao.DB))

	// 重制密码，随机16位的新密码
	randPwd := util.RandStringBytes(16)

	if err := user.SetPassword(randPwd); err != nil {
		return serializer.RespCode(e.ErrorWithEncryption, c)
	}

	if err := userDao.UpdateByID(user.ID, user); err != nil {
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}
	if u.Send == "true" {
		go func() { // 开启协程去发送邮件
			if err := service.SendUserRecoveredPwd(user.Email, randPwd); err != nil {
				service.Errorln("send recover email error,", err.Error())
			}
		}()
	}
	return serializer.RespSuccess(e.SuccessWithUpdatePwd, gin.H{"password": randPwd}, c)
}
