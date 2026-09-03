package guild

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetGuildTransferInfo(_ context.Context, req *guilddto.GetGuildTransferInfoReq) (*guilddto.GetGuildTransferInfoRes, error) {
	if req == nil || req.GuildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guild := guilddao.GetGuildByIdFromDB(req.GuildId)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	row := guilddao.GetGuildTransferInfo(req.GuildId)
	if row == nil {
		return &guilddto.GetGuildTransferInfoRes{
			Info: &guilddto.GuildTransferInfoItem{
				GuildId: strconv.FormatUint(req.GuildId, 10),
			},
		}, nil
	}
	return &guilddto.GetGuildTransferInfoRes{
		Info: toGuildTransferInfoItem(row),
	}, nil
}

func SaveGuildTransferInfo(_ context.Context, req *guilddto.SaveGuildTransferInfoReq) (*guilddto.SaveGuildTransferInfoRes, error) {
	if req == nil || req.GuildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guild := guilddao.GetGuildByIdFromDB(req.GuildId)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	row := liveentity.NewLiveGuildTransferInfo(req.GuildId)
	if existing := guilddao.GetGuildTransferInfo(req.GuildId); existing != nil {
		row.CreatedAt = existing.CreatedAt
	}
	row.Currency = currency
	row.PayeeName = strings.TrimSpace(req.PayeeName)
	row.BankName = strings.TrimSpace(req.BankName)
	row.AccountNo = strings.TrimSpace(req.AccountNo)
	row.BankCode = strings.TrimSpace(req.BankCode)
	row.Remark = strings.TrimSpace(req.Remark)
	row.UpdatedAt = time.Now()

	if err := guilddao.SaveGuildTransferInfo(row); err != nil {
		return nil, err
	}
	return &guilddto.SaveGuildTransferInfoRes{Success: true}, nil
}

func toGuildTransferInfoItem(row *liveentity.LiveGuildTransferInfo) *guilddto.GuildTransferInfoItem {
	if row == nil {
		return nil
	}
	updatedAt := ""
	if !row.UpdatedAt.IsZero() {
		updatedAt = row.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return &guilddto.GuildTransferInfoItem{
		GuildId:   strconv.FormatUint(row.ID, 10),
		Currency:  row.Currency,
		PayeeName: row.PayeeName,
		BankName:  row.BankName,
		AccountNo: row.AccountNo,
		BankCode:  row.BankCode,
		Remark:    row.Remark,
		UpdatedAt: updatedAt,
	}
}
