package convert

import "strconv"

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}
