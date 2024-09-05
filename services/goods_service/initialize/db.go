package initialize

import (
	"fmt"

	"github.com/zhaohaihang/goods_service/global"
	"github.com/zhaohaihang/goods_service/model"
	"go.uber.org/zap"

	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化DB
func InitDB() {
	var err error
	MysqlInfo := global.ServiceConfig.MySqlInfo
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
	// 可以在迁移时，设置表的字符集
	// err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").
	// 	AutoMigrate(&model.Banner{},&model.Category{},&model.Brand{},&model.Goods{},&model.GoodsCategoryBrand{})

	err = global.DB.AutoMigrate( &model.Banner{},&model.Category{},&model.Brand{},&model.Goods{},&model.GoodsCategoryBrand{})
	if err != nil {
		zap.S().Fatalw("db  AutoMigrate:", "err", err.Error())
	}
	zap.S().Infow("init goods db conn success")
}
