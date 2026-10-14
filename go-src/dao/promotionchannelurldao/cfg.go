package promotionchannelurldao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/promotionchannelurldto"
	"xr-game-server/entity/cms"
)

func GetById(id uint64) *entity.PromotionChannelUrl {
	if id == 0 {
		return nil
	}
	var row entity.PromotionChannelUrl
	if err := g.DB().Model(string(entity.TbPromotionChannelUrl)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByChannelCode(code int) *entity.PromotionChannelUrl {
	if code <= 0 {
		return nil
	}
	var row entity.PromotionChannelUrl
	if err := g.DB().Model(string(entity.TbPromotionChannelUrl)).Where("channel_code = ?", code).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.PromotionChannelUrl) error {
	_, err := g.DB().Model(string(entity.TbPromotionChannelUrl)).Save(row)
	return err
}

func Update(row *entity.PromotionChannelUrl) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbPromotionChannelUrl)).WherePri(id).Delete()
	return err
}

func GetList(req *promotionchannelurldto.PromotionChannelUrlListReq) (int, []*promotionchannelurldto.PromotionChannelUrlItem) {
	sql := `select id, channel_name as channelName, channel_code as channelCode, url, created_at as createdAt, updated_at as updatedAt from promotion_channel_urls where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*promotionchannelurldto.PromotionChannelUrlItem, 0)

	key := strings.TrimSpace(req.Key)
	if key != "" {
		sql += ` and (channel_name like ? or cast(channel_code as char) like ? or url like ?) `
		like := "%" + key + "%"
		param = append(param, like, like, like)
	}
	sql += ` order by channel_code asc, id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
