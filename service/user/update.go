package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/conf"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
	"strconv"
)

// UpdateEmail
// @Description: 修改用户绑定的邮箱
// @receiver u *RegisterAndUpdateService
// @param c *gin.Context
// @return interface{}
func (u *RegisterAndUpdateService) UpdateEmail(c *gin.Context) serializer.Response {
	// 查询redis判断邮箱验证码是否相同
	emailKey := cache.EmailCaptchaKey("update", u.Email)
	if captchaInRedis, err := cache.RedisClient.Get(emailKey).Result(); err != nil {
		service.Debugln("user UpdateEmail redis err,", err.Error())
		return serializer.RespCode(e.InvalidWithCaptcha, c) // redis错误
	} else {
		if captchaInRedis != u.EmailCaptchaValue {
			service.Infoln("user UpdateEmail captcha not equal", u.Email)
			return serializer.RespCode(e.InvalidWithCaptcha, c) // 两次的验证码不相等
		}
	}
	// dao层更新用户的email
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByUsername(u.Username)
	if err != nil {
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}
	// 密码错误，就直接返回
	if !user.CheckPassword(u.Password) {
		return serializer.RespCode(e.InvalidWithPassword, c)
	}
	// 如果邮箱未变，直接返回
	if user.Email == u.Email {
		return serializer.RespCode(e.InvalidWithSameEmail, c)
	}
	// 如果数据库中存在邮箱相同的用户，直接返回
	if _, te := userDao.GetUserByEmail(u.Email); te == nil {
		return serializer.RespCode(e.InvalidWithSameEmail, c)
	}
	// 尝试修改邮箱
	if err := userDao.UpdateEmailByUsername(u.Username, u.Email); err != nil {
		service.Errorln("user UpdateEmail mysql,", err.Error())
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}
	user.Email = u.Email // 直接修改，就不去数据库查询咯
	// 从redis删除邮箱验证码
	if err := cache.RedisClient.Del(emailKey).Err(); err != nil {
		service.Errorln("user UpdateEmail redis del,", err.Error())
	}

	// 修改邮箱成功
	return serializer.RespSuccess(e.SuccessWithUpdateEmail, serializer.BuildUser(user), c)
}

// UpdatePassword
// @Description: 修改用户密码
// @receiver u *RegisterAndUpdateService
// @param c *gin.Context
// @return serializer.Response
func (u *RegisterAndUpdateService) UpdatePassword(c *gin.Context) serializer.Response {
	// 查询redis判断邮箱验证码是否相同
	emailKey := cache.EmailCaptchaKey("password", u.Email)
	if captchaInRedis, err := cache.RedisClient.Get(emailKey).Result(); err != nil {
		service.Debugln("user UpdatePassword redis err,", err.Error())
		return serializer.RespCode(e.InvalidWithCaptcha, c) // redis错误
	} else {
		if captchaInRedis != u.EmailCaptchaValue {
			service.Infoln("user UpdatePassword captcha not equal", u.Email)
			return serializer.RespCode(e.InvalidWithCaptcha, c) // 两次的验证码不相等
		}
	}
	// dao层修改用户密码
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByUsername(u.Username)
	if err != nil || user.Email != u.Email { // 查无此人或者邮箱不匹配就返回
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}
	hashedPwd, err := model.SetPassword(u.Password)
	if err != nil {
		service.Errorln("user UpdatePassword pwd hash error,", err.Error())
		return serializer.RespCode(e.ErrorWithEncryption, c)
	}
	if err := userDao.UpdatePwdByUsername(u.Username, hashedPwd); err != nil {
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}
	// 从redis删除邮箱验证码
	if err := cache.RedisClient.Del(emailKey).Err(); err != nil {
		service.Errorln("user UpdateEmail redis del,", err.Error())
	}
	// 修改密码后需要重新登录，删除所有token
	DeleteAccessToken(user.ID)
	_ = DeleteRefreshToken(user.ID, dao.NewTokenDaoByDB(userDao.DB))

	return serializer.RespSuccess(e.SuccessWithUpdatePwd, nil, c)
}

func (u *UpdateService) Update(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	file, err := c.FormFile("file")
	flag := false
	if err == nil {
		if file.Size > conf.ImgMaxSize { // 图片过大
			return serializer.RespCode(e.InvalidWithImgSize, c)
		}
		// 上传图片到本地
		imgPath := conf.ImgPath + "user/" + strconv.FormatInt(user.ID, 10) + "/" + file.Filename
		if err := c.SaveUploadedFile(file, imgPath); err != nil {
			return serializer.RespCode(e.InvalidWithUpdateUser, c)
		}
		user.Avatar = "user/" + strconv.FormatInt(user.ID, 10) + "/" + file.Filename
		flag = true
	}
	if u.Sign != "" {
		user.Sign = u.Sign
		flag = true
	}
	if u.Username != "" {
		user.Username = u.Username
		flag = true
	}
	if flag {
		if err := userDao.UpdateByID(user.ID, user); err != nil {
			service.Errorln("update user error,", err.Error())
			return serializer.RespCode(e.InvalidWithUpdateUser, c)
		}
	}
	return serializer.RespSuccess(e.SuccessWithUpdate, serializer.BuildUser(user), c)
}
