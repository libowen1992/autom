package app

import (
	"sobot/core"
)

// AppManager 服务管理结构体
type AppManager struct {
	core.CommonData
}

func NewAppManager() *AppManager {
	am := &AppManager{
		CommonData: core.CommonData{
			Name:      "app",
			HumanName: "服务管理",
			Features: []core.Project{
				NewAppLog(),
				NewUpdate(),
				NewRoll(),
				NewHistory(),
			},
		},
	}
	am.Init()
	return am
}

// Init 初始化apps
func (am *AppManager) Init() {
	var appNames []string
	appNames = append(Wars, Jars...)
	appNames = append(appNames, Tars...)
	appNames = append(appNames, TarsNew...)
	for _, name := range appNames {
		app := &App{
			Name: name,
		}
		app.SetPKGName()
		app.SetPKGPath()
		app.SetStartCmd()
		app.SetStopCmd()
		app.SetNodes()
		app.SetLogPath()
		Apps = append(Apps, app)
	}
}
