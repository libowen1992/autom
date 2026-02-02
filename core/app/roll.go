package app

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"sobot/core"
	"sobot/global"
	"sobot/pkg/userInput"
	"strconv"
	"strings"
	"time"
)

var (
	ChoiceRollNumGuide = "请选择要回滚的数字序号,b返回,q退出"
)

type Roll struct {
	core.CommonData
}

func NewRoll() *Roll {
	return &Roll{core.CommonData{
		Name:      "roll",
		HumanName: "回滚",
		Features:  nil,
	}}
}

func (r *Roll) Run() error {
	for {
		records, err := GetDeployRecords()
		if err != nil {
			return err
		}
		RenderRollMenu(records)

		choice, err := userInput.UserString(ChoiceRollNumGuide)
		if err != nil {
			global.Logger.Error(err)
			return err
		}
		number, err := core.ChoiceJudge(choice)
		if err != nil {
			if errors.Is(err, core.ErrBack) {
				return nil
			} else {
				global.Logger.Error(err)
				continue
			}
		}

		if number > len(records) {
			global.Logger.Errorf("请在范围内选择")
			continue
		}

		rawRollRecord := records[number-1]
		deployId := rawRollRecord[0]
		appName := rawRollRecord[2]
		global.Logger.Warnf("您选择回滚的服务为: %s, 回滚的ID: %s", appName, deployId)
		c, err := userInput.UserString(ContinueGuide)
		if err != nil {
			return err
		}
		if strings.ToUpper(c) != "Y" {
			return nil
		}

		global.Logger.Infof("%s: 开始回滚，%s", appName, deployId)
		startTime := time.Now()
		script := filepath.Join("files", deployId, fmt.Sprintf("%s_rollback.sh", appName))
		cmd := exec.Command("/bin/bash", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errors.WithMessage(err, string(output))
		}
		cost := time.Since(startTime).String()
		global.Logger.Infof("%s: 完成回滚，耗时%s", appName, cost)
		return nil
	}

}

func GetDeployRecords() ([][]string, error) {
	drf, err := os.OpenFile(DeployRecords, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer drf.Close()

	reader := csv.NewReader(drf)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) > MaxRecords {
		records = records[len(records)-MaxRecords:]
	}

	return Reverse(records), nil
}

func RenderRollMenu(records [][]string) {
	headers := []string{"序号", "ID", "类型", "服务", "耗时", "节点", "更新包"}
	core.Render(headers, RollRecordsAddNum(records))
}

func RollRecordsAddNum(records [][]string) [][]string {
	var ret [][]string
	for index, tmpRecord := range records {
		var record []string
		record = append(record, strconv.Itoa(index+1))
		record = append(record, tmpRecord...)

		ret = append(ret, record)
	}
	return ret
}
