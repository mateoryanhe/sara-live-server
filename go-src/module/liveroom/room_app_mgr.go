package liveroom

import (
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/userstatus"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/aliyunmoderation"
	"xr-game-server/module/upload"
)

func normalizeLiveRoomCategory(category uint8) uint8 {
	if category == liveentity.LiveRoomCategoryGame || category == liveentity.LiveRoomCategoryPrivate {
		return category
	}
	return liveentity.LiveRoomCategoryHot
}

func applyRoomPricing(cfg *liveentity.LiveRoomCfg, ticket, billing float64) {
	if cfg == nil {
		return
	}
	if cfg.Ticket != ticket {
		cfg.SetTicket(ticket)
	}
	if cfg.Billing != billing {
		cfg.SetBilling(billing)
	}
}

func applyPrivateInviteType(cfg *liveentity.LiveRoomCfg, privateInviteType uint8) {
	if cfg == nil {
		return
	}
	if privateInviteType == 0 {
		privateInviteType = liveentity.DefaultPrivateInviteType(cfg.Category)
	}
	if cfg.PrivateInviteType != privateInviteType {
		cfg.SetPrivateInviteType(privateInviteType)
	}
}

func applyLiveRoomGameRecommends(roomID uint64, category uint8, gameCodes []string) error {
	if category != liveentity.LiveRoomCategoryGame {
		return nil
	}
	return SyncLiveRoomGameRecommendList(roomID, gameCodes)
}

// CreateRoom 创建直播间
// 业务规则:
//  1. 调用者必须已是主播(UserInfo.UserType 为普通主播或机器人主播)
//  2. 同一主播只能拥有一个直播间(再次调用直接返回已有信息)
func CreateRoom(ctx context.Context, req *liveroomdto.CreateLiveRoomReq) (res *liveroomdto.CreateLiveRoomRes, err error) {
	anchorId := httpserver.GetAuthId(ctx)
	logCreateRoomAppUpload(ctx, anchorId, req)

	user := userinfodao.GetUserInfoByUserId(anchorId)
	if user == nil || !user.IsAnchor() {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}
	if err := aliyunmoderation.RequireTextCompliant(aliyunmoderation.SceneComment, req.Title, req.Notice); err != nil {
		return nil, err
	}
	//if err := validateLiveRoomTag(req.TagId); err != nil {
	//	return nil, err
	//}

	coverName := ""
	if req.Cover != nil && req.Cover.Size > 0 {
		coverName, err = upload.UploadImageForApp(ctx, req.Cover)
		if err != nil {
			return nil, err
		}
	}

	category := normalizeLiveRoomCategory(req.Category)

	// 同一主播仅允许一个直播间(roomId == anchorId);CMS预创建的空直播间允许App完善资料
	if existing := liveroomdao.GetRoomById(anchorId); existing != nil {
		cfg := liveroomdao.GetLiveRoomCfg(existing.ID)
		if req.Title != "" && existing.Title != req.Title {
			existing.SetTitle(req.Title)
		}
		if coverName != "" && existing.Cover != coverName {
			existing.SetCover(coverName)
		}
		if req.Notice != "" && existing.Notice != req.Notice {
			existing.SetNotice(req.Notice)
		}
		if cfg != nil {
			if req.Category > 0 && cfg.Category != category {
				cfg.SetCategory(category)
			}
			if cfg.TagId != req.TagId {
				cfg.SetTagId(req.TagId)
			}
			applyRoomPricing(cfg, req.Ticket, req.Billing)
			applyPrivateInviteType(cfg, req.PrivateInviteType)
		}
		if err := applyLiveRoomGameRecommends(existing.ID, category, req.GameCodes); err != nil {
			return nil, err
		}
		return &liveroomdto.CreateLiveRoomRes{
			RoomId:  strconv.FormatUint(existing.ID, 10),
			GuildId: strconv.FormatUint(existing.GuildId, 10),
		}, nil
	}

	// 通过 syndb 异步入库,不直接 INSERT;LiveRoom.ID 复用主播用户ID
	room := liveentity.NewLiveRoom(
		anchorId,
		liveroomdao.GetAnchorGuildId(anchorId),
		coverName,
		req.Title,
		req.Notice,
	)
	liveroomdao.AddRoomToCache(room)
	cfg := liveroomdao.GetLiveRoomCfg(room.ID)
	if cfg != nil {
		cfg.SetCategory(category)
		cfg.SetTagId(req.TagId)
		cfg.SetTicket(req.Ticket)
		cfg.SetBilling(req.Billing)
		applyPrivateInviteType(cfg, req.PrivateInviteType)
	}
	if err := applyLiveRoomGameRecommends(room.ID, category, req.GameCodes); err != nil {
		return nil, err
	}

	return &liveroomdto.CreateLiveRoomRes{
		RoomId:  strconv.FormatUint(room.ID, 10),
		GuildId: strconv.FormatUint(room.GuildId, 10),
	}, nil
}

func logCreateRoomAppUpload(ctx context.Context, anchorId uint64, req *liveroomdto.CreateLiveRoomReq) {
	if req == nil {
		return
	}
	hasCover := req.Cover != nil && req.Cover.Size > 0
	coverFilename := ""
	var coverSize int64
	if hasCover {
		coverFilename = req.Cover.Filename
		coverSize = req.Cover.Size
	}
	g.Log().Infof(ctx,
		"CreateRoom app upload anchorId=%d title=%q notice=%q category=%d tagId=%d gameCodes=%v ticket=%.4f billing=%.4f privateInviteType=%d hasCover=%v coverFilename=%q coverSize=%d",
		anchorId, req.Title, req.Notice, req.Category, req.TagId, req.GameCodes, req.Ticket, req.Billing, req.PrivateInviteType,
		hasCover, coverFilename, coverSize,
	)
}

// loadOwnRoom 获取调用者(主播)自己的直播间;不存在则返回 LiveRoomNotExist
func loadOwnRoom(ctx context.Context) (*liveentity.LiveRoom, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	return room, nil
}

// UpdateCover 修改封面
func UpdateCover(ctx context.Context, req *liveroomdto.UpdateCoverReq) (*liveroomdto.UpdateCoverRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if room.Cover != req.Cover {
		room.SetCover(req.Cover)
	}
	return &liveroomdto.UpdateCoverRes{Success: true}, nil
}

// UpdateNotice 修改公告
func UpdateNotice(ctx context.Context, req *liveroomdto.UpdateNoticeReq) (*liveroomdto.UpdateNoticeRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if err := aliyunmoderation.RequireTextCompliant(aliyunmoderation.SceneComment, req.Notice); err != nil {
		return nil, err
	}
	if room.Notice != req.Notice {
		room.SetNotice(req.Notice)
		//liveroomdao.FlushRoomCache(room)
	}
	return &liveroomdto.UpdateNoticeRes{Success: true}, nil
}

func calcAge(birthday *time.Time) int {
	if birthday == nil || birthday.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	anniversary := time.Date(now.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(anniversary) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}

// ViewerCanUseLiveRoomCall 累计充值满 10 USD 可使用直播间 1v1 通话
func ViewerCanUseLiveRoomCall(userId uint64) bool {
	if userId == 0 {
		return false
	}
	stat := userinfodao.GetUserCumulativeStatByUserId(userId)
	return stat != nil && stat.TotalRecharge >= liveentity.LiveRoomCallMinTotalRechargeUSD
}

// CanInitiateLiveRoomCall 当前用户是否可向该直播间发起 1v1 通话(与 AllowCallIcon 规则一致)
func CanInitiateLiveRoomCall(room *liveentity.LiveRoom, cfg *liveentity.LiveRoomCfg, userId uint64) bool {
	return allowShowCallIcon(room, cfg, userId)
}

func allowShowCallIcon(room *liveentity.LiveRoom, cfg *liveentity.LiveRoomCfg, userId uint64) bool {
	if room == nil || cfg == nil || cfg.Category != liveentity.LiveRoomCategoryHot {
		return false
	}
	if userId == 0 || userId == room.ID {
		return false
	}

	inviteType := cfg.PrivateInviteType
	if inviteType == 0 {
		inviteType = liveentity.DefaultPrivateInviteType(cfg.Category)
	}
	switch inviteType {
	case liveentity.LiveRoomPrivateInviteAll:
		return true
	case liveentity.LiveRoomPrivateInviteVip:
		return ViewerCanUseLiveRoomCall(userId)
	case liveentity.LiveRoomPrivateInviteReject:
		return false
	default:
		return true
	}
}

// GetRoom 查询直播间(公开接口,任意登录用户可调用)
func GetRoom(ctx context.Context, req *liveroomdto.GetLiveRoomReq) (*liveroomdto.GetLiveRoomRes, error) {
	userId := httpserver.GetAuthId(ctx)
	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	cfg := liveroomdao.GetLiveRoomCfg(room.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	status := userstatus.LiveRoomStatusClosed
	if room.LiveRecordId > 0 {
		status = userstatus.LiveRoomStatusLive
	}

	res := &liveroomdto.GetLiveRoomRes{
		RoomId:            strconv.FormatUint(room.ID, 10),
		GuildId:           strconv.FormatUint(room.GuildId, 10),
		Title:             room.Title,
		Cover:             upload.GetUrlByName(room.Cover),
		Notice:            room.Notice,
		Status:            status,
		Category:          cfg.Category,
		TagId:             strconv.FormatUint(cfg.TagId, 10),
		TagName:           getRoomTagName(cfg.TagId),
		Ticket:            cfg.Ticket,
		Billing:           cfg.Billing,
		PrivateInviteType: cfg.PrivateInviteType,
		AllowCallIcon:     allowShowCallIcon(room, cfg, userId),
		CreateAt:          room.CreatedAt.Unix(),
		OnlineCount:       countAudienceInRoom(room.ID),
	}
	res.IsBotAnchor, res.CloudPlayerVideo = resolveBotAnchorRoomInfo(room.ID, cfg.CloudPlayerVideo)
	res.IsTest = cfg.IsTest
	if u := userinfodao.GetUserInfoByUserId(room.ID); u != nil {
		res.UserType = u.UserType
	}
	//判断一下房间类型
	if cfg.Category == liveentity.LiveRoomCategoryPrivate {
		clearFreeTime(userId, room.ID)
		//私密房免费时长
		pay := liveroomdao.GetLiveRoomBillingPay(userId, req.RoomId)

		res.FreeTime = pay.FreeTime
		res.TicketTime = pay.GetTicketTime()
		res.HasTicket = res.TicketTime > 0
	}

	if userId == room.ID {
		if income := liveroomdao.GetLiveRoomIncomeTotalFromCache(room.ID); income != nil {
			res.TotalVideoCallIncome = income.TotalVideoCallIncome
			res.TotalVideoCallTicketIncome = income.TotalVideoCallTicketIncome
			res.TotalVideoCallBillingIncome = income.TotalVideoCallBillingIncome
		}
	}

	return res, nil
}

// EnsureAnchorRoom 确保主播拥有直播间记录(CMS设为主播时预创建,App端后续可完善资料)
// guildId>0 时同步写入 live_room.guild_id(工会归属唯一来源)
func EnsureAnchorRoom(anchorId, guildId uint64) *liveentity.LiveRoom {
	if room := liveroomdao.GetRoomByAnchor(anchorId); room != nil {
		if guildId > 0 && room.GuildId != guildId {
			room.SetGuildId(guildId)
		}
		return room
	}
	// 已下架的直播间不在缓存中,直查 DB,避免重复创建
	if room := liveroomdao.GetRoomFromDB(anchorId); room != nil {
		if guildId > 0 && room.GuildId != guildId {
			room.SetGuildId(guildId)
		}
		if room.Status == liveentity.LiveRoomStatusOnShelf {
			liveroomdao.AddRoomToCache(room)
		}
		return room
	}
	room := liveentity.NewLiveRoom(anchorId, guildId, "", "", "")
	liveroomdao.AddRoomToCache(room)
	return room
}
