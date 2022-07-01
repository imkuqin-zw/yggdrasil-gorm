package driver

import "gorm.io/gorm"

var driverFactory = make(map[string]func(string) gorm.Dialector)

func RegisterFactory(name string, f func(string) gorm.Dialector) {
	driverFactory[name] = f
}

func GetFactory(name string) func(string) gorm.Dialector {
	f, _ := driverFactory[name]
	return f
}
