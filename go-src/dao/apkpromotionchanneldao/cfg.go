package apkpromotionchanneldao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/apkpromotionchanneldto"
	"xr-game-server/entity/cms"
)

func GetById(id uint64) *entity.ApkPromotionChannel {
	if id == 0 {
		return nil
	}
	var row entity.ApkPromotionChannel
	if err := g.DB().Model(string(entity.TbApkPromotionChannel)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByChannelCode(code string) *entity.ApkPromotionChannel {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil
	}
	var row entity.ApkPromotionChannel
	if err := g.DB().Model(string(entity.TbApkPromotionChannel)).Where("channel_code = ?", code).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.ApkPromotionChannel) error {
	_, err := g.DB().Model(string(entity.TbApkPromotionChannel)).Save(row)
	return err
}

func Update(row *entity.ApkPromotionChannel) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbApkPromotionChannel)).WherePri(id).Delete()
	return err
}

func ListAllForURLCache() []*entity.ApkPromotionChannel {
	var rows []*entity.ApkPromotionChannel
	_ = g.DB().Model(string(entity.TbApkPromotionChannel)).Scan(&rows)
	return rows
}

func GetList(req *apkpromotionchanneldto.ApkPromotionChannelListReq) (int, []*apkpromotionchanneldto.ApkPromotionChannelItem) {
	sql := `select id, channel_name as channelName, channel_code as channelCode,
apk_file_name as apkFileName, apk_download_url as apkDownloadUrl,
created_at as createdAt, updated_at as updatedAt from apk_promotion_channels where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*apkpromotionchanneldto.ApkPromotionChannelItem, 0)

	key := strings.TrimSpace(req.Key)
	if key != "" {
		sql += ` and (channel_name like ? or channel_code like ?) `
		like := "%" + key + "%"
		param = append(param, like, like)
	}
	sql += ` order by channel_code asc, id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
