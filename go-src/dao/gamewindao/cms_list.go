package gamewindao

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/gamequery"
	"xr-game-server/entity/game"
)

// CMSListFilter CMS 派彩记录查询条件
type CMSListFilter struct {
	UserId       uint64
	GameCode     string
	GameCodes    []string
	OrderId      string
	PlatformType string
	PageIndex    int
	PageSize     int
}

// CMSListItem CMS 派彩列表行(含 JOIN 出的下注金额)
type CMSListItem struct {
	ID           uint64    `json:"id"           orm:"id"`
	CreatedAt    time.Time `json:"createdAt"    orm:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt"    orm:"updated_at"`
	UserId       uint64    `json:"userId"       orm:"user_id"`
	GameCode     string    `json:"gameCode"     orm:"game_code"`
	NameEn       string    `json:"nameEn"       orm:"name_en"`
	Cover        string    `json:"cover"        orm:"cover"`
	Amount       float64   `json:"amount"       orm:"amount"`
	PlatformType string    `json:"platformType" orm:"platform_type"`
	OrderId      string    `json:"orderId"      orm:"order_id"`
	BetAmount    float64   `json:"betAmount"    orm:"bet_amount"`
}

// winBetOrderFallbackSQL 由派彩 order_id 推导下注 order_id 的 SQL 表达式.
// 厂商格式 {a}-{b}-{suffix}: 下注侧通常为 {b}-{b}-{suffix}.
const winBetOrderFallbackSQL = "CONCAT(" +
	"SUBSTRING_INDEX(SUBSTRING_INDEX(w.order_id,'-',2),'-',-1),'-'," +
	"SUBSTRING_INDEX(SUBSTRING_INDEX(w.order_id,'-',2),'-',-1),'-'," +
	"SUBSTRING_INDEX(w.order_id,'-',-1)" +
	")"

// CMSList CMS 分页查询游戏派彩记录(按 ID 倒序; LEFT JOIN 下注金额)
func CMSList(f *CMSListFilter) (int, []*CMSListItem) {
	list := make([]*CMSListItem, 0)
	if f == nil {
		return 0, list
	}
	pageIndex := f.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	ctx := gctx.New()
	base := applyCMSWinFilters(g.Model(string(entity.TbGameWinLog)+" w").Ctx(ctx), f)

	total, err := base.Clone().Count()
	if err != nil {
		return 0, list
	}

	betTb := string(entity.TbGameBetLog)
	m := base.Clone().
		LeftJoin(betTb+" b_exact", "b_exact.order_id = w.order_id").
		LeftJoin(betTb+" b_parent", "b_exact.id IS NULL AND b_parent.order_id = "+winBetOrderFallbackSQL).
		Fields(
			"w.id",
			"w.created_at",
			"w.updated_at",
			"w.user_id",
			"w.game_code",
			"w.name_en",
			"w.cover",
			"w.amount",
			"w.platform_type",
			"w.order_id",
			"COALESCE(b_exact.amount, b_parent.amount, 0) AS bet_amount",
		).
		Order("w.id desc").
		Limit(pageSize).
		Offset((pageIndex - 1) * pageSize)
	_ = m.Scan(&list)
	return total, list
}

func applyCMSWinFilters(m *gdb.Model, f *CMSListFilter) *gdb.Model {
	m = m.Where("w."+string(entity.GameWinLogAmount)+" > ?", 0)
	if f.UserId > 0 {
		m = m.Where("w."+string(entity.GameWinLogUserId)+" = ?", f.UserId)
	}
	m = gamequery.ApplyGameCodeFilter(m, "w."+string(entity.GameWinLogGameCode), f.GameCode, f.GameCodes)
	if orderId := strings.TrimSpace(f.OrderId); orderId != "" {
		m = m.Where("w."+string(entity.GameWinLogOrderId)+" LIKE ?", "%"+orderId+"%")
	}
	if platformType := strings.TrimSpace(f.PlatformType); platformType != "" {
		m = m.Where("UPPER(w."+string(entity.GameWinLogPlatformType)+") = ?", strings.ToUpper(platformType))
	}
	return m
}
