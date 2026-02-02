package app

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sobot/core"
	"sobot/global"
	"sobot/pkg/convert"
	"sobot/pkg/userInput"
	"strings"
	"sync"
	"time"
)

const (
	Guide         = "请输入要更新的包的url, 多个包使用空格分开,返回b"
	ContinueGuide = "确认请输入y,否则输入n"
)

type Update struct {
	core.CommonData
}

func NewUpdate() *Update {
	return &Update{core.CommonData{
		Name:      "update",
		HumanName: "更新",
		Features:  nil,
	}}
}

func (u *Update) Run() error {
	var wg sync.WaitGroup
	userStrings, err := userInput.UserString(Guide)
	if err != nil {
		return err
	}

	if userStrings == "q" || userStrings == "quit" {
		os.Exit(1)
	} else if userStrings == "b" {
		return nil
	}

	deployApps := ParseDeployApps(userStrings)
	if len(deployApps) == 0 {
		return errors.New("没有服务需要部署")
	}
	renderDeploysStart(deployApps)
	c, err := userInput.UserString(ContinueGuide)
	if err != nil {
		return err
	}
	if strings.ToUpper(c) != "Y" {
		return nil
	}

	var now = time.Now()
	deployId := convert.TimeToDateTimeFormat(now)
	dt := &AppDeployTask{
		DeployId:      deployId,
		StartTime:     now,
		Deploys:       make(chan *DeployApp, len(deployApps)),
		DeploysResult: make(chan *DeployApp, len(deployApps)),
	}
	for _, deployApp := range deployApps {
		if err := deployApp.GetDeployScript(deployId); err != nil {
			global.Logger.Error(err)
			continue
		}
		if err := deployApp.GetRollBackScript(deployId); err != nil {
			global.Logger.Error(err)
			continue
		}
		dt.Deploys <- deployApp
	}
	close(dt.Deploys)

	for deploy := range dt.Deploys {
		wg.Add(1)
		go func(deploy *DeployApp) error {
			defer wg.Done()
			global.Logger.Infof("%s: 开始部署", deploy.Name)
			cmd := exec.Command("/bin/bash", deploy.DeployScript)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return errors.WithMessage(err, string(output))
			}
			deploy.EndTime = time.Now()
			dt.DeploysResult <- deploy

			global.Logger.Infof("%s: 完成部署", deploy.Name)
			return nil
		}(deploy)
	}
	wg.Wait()

	close(dt.DeploysResult)
	global.Logger.Infof("所有服务完成部署，耗时 %s", time.Since(dt.StartTime).String())
	if err := renderDeploysEndWriteRecords(dt.DeploysResult, dt.StartTime); err != nil {
		global.Logger.Error(err)
		return err
	}

	return nil
}

func (u *Update) RunNew() error {
	for {
		var wg sync.WaitGroup
		urls, err := userInput.UserString(Guide)
		if err != nil {
			return err
		}
		deployApps := ParseDeployApps(urls)
		if len(deployApps) == 0 {
			return errors.New("没有服务需要部署")
		}
		renderDeploysStart(deployApps)
		c, err := userInput.UserString(ContinueGuide)
		if err != nil {
			return err
		}
		if strings.ToUpper(c) != "Y" {
			return nil
		}

		var now = time.Now()
		deployId := convert.TimeToDateTimeFormat(now)
		dt := &AppDeployTask{
			DeployId:      deployId,
			StartTime:     now,
			Deploys:       make(chan *DeployApp, len(deployApps)),
			DeploysResult: make(chan *DeployApp, len(deployApps)),
		}
		for _, deployApp := range deployApps {
			if err := deployApp.GetDeployScript(deployId); err != nil {
				global.Logger.Error(err)
				continue
			}
			if err := deployApp.GetRollBackScript(deployId); err != nil {
				global.Logger.Error(err)
				continue
			}
			dt.Deploys <- deployApp
		}
		close(dt.Deploys)

		for deploy := range dt.Deploys {
			wg.Add(1)
			go func(deploy *DeployApp) error {
				defer wg.Done()
				global.Logger.Infof("%s: 开始部署", deploy.Name)
				cmd := exec.Command("/bin/bash", deploy.DeployScript)
				output, err := cmd.CombinedOutput()
				if err != nil {
					return errors.WithMessage(err, string(output))
				}
				deploy.EndTime = time.Now()
				dt.DeploysResult <- deploy

				global.Logger.Infof("%s: 完成部署", deploy.Name)
				return nil
			}(deploy)
		}
		wg.Wait()

		close(dt.DeploysResult)
		global.Logger.Infof("所有服务完成部署，耗时 %s", time.Since(dt.StartTime).String())
		if err := renderDeploysEndWriteRecords(dt.DeploysResult, dt.StartTime); err != nil {
			global.Logger.Error(err)
			return err
		}

		return nil
	}
}

func ParseDeployApps(urls string) (deployApps []*DeployApp) {
	urlSplits := strings.Fields(urls)
	if len(urlSplits) == 0 {
		return
	}
	for _, url := range urlSplits {
		if err := CheckUrl(url); err != nil {
			global.Logger.Error(err)
			continue
		}
		split := strings.Split(strings.TrimSpace(url), "/")
		pkgName := split[len(split)-1]
		app, err := GetApp(pkgName)
		if err != nil {
			global.Logger.Warnf("没有对应的服务%s", pkgName)
			continue
		}
		if len(app.Nodes) == 0 {
			global.Logger.Warnf("%s服务没有对应的节点可部署，请检查配置", app.Name)
			continue
		}
		deployApp := &DeployApp{
			App:    app,
			PKGURL: url,
		}

		deployApps = append(deployApps, deployApp)
	}
	return
}

func CheckUrl(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("get pkg url not 200 %s, please check network", url)
	}
	return nil
}

func renderDeploysStart(apps []*DeployApp) {
	headers := []string{"服务", "节点", "更新包"}
	var data [][]string
	for _, app := range apps {
		data = append(data, []string{app.Name,
			strings.Join(app.Nodes, "\n"), app.PKGURL})
	}
	core.Render(headers, data)
}

func renderDeploysEndWriteRecords(results chan *DeployApp, startTime time.Time) error {
	headers := []string{"ID", "类型", "服务", "耗时", "节点", "更新包"}
	var data [][]string
	for ret := range results {
		data = append(data, []string{convert.TimeToDateTimeFormat(startTime),
			DeployType,
			ret.Name,
			ret.EndTime.Sub(startTime).String(),
			strings.Join(ret.Nodes, "\n"),
			ret.PKGURL})
	}

	drf, err := os.OpenFile(DeployRecords, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer drf.Close()

	w := csv.NewWriter(drf)
	w.WriteAll(data)

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	core.Render(headers, data)
	return nil
}
