package staticcachecfgdao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/staticcachecfgdto"
	"xr-game-server/entity/sys"
)

func GetByID(id uint64) *entity.StaticCacheRule {
	if id == 0 {
		return nil
	}
	var row entity.StaticCacheRule
	if err := g.DB().Model(string(entity.TbStaticCacheRule)).Where("id = ?", id).Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func GetByFileName(fileName string) *entity.StaticCacheRule {
	fileName = strings.ToLower(strings.TrimSpace(fileName))
	if fileName == "" {
		return nil
	}
	var row entity.StaticCacheRule
	if err := g.DB().Model(string(entity.TbStaticCacheRule)).Where("file_name = ?", fileName).Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func Save(row *entity.StaticCacheRule) error {
	_, err := g.DB().Model(string(entity.TbStaticCacheRule)).Save(row)
	return err
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbStaticCacheRule)).WherePri(id).Delete()
	return err
}

func CountAll() int {
	n, _ := g.DB().Model(string(entity.TbStaticCacheRule)).Count()
	return n
}

func ListAllFileNames() []string {
	rows := make([]*entity.StaticCacheRule, 0)
	_ = g.DB().Model(string(entity.TbStaticCacheRule)).Fields("file_name").Order("id asc").Scan(&rows)
	ret := make([]string, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(row.FileName))
		if name != "" {
			ret = append(ret, name)
		}
	}
	return ret
}

func GetList(req *staticcachecfgdto.StaticCacheRuleListReq) (int, []*staticcachecfgdto.StaticCacheRuleItem) {
	sql := `select id, file_name, remark, created_at, updated_at from static_cache_rules where 1=1 `
	params := make([]any, 0)
	key := strings.TrimSpace(req.Key)
	if key != "" {
		sql += ` and (file_name like ? or remark like ?) `
		like := "%" + key + "%"
		params = append(params, like, like)
	}
	sql += ` order by id desc`
	ctx := gctx.New()
	countSQL := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSQL, params)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	ret := make([]*staticcachecfgdto.StaticCacheRuleItem, 0)
	g.DB().GetScan(ctx, &ret, sql, params)
	return total, ret
}
