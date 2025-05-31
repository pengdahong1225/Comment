package db

import (
	"Comment/module/logger"
	"Comment/module/settings"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlSession() (*gorm.DB, error) {
	newLogger, err := logger.NewOrmLogger()
	if err != nil {
		logrus.Errorf("new orm logger error: %v", err)
	}

	cfg := settings.Instance().MysqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User,
		cfg.Pwd, cfg.Host, cfg.Port, cfg.Db)

	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		// SkipDefaultTransaction: true, //全局禁用默认事务
	})
	if err != nil {
		return nil, err
	}

	return db_, nil
}
