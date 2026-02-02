package host

import "sobot/core"

type Host struct {
	core.CommonData
}

func New() *Host {
	return &Host{core.CommonData{
		Name:      "host",
		HumanName: "主机管理",
		Features: []core.Project{
			NewHostTerminal(),
		},
	}}
}
