package simulatordevicewhitelistdao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/simulatordevicewhitelistdto"
	"xr-game-server/entity/cms"
)

func GetById(id uint64) *entity.SimulatorDeviceWhitelist {
	if id == 0 {
		return nil
	}
	var row entity.SimulatorDeviceWhitelist
	if err := g.DB().Model(string(entity.TbSimulatorDeviceWhitelist)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByDeviceId(deviceId string) *entity.SimulatorDeviceWhitelist {
	deviceId = strings.TrimSpace(deviceId)
	if deviceId == "" {
		return nil
	}
	var row entity.SimulatorDeviceWhitelist
	if err := g.DB().Model(string(entity.TbSimulatorDeviceWhitelist)).Where("device_id = ?", deviceId).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.SimulatorDeviceWhitelist) error {
	_, err := g.DB().Model(string(entity.TbSimulatorDeviceWhitelist)).Save(row)
	return err
}

func Update(row *entity.SimulatorDeviceWhitelist) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbSimulatorDeviceWhitelist)).WherePri(id).Delete()
	return err
}

// ListAllDeviceIds 全部设备码,用于内存缓存
func ListAllDeviceIds() []string {
	rows := make([]*entity.SimulatorDeviceWhitelist, 0)
	_ = g.DB().Model(string(entity.TbSimulatorDeviceWhitelist)).
		Fields("device_id").
		Order("id asc").
		Scan(&rows)
	ret := make([]string, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		id := strings.TrimSpace(row.DeviceId)
		if id == "" {
			continue
		}
		ret = append(ret, id)
	}
	return ret
}

func GetList(req *simulatordevicewhitelistdto.SimulatorDeviceWhitelistListReq) (int, []*simulatordevicewhitelistdto.SimulatorDeviceWhitelistItem) {
	sql := `select id, device_id as deviceId, created_at, updated_at from simulator_device_whitelists where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*simulatordevicewhitelistdto.SimulatorDeviceWhitelistItem, 0)

	key := strings.TrimSpace(req.Key)
	if key != "" {
		sql += ` and device_id like ? `
		param = append(param, "%"+key+"%")
	}
	sql += ` order by id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
