package database

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	isInit bool
	Gorm   *gorm.DB
	err    error
)

type mysqlInfo struct {
	host     string
	user     string
	password string
	port     string
	dbName   string
}

func NewMysqlConInfo(host, user, password, port, dbName string) *mysqlInfo {
	return &mysqlInfo{
		host:     host,
		user:     user,
		password: password,
		port:     port,
		dbName:   dbName,
	}
}

func init() {
	if isInit {
		return
	}
	mysqlConInfo := NewMysqlConInfo(
		"127.0.0.1",
		"root",
		"chenteng",
		"3306",
		"demo",
	)
	g, err := CreateDB(mysqlConInfo)
	if err != nil {
		log.Fatal("初始化数据库失败", err)
	}
	Gorm = g
}

func CreateDB(info *mysqlInfo) (*gorm.DB, error) {
	var g *gorm.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local",
		info.user,
		info.password,
		info.host,
		info.port,
		info.dbName,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	g, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	//连接池最大允许的空闲连接数
	sqlDB, err := g.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(100)
	//设置最大连接数
	sqlDB.SetMaxIdleConns(20)
	//设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(60 * time.Second)
	isInit = true
	return g, nil
}
