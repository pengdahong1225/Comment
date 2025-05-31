package logger

import (
	"fmt"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func NewOrmLogger() (logger.Interface, error) {
	path := fmt.Sprintf("%s/orm.log", "./log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	writer := io.MultiWriter(os.Stdout, file)
	newLogger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	return newLogger, nil
}
