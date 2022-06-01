package libs

import(
	"fmt"
	"time"
)

func Timer(start *time.Time) string {
	duration := time.Since(*start)
	if duration.Hours() >= 1 {
		return fmt.Sprintf("%dh%dm%ds", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
	} else if duration.Minutes() >= 1 {
		return fmt.Sprintf("%dm%ds", int(duration.Minutes())%60, int(duration.Seconds())%60)
	} else {
		return fmt.Sprintf("%.2fs", duration.Seconds())
	}
}