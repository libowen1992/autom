package app

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"sobot/core"
	"sobot/global"
	"sobot/pkg/convert"
	"sobot/pkg/userInput"
)

var (
	appLogGuide   = "请输入要查看的服务名称，支持模糊匹配,b返回"
	choicAppGuide = "请选择app序号"
	CHoiceHost    = "请选择主机序号"
	ErrExec       = "执行命令失败"
)

type AppLog struct {
	core.CommonData
}

func NewAppLog() *AppLog {
	return &AppLog{core.CommonData{
		Name:      "log",
		HumanName: "查看日志",
		Features:  nil,
	}}
}

func (ag *AppLog) Run() error {
	for {
		fuzzName, err := userInput.UserString(appLogGuide)
		if fuzzName == "q" || fuzzName == "quit" {
			os.Exit(1)
		} else if fuzzName == "b" {
			return nil
		}
		if err != nil {
			global.Logger.Error(err)
			continue
		}
		apps := GetAppsByNameFuzzy(fuzzName)
		if len(apps) == 0 {
			global.Logger.Warnf("%s服务没有找到", fuzzName)
			continue
		}
		if len(apps) == 1 {
			if err := singleAppStep(apps[0]); err != nil {
				if errors.Is(err, core.ErrBack) || errors.Is(err, core.ErrContinue) {
					continue
				} else {
					global.Logger.Error(err)
					return err
				}
			}
		}

		RenderMulApps(apps)
		choice, err := userInput.UserString(choicAppGuide)
		if err != nil {
			global.Logger.Error(err)
			return err
		}
		appNum, err := core.ChoiceJudge(choice)
		if err != nil {
			if errors.Is(err, core.ErrBack) {
				return nil
			} else {
				global.Logger.Error(err)
				continue
			}
		}
		if appNum > len(apps) {
			continue
		}

		app := apps[appNum-1]
		if err := singleAppStep(app); err != nil {
			if errors.Is(err, core.ErrBack) || errors.Is(err, core.ErrContinue) {
				continue
			} else {
				global.Logger.Error(err)
				return err
			}
		}

		return nil
	}
}

func RenderMulApps(apps []*App) {
	headers := []string{"序号", "服务"}
	var data [][]string
	for index, app := range apps {
		item := []string{convert.IntToStr(index + 1), app.Name}
		data = append(data, item)
	}
	core.Render(headers, data)
}

func RenderAppMulHost(app *App) {
	headers := []string{"序号", "主机"}
	var data [][]string
	for index, node := range app.Nodes {
		item := []string{convert.IntToStr(index + 1), node}
		data = append(data, item)
	}
	core.Render(headers, data)
}

func SSHAppLog(host string, app *App) error {
	cmdStr := fmt.Sprintf("ssh -q -o StrictHostKeyChecking=no %s@'%s' -p %d -t 'cd %s;bash --login'",
		global.HostSetting.User,
		host, global.HostSetting.Port, app.LogPath)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func singleAppStep(app *App) error {
	if len(app.Nodes) == 0 {
		global.Logger.Warnf("%s服务没有节点", app.Name)
		return core.ErrContinue
	} else if len(app.Nodes) == 1 {
		if err := SSHAppLog(app.Nodes[0], app); err != nil {
			return err
		}
	} else {
		RenderAppMulHost(app)
		choice, err := userInput.UserString(CHoiceHost)
		if err != nil {
			global.Logger.Error(err)
			return err
		}
		HostNum, err := core.ChoiceJudge(choice)
		if err != nil {
			if errors.Is(err, core.ErrBack) {
				return core.ErrBack
			} else {
				global.Logger.Error(err)
				return core.ErrContinue
			}
		}
		if HostNum > len(app.Nodes) {
			return core.ErrContinue
		}

		if err := SSHAppLog(app.Nodes[HostNum-1], app); err != nil {
			return err
		}
	}

	return nil
}
