package guild

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/liverevenuesharecfg"
)

func genGuildId() uint64 {
	return snowflake.GetId()
}

// resolveGuildListVisibility 非管理员仅看可见性表中的工会;管理员/超管看全部
func resolveGuildListVisibility(ctx context.Context) (visibleGuildIds []uint64, filterByVisibility bool) {
	return cmsvis.ResolveGuildListVisibility(ctx)
}

// GetGuildList 获取直播工会列表
func GetGuildList(ctx context.Context, req *guilddto.GuildListReq) (res *httpserver.CMSQueryResp, err error) {
	visibleIds, filter := resolveGuildListVisibility(ctx)
	total, guilds := guilddao.GetGuildList(req, visibleIds, filter)
	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  guilds,
	}, nil
}

// GetGuildListForVisibility 可见性管理页拉取全部上架工会(不按可见性表过滤)
func GetGuildListForVisibility(_ context.Context, req *guilddto.GuildListForVisibilityReq) (res *httpserver.CMSQueryResp, err error) {
	listReq := &guilddto.GuildListReq{
		CMSQueryReq: req.CMSQueryReq,
		Name:        req.Name,
	}
	total, guilds := guilddao.GetGuildList(listReq, nil, false)
	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  guilds,
	}, nil
}

// GetOffShelfGuildList 获取已下架工会列表(垃圾库,直查DB)
func GetOffShelfGuildList(ctx context.Context, req *guilddto.OffShelfGuildListReq) (res *httpserver.CMSQueryResp, err error) {
	visibleIds, filter := resolveGuildListVisibility(ctx)
	total, guilds := guilddao.GetOffShelfGuildList(req, visibleIds, filter)
	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  guilds,
	}, nil
}

// CreateGuild 创建直播工会
func CreateGuild(ctx context.Context, req *guilddto.CreateGuildReq) (res *guilddto.CreateGuildRes, err error) {
	if existing := guilddao.GetGuildByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}
	if err = validateGuildLeader(req.LeaderId); err != nil {
		return nil, err
	}
	if req.GuildType != liveentity.LiveGuildTypeNormal && req.GuildType != liveentity.LiveGuildTypeCoinMerchant {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	sharePercent, err := resolveGuildSharePercent(req.GuildType, req.SharePercent, 0, false)
	if err != nil {
		return nil, err
	}

	guild := liveentity.NewLiveGuild(
		genGuildId(),
		req.Name,
		req.LeaderId,
		resolveLeaderName(req.LeaderId),
		req.Description,
	)
	guild.GuildType = req.GuildType
	guild.SharePercent = sharePercent
	creatorId := httpserver.GetAuthId(ctx)
	if creatorId > 0 {
		guild.CreatorId = creatorId
		guild.CreatorName = resolveLeaderName(creatorId)
	}
	if err = guilddao.CreateGuild(guild); err != nil {
		return nil, err
	}
	// 创建人默认可见
	if creatorId > 0 {
		_ = guilddao.AddGuildVisibility(guild.ID, creatorId)
	}

	return &guilddto.CreateGuildRes{ID: strconv.FormatUint(guild.ID, 10)}, nil
}

// UpdateGuild 更新直播工会
func UpdateGuild(ctx context.Context, req *guilddto.UpdateGuildReq) (res *guilddto.UpdateGuildRes, err error) {
	guild := guilddao.GetGuildById(req.ID)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	if existing := guilddao.GetGuildByName(req.Name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}
	if err = validateGuildLeader(req.LeaderId); err != nil {
		return nil, err
	}
	if req.GuildType != liveentity.LiveGuildTypeNormal && req.GuildType != liveentity.LiveGuildTypeCoinMerchant {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	sharePercent, err := resolveGuildSharePercent(req.GuildType, req.SharePercent, guild.SharePercent, true)
	if err != nil {
		return nil, err
	}

	guild.Name = req.Name
	guild.LeaderId = req.LeaderId
	guild.LeaderName = resolveLeaderName(req.LeaderId)
	guild.Description = req.Description
	guild.GuildType = req.GuildType
	guild.SharePercent = sharePercent
	if err = guilddao.UpdateGuild(guild); err != nil {
		return nil, err
	}

	return &guilddto.UpdateGuildRes{Success: true}, nil
}

// resolveGuildSharePercent 普通工会强制用全局工会分佣;币商工会才允许自定义单一分佣比例
func resolveGuildSharePercent(
	guildType uint8,
	shareIn *float64,
	existing float64,
	useExisting bool,
) (sharePercent float64, err error) {
	if guildType != liveentity.LiveGuildTypeCoinMerchant {
		return liverevenuesharecfg.ResolveGuildSharePercent(), nil
	}
	if useExisting {
		sharePercent = existing
	} else {
		sharePercent = liverevenuesharecfg.ResolveGuildSharePercent()
	}
	if shareIn != nil {
		if *shareIn < 0 || *shareIn > 100 {
			return 0, errercode.CreateCode(errercode.InvalidParam)
		}
		sharePercent = *shareIn
	}
	return sharePercent, nil
}

// DeleteGuild 下架直播工会(同步下架工会下全部主播间,并归档清零工会未结算收益/日有效次数)
func DeleteGuild(ctx context.Context, req *guilddto.DeleteGuildReq) (res *guilddto.DeleteGuildRes, err error) {
	if guilddao.GetGuildById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	liveroom.OffShelfGuildLiveRooms(ctx, req.ID)
	liveroomdao.ArchiveAndClearGuildUnsettledIncome(req.ID)
	liveroomdao.ClearRecentUnsettledDailyGuildLiveDuration(req.ID)
	if err = guilddao.DeleteGuild(req.ID); err != nil {
		return nil, err
	}
	return &guilddto.DeleteGuildRes{Success: true}, nil
}

// OnShelfGuild 上架直播工会(同步上架工会下全部主播间)
func OnShelfGuild(ctx context.Context, req *guilddto.OnShelfGuildReq) (res *guilddto.OnShelfGuildRes, err error) {
	guild := guilddao.GetGuildByIdFromDB(req.ID)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if guild.Status == liveentity.LiveGuildStatusOnShelf {
		return &guilddto.OnShelfGuildRes{Success: true}, nil
	}
	// 上架前校验名称是否与已上架工会冲突
	if existing := guilddao.GetGuildByName(guild.Name); existing != nil && existing.ID != guild.ID {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}
	if err = guilddao.OnShelfGuild(req.ID); err != nil {
		return nil, err
	}
	liveroom.OnShelfGuildLiveRooms(ctx, req.ID)
	return &guilddto.OnShelfGuildRes{Success: true}, nil
}

func resolveLeaderName(leaderId uint64) string {
	if leaderId == 0 {
		return ""
	}
	user := cmsuserdao.GetCMSUserById(leaderId)
	if user == nil {
		return ""
	}
	return user.Name
}

func validateGuildLeader(leaderId uint64) error {
	if leaderId == 0 {
		return nil
	}
	user := cmsuserdao.GetCMSUserById(leaderId)
	if user == nil {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if user.Status != 1 {
		return errercode.CreateCode(errercode.GuildLeaderDisabled)
	}
	return nil
}
