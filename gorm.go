package xgorm

import (
	"context"
	"strings"
	"time"

	"github.com/imkuqin-zw/yggdrasil-gorm/driver"
	"github.com/imkuqin-zw/yggdrasil-gorm/plugin"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/log"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DB = gorm.DB

var (
	// ErrRecordNotFound record not found error
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	// ErrNotImplemented not implemented
	ErrNotImplemented = gorm.ErrNotImplemented
	// ErrMissingWhereClause missing where clause
	ErrMissingWhereClause = gorm.ErrMissingWhereClause
	// ErrUnsupportedRelation unsupported relations
	ErrUnsupportedRelation = gorm.ErrUnsupportedRelation
	// ErrPrimaryKeyRequired primary keys required
	ErrPrimaryKeyRequired = gorm.ErrPrimaryKeyRequired
	// ErrModelValueRequired model value required
	ErrModelValueRequired = gorm.ErrModelValueRequired
	// ErrInvalidData unsupported data
	ErrInvalidData = gorm.ErrInvalidData
	// ErrUnsupportedDriver unsupported driver
	ErrUnsupportedDriver = gorm.ErrUnsupportedDriver
	// ErrRegistered registered
	ErrRegistered = gorm.ErrRegistered
	// ErrInvalidField invalid field
	ErrInvalidField = gorm.ErrInvalidField
	// ErrEmptySlice empty slice found
	ErrEmptySlice = gorm.ErrEmptySlice
	// ErrDryRunModeUnsupported dry run mode unsupported
	ErrDryRunModeUnsupported = gorm.ErrDryRunModeUnsupported
	// ErrInvalidDB invalid db
	ErrInvalidDB = gorm.ErrInvalidDB
	// ErrInvalidValue invalid value
	ErrInvalidValue = gorm.ErrInvalidValue
	// ErrInvalidValueOfLength invalid values do not match length
	ErrInvalidValueOfLength = gorm.ErrInvalidValueOfLength
	// ErrPreloadNotAllowed preload is not allowed when count is used
	ErrPreloadNotAllowed = gorm.ErrPreloadNotAllowed
)

func Open(config *Config) *DB {
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

	if len(config.Plugins) > 0 {
		plugins := strings.Split(config.Plugins, ",")
		for _, name := range plugins {
			if err := db.Use(plugin.GetPlugin(name, config.Name)); err != nil {
				log.Fatalf("fault to use plugin, error: %+v", err)
				return nil
			}
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

func NewDB(name string) *DB {
	c := new(Config)
	if err := config.Get("gorm." + name).Scan(c); err != nil {
		log.Fatalf("fault to load gorm config, error: %s", err.Error())
	}
	plugins := config.Get("gorm.global.plugins").String("")
	if len(plugins) > 0 && len(c.Plugins) > 0 {
		c.Plugins = plugins + "," + c.Plugins
	}
	return Open(c)
}
