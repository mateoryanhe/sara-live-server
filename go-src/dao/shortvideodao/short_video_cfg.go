package shortvideodao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

var cfg *entity.ShortVideoCfg

func InitShortVideoCfgDao() {
	g.DB().Model(string(entity.TbShortVideoCfg)).Order("id asc").Limit(1).Scan(&cfg)
}

func Get() *entity.ShortVideoCfg {
	return cfg
}

func Save(row *entity.ShortVideoCfg) error {
	_, err := g.DB().Model(string(entity.TbShortVideoCfg)).Save(row)
	cfg = row
	return err
}
