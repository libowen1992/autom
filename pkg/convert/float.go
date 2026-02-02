package convert

import (
	"strconv"
)

func Float64ToPercentString(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}

func Float64ToPercentFloat64(num float64) float64 {
	return StrToFloat64(Float64ToPercentString(num))
}
