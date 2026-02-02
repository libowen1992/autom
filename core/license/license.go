package license

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sobot/core"
	"sobot/core/app"
)

type License struct {
	core.CommonData
}

func New() *License {
	return &License{core.CommonData{
		Name:      "license",
		HumanName: "获取序列号",
	}}
}

func (l *License) Run() error {
	headers := []string{"主机", "序列号"}
	var data [][]string
	basicLogin, _ := app.GetAppByName("basic-login")
	for _, node := range basicLogin.Nodes {
		apiUrl := fmt.Sprintf("http://%s:8003/basic-login/getSerialNo/4?prefix=%s", node, node)
		resp, err := http.Get(apiUrl)
		if err != nil {
			data = append(data, []string{node, err.Error()})
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			data = append(data, []string{node, fmt.Sprintf("%s 返回值非200，请检查接口", apiUrl)})
		}
		resData := ResData{}
		byteData, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(byteData, &resData); err != nil {
			data = append(data, []string{node, err.Error()})
		}
		data = append(data, []string{node, resData.Item})
	}
	core.Render(headers, data)
	return nil
}

type ResData struct {
	Item    string `json:"item"`
	RetCode string `json:"retCode"`
}
