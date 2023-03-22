package conf

import (
	ginI18n "github.com/fishjar/gin-i18n"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/pkg/wlog"
	"gopkg.in/ini.v1"
	"strings"
)

var (
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
	DockerServerIP  string
	FlagEnv         string
	FlagPrefix      string
	MaxTeamCount    int
	WebBoxStartPort int
	PwnBoxStartPort int
)

// Init
// @Description: 从config.ini读取配置
func Init(confPath string) {
	if confPath == "" {
		confPath = "./conf/config.ini"
	}
	file, err := ini.Load(confPath)
	if err != nil {
		panic(err)
	}
	// load config file
	loadServer(file)
	loadMySQL(file)
	loadRedis(file)
	loadEmail(file)
	loadImage(file)
	loadI18n(file)
	loadLog(file)
	loadPage(file)
	loadFlag(file)
	// logger init
	wlog.InitLogger(LogPath, LogLevel, LogKeepCount)
	// mysql read（主）
	pathRead := strings.Join([]string{DBUser, ":", DBPwd, "@tcp(", DBHost, ":", DBPort, ")/", DBName, "?charset=utf8mb4&parseTime=true"}, "")
	// mysql write（从）
	pathWrite := strings.Join([]string{DBUser, ":", DBPwd, "@tcp(", DBHost, ":", DBPort, ")/", DBName, "?charset=utf8mb4&parseTime=true"}, "")
	// mysql init
	dao.InitDatabase(pathRead, pathWrite)
	// redis init
	cache.InitDatabase(RedisHost+":"+RedisPort, RedisName, RedisPwd)
}

// loadServer
// @Description: 获取server的config
// @param file *ini.File
func loadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = ":" + file.Section("service").Key("HttpPort").String() // 加上:前缀
}

// loadMySQL
// @Description: 获取MySQL的config
// @param file *ini.File
func loadMySQL(file *ini.File) {
	DBHost = file.Section("mysql").Key("DBHost").String()
	DBPort = file.Section("mysql").Key("DBPort").String()
	DBUser = file.Section("mysql").Key("DBUser").String()
	DBPwd = file.Section("mysql").Key("DBPwd").String()
	DBName = file.Section("mysql").Key("DBName").String()
}

// loadRedis
// @Description: 获取redis的config
// @param file *ini.File
func loadRedis(file *ini.File) {
	RedisHost = file.Section("redis").Key("RedisHost").String()
	RedisPort = file.Section("redis").Key("RedisPort").String()
	RedisPwd = file.Section("redis").Key("RedisPwd").String()
	RedisName = file.Section("redis").Key("RedisName").String()
}

// loadEmail
// @Description: 获取email的config
// @param file *ini.File
func loadEmail(file *ini.File) {
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpPort, _ = file.Section("email").Key("SmtpPort").Int()
	SmtpEmail = file.Section("email").Key("SmtpEmail").String()
	SmtpToken = file.Section("email").Key("SmtpToken").String()
	SmtpSendName = file.Section("email").Key("SmtpSendName").String()
}

// loadImage
// @Description: 获取image的config
// @param file *ini.File
func loadImage(file *ini.File) {
	ImgHost = file.Section("image").Key("ImgHost").String()
	ImgPort = file.Section("image").Key("ImgPort").String()
	ImgMaxSize, _ = file.Section("image").Key("ImgMaxSize").Int64()
	ImgPath = file.Section("image").Key("ImgPath").String()
	ImageTypes = file.Section("image").Key("ImageTypes").Strings(",")
}

// loadI18n
// @Description: 获取国际化配置
// @param file *ini.File
func loadI18n(file *ini.File) {
	DefaultLang = file.Section("i18n").Key("DefaultLang").String()
	SupportLang = file.Section("i18n").Key("SupportLang").String()
	LangFilePath = file.Section("i18n").Key("LangFilePath").String()
	ginI18n.LocalizerInit(DefaultLang, SupportLang, LangFilePath) // 初始化国际化配置
}

// loadLog
// @Description: 获取log配置
// @param file *ini.File
func loadLog(file *ini.File) {
	LogPath = file.Section("log").Key("LogPath").String()
	LogLevel, _ = file.Section("log").Key("LogLevel").Int()
	LogKeepCount, _ = file.Section("log").Key("LogKeepCount").Int()
}

// loadPage
// @Description: 配置默认分页显示数量
// @param file *ini.File
func loadPage(file *ini.File) {
	PageSize, _ = file.Section("page").Key("PageSize").Uint()
}

// loadFlag
// @Description: 配置flag
// @param file *ini.File
func loadFlag(file *ini.File) {
	DockerServerIP = file.Section("flag").Key("DockerServerIP").String()
	FlagEnv = file.Section("flag").Key("FlagEnv").String()
	FlagPrefix = file.Section("flag").Key("FlagPrefix").String()
	MaxTeamCount, _ = file.Section("flag").Key("MaxTeamCount").Int()
	WebBoxStartPort, _ = file.Section("flag").Key("WebBoxStartPort").Int()
	PwnBoxStartPort, _ = file.Section("flag").Key("PwnBoxStartPort").Int()
}
