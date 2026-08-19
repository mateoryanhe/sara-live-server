package cfg

import (
	"runtime/debug"
)

const (
	mib                 = 1024 * 1024
	defaultMemoryLimitM = 300
)

// ApplyMemoryLimit 设置 Go 堆内存软上限(MB),limitM<=0 时使用 defaultMemoryLimitM.
func ApplyMemoryLimit(limitM int) {
	if limitM <= 0 {
		limitM = defaultMemoryLimitM
	}
	limitBytes := int64(limitM) * mib
	debug.SetMemoryLimit(limitBytes)
}
