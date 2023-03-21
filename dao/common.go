package dao

import (
	"context"
	"go_awd/pkg/wlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"time"
)

var _db *gorm.DB

// InitDatabase
// @Description: 初始化MySQL database
// @param pathRead string
// @param pathWrite string
func InitDatabase(pathRead, pathWrite string) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       pathRead,
		DefaultStringSize:         256,  // string类型字段默认长度
		DisableDatetimePrecision:  true, // 禁止datetime精度
		DontSupportRenameIndex:    true, // 如果需要重命名索引，需要把索引删除后再重建
		DontSupportRenameColumn:   true, // 用change重命名列，mysql8之前的数据库不支持
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: wlog.GormLogger, // 设置gorm的log
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		wlog.Logger.Errorln("mysql dao init error, ", err)
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)                 // 空闲连接池
	sqlDB.SetMaxOpenConns(100)                 // 打开
	sqlDB.SetConnMaxLifetime(30 * time.Second) // 空闲连接最多存活时间

	// 主从配置
	_db = db
	_ = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(pathWrite)},                      // 写操作
		Replicas: []gorm.Dialector{mysql.Open(pathRead), mysql.Open(pathRead)}, // 读操作
		Policy:   dbresolver.RandomPolicy{},
	}))

	// 数据迁移
	migration()
}

// NewDBClient
// @Description: 获取db对象
// @param ctx context.Context
// @return *gorm.DB
func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
