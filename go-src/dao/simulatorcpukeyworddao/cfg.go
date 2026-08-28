package simulatorcpukeyworddao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/simulatorcpukeyworddto"
	"xr-game-server/entity/cms"
)

func GetById(id uint64) *entity.SimulatorCpuKeyword {
	if id == 0 {
		return nil
	}
	var row entity.SimulatorCpuKeyword
	if err := g.DB().Model(string(entity.TbSimulatorCpuKeyword)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByKeyword(keyword string) *entity.SimulatorCpuKeyword {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}
	var row entity.SimulatorCpuKeyword
	if err := g.DB().Model(string(entity.TbSimulatorCpuKeyword)).Where("keyword = ?", keyword).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.SimulatorCpuKeyword) error {
	_, err := g.DB().Model(string(entity.TbSimulatorCpuKeyword)).Save(row)
	return err
}

func Update(row *entity.SimulatorCpuKeyword) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbSimulatorCpuKeyword)).WherePri(id).Delete()
	return err
}

func CountAll() int {
	n, _ := g.DB().Model(string(entity.TbSimulatorCpuKeyword)).Count()
	return n
}

// ListAllKeywords 全部关键词(小写),用于内存缓存
func ListAllKeywords() []string {
	rows := make([]*entity.SimulatorCpuKeyword, 0)
	_ = g.DB().Model(string(entity.TbSimulatorCpuKeyword)).
		Fields("keyword").
		Order("id asc").
		Scan(&rows)
	ret := make([]string, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		kw := strings.ToLower(strings.TrimSpace(row.Keyword))
		if kw == "" {
			continue
		}
		ret = append(ret, kw)
	}
	return ret
}

func GetList(req *simulatorcpukeyworddto.SimulatorCpuKeywordListReq) (int, []*simulatorcpukeyworddto.SimulatorCpuKeywordItem) {
	sql := `select id, keyword, remark, created_at, updated_at from simulator_cpu_keywords where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*simulatorcpukeyworddto.SimulatorCpuKeywordItem, 0)

	key := strings.TrimSpace(req.Key)
	if key != "" {
		sql += ` and (keyword like ? or remark like ?) `
		like := "%" + key + "%"
		param = append(param, like, like)
	}
	sql += ` order by id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
