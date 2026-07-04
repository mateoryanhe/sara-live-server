package cfg

import (
	"runtime/debug"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
	g.Log().Warningf(gctx.New(), "Go memoryLimit=%dM (%d bytes)", limitM, limitBytes)
}
