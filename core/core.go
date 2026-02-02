package core

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"os"
	"sobot/global"
	"sobot/pkg/userInput"
	"strconv"
	"strings"
)

const (
	Guide = "请输入数字选择,b返回,q退出"
)

var (
	ErrContinue = errors.New("输入有误")
	ErrBack     = errors.New("back return")
)

type Project interface {
	FullName() string
	Run() error
	Register(p Project) []Project
}

type CommonData struct {
	Name      string
	HumanName string
	Features  []Project
}

func (c *CommonData) FullName() string {
	return fmt.Sprintf("%s-%s", strings.Title(c.Name), c.HumanName)
}

func (c *CommonData) Register(p Project) []Project {
	c.Features = append(c.Features, p)
	return c.Features
}

func (c *CommonData) Render() {
	headers := []string{"序号", "项目"}
	var data [][]string
	for k, v := range c.Features {
		data = append(data, []string{strconv.FormatInt(int64(k+1), 10), v.FullName()})
	}
	Render(headers, data)
}

func (c CommonData) Run() error {
	for {
		c.Render()
		choice, err := userInput.UserString(Guide)
		if err != nil {
			continue
		}
		number, err := ChoiceJudge(choice)
		if err != nil {
			if errors.Is(err, ErrBack) {
				return nil
			} else {
				continue
			}
		}

		if number > len(c.Features) {
			continue
		}
		op := c.Features[number-1]
		if err := op.Run(); err != nil {
			global.Logger.Warn(err)
			continue
		}
	}
}

func Render(headers []string, data [][]string) {
	fmt.Println("粘贴注意~windows: Ctrl + Shift + v, Mac: Command + v")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetRowLine(true)
	var colors []tablewriter.Colors
	for i := 0; i < len(headers); i++ {
		colors = append(colors, tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor})
	}
	table.AppendBulk(data)
	table.Render()
}

func ChoiceJudge(choice string) (num int, err error) {
	if choice == "q" || choice == "quit" {
		os.Exit(1)
	} else if choice == "b" {
		return 0, ErrBack
	}
	num, err = strconv.Atoi(choice)
	if err != nil {
		return 0, errors.WithMessage(err, ErrContinue.Error())
	}

	return num, nil
}
