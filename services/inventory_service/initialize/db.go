package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/zhaohaihang/inventory_service/global"
	"github.com/zhaohaihang/inventory_service/model"
	"log"
	"time"
)

func InitDB() {
	var err error
	MySQL := global.ServiceConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", MySQL.User, MySQL.Password, MySQL.Host, MySQL.Port, MySQL.Name)

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

	err = global.DB.AutoMigrate(&model.Inventory{},&model.StockSellDetail{})
	if err != nil {
		zap.S().Errorw("global.DB.AutoMigrate", "err", err.Error())
		panic(err)
	}
	zap.S().Infow("init goods db conn success")
}
