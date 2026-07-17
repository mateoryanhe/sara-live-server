package cfg

import (
	"runtime/debug"
)

const (
	mib                 = 1024 * 1024
	defaultMemoryLimitM = 400
)

func applyMemoryLimit(limitM int) {
	if limitM <= 0 {
		limitM = defaultMemoryLimitM
	}
	limitBytes := int64(limitM) * mib
	debug.SetMemoryLimit(limitBytes)
}
