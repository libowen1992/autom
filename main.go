package main

import (
	"flag"
	"github.com/pkg/errors"
	"os/user"
	"sobot/core/mainMenu"
	"sobot/global"
	"sobot/pkg/logger"
	"sobot/setting"
)

func main() {
	// init logger
	global.Logger = logger.NewStdoutConsole()

	// init setting
	if err := SetupSetting(SetupCommand()); err != nil {
		global.Logger.Fatal(err)
	}
	// 判断用户
	if err := JudgeUser(); err != nil {
		global.Logger.Fatal(err)
	}

	m := mainMenu.New()
	m.Run()
}

func SetupCommand() (file string) {
	flag.StringVar(&file, "conf", global.LogFilePath, "配置文件路径")
	flag.Parse()
	return
}

func SetupSetting(file string) error {
	st := setting.New()
	if err := st.Init(file); err != nil {
		return errors.Wrap(err, "读取配置文件失败")
	}
	if err := st.SetSection("common", &global.CommSetting); err != nil {
		return err
	}
	if err := st.SetSection("hosts", &global.HostSetting); err != nil {
		return err
	}
	if err := st.SetSection("redis", &global.RedisSetting); err != nil {
		return err
	}
	if err := st.SetSection("mysql", &global.MySQLSetting); err != nil {
		return err
	}
	if err := st.SetSection("apps", &global.AppSetting); err != nil {
		return err
	}
	return nil
}

func JudgeUser() error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	if user.Username != global.HostSetting.User {
		return errors.Errorf("当前执行用户为%s, 请使用%s用户执行", user.Username, global.HostSetting.User)
	}
	return nil
}
