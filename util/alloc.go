package util

import "runtime"

// MemAlloc returns the current memory allocations of the Golang process for
// calculating the memory overhead a data structure
func MemAlloc() uint64 {
	var stats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&stats)
	return stats.Alloc
}
