package host

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"sobot/core"
	"sobot/core/app"
	"sobot/global"
	"sobot/pkg/userInput"
	"strconv"
	"strings"
)

var (
	Guide = "请选择主机序号"
)

type HostTerminal struct {
	core.CommonData
}

func NewHostTerminal() *HostTerminal {
	return &HostTerminal{core.CommonData{
		Name:      "terminal",
		HumanName: "主机终端",
		Features:  nil,
	}}
}

func (ht *HostTerminal) Run() error {
	for {
		headers := []string{"序号", "主机", "部署的服务"}
		var data [][]string
		for index, host := range global.HostSetting.IPS {
			appNames := app.GetAppsByHost(host)
			data = append(data, []string{strconv.FormatInt(int64(index+1), 10), host,
				Tmp(appNames)})
		}
		core.Render(headers, data)
		choice, err := userInput.UserString(Guide)
		if err != nil {
			continue
		}
		number, err := core.ChoiceJudge(choice)
		if err != nil {
			if errors.Is(err, core.ErrBack) {
				return nil
			} else {
				continue
			}
		}

		if number > len(global.HostSetting.IPS) {
			continue
		}
		ht.terminal(global.HostSetting.IPS[number-1])
	}
}

func Tmp(dest []string) string {
	var ret string
	if len(dest) <= 8 {
		ret = strings.Join(dest, ",")
	} else if len(dest) > 8 && len(dest) <= 16 {
		ret = strings.Join(dest[:8], ",") + "\n" + strings.Join(dest[8:], ",")
	} else if len(dest) > 16 && len(dest) <= 24 {
		ret = strings.Join(dest[:8], ",") + "\n" + strings.Join(dest[8:16], ",") + "\n" + strings.Join(dest[16:], ",")
	} else if len(dest) > 24 && len(dest) <= 32 {
		ret = strings.Join(dest[:8], ",") + "\n" + strings.Join(dest[8:16], ",") + "\n" + strings.Join(dest[16:24], ",") + "\n" + strings.Join(dest[24:], ",")
	} else {
		ret = strings.Join(dest[:8], ",") + "\n" + strings.Join(dest[8:16], ",") + "\n" + strings.Join(dest[16:24], ",") + "\n" + strings.Join(dest[24:32], ",") + "\n" + strings.Join(dest[32:], ",")
	}

	return ret
}

func (ht *HostTerminal) terminal(ip string) error {
	connect := fmt.Sprintf("ssh -q -o StrictHostKeyChecking=no %s@'%s' -p %d", global.HostSetting.User, ip, global.HostSetting.Port)
	cmd := exec.Command("bash", "-c", connect)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
