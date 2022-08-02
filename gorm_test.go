package xgorm

import (
	"testing"

	_ "github.com/imkuqin-zw/yggdrasil-gorm/driver/mysql"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/stretchr/testify/assert"
)

type Client struct {
	ID           uint32
	ClientID     string
	ClientSecret string
}

func TestNewDB(t *testing.T) {
	config.Set("gorm.test.driver", "mysql")
	config.Set("gorm.test.nameStrategy.singularTable", true)
	config.Set("gorm.test.dsn", "root:nihao123,./@tcp(sh-cynosdbmysql-grp-3r7tf1qw.sql.tencentcdb.com:21903)/uuid?charset=utf8mb4&parseTime=True&loc=Local")
	db := NewDB("test")
	client := &Client{}
	b, err := FindOne(db, client)
	assert.True(t, b)
	assert.Nil(t, err)
}
