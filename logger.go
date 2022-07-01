package xgorm

import (
	"context"
	"encoding/json"
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
	log.Infof(strings.TrimRight(s, "\n"), i...)
}

func (l *logger) Warn(ctx context.Context, s string, i ...interface{}) {
	log.Warnf(strings.TrimRight(s, "\n"), i...)
}

func (l *logger) Error(ctx context.Context, s string, i ...interface{}) {
	log.Errorf(strings.TrimRight(s, "\n"), i...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	cost := time.Since(begin)
	sql, rows := fc()
	fields := map[string]interface{}{
		"sql":  sql,
		"cost": cost.Seconds(),
		"rows": rows,
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		fields["err"] = err
		data, _ := json.Marshal(fields)
		log.Errorf("%s\t%s", "gorm", string(data))
		return
	}
	if l.slowThreshold < cost && log.Enable(types.LvWarn) {
		data, _ := json.Marshal(fields)
		log.Warnf("%s\t%s", "gorm", string(data))
		return
	}
	if log.Enable(types.LvDebug) {
		data, _ := json.Marshal(fields)
		log.Debugf("%s\t%s", "gorm", string(data))
	}
	return
}

func (l *logger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	return l
}
