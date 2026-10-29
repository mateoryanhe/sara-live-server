package liveroom

import (
	"xr-game-server/dao/liveroomdao"
	liveentity "xr-game-server/entity/live"
)

func voiceChatMicSeatCountForRoom(roomId uint64) int {
	return liveroomdao.VoiceChatMicSeatCountForRoom(roomId)
}

func voiceChatMicSeatCountForCfg(cfg *liveentity.LiveRoomCfg) int {
	if cfg == nil {
		return int(liveentity.DefaultVoiceChatMicSeatCount)
	}
	return cfg.EffectiveVoiceChatMicSeatCount()
}

func applyVoiceChatMicSeatCountToCfg(cfg *liveentity.LiveRoomCfg, category uint8, requested uint8) {
	if cfg == nil || normalizeLiveRoomCategory(category) != liveentity.LiveRoomCategoryVoiceChat {
		return
	}
	if requested > 0 {
		cfg.SetVoiceChatMicSeatCount(requested)
	}
}

func voiceChatMicSeatCountForAPI(cfg *liveentity.LiveRoomCfg) uint8 {
	if cfg == nil || normalizeLiveRoomCategory(cfg.Category) != liveentity.LiveRoomCategoryVoiceChat {
		return 0
	}
	return uint8(voiceChatMicSeatCountForCfg(cfg))
}
