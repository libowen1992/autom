package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sobot/core"
	"sobot/global"
	"strconv"
	"time"
)

type Processlist struct {
	core.CommonData
}

func NewProcesslist() *Processlist {
	return &Processlist{core.CommonData{
		Name:      "processlist",
		HumanName: "会话列表",
	}}
}

type Process struct {
	Id      int64         `db:"ID"`
	User    string        `db:"USER"`
	Host    string        `db:"HOST"`
	DB      string        `db:"DB"`
	Command string        `db:"COMMAND"`
	Time    time.Duration `db:"TIME"`
	State   string        `db:"STATE"`
	Info    string        `db:"INFO"`
}

func (p *Processlist) Run() error {
	headers := []string{"序号",
		"用户",
		"客户端",
		"数据库",
		"crud",
		"时间(秒)",
		"状态",
		"命令",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?charset=utf8&parseTime=True&loc=Local",
		global.MySQLSetting.User, global.MySQLSetting.Pass, global.MySQLSetting.IP, global.MySQLSetting.Port)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	var pl []*Process
	sqlStr := `select ID,
      					USER,
      					IFNULL(HOST, "") as HOST,
      					DB,
      					COMMAND,
      					TIME,
      					IFNULL(STATE, "") AS STATE,
      					IFNULL(INFO, "") AS INFO from PROCESSLIST
				WHERE DB NOT IN ("information_schema") AND COMMAND <> "Sleep"`
	if err = db.Select(&pl, sqlStr); err != nil {
		return err
	}

	var data [][]string
	for _, v := range pl {
		data = append(data, []string{strconv.FormatInt(v.Id, 10),
			v.User,
			v.Host,
			v.DB,
			v.Command,
			fmt.Sprintf("%f", v.Time.Seconds()),
			v.State,
			v.Info,
		})
	}

	core.Render(headers, data)

	return nil
}
