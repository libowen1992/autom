package mainMenu

import (
	"sobot/core"
	"sobot/core/app"
	"sobot/core/host"
	"sobot/core/license"
	"sobot/core/mysql"
	"sobot/core/redis"
)

type MainMenu struct {
	core.CommonData
}

func New() *MainMenu {
	var m = &MainMenu{}
	m.Register(app.NewAppManager())
	m.Register(host.New())
	m.Register(redis.New())
	m.Register(mysql.New())
	m.Register(license.New())
	return m
}
