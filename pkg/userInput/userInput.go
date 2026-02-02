package userInput

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sobot/pkg/convert"
	"strings"
)

var (
	DefaultGuide = "请输入"
)

// UserInputString Read user input os.stdin return a string
func UserString(guide string) (input string, err error) {
	if len(guide) == 0 {
		guide = DefaultGuide
	}
	fmt.Printf("%s>: ", guide)
	reader := bufio.NewReader(os.Stdin)
	input, err = reader.ReadString('\n')
	if err != nil {
		return "", errors.WithMessagef(err, "%s", input)
	}
	input = strings.TrimSpace(input)
	return
}

// UserInputNum Read user input os.stdin return number
func UserNum(guide string) (number int, err error) {
	input, err := UserString(guide)
	if err != nil {
		return
	}
	number, err = convert.StrToInt(input)
	if err != nil {
		return 0, errors.WithMessagef(err, "请输入数字:%s", input)
	}
	return
}
