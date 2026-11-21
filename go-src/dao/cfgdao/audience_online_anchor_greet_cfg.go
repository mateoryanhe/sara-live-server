package cfgdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	messageentity "xr-game-server/entity/message"
)

const audienceOnlineAnchorGreetCfgCacheKey = "audience_online_anchor_greet"

var audienceOnlineAnchorGreetCfgCacheMgr *cache.RowCache[*messageentity.AudienceOnlineAnchorGreetCfg]

func InitAudienceOnlineAnchorGreetCfgDao() {
	audienceOnlineAnchorGreetCfgCacheMgr = cache.NewRowCache[*messageentity.AudienceOnlineAnchorGreetCfg]()
}

func loadAudienceOnlineAnchorGreetCfgFromDB() *messageentity.AudienceOnlineAnchorGreetCfg {
	var row messageentity.AudienceOnlineAnchorGreetCfg
	if err := g.DB().Model(string(messageentity.TbAudienceOnlineAnchorGreetCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveAudienceOnlineAnchorGreetCfg(row *messageentity.AudienceOnlineAnchorGreetCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(messageentity.TbAudienceOnlineAnchorGreetCfg)).Save(row)
	return err
}

func ReloadAudienceOnlineAnchorGreetCfgCache() {
	if audienceOnlineAnchorGreetCfgCacheMgr == nil {
		return
	}
	audienceOnlineAnchorGreetCfgCacheMgr.PublishRow(gctx.New(), audienceOnlineAnchorGreetCfgCacheKey, loadAudienceOnlineAnchorGreetCfgFromDB())
}

func GetAudienceOnlineAnchorGreetCfgCached() *messageentity.AudienceOnlineAnchorGreetCfg {
	if audienceOnlineAnchorGreetCfgCacheMgr == nil {
		return loadAudienceOnlineAnchorGreetCfgFromDB()
	}
	return audienceOnlineAnchorGreetCfgCacheMgr.MustGetRow(gctx.New(), audienceOnlineAnchorGreetCfgCacheKey, func(ctx context.Context) (*messageentity.AudienceOnlineAnchorGreetCfg, error) {
		return loadAudienceOnlineAnchorGreetCfgFromDB(), nil
	})
}