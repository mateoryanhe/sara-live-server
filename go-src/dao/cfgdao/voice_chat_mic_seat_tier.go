package cfgdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/voicechatmicseattierdto"
	liveentity "xr-game-server/entity/live"
)

func GetVoiceChatMicSeatTierById(id uint64) *liveentity.VoiceChatMicSeatTier {
	var row liveentity.VoiceChatMicSeatTier
	err := g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier)).Where("id = ?", id).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func GetVoiceChatMicSeatTierBySeatCount(seatCount int) *liveentity.VoiceChatMicSeatTier {
	var row liveentity.VoiceChatMicSeatTier
	err := g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier)).Where("seat_count = ?", seatCount).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func CreateVoiceChatMicSeatTier(row *liveentity.VoiceChatMicSeatTier) error {
	_, err := g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier)).Save(row)
	return err
}

func DeleteVoiceChatMicSeatTier(id uint64) error {
	_, err := g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier)).WherePri(id).Delete()
	return err
}

func GetOnShelfVoiceChatMicSeatTiers() []*liveentity.VoiceChatMicSeatTier {
	ret := make([]*liveentity.VoiceChatMicSeatTier, 0)
	_ = g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier)).
		Where("status = ?", liveentity.VoiceChatMicSeatTierStatusOnShelf).
		Order("sort desc, id desc").
		Scan(&ret)
	return ret
}

func GetAllVoiceChatMicSeatTiersForSync() []*liveentity.VoiceChatMicSeatTier {
	ret := make([]*liveentity.VoiceChatMicSeatTier, 0)
	_ = g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier)).
		Order("sort desc, seat_count asc, id asc").
		Scan(&ret)
	return ret
}

func GetVoiceChatMicSeatTierList(req *voicechatmicseattierdto.VoiceChatMicSeatTierListReq) (int, []*voicechatmicseattierdto.VoiceChatMicSeatTierListRes) {
	sql := `select id, seat_count, sort, status, created_at, updated_at
		from voice_chat_mic_seat_tiers where 1=1`
	param := make([]any, 0)
	if req.SeatCount > 0 {
		sql += ` and seat_count = ?`
		param = append(param, req.SeatCount)
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = 0`
	case 2:
		sql += ` and status = 1`
	}
	sql += ` order by sort desc, seat_count asc, id desc`
	ctx := gctx.New()
	total, _ := g.DB().GetCount(ctx, str.GetCountSQL(sql), param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())

	type row struct {
		ID        uint64 `json:"id"`
		SeatCount int    `json:"seat_count"`
		Sort      int    `json:"sort"`
		Status    uint8  `json:"status"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	rows := make([]*row, 0)
	_ = g.DB().GetScan(ctx, &rows, sql, param)
	list := make([]*voicechatmicseattierdto.VoiceChatMicSeatTierListRes, 0, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		list = append(list, &voicechatmicseattierdto.VoiceChatMicSeatTierListRes{
			ID:        strconv.FormatUint(r.ID, 10),
			SeatCount: r.SeatCount,
			Sort:      r.Sort,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	return total, list
}
