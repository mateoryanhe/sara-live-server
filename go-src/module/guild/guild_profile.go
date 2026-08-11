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

// GetMyGuildProfile 获取当前 CMS 用户作为会长关联的工会基础信息
func GetMyGuildProfile(ctx context.Context, _ *guilddto.GetMyGuildProfileReq) (*guilddto.GetMyGuildProfileRes, error) {
	guild, err := getGuildByCMSUser(ctx)
	if err != nil {
		return nil, err
	}
	return toMyGuildProfileRes(guild), nil
}

// UpdateMyGuildProfile 更新当前 CMS 用户作为会长关联的工会基础信息
func UpdateMyGuildProfile(ctx context.Context, req *guilddto.UpdateMyGuildProfileReq) (*guilddto.UpdateMyGuildProfileRes, error) {
	guild, err := getGuildByCMSUser(ctx)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if existing := guilddao.GetGuildByName(name); existing != nil && existing.ID != guild.ID {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}

	guild.Name = name
	guild.BankCard = strings.TrimSpace(req.BankCard)
	guild.Contact = strings.TrimSpace(req.Contact)
	guild.Description = strings.TrimSpace(req.Description)

	if err = guilddao.UpdateGuild(guild); err != nil {
		return nil, err
	}
	guilddao.RemoveGuildCache(guild.ID)

	return &guilddto.UpdateMyGuildProfileRes{Success: true}, nil
}

func getGuildByCMSUser(ctx context.Context) (*entity.LiveGuild, error) {
	cmsUserId := httpserver.GetAuthId(ctx)
	if cmsUserId == 0 {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	guild := guilddao.GetGuildByLeaderId(cmsUserId)
	if guild == nil || guild.ID == 0 {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	return guild, nil
}

func toMyGuildProfileRes(guild *entity.LiveGuild) *guilddto.GetMyGuildProfileRes {
	if guild == nil {
		return &guilddto.GetMyGuildProfileRes{}
	}
	updatedAt := ""
	if !guild.UpdatedAt.IsZero() {
		updatedAt = guild.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return &guilddto.GetMyGuildProfileRes{
		ID:          strconv.FormatUint(guild.ID, 10),
		Name:        guild.Name,
		BankCard:    guild.BankCard,
		Contact:     guild.Contact,
		Description: guild.Description,
		UpdatedAt:   updatedAt,
	}
}
