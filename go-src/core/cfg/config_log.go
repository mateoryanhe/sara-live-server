package cfg

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/core/xrjson"
)

func logConfigContent() {
	ctx := gctx.New()
	if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		filePath, err := adapter.GetFilePath()
		if err == nil && filePath != "" {
			content := gfile.GetContents(filePath)
			if content != "" {
				g.Log().Warningf(ctx, "config.yaml (%s):\n%s", filePath, content)
				return
			}
		}
	}
	data, err := g.Cfg().Data(ctx)
	if err != nil {
		g.Log().Error(ctx, "无法读取配置内容")
		return
	}
	g.Log().Warningf(ctx, "config:\n%s", xrjson.MustMarshalIndent(data))
}
