package mysql

import (
	"Comment/module/db"
	"gorm.io/gorm"
)

var (
	DBSession *gorm.DB
	err       error
)

func Init() error {
	DBSession, err = db.NewMysqlSession()
	if err != nil {
		return err
	}
	return nil
}
