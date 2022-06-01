package libs

import(
	"fmt"
	"time"
)

func Timer(start *time.Time) string {
	duration := time.Since(*start)
	if duration < 1000 {
		return fmt.Sprintf("%d ms", duration)
	}
	if duration < 1000*60 {
		return fmt.Sprintf("%.2f s", float64(duration)/1000)
	}
	if duration < 1000*60*60 {
		return fmt.Sprintf("%.2f min", float64(duration)/1000/60)
	}
	return fmt.Sprintf("%.2f h", float64(duration)/1000/60/60)
}