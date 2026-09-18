package staticcachecfg

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/staticcachecfgdao"
	"xr-game-server/dto/staticcachecfgdto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

var defaultNoCacheFileNames = []string{
	"index.html",
	"version.js",
}

var noCacheFileNameCache atomic.Value // map[string]struct{}

func Init() {
	seedDefaultsIfEmpty()
	reloadCache()
	httpserver.SetStaticNoCacheMatcher(MatchNoCacheFile)
}

func seedDefaultsIfEmpty() {
	if staticcachecfgdao.CountAll() > 0 {
		return
	}
	for _, fileName := range defaultNoCacheFileNames {
		_ = staticcachecfgdao.Save(&entity.StaticCacheRule{
			FileName: fileName,
			Remark:   "default",
		})
	}
}

func reloadCache() {
	snapshot := make(map[string]struct{})
	for _, fileName := range staticcachecfgdao.ListAllFileNames() {
		if normalized := normalizeFileName(fileName); normalized != "" {
			snapshot[normalized] = struct{}{}
		}
	}
	noCacheFileNameCache.Store(snapshot)
}

// MatchNoCacheFile 仅从进程缓存匹配文件名,请求阶段不访问数据库.
func MatchNoCacheFile(filePath string) bool {
	fileName := normalizeFileName(filepath.Base(strings.TrimSpace(filePath)))
	if fileName == "" {
		return false
	}
	raw := noCacheFileNameCache.Load()
	if raw == nil {
		return false
	}
	snapshot, ok := raw.(map[string]struct{})
	if !ok {
		return false
	}
	_, ok = snapshot[fileName]
	return ok
}

func normalizeFileName(fileName string) string {
	fileName = strings.ToLower(strings.TrimSpace(fileName))
	if fileName == "" || fileName == "." || fileName == ".." || strings.ContainsAny(fileName, `/\\`) {
		return ""
	}
	return fileName
}

func GetList(_ context.Context, req *staticcachecfgdto.StaticCacheRuleListReq) (*httpserver.CMSQueryResp, error) {
	total, list := staticcachecfgdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *staticcachecfgdto.CreateStaticCacheRuleReq) (*staticcachecfgdto.CreateStaticCacheRuleRes, error) {
	fileName := normalizeFileName(req.FileName)
	if fileName == "" || staticcachecfgdao.GetByFileName(fileName) != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.StaticCacheRule{FileName: fileName, Remark: strings.TrimSpace(req.Remark)}
	if err := staticcachecfgdao.Save(row); err != nil {
		return nil, err
	}
	reloadCache()
	return &staticcachecfgdto.CreateStaticCacheRuleRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *staticcachecfgdto.UpdateStaticCacheRuleReq) (*staticcachecfgdto.UpdateStaticCacheRuleRes, error) {
	row := staticcachecfgdao.GetByID(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	fileName := normalizeFileName(req.FileName)
	if fileName == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := staticcachecfgdao.GetByFileName(fileName); existing != nil && existing.ID != row.ID {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.FileName = fileName
	row.Remark = strings.TrimSpace(req.Remark)
	if err := staticcachecfgdao.Save(row); err != nil {
		return nil, err
	}
	reloadCache()
	return &staticcachecfgdto.UpdateStaticCacheRuleRes{Success: true}, nil
}

func Delete(_ context.Context, req *staticcachecfgdto.DeleteStaticCacheRuleReq) (*staticcachecfgdto.DeleteStaticCacheRuleRes, error) {
	if staticcachecfgdao.GetByID(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := staticcachecfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadCache()
	return &staticcachecfgdto.DeleteStaticCacheRuleRes{Success: true}, nil
}
