package redis

import (
	"sobot/core"
)

type Redis struct {
	core.CommonData
}

func New() *Redis {
	return &Redis{core.CommonData{
		Name:      "redis",
		HumanName: "redis管理",
		Features:  []core.Project{NewRedisTernimal(), NewSlowLogManager()},
	}}
}
