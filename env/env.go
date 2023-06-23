package env

import "runtime"

func GetCpuCores() int {
	return runtime.NumCPU()
}
