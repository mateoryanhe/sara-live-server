package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoom db.TbName = "live_rooms"
)

const (
	LiveRoomGuildId db.TbCol = "guild_id"
	LiveRoomTitle   db.TbCol = "title"
	LiveRoomCover   db.TbCol = "cover"
	LiveRoomNotice  db.TbCol = "notice"
	LiveRoomStatus  db.TbCol = "status"
)

// 直播间状态
const (
	LiveRoomStatusClosed uint8 = 0 // 未开播/已下播
	LiveRoomStatusLive   uint8 = 1 // 直播中
)

// LiveRoom 直播间(LiveRoom.ID 与 UserInfo.ID 均为主播用户ID,每个主播仅一个直播间)
type LiveRoom struct {
	migrate.OneModel
	GuildId uint64 `gorm:"index;default:0;comment:所属工会ID" json:"guildId"`
	Title   string `gorm:"size:128;default:'';comment:直播间标题" json:"title"`
	Cover   string `gorm:"size:255;default:'';comment:封面图URL" json:"cover"`
	Notice  string `gorm:"size:512;default:'';comment:公告" json:"notice"`
	Status  uint8  `gorm:"default:0;comment:状态(0未开播,1直播中)" json:"status"`
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
	r.SetStatus(LiveRoomStatusClosed)
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

func (r *LiveRoom) SetCover(v string) {
	r.Cover = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomCover, &syndb.ColData{
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

func (r *LiveRoom) SetStatus(v uint8) {
	r.Status = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoom, LiveRoomStatus, &syndb.ColData{
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
	syndb.RegQuickWithLarge(TbLiveRoom, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbLiveRoom, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbLiveRoom, LiveRoomGuildId)
	syndb.RegQuickWithLarge(TbLiveRoom, LiveRoomTitle)
	syndb.RegQuickWithLarge(TbLiveRoom, LiveRoomCover)
	syndb.RegQuickWithLarge(TbLiveRoom, LiveRoomNotice)
	syndb.RegQuickWithLarge(TbLiveRoom, LiveRoomStatus)
	migrate.AutoMigrate(&LiveRoom{})
}
