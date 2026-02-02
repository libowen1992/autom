package convert

import (
	"fmt"
	"time"
)

const (
	day            = 24 * time.Hour
	dateTimeLayout = "20060102150405"
	dateLayout     = "2006-01-02"
)

func TimeToDateTimeFormat(t time.Time) string {
	return t.Format(dateTimeLayout)
}

func TimeToDateFormat(t time.Time) string {
	return t.Format(dateLayout)
}

func SecondHumanTime(t string) string {
	d, _ := time.ParseDuration(t)
	if d < day {
		return d.String()
	}
	days := d / day
	return fmt.Sprintf("%ddays", days)
}

func UnixTimeToDataTimeFormat(s, ns int64) string {
	time := time.Unix(s, ns)
	return TimeToDateTimeFormat(time)
}
