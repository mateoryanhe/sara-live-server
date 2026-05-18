package guilddao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
)

// GetGuildById 根据ID获取工会
func GetGuildById(id uint64) *entity.LiveGuild {
	var guild entity.LiveGuild
	err := g.DB().Model(string(entity.TbLiveGuild)).Where("id = ?", id).Scan(&guild)
	if err != nil {
		return nil
	}
	return &guild
}

// GetGuildByName 根据名称获取工会
func GetGuildByName(name string) *entity.LiveGuild {
	var guild entity.LiveGuild
	err := g.DB().Model(string(entity.TbLiveGuild)).Where("name = ?", name).Scan(&guild)
	if err != nil {
		return nil
	}
	return &guild
}

// CreateGuild 创建工会
func CreateGuild(guild *entity.LiveGuild) error {
	_, err := g.DB().Model(string(entity.TbLiveGuild)).Save(guild)
	if err == nil {
		return nil
	}
	return err
}

// UpdateGuild 更新工会
func UpdateGuild(guild *entity.LiveGuild) error {
	return CreateGuild(guild)
}

// DeleteGuild 删除工会
func DeleteGuild(id uint64) error {
	_, err := g.DB().Model(string(entity.TbLiveGuild)).WherePri(id).Delete()
	if err == nil {
		return nil
	}
	return err
}

// GetGuildList 获取工会列表
func GetGuildList(req *guilddto.GuildListReq) (int, []*guilddto.GuildListRes) {
	sql := `select g.id, g.name, g.leader_id, g.contact, g.description, g.status, g.created_at, g.updated_at
            from live_guilds g
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*guilddto.GuildListRes, 0)

	if req.Name != "" {
		sql += ` and g.name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Name))
	}

	sql += ` order by g.created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
