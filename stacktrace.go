package betterr

import "runtime"

func GetCallers(skip int) []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	return pcs[0:n]
}
