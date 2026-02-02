package mysql

import (
	"fmt"
	"os"
	"os/exec"
	"sobot/core"
	"sobot/global"
)

type MySQLTernimal struct {
	core.CommonData
}

func NewTernimal() *MySQLTernimal {
	return &MySQLTernimal{core.CommonData{
		Name:      "ternimal",
		HumanName: "MySQL终端",
		Features:  nil,
	}}
}

func (rt *MySQLTernimal) Run() error {
	connect := fmt.Sprintf("mysql -h %s -u%s -p%d -p'%s'",
		global.MySQLSetting.IP, global.MySQLSetting.User,
		global.MySQLSetting.Port, global.MySQLSetting.Pass)
	cmd := exec.Command("bash", "-c", connect)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
