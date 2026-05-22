package globalcfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"xr-game-server/core/event"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

var cacheMap = make(map[entity.GlobalCfgModule]map[string]*entity.GlobalCfg)

func InitGlobalCfg() {
	tempMap := make(map[entity.GlobalCfgModule]map[string]*entity.GlobalCfg)
	lst := make([]*entity.GlobalCfg, 0)
	g.Model(string(entity.TbGlobalCfg)).Unscoped().Scan(&lst)
	for _, v := range lst {
		tempMapData, ok := tempMap[entity.GlobalCfgModule(v.Module)]
		if !ok {
			tempMapData = make(map[string]*entity.GlobalCfg)
			tempMap[entity.GlobalCfgModule(v.Module)] = tempMapData
		}
		tempMapData[v.Key] = v
	}
	cacheMap = tempMap
	event.Pub(gameevent.GlobalCfgEvent, nil)
}

func Save(data *entity.GlobalCfg) {
	g.Model(string(entity.TbGlobalCfg)).Save(data)
	InitGlobalCfg()
}

func DelById(id uint64) {
	g.Model(string(entity.TbGlobalCfg)).Unscoped().WherePri(id).Delete()
	InitGlobalCfg()
}

func GetCfgList(module, moduleName string) []*entity.GlobalCfg {
	source := getCfgByModule(module)
	if moduleName == "" {
		return source
	}
	keyword := strings.ToLower(strings.TrimSpace(moduleName))
	ret := make([]*entity.GlobalCfg, 0, len(source))
	for _, v := range source {
		if strings.Contains(strings.ToLower(v.ModuleName), keyword) ||
			strings.Contains(strings.ToLower(v.Module), keyword) {
			ret = append(ret, v)
		}
	}
	return ret
}

func getCfgByModule(module string) []*entity.GlobalCfg {
	// 如果module为空字符串，返回所有配置
	if module == "" {
		ret := make([]*entity.GlobalCfg, 0)
		for _, moduleMap := range cacheMap {
			for _, v := range moduleMap {
				ret = append(ret, v)
			}
		}
		return ret
	}

	// 否则返回指定模块的配置
	data, ok := cacheMap[entity.GlobalCfgModule(module)]
	if !ok {
		return make([]*entity.GlobalCfg, 0)
	}
	ret := make([]*entity.GlobalCfg, 0)
	for _, v := range data {
		ret = append(ret, v)
	}
	return ret
}

func GetUint64(module entity.GlobalCfgModule, key string, defaultVal uint64) uint64 {
	mapData, ok := cacheMap[module]
	if !ok {
		return defaultVal
	}
	val, ok := mapData[key]
	if !ok {
		return defaultVal
	}
	return gconv.Uint64(val.Value)
}

func GetStr(module entity.GlobalCfgModule, key string, defaultVal string) string {
	mapData, ok := cacheMap[module]
	if !ok {
		return defaultVal
	}
	val, ok := mapData[key]
	if !ok {
		return defaultVal
	}
	return val.Value
}

func GetBool(module entity.GlobalCfgModule, key string, defaultVal bool) bool {
	mapData, ok := cacheMap[module]
	if !ok {
		return defaultVal
	}
	val, ok := mapData[key]
	if !ok {
		return defaultVal
	}
	return gconv.Bool(val.Value)
}
