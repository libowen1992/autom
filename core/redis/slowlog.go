package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sobot/core"
	"sobot/global"
	"strconv"
)

type SlowLogManager struct {
	core.CommonData
	db redis.Conn
}

func NewSlowLogManager() *SlowLogManager {
	return &SlowLogManager{
		CommonData: core.CommonData{
			Name:      "slowlog",
			HumanName: "慢日志",
		},
	}
}

func (sm *SlowLogManager) Run() (err error) {
	if err = sm.setup(); err != nil {
		return
	}

	defer sm.db.Close()
	result, err := sm.db.Do("SLOWLOG", "GET")
	if err != nil {
		return
	}
	slowLogs, err := redis.SlowLogs(result, err)
	if err != nil {
		return
	}

	headers := []string{"ID", "开始时间", "耗时(秒)", "命令", "客户端地址"}
	data := make([][]string, 0, len(slowLogs))
	for _, v := range slowLogs {
		idStr := strconv.FormatInt(v.ID, 10)
		startTime := v.Time.String()
		executionTime := fmt.Sprintf("%f", v.ExecutionTime.Seconds())
		args := fmt.Sprintf("%s", v.Args)
		data = append(data, []string{idStr, startTime, executionTime, args, v.ClientAddr})
	}
	core.Render(headers, data)
	return nil
}

func (sm *SlowLogManager) setup() error {
	password := redis.DialPassword(global.RedisSetting.Pass)
	addrStr := fmt.Sprintf("%s:%d", global.RedisSetting.IP, global.RedisSetting.Port)
	conn, err := redis.Dial("tcp", addrStr, password)
	sm.db = conn
	if err != nil {
		return err
	}
	return nil
}
