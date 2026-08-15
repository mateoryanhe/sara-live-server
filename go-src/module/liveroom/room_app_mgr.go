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
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/aliyunmoderation"
	"xr-game-server/module/upload"
)

func normalizeLiveRoomCategory(category uint8) uint8 {
	if category == entity.LiveRoomCategoryGame || category == entity.LiveRoomCategoryPrivate {
		return category
	}
	return entity.LiveRoomCategoryHot
}

func applyRoomPricing(room *entity.LiveRoom, ticket, billing float64) {
	if room.Ticket != ticket {
		room.SetTicket(ticket)
	}
	if room.Billing != billing {
		room.SetBilling(billing)
	}
}

func applyPrivateInviteType(room *entity.LiveRoom, privateInviteType uint8) {
	if privateInviteType == 0 {
		privateInviteType = entity.DefaultPrivateInviteType(room.Category)
	}
	if room.PrivateInviteType != privateInviteType {
		room.SetPrivateInviteType(privateInviteType)
	}
}

func applyLiveRoomGameRecommends(roomID uint64, category uint8, gameCodes []string) error {
	if category != entity.LiveRoomCategoryGame {
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
		if req.Title != "" && existing.Title != req.Title {
			existing.SetTitle(req.Title)
		}
		if coverName != "" && existing.Cover != coverName {
			existing.SetCover(coverName)
		}
		if req.Notice != "" && existing.Notice != req.Notice {
			existing.SetNotice(req.Notice)
		}
		if req.Category > 0 && existing.Category != category {
			existing.SetCategory(category)
		}
		if existing.TagId != req.TagId {
			existing.SetTagId(req.TagId)
		}
		applyRoomPricing(existing, req.Ticket, req.Billing)
		applyPrivateInviteType(existing, req.PrivateInviteType)
		if err := applyLiveRoomGameRecommends(existing.ID, category, req.GameCodes); err != nil {
			return nil, err
		}
		markLiveRoomCreated(user)
		return &liveroomdto.CreateLiveRoomRes{
			RoomId:  strconv.FormatUint(existing.ID, 10),
			GuildId: strconv.FormatUint(existing.GuildId, 10),
		}, nil
	}

	// 通过 syndb 异步入库,不直接 INSERT;LiveRoom.ID 复用主播用户ID
	room := entity.NewLiveRoom(
		anchorId,
		user.GuildId,
		coverName,
		req.Title,
		req.Notice,
	)
	room.SetCategory(category)
	room.SetTagId(req.TagId)
	room.SetTicket(req.Ticket)
	room.SetBilling(req.Billing)
	applyPrivateInviteType(room, req.PrivateInviteType)

	liveroomdao.AddRoomToCache(room)
	if err := applyLiveRoomGameRecommends(room.ID, category, req.GameCodes); err != nil {
		return nil, err
	}
	markLiveRoomCreated(user)

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
func loadOwnRoom(ctx context.Context) (*entity.LiveRoom, error) {
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

// markLiveRoomCreated App端创建/完善直播间后,标记 user_infos.has_live_room = true
func markLiveRoomCreated(user *entity.UserInfo) {
	if user == nil || user.HasLiveRoom {
		return
	}
	user.SetHasLiveRoom(true)
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

func allowShowCallIcon(room *entity.LiveRoom, userId uint64) bool {
	if room == nil || room.Category != entity.LiveRoomCategoryHot {
		return false
	}
	if userId == 0 || userId == room.ID {
		return false
	}

	inviteType := room.PrivateInviteType
	if inviteType == 0 {
		inviteType = entity.DefaultPrivateInviteType(room.Category)
	}
	switch inviteType {
	case entity.LiveRoomPrivateInviteAll:
		return true
	case entity.LiveRoomPrivateInviteVip:
		user := userinfodao.GetUserInfoByUserId(userId)
		return user != nil && user.VipLevel > 0
	case entity.LiveRoomPrivateInviteReject:
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
		Category:          room.Category,
		TagId:             strconv.FormatUint(room.TagId, 10),
		TagName:           getRoomTagName(room.TagId),
		Ticket:            room.Ticket,
		Billing:           room.Billing,
		PrivateInviteType: room.PrivateInviteType,
		AllowCallIcon:     allowShowCallIcon(room, userId),
		CreateAt:          room.CreatedAt.Unix(),
		OnlineCount:       countAudienceInRoom(room.ID),
	}
	res.IsBotAnchor, res.CloudPlayerVideo = resolveBotAnchorRoomInfo(room.ID, room.CloudPlayerVideo)
	res.IsTest = room.IsTest
	if u := userinfodao.GetUserInfoByUserId(room.ID); u != nil {
		res.UserType = u.UserType
	}
	//判断一下房间类型
	if room.Category == entity.LiveRoomCategoryPrivate {
		clearFreeTime(userId, room.ID)
		//私密房免费时长
		pay := liveroomdao.GetLiveRoomBillingPay(userId, req.RoomId)

		res.FreeTime = pay.FreeTime
		res.TicketTime = pay.GetTicketTime()
		res.HasTicket = res.TicketTime > 0
	}

	if userId == room.ID {
		res.TotalVideoCallIncome = room.TotalVideoCallIncome
		res.TotalVideoCallTicketIncome = room.TotalVideoCallTicketIncome
		res.TotalVideoCallBillingIncome = room.TotalVideoCallBillingIncome
	}

	return res, nil
}

// EnsureAnchorRoom 确保主播拥有直播间记录(CMS设为主播时预创建,App端后续可完善资料)
func EnsureAnchorRoom(anchorId, guildId uint64) *entity.LiveRoom {
	if room := liveroomdao.GetRoomByAnchor(anchorId); room != nil {
		return room
	}
	// 已下架的直播间不在缓存中,直查 DB,避免重复创建
	if room := liveroomdao.GetRoomFromDB(anchorId); room != nil {
		return room
	}
	room := entity.NewLiveRoom(anchorId, guildId, "", "", "")
	liveroomdao.AddRoomToCache(room)
	return room
}
