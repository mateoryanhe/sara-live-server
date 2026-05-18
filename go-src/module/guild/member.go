package guild

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func genMemberId() uint64 {
	return snowflake.GetId()
}

// ApplyJoin 申请加入工会
func ApplyJoin(ctx context.Context, req *guilddto.ApplyJoinGuildReq) (res *guilddto.ApplyJoinGuildRes, err error) {
	userId := httpserver.GetAuthId(ctx)

	// 工会必须存在
	g := guilddao.GetGuildByIdCached(req.GuildId)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	// 已存在 申请中 / 已加入 的记录则不允许重复申请
	exist := guilddao.FindMemberInGuild(req.GuildId, userId,
		entity.GuildMemberStatusPending, entity.GuildMemberStatusApproved)
	if exist != nil {
		return nil, errercode.CreateCode(errercode.GuildApplyExist)
	}

	// 构造新成员对象,所有字段写入通过 syndb 异步入库,不直接 Insert
	m := entity.NewLiveGuildMember(genMemberId(), req.GuildId, userId)
	// 加入缓存
	guilddao.AddMemberToCache(m)

	return &guilddto.ApplyJoinGuildRes{
		MemberId: strconv.FormatUint(m.ID, 10),
	}, nil
}

// ApproveMember 审批通过工会申请(会长操作);通过后写入 user_info.guild_id
func ApproveMember(ctx context.Context, req *guilddto.ApproveGuildMemberReq) (res *guilddto.ApproveGuildMemberRes, err error) {
	operatorId := httpserver.GetAuthId(ctx)

	m := guilddao.GetMember(req.MemberId)
	if m == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	g := guilddao.GetGuildByIdCached(m.GuildId)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if g.LeaderId != operatorId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}

	// 仅"申请中"状态可以被通过
	if m.Status != entity.GuildMemberStatusPending {
		return nil, errercode.CreateCode(errercode.GuildApplyExist)
	}

	// 1) 成员表状态置为已加入
	m.SetStatus(entity.GuildMemberStatusApproved)
	guilddao.FlushMemberCache(m)

	// 2) 将工会ID写入用户信息表(玩家成为主播,需要绑定工会)
	user := userinfodao.GetUserInfoByUserId(m.UserId)
	if user != nil {
		user.SetGuildId(m.GuildId)
	}

	return &guilddto.ApproveGuildMemberRes{Success: true}, nil
}

// RejectMember 拒绝工会申请(会长操作)
func RejectMember(ctx context.Context, req *guilddto.RejectGuildMemberReq) (res *guilddto.RejectGuildMemberRes, err error) {
	operatorId := httpserver.GetAuthId(ctx)

	m := guilddao.GetMember(req.MemberId)
	if m == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	// 权限校验: 仅会长可操作
	g := guilddao.GetGuildByIdCached(m.GuildId)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if g.LeaderId != operatorId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}

	// 仅"申请中"的记录可被拒绝
	if m.Status != entity.GuildMemberStatusPending {
		return nil, errercode.CreateCode(errercode.GuildApplyExist)
	}

	m.SetStatus(entity.GuildMemberStatusRejected)
	guilddao.FlushMemberCache(m)
	// 列表缓存中的对象是同一指针,状态已就地变更,无需重建

	return &guilddto.RejectGuildMemberRes{Success: true}, nil
}

// KickMember 剔除工会成员(会长操作)
func KickMember(ctx context.Context, req *guilddto.KickGuildMemberReq) (res *guilddto.KickGuildMemberRes, err error) {
	operatorId := httpserver.GetAuthId(ctx)

	m := guilddao.GetMember(req.MemberId)
	if m == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	g := guilddao.GetGuildByIdCached(m.GuildId)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if g.LeaderId != operatorId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	// 不可剔除自己
	if m.UserId == operatorId {
		return nil, errercode.CreateCode(errercode.GuildKickSelf)
	}
	// 仅"已加入"的成员可被剔除
	if m.Status != entity.GuildMemberStatusApproved {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	m.SetStatus(entity.GuildMemberStatusKicked)
	guilddao.FlushMemberCache(m)

	// 被剔除后清空该用户的所属工会
	if user := userinfodao.GetUserInfoByUserId(m.UserId); user != nil && user.GuildId == m.GuildId {
		user.SetGuildId(0)
	}

	return &guilddto.KickGuildMemberRes{Success: true}, nil
}

// GetMemberList 查询工会成员列表(仅返回已加入的成员)
func GetMemberList(ctx context.Context, req *guilddto.GuildMemberListReq) (res *guilddto.GuildMemberListRes, err error) {
	if g := guilddao.GetGuildByIdCached(req.GuildId); g == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	all := guilddao.GetMembersByGuild(req.GuildId)
	items := make([]*guilddto.GuildMemberItem, 0, len(all))
	for _, m := range all {
		if m.Status != entity.GuildMemberStatusApproved {
			continue
		}
		items = append(items, buildMemberItem(m))
	}
	return &guilddto.GuildMemberListRes{List: items}, nil
}

// GetPendingMemberList 获取待审批的入会申请列表(仅会长可调用)
func GetPendingMemberList(ctx context.Context, req *guilddto.PendingMemberListReq) (res *guilddto.PendingMemberListRes, err error) {
	operatorId := httpserver.GetAuthId(ctx)

	g := guilddao.GetGuildByIdCached(req.GuildId)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if g.LeaderId != operatorId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}

	all := guilddao.GetMembersByGuild(req.GuildId)
	items := make([]*guilddto.GuildMemberItem, 0)
	for _, m := range all {
		if m.Status != entity.GuildMemberStatusPending {
			continue
		}
		items = append(items, buildMemberItem(m))
	}
	return &guilddto.PendingMemberListRes{List: items}, nil
}

// buildMemberItem 统一构造成员返回项
func buildMemberItem(m *entity.LiveGuildMember) *guilddto.GuildMemberItem {
	return &guilddto.GuildMemberItem{
		MemberId:  strconv.FormatUint(m.ID, 10),
		GuildId:   strconv.FormatUint(m.GuildId, 10),
		UserId:    strconv.FormatUint(m.UserId, 10),
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Format(time.DateTime),
		UpdatedAt: m.UpdatedAt.Format(time.DateTime),
	}
}
