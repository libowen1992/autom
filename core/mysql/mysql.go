package mysql

import (
	"sobot/core"
)

type MySQL struct {
	core.CommonData
}

func New() *MySQL {
	return &MySQL{core.CommonData{
		Name:      "mysql",
		HumanName: "mysql管理",
		Features:  []core.Project{NewTernimal(), NewProcesslist()},
	}}
}
