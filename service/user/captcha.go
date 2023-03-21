package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go_awd/cache"
	"go_awd/pkg/e"
	"go_awd/pkg/util"
	"go_awd/serializer"
	"go_awd/service"
	"time"
)

var ImgCaptchaContainer = base64Captcha.NewCaptcha(
	base64Captcha.NewDriverDigit(80, 160, 5, 0.7, 80),
	cache.NewCaptchaMemStore(3*time.Minute),
)

// GenEmailCaptcha
// @Description: 生成邮箱验证码并发送
// @receiver em *EmailCaptchaValidateService
// @param c *gin.Context
// @return serializer.Response
func (em *EmailCaptchaValidateService) GenEmailCaptcha(c *gin.Context, limit int, tm time.Duration) serializer.Response {
	go func() {
		captchaKey := util.RandStringBytes(5)
		key := cache.EmailSendCountKey(em.Type, em.Email)
		// 先读已经发了多少次
		if curCount, err := cache.RedisClient.Get(key).Int(); err == nil && curCount >= limit {
			service.Infoln("email request too many from", em.Email)
			return // 访问频率太高
		}

		if em.Type == "register" { // 未注册用户注册时绑定邮箱
			if err := service.SendUserRegisterEmail(em.Email, captchaKey); err != nil {
				service.Errorln("email captcha send error,", err.Error())
				return
			}
		} else if em.Type == "update" { // 已注册用户修改邮箱
			if err := service.SendUserUpdateEmail(em.Email, captchaKey); err != nil {
				service.Errorln("email captcha send error,", err.Error())
				return
			}
		} else if em.Type == "password" { // 已注册用户修改密码
			if err := service.SendUserUpdatePwd(em.Email, captchaKey); err != nil {
				service.Errorln("email captcha send error,", err.Error())
				return
			}
		} else if em.Type == "recover" {
			if err := service.SendUserRecoverPwd(em.Email, captchaKey); err != nil {
				service.Errorln("email captcha send error,", err.Error())
				return
			}
		}

		// 发送成功才存入redis
		if err := cache.RedisClient.Set(cache.EmailCaptchaKey(em.Type, em.Email), captchaKey, 5*time.Minute).Err(); err != nil {
			service.Errorln("email captcha store to redis error,", err.Error())
			return
		}
		// 更新发送的次数
		count, err := cache.RedisClient.Incr(key).Result()
		if err != nil {
			service.Errorln("email captcha count store to redis error,", err.Error())
			return
		}
		if count == 1 { // 第一次发送，设置过期时间
			if err := cache.RedisClient.Expire(key, tm).Err(); err != nil {
				service.Errorln("email captcha count set expire to redis error,", err.Error())
				return
			}
		}
	}()

	// 无论怎样都返回成功状态
	return serializer.RespSuccess(e.SuccessWithGenCaptcha, nil, c)
}

// GenCaptcha
// @Description: 生成验证码
// @receiver service *CaptchaValidateService
// @param c *gin.Context
// @return serializer.Response
func (service *CaptchaValidateService) GenCaptcha(c *gin.Context) serializer.Response {
	id, b64s, err := ImgCaptchaContainer.Generate() // 3分钟的图形验证码有效期
	if err != nil {
		return serializer.RespErr(e.ErrorWithGenCaptcha, err, c)
	}
	return serializer.RespSuccess(e.SuccessWithGenCaptcha, gin.H{"captcha_key": id, "captcha_b64": b64s}, c)
}
