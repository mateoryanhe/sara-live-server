package guild

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// GetMyGuildProfile 获取当前 CMS 用户作为会长管理的全部工会基础信息
func GetMyGuildProfile(ctx context.Context, _ *guilddto.GetMyGuildProfileReq) (*guilddto.GetMyGuildProfileRes, error) {
	cmsUserId, err := getCMSUserId(ctx)
	if err != nil {
		return nil, err
	}
	rows := guilddao.ListGuildsByLeaderId(cmsUserId)
	list := make([]*guilddto.MyGuildProfileItem, 0, len(rows))
	for _, row := range rows {
		if item := toMyGuildProfileItem(row); item != nil {
			list = append(list, item)
		}
	}
	return &guilddto.GetMyGuildProfileRes{List: list}, nil
}

// UpdateMyGuildProfile 更新当前 CMS 用户作为会长管理的指定工会基础信息
func UpdateMyGuildProfile(ctx context.Context, req *guilddto.UpdateMyGuildProfileReq) (*guilddto.UpdateMyGuildProfileRes, error) {
	guild, err := getGuildOwnedByCMSUser(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if existing := guilddao.GetGuildByName(name); existing != nil && existing.ID != guild.ID {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}

	oldName := guild.Name
	guild.SetName(name)
	guild.SetBankCard(strings.TrimSpace(req.BankCard))
	guild.SetContact(strings.TrimSpace(req.Contact))
	guild.SetDescription(strings.TrimSpace(req.Description))
	guilddao.UpdateGuild(guild, oldName, guild.LeaderId)

	return &guilddto.UpdateMyGuildProfileRes{Success: true}, nil
}

func getCMSUserId(ctx context.Context) (uint64, error) {
	cmsUserId := httpserver.GetAuthId(ctx)
	if cmsUserId == 0 {
		return 0, errercode.CreateCode(errercode.NoPermission)
	}
	return cmsUserId, nil
}

func getGuildOwnedByCMSUser(ctx context.Context, guildId uint64) (*entity.LiveGuild, error) {
	if guildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	cmsUserId, err := getCMSUserId(ctx)
	if err != nil {
		return nil, err
	}
	guild := guilddao.GetGuildById(guildId)
	if guild == nil || guild.ID == 0 {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if guild.LeaderId != cmsUserId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	return guild, nil
}

func toMyGuildProfileItem(guild *entity.LiveGuild) *guilddto.MyGuildProfileItem {
	if guild == nil || guild.ID == 0 {
		return nil
	}
	updatedAt := ""
	if !guild.UpdatedAt.IsZero() {
		updatedAt = guild.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return &guilddto.MyGuildProfileItem{
		ID:          strconv.FormatUint(guild.ID, 10),
		Name:        guild.Name,
		BankCard:    guild.BankCard,
		Contact:     guild.Contact,
		Description: guild.Description,
		UpdatedAt:   updatedAt,
	}
}
