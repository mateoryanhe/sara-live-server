package recharge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

const rtdnMessageDedupTTL = 7 * 24 * time.Hour

var rtdnProcessedCache = gcache.New()
var rtdnProcessedMu sync.Mutex

type pubSubPushRequest struct {
	Message struct {
		Data        string `json:"data"`
		MessageID   string `json:"messageId"`
		PublishTime string `json:"publishTime"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type developerNotification struct {
	Version                    string                      `json:"version"`
	PackageName                string                      `json:"packageName"`
	EventTimeMillis            string                      `json:"eventTimeMillis"`
	OneTimeProductNotification *oneTimeProductNotification `json:"oneTimeProductNotification"`
	TestNotification           *testNotification           `json:"testNotification"`
}

type oneTimeProductNotification struct {
	Version          string `json:"version"`
	NotificationType int    `json:"notificationType"`
	PurchaseToken    string `json:"purchaseToken"`
	Sku              string `json:"sku"`
}

type testNotification struct {
	Version string `json:"version"`
}

// HandleGooglePlayRTDN 处理 Google Play RTDN Pub/Sub Push 回调
func HandleGooglePlayRTDN(ctx context.Context, body []byte) error {
	var push pubSubPushRequest
	if err := json.Unmarshal(body, &push); err != nil {
		logGooglePlayError(ctx, "rtdn invalid pubsub body err=%v", err)
		return nil
	}
	if push.Message.MessageID != "" && isRTDNMessageProcessed(push.Message.MessageID) {
		logGooglePlayInfo(ctx, "rtdn duplicate messageId=%s", push.Message.MessageID)
		return nil
	}
	if push.Message.Data == "" {
		logGooglePlayError(ctx, "rtdn empty message data")
		return nil
	}
	raw, err := base64.StdEncoding.DecodeString(push.Message.Data)
	if err != nil {
		logGooglePlayError(ctx, "rtdn decode data failed err=%v", err)
		return nil
	}
	var notification developerNotification
	if err := json.Unmarshal(raw, &notification); err != nil {
		logGooglePlayError(ctx, "rtdn decode notification failed err=%v", err)
		return nil
	}
	if notification.TestNotification != nil {
		logGooglePlayInfo(ctx, "rtdn test notification version=%s", notification.TestNotification.Version)
		markRTDNMessageProcessed(push.Message.MessageID)
		return nil
	}
	otp := notification.OneTimeProductNotification
	if otp == nil {
		logGooglePlayInfo(ctx, "rtdn ignored: no one-time product notification")
		markRTDNMessageProcessed(push.Message.MessageID)
		return nil
	}
	switch otp.NotificationType {
	case googleOneTimeProductPurchased:
		err = handleGoogleOneTimeProductPurchased(ctx, notification.PackageName, otp.Sku, otp.PurchaseToken)
	default:
		logGooglePlayInfo(ctx, "rtdn ignored notificationType=%d sku=%s", otp.NotificationType, otp.Sku)
	}
	if err != nil {
		return err
	}
	markRTDNMessageProcessed(push.Message.MessageID)
	return nil
}

func isRTDNMessageProcessed(messageId string) bool {
	if messageId == "" {
		return false
	}
	v, err := rtdnProcessedCache.Get(gctx.New(), messageId)
	return err == nil && !v.IsNil()
}

func markRTDNMessageProcessed(messageId string) {
	if messageId == "" {
		return
	}
	rtdnProcessedMu.Lock()
	defer rtdnProcessedMu.Unlock()
	_ = rtdnProcessedCache.Set(gctx.New(), messageId, 1, rtdnMessageDedupTTL)
}
