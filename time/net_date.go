package time

import (
	"fmt"
	go_time "time"
)

// ToNetLongDateString converts time to .NET ToLongDateString() format.
func ToNetLongDateString(time go_time.Time) string {
	return fmt.Sprintf("%s, %s %d, %d", time.Weekday(), time.Month(), time.Day(), time.Year())
}
