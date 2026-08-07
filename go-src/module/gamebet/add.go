package gamebet

import (
	"strings"

	"xr-game-server/constants/gameplatform"
	"xr-game-server/core/event"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

// AddGameBetLog 新增游戏下注记录(quick 缓冲写库,并插入列表缓存头部).
func AddGameBetLog(userId uint64, gameCode, nameEn, cover string, platformType gameplatform.Platform, orderId string, amount float64) *entity.GameBetLog {
	if userId == 0 || !gameplatform.IsValid(platformType) {
		return nil
	}
	liveRoomId, liveRecordId := gamebetdao.ResolveAudienceLiveContext(userId)
	row := entity.NewGameBetLog(
		userId,
		strings.TrimSpace(gameCode),
		strings.TrimSpace(nameEn),
		strings.TrimSpace(cover),
		platformType,
		strings.TrimSpace(orderId),
		amount,
		liveRoomId,
		liveRecordId,
	)
	if amount > 0 {
		gamebetdao.PrependGameBetToAppListCache(userId, row)
		if liveRecordId > 0 {
			event.Pub(gameevent.GameBetCreatedEvent, gameevent.NewGameBetCreatedEventData(userId, liveRecordId, amount))
		}
	}
	return row
}
