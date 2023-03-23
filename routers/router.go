package routers

import (
	ginI18n "github.com/fishjar/gin-i18n"
	"github.com/gin-gonic/gin"
	api "go_awd/api/v1"
	"go_awd/conf"
	"go_awd/middleware/auth"
	"go_awd/middleware/cors"
	lm "go_awd/middleware/limit"
	"go_awd/pkg/e"
	"go_awd/pkg/wlog"
	"go_awd/serializer"
	"net/http"
	"runtime"
	"time"
)

// NewRouter
// @Description: 设置gin router
// @return *gin.Engine
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(wlog.GinLogger, gin.Recovery())       // gin的log中间件
	router.Use(ginI18n.GinLocalizer())               // i18n国际化中间件
	router.Use(cors.CorsMiddleware())                // cors中间件
	router.StaticFS("/static", http.Dir("./static")) // 设置fs路径

	v1 := router.Group("api/v1")                          // v1版本的api
	v1.Use(lm.LimitMiddleware(200, 5*time.Minute, "all")) // 5分钟最多请求200次
	{
		if conf.AppMode == "debug" {
			// ping test
			v1.GET("ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, serializer.RespSuccess(e.Success, runtime.GOOS, c))
			})
		}
	}
	// 无需权限校验的接口
	unAuthed := v1.Group("/")
	unAuthedUser := unAuthed.Group("user")
	unAuthedUser.Use(auth.FileTypeMiddleware()) // 过滤文件后缀
	{
		// 用户注册
		unAuthedUser.POST("register",
			lm.LimitMiddleware(30, time.Minute, "register"),
			api.UserRegister)
		// 用户登录
		unAuthedUser.POST("login",
			lm.LimitMiddleware(30, time.Minute, "login"),
			api.UserLogin)
		// 图像验证码,限制一分钟只能30次
		unAuthedUser.GET("captcha",
			lm.LimitMiddleware(30, time.Minute, "captcha"),
			api.GenCaptcha)
		// 发送邮箱验证码
		unAuthedUser.POST("email",
			lm.LimitMiddleware(30, time.Minute, "email"),
			api.SendEmailCaptcha)
		// 重制密码
		unAuthedUser.POST("recover",
			lm.LimitMiddleware(30, time.Minute, "recover"),
			api.RecoverUserPwd)
		// 查询用户信息
		unAuthedUser.GET(":id", api.ShowUserInfoByID)
	}

	// 以下api需要权限校验
	authed := v1.Group("/")
	authed.Use(auth.JWT())
	// 需要权限的user接口
	authedUser := authed.Group("user")
	{
		authedUser.GET("logout", api.UserLogout) // 退出登录
		authedUser.PUT("email",                  // 修改用户邮箱
			lm.LimitMiddleware(30, time.Minute, "update_email"),
			api.UpdateUserEmail)
		authedUser.PUT("password", // 修改用户密码
			lm.LimitMiddleware(30, time.Minute, "update_password"),
			api.UpdateUserPwd)
		authedUser.PUT("update", // 更新用户信息
			lm.LimitMiddleware(30, time.Minute, "update_user"),
			api.UpdateUserInfo)
		authedUser.GET("self", api.ShowUserInfo) // 展示当前用户信息
	}
	// 需要权限的team接口
	authedTeam := authed.Group("team")
	authedTeam.Use(auth.FileTypeMiddleware()) // 过滤文件后缀
	{
		authedTeam.POST("create", api.CreateTeam)                            // 创建团队
		authedTeam.POST("apply", api.ApplyTeam)                              // 申请入队
		authedTeam.DELETE("cancel", api.CancelApplyTeam)                     // 取消申请入队
		authedTeam.DELETE("leave", api.LeaveTeam)                            // 离开团队
		authedTeam.DELETE("dismiss", auth.Role("leader"), api.DismissTeam)   // 解散团队
		authedTeam.DELETE("transfer", auth.Role("leader"), api.TransferTeam) // 转让团队队长
		authedTeam.POST("accept", auth.Role("leader"), api.AcceptTeam)       // 同意入队申请
		authedTeam.POST("reject", auth.Role("leader"), api.RejectTeam)       // 拒绝入队申请
		authedTeam.PUT("update",                                             // 更新团队信息
			auth.Role("leader"),
			lm.LimitMiddleware(30, time.Minute, "update_user"),
			api.UpdateTeam)
	}

	adminAuthed := authed.Group("admin")
	adminAuthed.Use(auth.Role("admin"))
	{
		adminAuthed.GET("challenge", api.ShowChallenge)            // 获取题目列表
		adminAuthed.POST("challenge", api.CreateOrUpdateChallenge) // 上传题目
		adminAuthed.PUT("challenge/:id", api.UpdateChallengeInfo)  // 修改题目
		adminAuthed.POST("challenge/:id",                          // 开启题目测试
			lm.LimitMiddleware(1, 5*time.Second, "start_test_challenge"),
			api.StartTestChallenge)
		adminAuthed.DELETE("challenge/:id") // 删除题目环境
	}
	return router
}
