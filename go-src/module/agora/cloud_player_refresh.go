package agora

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/entity/live"
)

const cloudPlayerTokenRefreshStartupDelay = time.Minute

type cloudPlayerTokenState struct {
	playerId      string
	tokenExpireAt time.Time
	timerEntry    *gtimer.Entry
}

var cloudPlayerTokenStates sync.Map // anchorId uint64 -> *cloudPlayerTokenState

func initCloudPlayerTokenRefresh() {
	xrtimer.AddOnce(gctx.New(), cloudPlayerTokenRefreshStartupDelay, restoreCloudPlayerTokenRefreshTimers)
}

func restoreCloudPlayerTokenRefreshTimers(ctx context.Context) {
	for _, room := range liveroomdao.ListActiveCloudPlayerRooms() {
		if room == nil || room.ID == 0 || strings.TrimSpace(room.CloudPlayerId) == "" {
			continue
		}
		expireAt := resolveCloudPlayerTokenExpireAt(room)
		ScheduleCloudPlayerTokenRefresh(room.ID, room.CloudPlayerId, expireAt)
		g.Log().Infof(ctx, "restore cloud player token timer anchorId=%d playerId=%s expireAt=%s",
			room.ID, room.CloudPlayerId, expireAt.Format(time.RFC3339))
	}
}

func resolveCloudPlayerTokenExpireAt(room *entity.LiveRoom) time.Time {
	if room != nil && room.CloudPlayerTokenExpireAt != nil && !room.CloudPlayerTokenExpireAt.IsZero() {
		return *room.CloudPlayerTokenExpireAt
	}
	return estimateCloudPlayerTokenExpireAt(room)
}

func estimateCloudPlayerTokenExpireAt(room *entity.LiveRoom) time.Time {
	expireSeconds := getChannelTokenExpireSeconds()
	if room != nil && room.LiveRecordId > 0 {
		if record := liveroomdao.GetLiveRecordById(room.LiveRecordId); record != nil && !record.StartTime.IsZero() {
			return record.StartTime.Add(time.Duration(expireSeconds) * time.Second)
		}
	}
	return time.Now().Add(time.Duration(expireSeconds) * time.Second)
}

func cloudPlayerTokenRefreshAheadDuration() time.Duration {
	return time.Duration(tokenRefreshAheadGapSeconds) * time.Second
}

// ScheduleCloudPlayerTokenRefresh 按 token 到期时间前 2 小时安排续期任务(每个机器人独立计时)
func ScheduleCloudPlayerTokenRefresh(anchorId uint64, playerId string, tokenExpireAt time.Time) {
	playerId = strings.TrimSpace(playerId)
	if anchorId == 0 || playerId == "" || tokenExpireAt.IsZero() {
		return
	}

	CancelCloudPlayerTokenRefresh(anchorId)

	refreshAt := tokenExpireAt.Add(-cloudPlayerTokenRefreshAheadDuration())
	delay := time.Until(refreshAt)
	if delay <= 0 {
		xrtimer.AddOnce(gctx.New(), time.Millisecond, func(ctx context.Context) {
			refreshCloudPlayerTokenForAnchor(ctx, anchorId)
		})
		cloudPlayerTokenStates.Store(anchorId, &cloudPlayerTokenState{
			playerId:      playerId,
			tokenExpireAt: tokenExpireAt,
		})
		return
	}

	entry := xrtimer.AddOnce(gctx.New(), delay, func(ctx context.Context) {
		refreshCloudPlayerTokenForAnchor(ctx, anchorId)
	})
	cloudPlayerTokenStates.Store(anchorId, &cloudPlayerTokenState{
		playerId:      playerId,
		tokenExpireAt: tokenExpireAt,
		timerEntry:    entry,
	})
}

// CancelCloudPlayerTokenRefresh 取消机器人云播放器 token 续期任务
func CancelCloudPlayerTokenRefresh(anchorId uint64) {
	if anchorId == 0 {
		return
	}
	if v, ok := cloudPlayerTokenStates.LoadAndDelete(anchorId); ok {
		if state, ok := v.(*cloudPlayerTokenState); ok && state.timerEntry != nil {
			state.timerEntry.Close()
		}
	}
}

func refreshCloudPlayerTokenForAnchor(ctx context.Context, anchorId uint64) {
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil || room.LiveRecordId == 0 {
		CancelCloudPlayerTokenRefresh(anchorId)
		return
	}

	playerId := strings.TrimSpace(room.CloudPlayerId)
	if playerId == "" {
		CancelCloudPlayerTokenRefresh(anchorId)
		return
	}
	restoreCloudPlayerSequence(playerId)

	expireAt, err := updateBotAnchorCloudPlayerToken(ctx, anchorId, playerId)
	if err == nil {
		expireTime := time.Unix(expireAt, 0)
		room.SetCloudPlayerTokenExpireAt(&expireTime)
		liveroomdao.FlushRoomCache(room)
		ScheduleCloudPlayerTokenRefresh(anchorId, playerId, expireTime)
		return
	}
	if !isCloudPlayerNotFound(err) {
		g.Log().Warningf(ctx, "refresh cloud player token failed anchorId=%d playerId=%s err=%v", anchorId, playerId, err)
		return
	}

	g.Log().Warningf(ctx, "cloud player missing, recreating anchorId=%d playerId=%s", anchorId, playerId)
	unregisterCloudPlayerSequence(playerId)
	newPlayerId, newExpireAt, createErr := StartBotAnchorCloudPlayer(ctx, anchorId, room.CloudPlayerVideo)
	if createErr != nil {
		xrlog.ErrorWithErr(ctx, logSourceAgoraCloudPlayer,
			fmt.Sprintf("recreate failed anchorId=%d oldPlayerId=%s", anchorId, playerId),
			createErr)
		return
	}
	room.SetCloudPlayerId(newPlayerId)
	newExpireTime := time.Unix(newExpireAt, 0)
	room.SetCloudPlayerTokenExpireAt(&newExpireTime)
	liveroomdao.FlushRoomCache(room)
	ScheduleCloudPlayerTokenRefresh(anchorId, newPlayerId, newExpireTime)
	g.Log().Infof(ctx, "cloud player recreated anchorId=%d playerId=%s", anchorId, newPlayerId)
}

func updateBotAnchorCloudPlayerToken(ctx context.Context, anchorId uint64, playerId string) (int64, error) {
	playerId = strings.TrimSpace(playerId)
	if anchorId == 0 || playerId == "" {
		return 0, nil
	}

	cfg := getAgoraCfgCache()
	if err := validateAgoraRestCfg(cfg); err != nil {
		return 0, err
	}

	channelName := strconv.FormatUint(anchorId, 10)
	token, expireAt, err := BuildChannelToken(anchorId, channelName, agoradto.RTCRolePublisher)
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("/%s/v1/projects/%s/cloud-player/players/%s", cfg.CloudPlayerRegion, cfg.AppId, playerId)
	query := fmt.Sprintf("sequence=%d", nextCloudPlayerSequence(playerId))
	requestID := guid.S()

	client := newAgoraRestClient(cfg)
	client.SetHeader("X-Request-ID", requestID)
	resp, err := client.Patch(ctx, agoraRestBaseURL+path+"?"+query, cloudPlayerUpdateReq{
		Player: cloudPlayerUpdateBody{Token: token},
	})
	if err != nil {
		return 0, err
	}
	defer resp.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		g.Log().Infof(ctx, "refresh cloud player token ok anchorId=%d playerId=%s requestId=%s expireAt=%d",
			anchorId, playerId, requestID, expireAt)
		return expireAt, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return 0, errCloudPlayerNotFound
	}

	respBody := resp.ReadAll()
	xrlog.Error(ctx, logSourceAgoraCloudPlayer, fmt.Sprintf(
		"refresh token failed anchorId=%d playerId=%s requestId=%s status=%d path=%s body=%s",
		anchorId, playerId, requestID, resp.StatusCode, path, string(respBody),
	))
	return 0, fmt.Errorf("refresh cloud player token status=%d", resp.StatusCode)
}

var errCloudPlayerNotFound = fmt.Errorf("cloud player not found")

func isCloudPlayerNotFound(err error) bool {
	return err == errCloudPlayerNotFound
}
