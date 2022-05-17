package libs

import "runtime"

func MaxCPU() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
			return maxProcs
	}
	return numCPU
}