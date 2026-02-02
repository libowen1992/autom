package global

import (
	"go.uber.org/zap"
	"sobot/setting"
)

var (
	LogFilePath   = "./conf.yml"
	SnowflakeTime = "2021-04-02"
	HostSetting   *setting.Hosts
	MySQLSetting  *setting.MySQL
	RedisSetting  *setting.Redis
	AppSetting    setting.Apps
	CommSetting   setting.Common
	//ESSetting *conf.ES
	//ServiceSetting *conf.Services
	Logger *zap.SugaredLogger
)
