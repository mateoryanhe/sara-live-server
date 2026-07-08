package call

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/errercode"
)

// GetCallConfirmStatus 查询对方应答确认状态
func GetCallConfirmStatus(ctx context.Context, req *calldto.GetCallConfirmStatusReq) (*calldto.GetCallConfirmStatusRes, error) {
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

	var peerConfirmAt int64
	if t := order.PeerConfirmTime(userId); t != nil {
		peerConfirmAt = t.Unix()
	}

	return &calldto.GetCallConfirmStatusRes{
		Success:            true,
		OrderId:            strconv.FormatUint(order.ID, 10),
		SelfConfirmed:      order.HasUserConfirmed(userId),
		PeerConfirmed:      order.HasPeerConfirmed(userId),
		PeerConfirmAt:      peerConfirmAt,
		AnswerConfirmCount: order.AnswerConfirmCount(),
		CallStarted:        order.IsCallStarted(),
	}, nil
}
