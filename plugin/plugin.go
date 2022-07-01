package plugin

import "gorm.io/gorm"

type Factory func(dbName string) gorm.Plugin

var plugins = make(map[string]func(dbName string) gorm.Plugin)

func RegisterPluginFactory(name string, f Factory) {
	plugins[name] = f
}

func Get5PluginFactory(name string) Factory {
	f, _ := plugins[name]
	return f
}

func GetPlugin(name, dbName string) gorm.Plugin {
	f, ok := plugins[name]
	if !ok {
		return nil
	}
	return f(dbName)
}
