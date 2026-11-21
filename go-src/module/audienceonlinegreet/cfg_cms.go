package audienceonlinegreet

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/audienceonlinegreetdto"
	messageentity "xr-game-server/entity/message"
	"xr-game-server/errercode"
)

func GetCfg(_ context.Context, _ *audienceonlinegreetdto.GetAudienceOnlineAnchorGreetCfgReq) (*audienceonlinegreetdto.GetAudienceOnlineAnchorGreetCfgRes, error) {
	cfg := cfgdao.GetAudienceOnlineAnchorGreetCfgCached()
	if cfg == nil {
		return &audienceonlinegreetdto.GetAudienceOnlineAnchorGreetCfgRes{
			Cfg: defaultCfgItem(),
		}, nil
	}
	return &audienceonlinegreetdto.GetAudienceOnlineAnchorGreetCfgRes{Cfg: toCfgItem(cfg)}, nil
}

func SaveCfg(_ context.Context, req *audienceonlinegreetdto.SaveAudienceOnlineAnchorGreetCfgReq) (*audienceonlinegreetdto.SaveAudienceOnlineAnchorGreetCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	// CMS 不维护招呼文案，保存时清空库内列表，运行时走 default_greet_lines.go 内置 30 条/语种。
	anchorCount := normalizeAnchorCount(req.AnchorCount)
	intervalMin, intervalMax := normalizeGreetIntervalSec(req.IntervalMinSec, req.IntervalMaxSec)

	existing := cfgdao.GetAudienceOnlineAnchorGreetCfgCached()
	row := &messageentity.AudienceOnlineAnchorGreetCfg{
		Enabled:        req.Enabled,
		AnchorCount:    anchorCount,
		IntervalMinSec: intervalMin,
		IntervalMaxSec: intervalMax,
		ContentsEn:     messageentity.EncodeAudienceGreetContentListJSON(nil),
		ContentsEs:     messageentity.EncodeAudienceGreetContentListJSON(nil),
		ContentsPt:     messageentity.EncodeAudienceGreetContentListJSON(nil),
		ContentsHi:     messageentity.EncodeAudienceGreetContentListJSON(nil),
		ContentsId:     messageentity.EncodeAudienceGreetContentListJSON(nil),
		ContentEn:      "",
		Content:        "",
		ContentEs:      "",
		ContentPt:      "",
		ContentHi:      "",
		ContentId:      "",
		ContentZhCN:    "",
		ContentZhTW:    "",
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = req.ID
		row.CreatedAt = existing.CreatedAt
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveAudienceOnlineAnchorGreetCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadAudienceOnlineAnchorGreetCfgCache()
	return &audienceonlinegreetdto.SaveAudienceOnlineAnchorGreetCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func normalizeAnchorCount(v uint8) uint8 {
	if v == 0 {
		return messageentity.DefaultAudienceOnlineAnchorGreetAnchorCount
	}
	if v > messageentity.AudienceOnlineAnchorGreetAnchorCountMax {
		return messageentity.AudienceOnlineAnchorGreetAnchorCountMax
	}
	return v
}

func normalizeGreetIntervalSec(minSec, maxSec uint32) (uint32, uint32) {
	if minSec == 0 && maxSec == 0 {
		return messageentity.DefaultAudienceOnlineAnchorGreetIntervalMinSec,
			messageentity.DefaultAudienceOnlineAnchorGreetIntervalMaxSec
	}
	if minSec > messageentity.AudienceOnlineAnchorGreetIntervalSecMax {
		minSec = messageentity.AudienceOnlineAnchorGreetIntervalSecMax
	}
	if maxSec > messageentity.AudienceOnlineAnchorGreetIntervalSecMax {
		maxSec = messageentity.AudienceOnlineAnchorGreetIntervalSecMax
	}
	if maxSec < minSec {
		maxSec = minSec
	}
	return minSec, maxSec
}

func defaultCfgItem() *audienceonlinegreetdto.AudienceOnlineAnchorGreetCfgItem {
	return attachBuiltinMeta(&audienceonlinegreetdto.AudienceOnlineAnchorGreetCfgItem{
		Enabled:        false,
		AnchorCount:    messageentity.DefaultAudienceOnlineAnchorGreetAnchorCount,
		IntervalMinSec: messageentity.DefaultAudienceOnlineAnchorGreetIntervalMinSec,
		IntervalMaxSec: messageentity.DefaultAudienceOnlineAnchorGreetIntervalMaxSec,
	})
}

func toCfgItem(cfg *messageentity.AudienceOnlineAnchorGreetCfg) *audienceonlinegreetdto.AudienceOnlineAnchorGreetCfgItem {
	if cfg == nil {
		return defaultCfgItem()
	}
	intervalMin, intervalMax := cfg.EffectiveGreetIntervalSec()
	return attachBuiltinMeta(&audienceonlinegreetdto.AudienceOnlineAnchorGreetCfgItem{
		ID:             strconv.FormatUint(cfg.ID, 10),
		Enabled:        cfg.Enabled,
		AnchorCount:    cfg.EffectiveAnchorCount(),
		IntervalMinSec: intervalMin,
		IntervalMaxSec: intervalMax,
		CreatedAt:      formatTime(cfg.CreatedAt),
		UpdatedAt:      formatTime(cfg.UpdatedAt),
	})
}

func attachBuiltinMeta(item *audienceonlinegreetdto.AudienceOnlineAnchorGreetCfgItem) *audienceonlinegreetdto.AudienceOnlineAnchorGreetCfgItem {
	if item == nil {
		return nil
	}
	builtin := builtinGreetCfgContents()
	item.BuiltinContentsEn = cloneStrSlice(builtin.En)
	item.BuiltinContentsEs = cloneStrSlice(builtin.Es)
	item.BuiltinContentsPt = cloneStrSlice(builtin.Pt)
	item.BuiltinContentsHi = cloneStrSlice(builtin.Hi)
	item.BuiltinContentsId = cloneStrSlice(builtin.Id)
	item.UsingBuiltinContents = true
	return item
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}