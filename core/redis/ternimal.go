package redis

import (
	"fmt"
	"os"
	"os/exec"
	"sobot/core"
	"sobot/global"
)

type RedisTernimal struct {
	core.CommonData
}

func NewRedisTernimal() *RedisTernimal {
	return &RedisTernimal{core.CommonData{
		Name:      "ternimal",
		HumanName: "Redis终端",
		Features:  nil,
	}}
}

func (rt *RedisTernimal) Run() error {
	connect := fmt.Sprintf("redis-cli --raw -a '%s' -h %s -p %d",
		global.RedisSetting.Pass, global.RedisSetting.IP, global.RedisSetting.Port)
	cmd := exec.Command("bash", "-c", connect)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
