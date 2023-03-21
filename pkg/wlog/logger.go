package wlog

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginLog "go_awd/pkg/wlog/gin"
	gormLog "go_awd/pkg/wlog/gorm"
	uLog "go_awd/pkg/wlog/user"
)

var (
	GinLogger  gin.HandlerFunc
	GormLogger *gormLog.LoggerGorm
	Logger     *logrus.Logger // 分发给coder的log
)

// InitLogger
// @Description: 初始化所有的logger
// @param logPath string
// @param logLevel int
// @param logKeepCount int
func InitLogger(logPath string, logLevel, logKeepCount int) {
	opt := uLog.GetOpt()
	opt.LogPath = logPath
	opt.LogLevel = logrus.Level(logLevel)
	opt.KeepCount = logKeepCount
	opt.FileNamePrefix = "cmd.log"
	opt.ErrLogPrefix = "cmd.err.log"
	opt.LogPrefix = "CMD|"
	lg, err := uLog.New(opt)
	if err != nil {
		panic(err)
	}
	Logger = lg
	initGinLogger(logPath, logLevel, logKeepCount)
	initGormLogger(logPath, logLevel, logKeepCount)
}

// initGinLogger
// @Description: 初始化gin的logger
// @param logPath string
// @param logLevel int
// @param logKeepCount int
func initGinLogger(logPath string, logLevel, logKeepCount int) {
	opt := ginLog.GetOpt()
	opt.OptLogrus.LogPath = logPath
	opt.OptLogrus.LogLevel = logrus.Level(logLevel)
	opt.OptLogrus.KeepCount = logKeepCount
	opt.OptLogrus.FileNamePrefix = "gin.log"
	opt.OptLogrus.ErrLogPrefix = "gin.err.log"
	opt.OptLogrus.LogPrefix = "GIN|"
	// 不记录如下api的log
	opt.SkipRoute = map[string]struct{}{
		"/api/v1/ping": {},
	}
	_, gLog, err := ginLog.New(opt)
	if err != nil {
		panic(err)
	}
	GinLogger = gLog
}

// initGormLogger
// @Description: 初始化gorm的logger
// @param logPath string
// @param logLevel int
// @param logKeepCount int
func initGormLogger(logPath string, logLevel, logKeepCount int) {
	opt := gormLog.GetOpt()
	opt.OptLogrus.LogPath = logPath
	opt.OptLogrus.LogLevel = logrus.Level(logLevel)
	opt.OptLogrus.KeepCount = logKeepCount
	opt.OptLogrus.FileNamePrefix = "gorm.log"
	opt.OptLogrus.ErrLogPrefix = "gorm.err.log"
	opt.OptLogrus.LogPrefix = "GORM|"
	gLog, err := gormLog.New(opt)
	if err != nil {
		panic(err)
	}
	GormLogger = gLog
}
