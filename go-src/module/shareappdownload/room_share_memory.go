package shareappdownload

import (
	"strings"
	"sync/atomic"

	"xr-game-server/constants/lang"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/user"
)

type liveRoomShareSnapshot struct {
	ShareTextZhCN string
	ShareTextZhTW string
	ShareTextEn   string
	ShareTextEs   string
	ShareTextPt   string
	ShareTextHi   string
	ShareTextId   string
	AfLink        string
}

var (
	liveRoomShareCache         atomic.Value // *liveRoomShareSnapshot
	emptyLiveRoomShareSnapshot = &liveRoomShareSnapshot{}
)

func reloadLiveRoomShareMemory() {
	liveRoomShareCache.Store(toLiveRoomShareSnapshot(cfgdao.LoadLiveRoomShareCfg()))
}

func getLiveRoomShareSnapshot() *liveRoomShareSnapshot {
	v := liveRoomShareCache.Load()
	if v == nil {
		return emptyLiveRoomShareSnapshot
	}
	snap, ok := v.(*liveRoomShareSnapshot)
	if !ok || snap == nil {
		return emptyLiveRoomShareSnapshot
	}
	return snap
}

func toLiveRoomShareSnapshot(row *entity.LiveRoomShareCfg) *liveRoomShareSnapshot {
	if row == nil {
		return emptyLiveRoomShareSnapshot
	}
	return &liveRoomShareSnapshot{
		ShareTextZhCN: strings.TrimSpace(row.ShareTextZhCN),
		ShareTextZhTW: strings.TrimSpace(row.ShareTextZhTW),
		ShareTextEn:   strings.TrimSpace(row.ShareTextEn),
		ShareTextEs:   strings.TrimSpace(row.ShareTextEs),
		ShareTextPt:   strings.TrimSpace(row.ShareTextPt),
		ShareTextHi:   strings.TrimSpace(row.ShareTextHi),
		ShareTextId:   strings.TrimSpace(row.ShareTextId),
		AfLink:        strings.TrimSpace(row.AfLink),
	}
}

func liveRoomShareTextForLang(snap *liveRoomShareSnapshot, clientLang lang.Lang) string {
	cms := liveRoomShareTextFromCMS(snap, clientLang)
	if cms != "" {
		return cms
	}
	if def := DefaultLiveRoomShareText(clientLang); def != "" {
		return def
	}
	return defaultLiveRoomShareTextEn
}

func liveRoomShareTextFromCMS(snap *liveRoomShareSnapshot, clientLang lang.Lang) string {
	if snap == nil {
		return ""
	}
	var text string
	switch clientLang {
	case lang.LangES:
		text = snap.ShareTextEs
	case lang.LangPT:
		text = snap.ShareTextPt
	case lang.LangHI:
		text = snap.ShareTextHi
	case lang.LangID:
		text = snap.ShareTextId
	default:
		text = snap.ShareTextEn
	}
	text = strings.TrimSpace(text)
	if text != "" {
		return text
	}
	text = strings.TrimSpace(snap.ShareTextEn)
	return text
}