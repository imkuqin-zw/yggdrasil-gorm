package xgorm

import (
	"time"
)

const (
	defaultMaxIdleConn     = 10
	defaultMaxOpenConn     = 100
	defaultConnMaxLifetime = time.Second * 300
	defaultSlowThreshold   = time.Millisecond * 500
)

// conf options
type Config struct {
	// DSN: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	Name   string
	DSN    string
	Driver string

	PrepareStmt            bool
	DryRun                 bool
	SkipDefaultTransaction bool
	MaxIdleConn            int
	MaxOpenConn            int
	ConnMaxLifetime        time.Duration
	NameStrategy           struct {
		TablePrefix   string
		SingularTable bool
		NoLowerCase   bool
	}
	// 慢日志阈值
	SlowThreshold time.Duration

	Plugins string
}

// Check
func (config *Config) SetDefault() {
	if config.MaxIdleConn == 0 {
		config.MaxIdleConn = defaultMaxIdleConn
	}
	if config.MaxOpenConn == 0 {
		config.MaxOpenConn = defaultMaxOpenConn
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = defaultConnMaxLifetime
	}
	if config.SlowThreshold < time.Millisecond {
		config.SlowThreshold = defaultSlowThreshold
	}
}
