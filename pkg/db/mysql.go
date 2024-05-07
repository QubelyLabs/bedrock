package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var sql *gorm.DB

func InitSQL(host string, port int, user, password, name string) error {
	m, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
			user, password, host, port, name,
		), // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  false, // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    false, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   false, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,       // Don't include params in the SQL log
				Colorful:                  true,        // Disable color
			},
		),
	})

	if err != nil {
		return err
	}

	sql = m

	return nil
}

func SQL() *gorm.DB {
	return sql
}
