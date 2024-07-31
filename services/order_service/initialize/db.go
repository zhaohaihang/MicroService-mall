package initialize

import (
	"fmt"
	"log"
	"time"

	"github.com/zhaohaihang/order_service/global"
	"github.com/zhaohaihang/order_service/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	var err error
	MySQL := global.ServiceConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySQL.User, 
		MySQL.Password,
		MySQL.Host,
		MySQL.Port,
		MySQL.Name)
	// 创建日志文件
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
	err = global.DB.AutoMigrate(&model.OrderGoods{},&model.OrderInfo{},&model.ShoppingCart{})
	if err != nil {
		zap.S().Fatalw("db AutoMigrate : %s", "err", err.Error())
	}
	zap.S().Infow("db conn success")
}
