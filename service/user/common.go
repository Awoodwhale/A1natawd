package user

// RegisterAndUpdateService
// @Description: 注册与更新的service
type RegisterAndUpdateService struct {
	Username          string `json:"username" form:"username" binding:"required,lte=50" msg:"invalid_username"`
	Password          string `json:"password" form:"password" binding:"required,lte=50" msg:"invalid_password"`
	Email             string `json:"email" form:"email" binding:"required,email" msg:"invalid_email"`
	EmailCaptchaValue string `json:"captcha" form:"captcha" binding:"required,len=5" msg:"invalid_captcha"` // 邮箱验证码 	// 图形验证码的value
}

// LoginService
// @Description: 登录的service
type LoginService struct {
	Username     string `json:"username" form:"username" binding:"required,lte=50" msg:"invalid_username"`
	Password     string `json:"password" form:"password" binding:"required,lte=50" msg:"invalid_password"`
	CaptchaKey   string `json:"key" form:"key" binding:"required" msg:"invalid_captcha"`               // 验证码key
	CaptchaValue string `json:"captcha" form:"captcha" binding:"required,len=5" msg:"invalid_captcha"` // 图形验证码的value
}

// EmailCaptchaValidateService
// @Description: 邮箱验证的service
type EmailCaptchaValidateService struct {
	Email string `json:"email" form:"email" binding:"required,email" msg:"invalid_email"`
	// type == register 注册时绑定邮箱
	// type == update 已注册用户修改邮箱
	// type == password 已注册用户修改密码
	Type string `json:"type" form:"type" binding:"required,oneof=register update password recover" msg:"invalid_type"`
}

// CaptchaValidateService
// @Description: 验证码service
type CaptchaValidateService struct {
	CaptchaKey   string `json:"captcha_key" form:"captcha_key" binding:"omitempty" msg:"invalid_captcha_key"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value" binding:"omitempty,len=5" msg:"invalid_captcha"`
}

// UpdateService
// @Description: 用户更新的service
type UpdateService struct {
	Username string `json:"username" form:"username" binding:"omitempty,lte=50" msg:"invalid_username"`
	Sign     string `json:"sign" form:"sign" binding:"omitempty,lte=100" msg:"invalid_sign"`
}

// RecoverPwdService
// @Description: 找回密码service
type RecoverPwdService struct {
	Username          string `json:"username" form:"username" binding:"required,lte=50" msg:"invalid_username"`
	Email             string `json:"email" form:"email" binding:"required,email" msg:"invalid_email"`
	EmailCaptchaValue string `json:"captcha" form:"captcha" binding:"required,len=5" msg:"invalid_captcha"`
}

// EmptyService
// @Description: 空service，用来获取用户信息
type EmptyService struct{}
