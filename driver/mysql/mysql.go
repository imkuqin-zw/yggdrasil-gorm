package mysql

import (
	"github.com/imkuqin-zw/yggdrail-gorm/driver"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	driver.RegisterFactory("mysql", NewDialector)
}

func NewDialector(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}
