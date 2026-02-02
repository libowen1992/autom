package app

import (
	"encoding/csv"
	"os"
	"sobot/core"
)

type HistoryApp struct {
	core.CommonData
}

func NewHistory() *HistoryApp {
	return &HistoryApp{core.CommonData{
		Name:      "history",
		HumanName: "历史记录",
		Features:  nil,
	}}
}

func (ha *HistoryApp) Run() error {
	drf, err := os.OpenFile(DeployRecords, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer drf.Close()

	r := csv.NewReader(drf)

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	if len(records) > MaxRecords {
		records = records[len(records)-MaxRecords:]
	}

	headers := []string{"ID", "类型", "服务", "耗时", "节点", "更新包"}
	core.Render(headers, Reverse(records))
	return nil
}

func Reverse(s [][]string) [][]string {
	var d [][]string
	for i := len(s) - 1; i >= 0; i-- {
		d = append(d, s[i])
	}
	return d
}
