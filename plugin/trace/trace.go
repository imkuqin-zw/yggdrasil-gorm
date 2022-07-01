package trace

import (
	"fmt"

	"github.com/imkuqin-zw/yggdrasil-gorm/plugin"
	"gorm.io/gorm"
)

const (
	callBackBeforeName = "otel:before"
	callBackAfterName  = "otel:after"
	opCreate           = "INSERT"
	opQuery            = "SELECT"
	opDelete           = "DELETE"
	opUpdate           = "UPDATE"
)

func beforeName(name string) string {
	return callBackBeforeName + "_" + name
}

func afterName(name string) string {
	return callBackAfterName + "_" + name
}

type OtelPlugin struct {
	dbName string
}

func (op *OtelPlugin) Name() string {
	return "OpenTelemetryPlugin"
}

type registerCallback interface {
	Register(name string, fn func(*gorm.DB)) error
}

type gormHookFunc func(tx *gorm.DB)

func (op *OtelPlugin) Initialize(db *gorm.DB) error {
	registerHooks := []struct {
		callback registerCallback
		hook     gormHookFunc
		name     string
	}{
		// before hooks
		{db.Callback().Create().Before("gorm:before_create"), op.before(opCreate), beforeName("create")},
		{db.Callback().Query().Before("gorm:query"), op.before(opQuery), beforeName("query")},
		{db.Callback().Delete().Before("gorm:before_delete"), op.before(opDelete), beforeName("delete")},
		{db.Callback().Update().Before("gorm:before_update"), op.before(opUpdate), beforeName("update")},
		{db.Callback().Row().Before("gorm:row"), op.before(""), beforeName("row")},
		{db.Callback().Raw().Before("gorm:raw"), op.before(""), beforeName("raw")},

		// after hooks
		{db.Callback().Create().After("gorm:after_create"), op.after(opCreate), afterName("create")},
		{db.Callback().Query().After("gorm:after_query"), op.after(opQuery), afterName("select")},
		{db.Callback().Delete().After("gorm:after_delete"), op.after(opDelete), afterName("delete")},
		{db.Callback().Update().After("gorm:after_update"), op.after(opUpdate), afterName("update")},
		{db.Callback().Row().After("gorm:row"), op.after(""), afterName("row")},
		{db.Callback().Raw().After("gorm:raw"), op.after(""), afterName("raw")},
	}

	for _, h := range registerHooks {
		if err := h.callback.Register(h.name, h.hook); err != nil {
			return fmt.Errorf("register %s hook: %w", h.name, err)
		}
	}

	return nil
}

func newPlugin(dbName string) gorm.Plugin {
	return &OtelPlugin{dbName: dbName}
}

func init() {
	plugin.RegisterPluginFactory("trace", newPlugin)
}
