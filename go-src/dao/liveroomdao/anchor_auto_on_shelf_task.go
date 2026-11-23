package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	liveentity "xr-game-server/entity/live"
)

func InsertAnchorAutoOnShelfTask(row *liveentity.AnchorAutoOnShelfTask) error {
	if row == nil || row.ID == 0 {
		return nil
	}
	_, err := g.Model(string(liveentity.TbAnchorAutoOnShelfTask)).Ctx(gctx.New()).Data(row).Insert()
	return err
}

func CancelPendingAnchorAutoOnShelfTasks(anchorId uint64, now time.Time) (int64, error) {
	if anchorId == 0 {
		return 0, nil
	}
	res, err := g.Model(string(liveentity.TbAnchorAutoOnShelfTask)).Ctx(gctx.New()).
		Where(string(liveentity.AnchorAutoOnShelfTaskAnchorId)+" = ?", anchorId).
		Where(string(liveentity.AnchorAutoOnShelfTaskStatus)+" = ?", liveentity.AnchorAutoOnShelfTaskStatusPending).
		Data(g.Map{
			string(liveentity.AnchorAutoOnShelfTaskStatus):      liveentity.AnchorAutoOnShelfTaskStatusCancelled,
			string(liveentity.AnchorAutoOnShelfTaskCompletedAt): now,
			"updated_at": now,
		}).Update()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func LoadPendingAnchorAutoOnShelfTasks() []*liveentity.AnchorAutoOnShelfTask {
	ctx := gctx.New()
	rows := make([]*liveentity.AnchorAutoOnShelfTask, 0)
	err := g.Model(string(liveentity.TbAnchorAutoOnShelfTask)).Ctx(ctx).
		Where(string(liveentity.AnchorAutoOnShelfTaskStatus)+" = ?", liveentity.AnchorAutoOnShelfTaskStatusPending).
		OrderAsc(string(liveentity.AnchorAutoOnShelfTaskRestoreAt)).
		Scan(&rows)
	if err != nil {
		g.Log().Errorf(ctx, "load pending anchor auto on shelf tasks failed: %v", err)
		return nil
	}
	return rows
}

func GetAnchorAutoOnShelfTaskById(taskId uint64) *liveentity.AnchorAutoOnShelfTask {
	if taskId == 0 {
		return nil
	}
	row := &liveentity.AnchorAutoOnShelfTask{}
	err := g.Model(string(liveentity.TbAnchorAutoOnShelfTask)).Ctx(gctx.New()).
		Where("id = ?", taskId).
		Scan(row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return row
}

// FinishAnchorAutoOnShelfTask 仅当仍为待执行时标记完成，避免重复执行。
func FinishAnchorAutoOnShelfTask(taskId uint64, status uint8, completedAt time.Time) bool {
	if taskId == 0 {
		return false
	}
	res, err := g.Model(string(liveentity.TbAnchorAutoOnShelfTask)).Ctx(gctx.New()).
		Where("id = ?", taskId).
		Where(string(liveentity.AnchorAutoOnShelfTaskStatus)+" = ?", liveentity.AnchorAutoOnShelfTaskStatusPending).
		Data(g.Map{
			string(liveentity.AnchorAutoOnShelfTaskStatus):      status,
			string(liveentity.AnchorAutoOnShelfTaskCompletedAt): completedAt,
			"updated_at": completedAt,
		}).Update()
	if err != nil {
		g.Log().Errorf(gctx.New(), "finish anchor auto on shelf task id=%d err=%v", taskId, err)
		return false
	}
	affected, _ := res.RowsAffected()
	return affected > 0
}