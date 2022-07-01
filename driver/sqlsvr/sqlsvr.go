package sqlsvr

import (
	"github.com/imkuqin-zw/yggdrasil-gorm/driver"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	driver.RegisterFactory("sqlserver", NewDialector)
}

func NewDialector(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}
