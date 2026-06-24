package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoom db.TbName = "live_rooms"
)

const (
	LiveRoomGuildId                      db.TbCol = "guild_id"
	LiveRoomTitle                        db.TbCol = "title"
	LiveRoomCover                        db.TbCol = "cover"
	LiveRoomNotice                       db.TbCol = "notice"
	LiveRoomLiveId                       db.TbCol = "live_record_id"
	LiveRoomHeartTime                    db.TbCol = "heart_time"
	LiveRoomBan                          db.TbCol = "ban"
	LiveRoomBanApplyTime                 db.TbCol = "ban_apply_time"
	LiveRoomBanReason                    db.TbCol = "ban_reason"
	LiveRoomTotalIncome                  db.TbCol = "total_income"
	LiveRoomTotalGiftIncome              db.TbCol = "total_gift_income"
	LiveRoomTotalPaidDanmakuIncome       db.TbCol = "total_paid_danmaku_income"
	LiveRoomTotalPrivateRoomTicketIncome db.TbCol = "total_private_room_ticket_income"
	LiveRoomTotalPrivateRoomWatchIncome  db.TbCol = "total_private_room_watch_income"
	LiveRoomCategory                     db.TbCol = "category"
	LiveRoomTicket                       db.TbCol = "ticket"
	LiveRoomBilling                      db.TbCol = "billing"
)

const (
	LiveRoomCategoryHot     uint8 = 1 // hot
	LiveRoomCategoryGame    uint8 = 2 // game
	LiveRoomCategoryPrivate uint8 = 3 // 私密
)

// LiveRoom 直播间(LiveRoom.ID 与 UserInfo.ID 均为主播用户ID,每个主播仅一个直播间)
type LiveRoom struct {
	migrate.OneModel
	GuildId                      uint64     `gorm:"index;default:0;comment:所属工会ID" json:"guildId"`
	Title                        string     `gorm:"size:128;default:'';comment:直播间标题" json:"title"`
	Cover                        string     `gorm:"size:255;default:'';comment:封面图URL" json:"cover"`
	Notice                       string     `gorm:"size:512;default:'';comment:公告" json:"notice"`
	LiveRecordId                 uint64     `gorm:"default:0;comment:直播记录id" json:"liveRecordId"`
	HeartTime                    *time.Time `gorm:"comment:房间心跳状态,大于5分钟，判断下播" json:"heart_time"`
	Ban                          bool       `gorm:"default:0;comment:封禁状态" json:"ban"`
	BanApplyTime                 *time.Time `gorm:"comment:封禁截止时间" json:"banApplyTime"`
	BanReason                    string     `gorm:"size:512;default:'';comment:封禁原因" json:"banReason"`
	TotalIncome                  float64    `gorm:"default:0;comment:直播收益" json:"totalIncome"`
	TotalGiftIncome              float64    `gorm:"default:0;comment:累计礼物收益" json:"totalGiftIncome"`
	TotalPaidDanmakuIncome       float64    `gorm:"default:0;comment:累计付费弹幕收益" json:"totalPaidDanmakuIncome"`
	TotalPrivateRoomTicketIncome float64    `gorm:"default:0;comment:累计私密直播间门票收益" json:"totalPrivateRoomTicketIncome"`
	TotalPrivateRoomWatchIncome  float64    `gorm:"default:0;comment:累计私密房观看收益" json:"totalPrivateRoomWatchIncome"`
	Category                     uint8      `gorm:"default:1;comment:分类(1=hot,2=game,3=私密)" json:"category"`
	Ticket                       float64    `gorm:"type:decimal(10,4);default:0;comment:门票价格(钻石)" json:"ticket"`
	Billing                      float64    `gorm:"type:decimal(10,4);default:0;comment:计费价格(每分钟钻石)" json:"billing"`
}

// NewLiveRoom 构造内存对象,字段写入通过 syndb 异步入库
// anchorId 同时作为 LiveRoom 的主键 ID
func NewLiveRoom(anchorId, guildId uint64, title, cover, notice string) *LiveRoom {
	r := &LiveRoom{}
	r.ID = anchorId
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetGuildId(guildId)
	r.SetTitle(title)
	r.SetCover(cover)
	r.SetNotice(notice)
	r.SetCategory(LiveRoomCategoryHot)
	return r
}

func (r *LiveRoom) SetGuildId(v uint64) {
	r.GuildId = v
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomGuildId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetTitle(v string) {
	r.Title = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomTitle, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetHeartTime(v *time.Time) {
	r.HeartTime = v
	syndb.AddDataToLazyChan(TbLiveRoom, LiveRoomHeartTime, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetCover(v string) {
	r.Cover = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomCover, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomLiveId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetNotice(v string) {
	r.Notice = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomNotice, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBan(v bool) {
	r.Ban = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomBan, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBanApplyTime(v *time.Time) {
	r.BanApplyTime = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomBanApplyTime, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBanReason(v string) {
	r.BanReason = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomBanReason, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) AddTotalIncome(v float64) {
	r.TotalIncome = math.AddFloat64(r.TotalIncome, v)
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomTotalIncome, &syndb.ColData{
		IdVal: r.ID, ColVal: r.TotalIncome,
	})
}

func (r *LiveRoom) AddTotalGiftIncome(v float64) {
	if v <= 0 {
		return
	}
	r.TotalGiftIncome = math.AddFloat64(r.TotalGiftIncome, v)
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomTotalGiftIncome, &syndb.ColData{
		IdVal: r.ID, ColVal: r.TotalGiftIncome,
	})
}

func (r *LiveRoom) AddTotalPaidDanmakuIncome(v float64) {
	if v <= 0 {
		return
	}
	r.TotalPaidDanmakuIncome = math.AddFloat64(r.TotalPaidDanmakuIncome, v)
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomTotalPaidDanmakuIncome, &syndb.ColData{
		IdVal: r.ID, ColVal: r.TotalPaidDanmakuIncome,
	})
}

func (r *LiveRoom) AddTotalPrivateRoomTicketIncome(v float64) {
	if v <= 0 {
		return
	}
	r.TotalPrivateRoomTicketIncome = math.AddFloat64(r.TotalPrivateRoomTicketIncome, v)
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomTotalPrivateRoomTicketIncome, &syndb.ColData{
		IdVal: r.ID, ColVal: r.TotalPrivateRoomTicketIncome,
	})
}

func (r *LiveRoom) AddTotalPrivateRoomWatchIncome(v float64) {
	if v <= 0 {
		return
	}
	r.TotalPrivateRoomWatchIncome = math.AddFloat64(r.TotalPrivateRoomWatchIncome, v)
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomTotalPrivateRoomWatchIncome, &syndb.ColData{
		IdVal: r.ID, ColVal: r.TotalPrivateRoomWatchIncome,
	})
}

func (r *LiveRoom) SetCategory(v uint8) {
	if v != LiveRoomCategoryHot && v != LiveRoomCategoryGame && v != LiveRoomCategoryPrivate {
		v = LiveRoomCategoryHot
	}
	r.Category = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomCategory, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetTicket(v float64) {
	r.Ticket = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomTicket, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBilling(v float64) {
	r.Billing = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomBilling, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveRoom, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveRoom, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddDataToQuickChan(TbLiveRoom, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func initLiveRoom() {
	syndb.RegQuickWithMiddle(TbLiveRoom, db.CreatedAtName)
	syndb.RegQuickWithMiddle(TbLiveRoom, db.UpdatedAtName)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomGuildId)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomTitle)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomCover)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomNotice)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomLiveId)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomBan)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomBanApplyTime)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomBanReason)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomTotalIncome)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomTotalGiftIncome)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomTotalPaidDanmakuIncome)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomTotalPrivateRoomTicketIncome)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomTotalPrivateRoomWatchIncome)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomCategory)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomTicket)
	syndb.RegQuickWithMiddle(TbLiveRoom, LiveRoomBilling)

	syndb.RegLazyWithLarge(TbLiveRoom, LiveRoomHeartTime)

	migrate.AutoMigrate(&LiveRoom{})
}
