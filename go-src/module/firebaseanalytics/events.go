package firebaseanalytics

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/gameevent"
)

func initEvents() {
	event.Sub(gameevent.RegisterEvent, onRegisterEvent)
}

func onRegisterEvent(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	userId := val.UserId
	registeredAt := val.RegisteredAt
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		trackRegister(ctx, userId, registeredAt)
	})
}

func trackRegister(ctx context.Context, userId uint64, at time.Time) {
	if userId == 0 || !IsEnabled() {
		return
	}
	countryCode := resolveUserCountryCode(userId)
	if at.IsZero() {
		at = time.Now()
	}
	sendSignUpEvent(ctx, userId, countryCode, at)
}

func resolveUserCountryCode(userId uint64) string {
	location := userloginlocationdao.GetByUserId(userId)
	if location == nil {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(location.RegisterCountry))
}

func clientIDForUser(userId uint64) string {
	return "srv." + strconv.FormatUint(userId, 10)
}
