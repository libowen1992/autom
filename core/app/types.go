package app

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"sobot/global"
	"strings"
	"text/template"
	"time"
)

var (
	WarSuffix = "war"
	JarSuffix = "jar"
	TarSuffix = "tar.gz"
	Wars      = []string{"basic-login", "basic-customer", "basic-monitor", "basic-pass", "basic-set", "customer-open", "set-open",
		"boss-admin", "boss-service", "open-service", "kb-service", "robot-service", "robot-modeling", "ws-service", "ws-open",
		"ws-report", "ws-tenant", "ws-job", "chat-data", "chat-input", "chat-job", "chat-job-new", "chat-kwb", "chat-msg", "chat-report",
		"chat-sdk", "chat-third", "chat-visit", "chat-wb", "chat-web", "chat-wx", "tenant-service", "call-bill", "call-data",
		"call-ivr", "call-job", "call-open", "call-report", "comm-admin", "icall-admin", "icall-data", "icall-import", "icall-job",
		"icall-link", "icall-monitor", "icall-open", "icall-service", "link-service", "link-translate","chat-report-new","chat-admin-sdk","im-service",}

	Jars = []string{"basic-config-service", "basic-message", "customer-base", "job-service", "set-base6", "message-base",
		"qcs-base", "qcs-common", "kb-job", "robot-base", "robot-base-prv", "robot-base-suggest", "robot-base-word", "robot-base-learn",
		"robot-work-index", "robot-work-suggest", "robot-knowledge-optimize", "robot-inspection-index", "ws-base-es6", "chat-inner-im",
		"chat-inner-im-consumer", "chat-reply", "chat-reply-consumer", "chat-reply-job", "chat-report-base", "chat-set", "qcs-chat", "data-base-es6",
		"chat-analysis", "chat-offline-analysis", "chat-realtime-monitor", "offlineJobProcess-es6", "call-realtime-analysis", "ivr-call-analysis",
		"basic-public", "basic-system", "boss-base3", "boss-service3", "boss-web3", "gateway", "crm-base", "crm-channel-qywx",
		"crm-job", "crm-open", "crm-report", "crm-user-service", "qywx-base", "qywx-callback", "qywx-job", "qywx-report",
		"qywx-service", "call-basic", "call-data-base", "comm-cti-core", "comm-cti-message", "cti-asr", "cti-backend", "cti-cache",
		"cti-config", "cti-data", "cti-new-websocket", "cti-outbound", "cti-service", "cti-websocket", "cti-web-tts", "cti-xnumber",
		"qcs-call", "tms-data", "tms-main", "tms-job", "icall-base", "link-fs", "sobot-cloud-inspection-main", "sobot-cloud-inspection-open",
		"sobot-cloud-inspection-job", "sobot-cloud-eureka", "sobot-cloud-zuul", "job-hh", "cti-nxcapi", "qywx-material-service",
		"qywx-route-callback", "qywx-kf-callback", "qywx-kf-service", "zk-sync-data-base", "zk-sync-data-service", "cti-datasync", "comm-open", "call-realtime-analysis01",
		"qywx-marketing-job", "qywx-rule-engine", "zk-marketing-base", "zk-marketing-job", "zk-marketing-report", "zk-marketing-service",
		"crm-export-base", "crm-qywx-base", "callcenter-mid-query","chat-router","chat-facebook","chat-whatsapp","chat-googleplay","chat-third-platform","qywx-welcome","qywx-upgrade","zk-marketing-open","qywx-lining-lable-import","chat-push",
	}

	Tars = []string{"agent", "callcenter", "call", "console", "dashboard", "framev2", "h5v2", "inner_business", "pcv2",
		"platformpro", "robotsupportpro", "scrm_minipro", "ticketclient", "wechatbar", "licensemanagement", "admins", "monitor", "mailservices"}
	TarsNew       = []string{"account", "crm", "home", "task-center", "work-order", "zkwxm"}
	DeployRecords = "./deploy_records.csv"
	DeployType    = "update"
	RollType      = "rollback"
	MaxRecords    = 20
)

// GetAppSuffix 根据服务名称获取服务包后缀
func GetAppSuffix(appName string) string {
	for _, app := range Wars {
		if appName == app {
			return WarSuffix
		}
	}

	for _, app := range Jars {
		if appName == app {
			return JarSuffix
		}
	}

	for _, app := range Tars {
		if appName == app {
			return TarSuffix
		}
	}

	for _, app := range TarsNew {
		if appName == app {
			return TarSuffix
		}
	}

	return ""
}

// IsNewFront 检查是否是新前端包 安装包路径不一样
func IsZhikeFront(appName string) bool {
	for _, app := range TarsNew {
		if appName == app {
			return true
		}
	}
	return false
}

// 服务
type App struct {
	// 服务名称
	Name string
	// 服务包名称
	PKGName string
	// 服务包安装目录
	PKGPath string
	// 日志目录
	LogPath string
	// 服务启动命令
	StartCmd string
	// 服务关闭命令
	StopCmd string
	// 服务部署的节点
	Nodes []string
}

var Apps []*App

func GetApp(pkgName string) (app *App, err error) {
	for _, app := range Apps {
		if pkgName == app.PKGName {
			return app, nil
		}
	}
	return nil, errors.New("没有匹配的服务")
}

func GetAppByName(name string) (app *App, err error) {
	for _, app := range Apps {
		if name == app.Name {
			return app, nil
		}
	}
	return nil, errors.New("没有匹配的服务")
}

func GetAppsByHost(host string) []string {
	var appNames []string
	for _, app := range Apps {
		if eleInSlice(host, app.Nodes) {
			appNames = append(appNames, app.Name)
		}
	}
	return appNames
}

func GetAppsByNameFuzzy(name string) []*App {
	var apps []*App
	for _, app := range Apps {
		if strings.Contains(app.Name, name) {
			apps = append(apps, app)
		}
	}
	return apps
}

func eleInSlice(name string, dest []string) bool {
	for _, ele := range dest {
		if name == ele {
			return true
		}
	}
	return false
}

// GetSuffix 服务后缀
func (a *App) GetSuffix() string {
	return GetAppSuffix(a.Name)
}

// SetPKGName 服务包名称
func (a *App) SetPKGName() {
	a.PKGName = fmt.Sprintf("%s.%s", a.Name, a.GetSuffix())
}

// SetPKGPath 服务包所在路径
func (a *App) SetPKGPath() {
	switch a.GetSuffix() {
	case WarSuffix:
		a.PKGPath = filepath.Join(global.CommSetting.AppInstallPath, "tomcat", a.Name, "webapps")
	case JarSuffix:
		a.PKGPath = filepath.Join(global.CommSetting.AppInstallPath, "offlinetask", a.Name)
	case TarSuffix:
		if IsZhikeFront(a.Name) {
			a.PKGPath = "/data/static/console-platform"
		} else {
			a.PKGPath = "/data/static"
		}
	}
}

// SetStartCmd 设置服务启动命令
func (a *App) SetStartCmd() {
	switch a.GetSuffix() {
	case WarSuffix:
		a.StartCmd = fmt.Sprintf("cd %s && ./startup.sh start",
			filepath.Join(global.CommSetting.AppInstallPath, "tomcat", a.Name, "bin"))
	case JarSuffix:
		a.StartCmd = fmt.Sprintf("cd %s && ./start.sh start",
			filepath.Join(global.CommSetting.AppInstallPath, "offlinetask", a.Name))
	}
}

// SetStopCmd 设置服务关闭命令
func (a *App) SetStopCmd() {
	switch a.GetSuffix() {
	case WarSuffix:
		a.StopCmd = fmt.Sprintf("ps -ef|grep -w %s |grep -v grep | awk '{print \\$2}'|xargs kill -9",
			filepath.Join(global.CommSetting.AppInstallPath, "tomcat", a.Name, "conf"))
	case JarSuffix:
		a.StopCmd = fmt.Sprintf("cd %s && ./start.sh stop",
			filepath.Join(global.CommSetting.AppInstallPath, "offlinetask", a.Name))
	}
}

// SetStopCmd 设置服务节点
func (a *App) SetNodes() {
	for _, settingApp := range global.AppSetting {
		if settingApp.Name == a.Name {
			a.Nodes = settingApp.Nodes
		}
	}
}

// set Log Path
func (a *App) SetLogPath() {
	switch a.GetSuffix() {
	case WarSuffix:
		a.LogPath = filepath.Join(global.CommSetting.AppInstallPath, "tomcat", a.Name, "logs")
	case JarSuffix:
		a.LogPath = filepath.Join(global.CommSetting.AppInstallPath, "offlinetask", a.Name, "logs")
	}
}

// GetScriptTemplate
func (a *App) GetDeployScriptTemplate() string {
	switch a.GetSuffix() {
	case WarSuffix:
		return "./tmpl/tomcat_update.sh"
	case JarSuffix:
		return "./tmpl/jar_update.sh"
	case TarSuffix:
		return "./tmpl/front_update.sh"
	default:
		return ""
	}
}

func (a *App) GetRollbackScriptTemplate() string {
	switch a.GetSuffix() {
	case WarSuffix:
		return "./tmpl/tomcat_rollback.sh"
	case JarSuffix:
		return "./tmpl/jar_rollback.sh"
	case TarSuffix:
		return "./tmpl/front_rollback.sh"
	default:
		return ""
	}
}

// 服务部署结构体
type DeployApp struct {
	// 服务结构体
	*App
	// 服务更新包URL
	PKGURL string
	// 部署结束时间
	EndTime time.Time
	// 更新脚本
	DeployScript string
	// 回滚脚本
	RollBackScript string
}

func (da *DeployApp) GetDeployScript(deployId string) error {
	scriptBaseDir := filepath.Join("files", deployId)
	da.DeployScript = filepath.Join(scriptBaseDir, fmt.Sprintf("%s_update.sh", da.Name))
	tmplateScript := da.GetDeployScriptTemplate()
	if err := os.MkdirAll(scriptBaseDir, 0755); err != nil {
		return err
	}
	data := UpdateAppData{
		AppName:         da.Name,
		AppPkgName:      da.PKGName,
		AppPkgPath:      da.PKGPath,
		PkgDownLoadDir:  filepath.Join("files", deployId),
		PkgUrl:          da.PKGURL,
		Nodes:           fmt.Sprintf("(%s)", strings.Join(da.Nodes, " ")),
		RemoteTmpDir:    filepath.Join("/tmp", deployId),
		RemoteBackupDir: filepath.Join(global.CommSetting.AppBackupPath, deployId),
		StopCmd:         da.StopCmd,
		StartCmd:        da.StartCmd,
		User:            global.HostSetting.User,
		Port:            global.HostSetting.Port,
	}
	df, err := os.Create(da.DeployScript)
	if err != nil {
		return err
	}
	defer df.Close()
	tmpl, err := template.ParseFiles(tmplateScript)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(df, data); err != nil {
		return err
	}
	return nil
}

func (da *DeployApp) GetRollBackScript(deployId string) error {
	scriptBaseDir := filepath.Join("files", deployId)
	da.RollBackScript = filepath.Join(scriptBaseDir, fmt.Sprintf("%s_rollback.sh", da.Name))
	templateScript := da.GetRollbackScriptTemplate()
	data := UpdateAppData{
		AppName:         da.Name,
		AppPkgName:      da.PKGName,
		AppPkgPath:      da.PKGPath,
		PkgDownLoadDir:  filepath.Join("files", deployId),
		PkgUrl:          da.PKGURL,
		Nodes:           fmt.Sprintf("(%s)", strings.Join(da.Nodes, " ")),
		RemoteTmpDir:    filepath.Join("/tmp", deployId),
		RemoteBackupDir: filepath.Join(global.CommSetting.AppBackupPath, deployId),
		StopCmd:         da.StopCmd,
		StartCmd:        da.StartCmd,
		User:            global.HostSetting.User,
		Port:            global.HostSetting.Port,
	}
	df, err := os.Create(da.RollBackScript)
	if err != nil {
		return err
	}
	defer df.Close()
	tmpl, err := template.ParseFiles(templateScript)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(df, data); err != nil {
		return err
	}
	return nil
}

type UpdateAppData struct {
	User            string
	Port            int
	AppName         string
	AppPkgName      string
	AppPkgPath      string
	PkgDownLoadDir  string
	PkgUrl          string
	StopCmd         string
	StartCmd        string
	Nodes           string
	RemoteTmpDir    string
	RemoteBackupDir string
	Cost            time.Duration
}

// 服务部署总结构体
type AppDeployTask struct {
	// 部署ID
	DeployId string
	// 部署开始时间
	StartTime time.Time
	// 部署结束时间
	EndTime time.Time
	// 部署的服务
	Deploys chan *DeployApp
	// 部署结果
	DeploysResult chan *DeployApp
}

type DeployRecord struct {
	// 部署ID
	DeployId   string
	DeployType string
	AppName    string
	AppPkgName string
}
