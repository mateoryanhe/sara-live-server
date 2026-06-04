package sensitiveword

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func Init() {
	ctx := gctx.New()
	if err := loadAll(ctx); err != nil {
		g.Log().Errorf(ctx, "sensitiveword init failed: %v", err)
	}
}
