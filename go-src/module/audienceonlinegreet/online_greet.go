package audienceonlinegreet

import (
	"context"
	"crypto/rand"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/lang"
	"xr-game-server/core/event"
	"xr-game-server/core/push"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	messageentity "xr-game-server/entity/message"
	userentity "xr-game-server/entity/user"
	"xr-game-server/gameevent"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/message"
)

var audienceGreetSentDay sync.Map

func initOnlineGreet() {
	event.Sub(gameevent.LoginEvent, onLoginAudienceAnchorGreet)
}

func onLoginAudienceAnchorGreet(data any) {
	val, ok := data.(*gameevent.LoginEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	clientLang := val.ClientLang
	if clientLang == "" {
		clientLang = lang.DefaultLang
	}
	userId := val.UserId
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		trySendAudienceOnlineAnchorGreets(ctx, userId, clientLang)
	})
}

func trySendAudienceOnlineAnchorGreets(ctx context.Context, audienceId uint64, clientLang lang.Lang) {
	cfg := cfgdao.GetAudienceOnlineAnchorGreetCfgCached()
	if cfg == nil || !cfg.Enabled {
		return
	}
	contents := greetContentsForLang(cfg, clientLang)
	if len(contents) == 0 {
		return
	}
	if !isAudienceUser(audienceId) {
		return
	}
	if !markAudienceGreetSentToday(audienceId) {
		return
	}

	anchorIds := selectGreetAnchorIds(cfg)
	intervalMin, intervalMax := cfg.EffectiveGreetIntervalSec()
	for _, anchorId := range anchorIds {
		if !sleepCtx(ctx, randomQueueDelay(intervalMin, intervalMax)) {
			return
		}
		if anchorId == 0 || anchorId == audienceId {
			continue
		}
		if !push.IsOnline(anchorId) {
			continue
		}
		if liveroom.IsRoomOffShelf(liveroomdao.GetRoomByAnchor(anchorId)) {
			continue
		}
		plain := pickRandomGreetContent(contents)
		if plain == "" {
			continue
		}
		content := message.BuildPrivateTextMessageContent(anchorId, audienceId, plain)
		if content == "" {
			continue
		}
		message.DeliverPrivateMessage(anchorId, audienceId, content)
	}
}

func randomQueueDelay(minSec, maxSec uint32) time.Duration {
	if maxSec < minSec {
		maxSec = minSec
	}
	if minSec == maxSec {
		return time.Duration(minSec) * time.Second
	}
	span := int64(maxSec - minSec + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(span))
	if err != nil {
		return time.Duration(minSec) * time.Second
	}
	return time.Duration(minSec+uint32(n.Int64())) * time.Second
}

func sleepCtx(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return true
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func isAudienceUser(userId uint64) bool {
	user := userinfodao.FindUserInfoByUserId(userId)
	if user == nil {
		return false
	}
	return !userentity.UserTypeIsAnchor(user.UserType)
}

func isHumanAnchorUserType(userType uint8) bool {
	return userType == userentity.UserTypeAnchor || userType == userentity.UserTypeSeniorAnchor
}

func markAudienceGreetSentToday(userId uint64) bool {
	today := time.Now().Format("2006-01-02")
	if v, ok := audienceGreetSentDay.Load(userId); ok {
		if day, _ := v.(string); day == today {
			return false
		}
	}
	audienceGreetSentDay.Store(userId, today)
	return true
}

func selectGreetAnchorIds(cfg *messageentity.AudienceOnlineAnchorGreetCfg) []uint64 {
	onlineAnchors := listOnlineHumanAnchorIds()
	if len(onlineAnchors) == 0 {
		return nil
	}
	limit := int(cfg.EffectiveAnchorCount())
	sort.Slice(onlineAnchors, func(i, j int) bool {
		return oneToOneCallCount(onlineAnchors[i]) > oneToOneCallCount(onlineAnchors[j])
	})
	if len(onlineAnchors) > limit {
		onlineAnchors = onlineAnchors[:limit]
	}
	return onlineAnchors
}

func listOnlineHumanAnchorIds() []uint64 {
	out := make([]uint64, 0, 16)
	for _, userId := range push.OnlineUserIds() {
		if userId == 0 {
			continue
		}
		user := userinfodao.FindUserInfoByUserId(userId)
		if user == nil || !isHumanAnchorUserType(user.UserType) {
			continue
		}
		room := liveroomdao.GetRoomByAnchor(userId)
		if liveroom.IsRoomOffShelf(room) {
			continue
		}
		out = append(out, userId)
	}
	return out
}

func oneToOneCallCount(anchorId uint64) uint64 {
	row := liveroomdao.GetOneToOneRoom(anchorId)
	if row == nil {
		return 0
	}
	return row.CallCount
}