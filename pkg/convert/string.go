package convert

import (
	"strconv"
	"time"
)

func StringToBytes(s string) []byte {
	return []byte(s)
}

func StrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func StrDateToTime(str string) (t time.Time, err error) {
	return time.Parse(dateLayout, str)
}

func StrDateTimeToTime(str string) (t time.Time, err error) {
	return time.Parse(dateTimeLayout, str)
}

func StrToFloat64(str string) float64 {
	ret, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return ret
}
