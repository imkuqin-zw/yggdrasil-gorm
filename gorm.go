package xgorm

import (
	"context"
	"time"

	"github.com/imkuqin-zw/yggdrasil-gorm/driver"
	"github.com/imkuqin-zw/yggdrasil-gorm/plugin"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/log"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Open(config *Config) *gorm.DB {
	config.SetDefault()
	cfg := &gorm.Config{
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		Logger: &logger{
			slowThreshold: config.SlowThreshold,
		},
		DryRun:      config.DryRun,
		PrepareStmt: config.PrepareStmt,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: config.NameStrategy.SingularTable,
			TablePrefix:   config.NameStrategy.TablePrefix,
			NoLowerCase:   config.NameStrategy.NoLowerCase,
		},
	}
	f := driver.GetFactory(config.Driver)
	if f == nil {
		log.Fatalf("unknown gorm driver, name: %s", config.Driver)
		return nil
	}
	db, err := gorm.Open(f(config.DSN), cfg)
	if err != nil {
		log.Fatalf("fault to connect mysql, error: %+v", err)
		return nil
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil
	}
	sqlDb.SetMaxOpenConns(config.MaxOpenConn)
	sqlDb.SetMaxIdleConns(config.MaxIdleConn)
	sqlDb.SetMaxIdleConns(config.MaxIdleConn)
	sqlDb.SetConnMaxLifetime(config.ConnMaxLifetime)

	for _, name := range config.Plugins {
		if err := db.Use(plugin.GetPlugin(name, config.Name)); err != nil {
			log.Fatalf("fault to use plugin, error: %+v", err)
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	if err := sqlDb.PingContext(ctx); err != nil {
		log.Fatalf("fault to ping mysql, error: %+v", err)
		return nil
	}
	return db
}

func NewDB(name string) *gorm.DB {
	c := new(Config)
	if err := config.Get("gorm." + name).Scan(c); err != nil {
		log.Fatalf("fault to load gorm config, error: %s", err.Error())
	}
	plugins := config.Get("gorm.global.plugins").StringSlice([]string{})
	if len(plugins) > 0 {
		c.Plugins = append(plugins, c.Plugins...)
	}
	return Open(c)
}
