package call

import (
	"context"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// CallHeart 通话心跳,呼叫者与接听者每10秒调用一次
func CallHeart(ctx context.Context, req *calldto.CallHeartReq) (*calldto.CallHeartRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	order := calldao.GetOrderById(req.OrderId)
	if order == nil {
		return nil, errercode.CreateCode(errercode.CallOrderNonExist)
	}
	if order.CallerId != userId && order.ReceiverId != userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if order.HasEnded() || (!order.IsCalling() && !order.HasAnswered()) {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	now := time.Now()
	lastHeartTime := getCallUserHeartTime(userId)
	upsertCallUser(userId, order.ID, now)

	if order.IsCallStarted() {
		if !isCallHeartIntervalExceeded(lastHeartTime, now) {
			return &calldto.CallHeartRes{Success: true}, nil
		}
		addCallDurationOnHeart(order, now)
		if err := chargeLiveRoomCallBillingIfDue(order, now); err != nil {
			return nil, err
		}
		calldao.FlushOrderCache(order)
	}

	return &calldto.CallHeartRes{Success: true}, nil
}

const callHeartProcessInterval = 10 * time.Second

func getCallUserHeartTime(userId uint64) *time.Time {
	callUser := calldao.GetUserById(userId)
	if callUser == nil {
		return nil
	}
	return callUser.HeartTime
}

func isCallHeartIntervalExceeded(lastHeartTime *time.Time, now time.Time) bool {
	if lastHeartTime == nil || lastHeartTime.IsZero() {
		return false
	}
	return now.Sub(*lastHeartTime) > callHeartProcessInterval
}

// addCallDurationOnHeart 每次心跳增加10秒通话时长(按接听时间对齐,避免双方重复累加)
func addCallDurationOnHeart(order *entity.CallOrder, now time.Time) {
	if order == nil || order.AnswerTime == nil {
		return
	}
	elapsedSec := int(now.Sub(*order.AnswerTime).Seconds())
	if elapsedSec < 10 {
		return
	}
	targetDuration := uint32(elapsedSec / 10 * 10)
	if order.CallDuration+10 <= targetDuration {
		order.SetCallDuration(order.CallDuration + 10)
	}
}
