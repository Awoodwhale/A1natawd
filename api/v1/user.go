package v1

import (
	"github.com/gin-gonic/gin"
	"go_awd/pkg/e"
	"go_awd/serializer"
	"go_awd/service/user"
	"net/http"
	"strconv"
	"time"
)

// GenCaptcha
// @Description: 生成验证码
// @param c *gin.Context
func GenCaptcha(c *gin.Context) {
	var service user.CaptchaValidateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GenCaptcha(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UserRegister
// @Description: 用户注册
// @param c *gin.Context
func UserRegister(c *gin.Context) {
	var service user.RegisterAndUpdateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Register(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UserLogin
// @Description: 用户登录
// @param c *gin.Context
func UserLogin(c *gin.Context) {
	var service user.LoginService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Login(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UserLogout
// @Description: 退出登录
// @param c *gin.Context
func UserLogout(c *gin.Context) {
	var service user.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Logout(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// SendEmailCaptcha
// @Description: 发送注册邮箱验证码
// @param c *gin.Context
func SendEmailCaptcha(c *gin.Context) {
	var service user.EmailCaptchaValidateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GenEmailCaptcha(c, 10, 10*time.Minute) // 10分钟之内限制发10封邮件
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UpdateUserEmail
// @Description: 修改用户邮箱
// @param c *gin.Context
func UpdateUserEmail(c *gin.Context) {
	var service user.RegisterAndUpdateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateEmail(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UpdateUserPwd
// @Description: 修改用户密码
// @param c *gin.Context
func UpdateUserPwd(c *gin.Context) {
	var service user.RegisterAndUpdateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdatePassword(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// RecoverUserPwd
// @Description: 找回密码
// @param c *gin.Context
func RecoverUserPwd(c *gin.Context) {
	var service user.RecoverPwdService
	if err := c.ShouldBind(&service); err == nil {
		res := service.RecoverPwd(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UpdateUserInfo
// @Description: 更新username、sign、avatar
// @param c *gin.Context
func UpdateUserInfo(c *gin.Context) {
	var service user.UpdateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Update(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// ShowUserInfoByID
// @Description: 获取id对应的用户
// @param c *gin.Context
func ShowUserInfoByID(c *gin.Context) {
	var service user.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, serializer.RespCode(e.Invalid, c))
			return
		}
		res := service.ShowByID(c, id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// ShowUserInfo
// @Description: 展示用户信息
// @param c *gin.Context
func ShowUserInfo(c *gin.Context) {
	var service user.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		res := service.ShowSelf(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}
