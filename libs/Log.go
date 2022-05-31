package libs

import (
	"os"
	"fmt"
)

func logDivider(log *os.File, message string, program string, params []string) {
	fmt.Fprintf(log, "\n\n\n==================== %s ====================\n", message)
	fmt.Fprintf(log, "Using %s with preset: %s\n", program, params)
	fmt.Fprintf(log, "\n\n\n")
}