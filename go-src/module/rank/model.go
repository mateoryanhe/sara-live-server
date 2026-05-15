package rank

import (
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/core/math"
)

type ObjModel struct {
	Id         uint64    `json:"id"`
	UpdateTime time.Time `json:"updateTime"`
	Val        uint64    `json:"val"`
}

func (receiver *ObjModel) Add(val uint64) {
	ret := math.ChkUInt64(receiver.Val, val)
	receiver.Val += ret
	receiver.UpdateTime = time.Now()
}
func (receiver *ObjModel) Reduce(val uint64) {
	if receiver.Val >= val {
		receiver.Val -= val
	} else {
		receiver.Val = common.Zero
	}
	receiver.UpdateTime = time.Now()
}

func NewObjModel(id uint64, val uint64, time time.Time) *ObjModel {
	return &ObjModel{
		Id:         id,
		UpdateTime: time,
		Val:        val,
	}
}

// 结算排行榜事件的结构体
type SettlementRankEventData struct {
	RankId uint32
}

func NewSettlementRankEventData(rankId uint32) *SettlementRankEventData {
	return &SettlementRankEventData{
		RankId: rankId,
	}
}

type UpEventData struct {
	Id     uint64
	TypeId uint32
}

func NewUpEventData(roleId uint64, typeId uint32) *UpEventData {
	return &UpEventData{
		Id:     roleId,
		TypeId: typeId,
	}
}

type DownEventData struct {
	Id       uint64
	LockTime *time.Time
	TypeId   uint32
}

func NewDownEventData(roleId uint64, lockTime *time.Time, typeId uint32) *DownEventData {
	return &DownEventData{
		Id:       roleId,
		LockTime: lockTime,
		TypeId:   typeId,
	}
}
