package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	"go_awd/serializer"
	"go_awd/service"
)

// Register
// @Description: 用户注册
// @receiver u *RegisterAndUpdateService
// @param c *gin.Context
// @return serializer.Response
func (u *RegisterAndUpdateService) Register(c *gin.Context) serializer.Response {
	// 查询redis判断邮箱验证码是否相同
	emailKey := cache.EmailCaptchaKey("register", u.Email)
	if captchaInRedis, err := cache.RedisClient.Get(emailKey).Result(); err != nil {
		service.Debugln("user Register redis err,", err.Error())
		return serializer.RespCode(e.InvalidWithCaptcha, c) // redis错误
	} else {
		if captchaInRedis != u.EmailCaptchaValue {
			service.Infoln("user register email captcha not equal,", u.Email)
			return serializer.RespCode(e.InvalidWithCaptcha, c) // 两次的验证码不相等
		}
	}

	// 查看用户是否存在
	userDao := dao.NewUserDao(c)
	if _, exist := userDao.ExistOrNotByUserNameOrEmail(u.Username, u.Email); exist {
		service.Debugln("user Register exist username or email")
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}

	// 不存在就添加用户
	user := &model.User{
		Username: u.Username,
		Email:    u.Email,
		Role:     model.UserRole, // 普通用户
		Avatar:   "avatar.png",
		Money:    0, // 初始金币为0
		Score:    0, // 初始积分为0
		Sign:     "这个人很懒，什么都没留下",
	}
	if err := user.SetPassword(u.Password); err != nil {
		service.Errorln("user Register encryption,", err.Error())
		return serializer.RespCode(e.ErrorWithEncryption, c)
	}

	// 创建用户，写入mysql
	if err := userDao.CreateUser(user); err != nil {
		service.Errorln("user Register mysql,", err.Error())
		return serializer.RespCode(e.InvalidWithCreateUser, c)
	}

	// 从redis删除邮箱验证码
	if err := cache.RedisClient.Del(emailKey).Err(); err != nil {
		service.Errorln("user Register redis del,", err.Error())
	}

	// 注册成功
	return serializer.RespSuccess(e.SuccessWithRegister, serializer.BuildUser(user), c)
}
