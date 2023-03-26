package conf

import (
	ginI18n "github.com/fishjar/gin-i18n"
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/pkg/wlog"
	"gopkg.in/ini.v1"
	"time"
)

var (
	file *ini.File
	// service
	AppMode  string
	HttpPort string
	// mysql
	DBHost string
	DBPort string
	DBUser string
	DBPwd  string
	DBName string
	// redis
	RedisHost string
	RedisPort string
	RedisPwd  string
	RedisName string
	// email
	SmtpHost     string
	SmtpPort     int
	SmtpEmail    string
	SmtpToken    string
	SmtpSendName string
	// img
	ImgHost    string
	ImgPort    string
	ImgMaxSize int64
	ImgPath    string
	ImageTypes []string
	// i18n
	DefaultLang  string
	SupportLang  string
	LangFilePath string
	// log
	LogPath      string
	LogLevel     int // 5 debug 2 error
	LogKeepCount int // log保存时间/天
	// page
	PageSize uint
	// flag
	DockerServerIP     string
	FlagEnv            string
	SSHUsernameEnv     string
	SSHPasswordEnv     string
	SSHDefaultUsername string
	SSHDefaultPassword string
	FlagPrefix         string
	ContainerExistTime time.Duration
	MaxTeamCount       int
	SSHStartPort       int
	WebBoxStartPort    int
	PwnBoxStartPort    int
)

// Init
// @Description: 从config.ini读取配置
func Init(confPath string) {
	if confPath == "" {
		confPath = "./conf/config.ini"
	}
	f, err := ini.Load(confPath)
	if err != nil {
		panic(err)
	}
	file = f
	// load config file
	loadServer()
	loadMySQL()
	loadRedis()
	loadEmail()
	loadImage()
	loadI18n()
	loadLog()
	loadPage()
	loadFlag()
	// logger init
	wlog.InitLogger(LogPath, LogLevel, LogKeepCount)
	// mysql read（主）
	pathRead := DBUser + ":" + DBPwd + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8mb4&parseTime=true"
	// mysql write（从）
	pathWrite := DBUser + ":" + DBPwd + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8mb4&parseTime=true"
	// mysql init
	dao.InitDatabase(pathRead, pathWrite)
	// redis init
	cache.InitDatabase(RedisHost+":"+RedisPort, RedisName, RedisPwd)
	// set gin mode
	gin.SetMode(AppMode)
}

// loadServer
// @Description: 获取server的config
// @param
func loadServer() {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = ":" + file.Section("service").Key("HttpPort").String() // 加上:前缀
}

// loadMySQL
// @Description: 获取MySQL的config
// @param
func loadMySQL() {
	DBHost = file.Section("mysql").Key("DBHost").String()
	DBPort = file.Section("mysql").Key("DBPort").String()
	DBUser = file.Section("mysql").Key("DBUser").String()
	DBPwd = file.Section("mysql").Key("DBPwd").String()
	DBName = file.Section("mysql").Key("DBName").String()
}

// loadRedis
// @Description: 获取redis的config
// @param
func loadRedis() {
	RedisHost = file.Section("redis").Key("RedisHost").String()
	RedisPort = file.Section("redis").Key("RedisPort").String()
	RedisPwd = file.Section("redis").Key("RedisPwd").String()
	RedisName = file.Section("redis").Key("RedisName").String()
}

// loadEmail
// @Description: 获取email的config
// @param
func loadEmail() {
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpPort, _ = file.Section("email").Key("SmtpPort").Int()
	SmtpEmail = file.Section("email").Key("SmtpEmail").String()
	SmtpToken = file.Section("email").Key("SmtpToken").String()
	SmtpSendName = file.Section("email").Key("SmtpSendName").String()
}

// loadImage
// @Description: 获取image的config
// @param
func loadImage() {
	ImgHost = file.Section("image").Key("ImgHost").String()
	ImgPort = file.Section("image").Key("ImgPort").String()
	ImgMaxSize, _ = file.Section("image").Key("ImgMaxSize").Int64()
	ImgPath = file.Section("image").Key("ImgPath").String()
	ImageTypes = file.Section("image").Key("ImageTypes").Strings(",")
}

// loadI18n
// @Description: 获取国际化配置
// @param
func loadI18n() {
	DefaultLang = file.Section("i18n").Key("DefaultLang").String()
	SupportLang = file.Section("i18n").Key("SupportLang").String()
	LangFilePath = file.Section("i18n").Key("LangFilePath").String()
	ginI18n.LocalizerInit(DefaultLang, SupportLang, LangFilePath) // 初始化国际化配置
}

// loadLog
// @Description: 获取log配置
// @param
func loadLog() {
	LogPath = file.Section("log").Key("LogPath").String()
	LogLevel, _ = file.Section("log").Key("LogLevel").Int()
	LogKeepCount, _ = file.Section("log").Key("LogKeepCount").Int()
}

// loadPage
// @Description: 配置默认分页显示数量
// @param
func loadPage() {
	PageSize, _ = file.Section("page").Key("PageSize").Uint()
}

// loadFlag
// @Description: 配置flag
// @param
func loadFlag() {
	DockerServerIP = file.Section("flag").Key("DockerServerIP").String()
	FlagEnv = file.Section("flag").Key("FlagEnv").String()
	SSHUsernameEnv = file.Section("flag").Key("SSHUsernameEnv").String()
	SSHPasswordEnv = file.Section("flag").Key("SSHPasswordEnv").String()
	SSHDefaultUsername = file.Section("flag").Key("SSHDefaultUsername").String()
	SSHDefaultPassword = file.Section("flag").Key("SSHDefaultPassword").String()
	FlagPrefix = file.Section("flag").Key("FlagPrefix").String()
	sec, _ := file.Section("flag").Key("ContainerExistTime").Int64()
	ContainerExistTime = time.Second * time.Duration(sec)
	MaxTeamCount, _ = file.Section("flag").Key("MaxTeamCount").Int()
	SSHStartPort, _ = file.Section("flag").Key("SSHStartPort").Int()
	WebBoxStartPort, _ = file.Section("flag").Key("WebBoxStartPort").Int()
	PwnBoxStartPort, _ = file.Section("flag").Key("PwnBoxStartPort").Int()
}
