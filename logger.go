package xgorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imkuqin-zw/yggdrasil/pkg/log"
	"github.com/imkuqin-zw/yggdrasil/pkg/types"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

type logger struct {
	slowThreshold time.Duration
}

func (l *logger) Info(ctx context.Context, s string, i ...interface{}) {
	log.InfoFiled(fmt.Sprintf(strings.TrimRight(s, "\n"), i...), log.Context(ctx))
}

func (l *logger) Warn(ctx context.Context, s string, i ...interface{}) {
	log.WarnFiled(fmt.Sprintf(strings.TrimRight(s, "\n"), i...), log.Context(ctx))
}

func (l *logger) Error(ctx context.Context, s string, i ...interface{}) {
	log.ErrorFiled(fmt.Sprintf(strings.TrimRight(s, "\n"), i...), log.Context(ctx))
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	cost := time.Since(begin)
	sql, rows := fc()
	fields := []log.Field{
		log.String("sql", sql),
		log.Duration("cost", cost),
		log.Int64("rows", rows),
		log.Context(ctx),
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.ErrorFiled("gorm", append(fields, log.Err(err))...)
		return
	}
	if l.slowThreshold < cost && log.Enable(types.LvWarn) {
		log.WarnFiled("gorm", fields...)
		return
	}
	if log.Enable(types.LvDebug) {
		log.DebugFiled("gorm", fields...)
	}
	return
}

func (l *logger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	return l
}
