package initialize

import (
	"fmt"
	"log"
	"time"

	"github.com/zhaohaihang/user_service/global"
	"github.com/zhaohaihang/user_service/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	MysqlInfo := global.ServiceConfig.MysqlInfo
	user := MysqlInfo.User
	password := MysqlInfo.Password
	name := MysqlInfo.Name
	host := MysqlInfo.Host
	port := MysqlInfo.Port
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	newLogger := logger.New(log.New(logFileWriter, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
		Colorful:      false,
	})
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		zap.S().Fatalw("gorm open dsn failed: %s", "err", err.Error())
	}
	err = global.DB.AutoMigrate(&model.User{})
	if err != nil {
		zap.S().Fatalw("db  AutoMigrate : %s", "err", err.Error())
	}
	zap.S().Infow("db conn success")
}
