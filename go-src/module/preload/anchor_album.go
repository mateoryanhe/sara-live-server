package preload

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/anchorphotodao"
	"xr-game-server/dao/anchorvideodao"
)

func preloadAnchorAlbum() {
	ctx := gctx.New()
	videoAnchors, videoRows := anchorvideodao.PreloadAnchorVideoToCache()
	photoAnchors, photoRows := anchorphotodao.PreloadAnchorPhotoToCache()
	if videoAnchors == 0 && photoAnchors == 0 {
		g.Log().Info(ctx, "preload anchor album skipped: no active photos/videos")
		return
	}
	g.Log().Infof(ctx, "preload anchor album done, videos anchors=%d rows=%d, photos anchors=%d rows=%d",
		videoAnchors, videoRows, photoAnchors, photoRows)
}