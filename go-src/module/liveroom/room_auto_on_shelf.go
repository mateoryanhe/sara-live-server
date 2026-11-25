package liveroom

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gtimer"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

const (
	anchorAutoOnShelfMinMinutes uint32 = 1
	anchorAutoOnShelfMaxMinutes uint32 = 43200 // 30 天
)

var anchorAutoOnShelfEntries sync.Map // anchorId uint64 -> *gtimer.Entry

func initAnchorAutoOnShelfWatch() {
	restorePendingAnchorAutoOnShelfTasks()
}

func restorePendingAnchorAutoOnShelfTasks() {
	tasks := liveroomdao.LoadPendingAnchorAutoOnShelfTasks()
	if len(tasks) == 0 {
		return
	}
	for _, task := range tasks {
		if task == nil || task.ID == 0 || task.AnchorId == 0 {
			continue
		}
		scheduleAnchorAutoOnShelfTask(task.ID, task.AnchorId, task.RestoreAt)
	}
}

// OffShelfLiveRoomWithAutoOnShelf CMS 下架主播并在指定分钟后自动上架
func OffShelfLiveRoomWithAutoOnShelf(ctx context.Context, req *accountdto.OffShelfLiveRoomWithAutoOnShelfReq) (*accountdto.OffShelfLiveRoomWithAutoOnShelfRes, error) {
	if req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	minutes := req.AutoOnShelfMinutes
	if minutes < anchorAutoOnShelfMinMinutes || minutes > anchorAutoOnShelfMaxMinutes {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	now := time.Now()
	if n, err := liveroomdao.CancelPendingAnchorAutoOnShelfTasks(req.AnchorId, now); err != nil {
		return nil, errercode.CreateCode(errercode.SysError)
	} else if n > 0 {
		cancelAnchorAutoOnShelfTimer(req.AnchorId)
	}

	if _, err := doOffShelfLiveRoom(ctx, req.AnchorId); err != nil {
		return nil, err
	}
	RefreshRoomListCache(ctx)

	restoreAt := now.Add(time.Duration(minutes) * time.Minute)
	task := &entity.AnchorAutoOnShelfTask{}
	task.ID = snowflake.GetId()
	task.CreatedAt = now
	task.UpdatedAt = now
	task.AnchorId = req.AnchorId
	task.OffShelfMinutes = minutes
	task.OffShelfAt = now
	task.RestoreAt = restoreAt
	task.Status = entity.AnchorAutoOnShelfTaskStatusPending
	task.OperatorId = httpserver.GetAuthId(ctx)
	if err := liveroomdao.InsertAnchorAutoOnShelfTask(task); err != nil {
		g.Log().Errorf(ctx, "insert anchor auto on shelf task anchorId=%d err=%v", req.AnchorId, err)
		return nil, errercode.CreateCode(errercode.SysError)
	}
	scheduleAnchorAutoOnShelfTask(task.ID, req.AnchorId, restoreAt)

	return &accountdto.OffShelfLiveRoomWithAutoOnShelfRes{
		Success:       true,
		TaskId:        task.ID,
		Status:        entity.LiveRoomStatusOffShelf,
		RestoreAt:     restoreAt.Format(time.RFC3339),
		RestoreAtUnix: restoreAt.UnixMilli(),
	}, nil
}

func scheduleAnchorAutoOnShelfTask(taskId, anchorId uint64, restoreAt time.Time) {
	if taskId == 0 || anchorId == 0 || restoreAt.IsZero() {
		return
	}
	cancelAnchorAutoOnShelfTimer(anchorId)

	delay := time.Until(restoreAt)
	if delay <= 0 {
		delay = time.Millisecond
	}
	entry := xrtimer.AddOnce(gctx.New(), delay, func(ctx context.Context) {
		executeAnchorAutoOnShelfTask(ctx, taskId, anchorId)
	})
	anchorAutoOnShelfEntries.Store(anchorId, entry)
}

func cancelAnchorAutoOnShelfTimer(anchorId uint64) {
	if anchorId == 0 {
		return
	}
	if v, ok := anchorAutoOnShelfEntries.LoadAndDelete(anchorId); ok {
		if entry, ok := v.(*gtimer.Entry); ok && entry != nil {
			entry.Close()
		}
	}
}

func cancelPendingAnchorAutoOnShelfForAnchor(anchorId uint64) {
	if anchorId == 0 {
		return
	}
	now := time.Now()
	n, err := liveroomdao.CancelPendingAnchorAutoOnShelfTasks(anchorId, now)
	if err != nil {
		g.Log().Errorf(gctx.New(), "cancel pending auto on shelf anchorId=%d err=%v", anchorId, err)
		return
	}
	if n > 0 {
		cancelAnchorAutoOnShelfTimer(anchorId)
	}
}

func executeAnchorAutoOnShelfTask(ctx context.Context, taskId, anchorId uint64) {
	anchorAutoOnShelfEntries.Delete(anchorId)

	lockKey := anchorAutoOnShelfLockKey(anchorId)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	task := liveroomdao.GetAnchorAutoOnShelfTaskById(taskId)
	if task == nil || task.Status != entity.AnchorAutoOnShelfTaskStatusPending {
		return
	}

	now := time.Now()
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil {
		room = liveroomdao.GetRoomFromDB(anchorId)
	}
	if room != nil && room.Status == entity.LiveRoomStatusOnShelf {
		liveroomdao.FinishAnchorAutoOnShelfTask(taskId, entity.AnchorAutoOnShelfTaskStatusCancelled, now)
		return
	}

	if _, err := doOnShelfLiveRoom(ctx, anchorId); err != nil {
		g.Log().Errorf(ctx, "auto on shelf live room failed taskId=%d anchorId=%d err=%v", taskId, anchorId, err)
		return
	}
	RefreshRoomListCache(ctx)
	if !liveroomdao.FinishAnchorAutoOnShelfTask(taskId, entity.AnchorAutoOnShelfTaskStatusDone, now) {
		return
	}
	g.Log().Infof(ctx, "anchor auto on shelf done taskId=%d anchorId=%d", taskId, anchorId)
}

func anchorAutoOnShelfLockKey(anchorId uint64) string {
	return fmt.Sprintf("anchor_auto_on_shelf:%d", anchorId)
}