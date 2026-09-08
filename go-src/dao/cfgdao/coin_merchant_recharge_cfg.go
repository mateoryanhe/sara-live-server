package cfgdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/coinmerchantrechargecfgdto"
	"xr-game-server/entity/recharge"
)

func GetCoinMerchantRechargeCfgById(id uint64) *entity.CoinMerchantRechargeCfg {
	var cfg entity.CoinMerchantRechargeCfg
	err := g.DB().Model(string(entity.TbCoinMerchantRechargeCfg)).Where("id = ?", id).Scan(&cfg)
	if err != nil || cfg.ID == 0 {
		return nil
	}
	return &cfg
}

func GetCoinMerchantRechargeCfgByName(name string) *entity.CoinMerchantRechargeCfg {
	var cfg entity.CoinMerchantRechargeCfg
	err := g.DB().Model(string(entity.TbCoinMerchantRechargeCfg)).Where("name = ?", name).Scan(&cfg)
	if err != nil || cfg.ID == 0 {
		return nil
	}
	return &cfg
}

func CreateCoinMerchantRechargeCfg(cfg *entity.CoinMerchantRechargeCfg) error {
	_, err := g.DB().Model(string(entity.TbCoinMerchantRechargeCfg)).Save(cfg)
	return err
}

func UpdateCoinMerchantRechargeCfg(cfg *entity.CoinMerchantRechargeCfg) error {
	return CreateCoinMerchantRechargeCfg(cfg)
}

func DeleteCoinMerchantRechargeCfg(id uint64) error {
	_, err := g.DB().Model(string(entity.TbCoinMerchantRechargeCfg)).WherePri(id).Delete()
	return err
}

func UpdateCoinMerchantRechargeCfgStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbCoinMerchantRechargeCfg)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetOnShelfCoinMerchantRechargeCfg() []*entity.CoinMerchantRechargeCfg {
	ret := make([]*entity.CoinMerchantRechargeCfg, 0)
	_ = g.DB().Model(string(entity.TbCoinMerchantRechargeCfg)).
		Where("status = ?", entity.CoinMerchantRechargeCfgStatusOnShelf).
		Order("price asc, id asc").
		Scan(&ret)
	return ret
}

func GetCoinMerchantRechargeCfgList(req *coinmerchantrechargecfgdto.CoinMerchantRechargeCfgListReq) (int, []*coinmerchantrechargecfgdto.CoinMerchantRechargeCfgListRes) {
	sql := `select id, name, price, gold, status, created_at, updated_at
		from coin_merchant_recharge_cfgs where 1=1`
	param := make([]any, 0)
	if req.Name != "" {
		sql += ` and name like ?`
		param = append(param, "%"+req.Name+"%")
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = 0`
	case 2:
		sql += ` and status = 1`
	}
	sql += ` order by price asc, id desc`
	ctx := gctx.New()
	total, _ := g.DB().GetCount(ctx, str.GetCountSQL(sql), param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())

	type row struct {
		ID        uint64  `json:"id"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Gold      uint64  `json:"gold"`
		Status    uint8   `json:"status"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}
	rows := make([]*row, 0)
	_ = g.DB().GetScan(ctx, &rows, sql, param)
	list := make([]*coinmerchantrechargecfgdto.CoinMerchantRechargeCfgListRes, 0, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		list = append(list, &coinmerchantrechargecfgdto.CoinMerchantRechargeCfgListRes{
			ID:        strconv.FormatUint(r.ID, 10),
			Name:      r.Name,
			Price:     r.Price,
			Gold:      r.Gold,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	return total, list
}
