package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	liverevenueconst "xr-game-server/constants/liverevenue"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var revenueLogCacheMgr *cache.RowCache[*entity.LiveRevenueLog]

func initRevenueLogDao() {
	revenueLogCacheMgr = cache.NewRowCache[*entity.LiveRevenueLog]()
}

// GetRevenueLogById 按主键查询直播收益流水(走缓存);数据库不存在则返回新实例
func GetRevenueLogById(id uint64) *entity.LiveRevenueLog {
	if id == 0 || revenueLogCacheMgr == nil {
		return nil
	}
	v := revenueLogCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.LiveRevenueLog, error) {
		var ret *entity.LiveRevenueLog
		err := g.DB().Model(string(entity.TbLiveRevenueLog)).WherePri(id).Scan(&ret)
		if err != nil || ret == nil {
			return entity.NewLiveRevenueLog(id), nil
		}
		return ret, nil
	})
	return v
}

// FindLatestUnrefundedVideoCallBillingLog 查询通话订单最近一条未退款的分钟计费流水
func FindLatestUnrefundedVideoCallBillingLog(orderId, callerId uint64) *entity.LiveRevenueLog {
	if orderId == 0 || callerId == 0 {
		return nil
	}
	ctx := gctx.New()
	var row entity.LiveRevenueLog
	err := g.DB().Model(string(entity.TbLiveRevenueLog)).Ctx(ctx).
		Where(string(entity.LiveRevenueLogBizId)+" = ?", orderId).
		Where(string(entity.LiveRevenueLogSenderId)+" = ?", callerId).
		Where(string(entity.LiveRevenueLogRevenueType)+" = ?", uint8(liverevenueconst.LiveRoomVideoCallBilling)).
		Where(string(entity.LiveRevenueLogStatus)+" = ?", entity.LiveRevenueLogStatusNormal).
		Where(string(entity.LiveRevenueLogTotalAmount)+" > ?", 0).
		OrderDesc("id").
		Limit(1).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	log := GetRevenueLogById(row.ID)
	if log == nil {
		return nil
	}
	log.ID = row.ID
	log.RevenueType = row.RevenueType
	log.RoomId = row.RoomId
	log.LiveRecordId = row.LiveRecordId
	log.SenderId = row.SenderId
	log.ReceiverId = row.ReceiverId
	log.BizId = row.BizId
	log.Count = row.Count
	log.UnitPrice = row.UnitPrice
	log.TotalAmount = row.TotalAmount
	log.Status = row.Status
	log.CreatedAt = row.CreatedAt
	log.UpdatedAt = row.UpdatedAt
	return log
}
