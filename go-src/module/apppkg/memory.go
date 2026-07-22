package apppkg

import (
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"xr-game-server/dao/apppkgdao"
	"xr-game-server/dto/apppkgdto"
	"xr-game-server/entity"
)

const appPkgTimeLayout = "2006-01-02 15:04:05"

type appPkgSnapshot struct {
	byID          map[uint64]*entity.AppPkg
	byPackageName map[string]*entity.AppPkg
	list          []*entity.AppPkg
}

var (
	appPkgCache     atomic.Value // *appPkgSnapshot
	emptyAppPkgList = make([]*entity.AppPkg, 0)
)

func Init() {
	reloadAppPkgMemory()
}

func reloadAppPkgMemory() {
	rows := apppkgdao.GetAll()
	byID := make(map[uint64]*entity.AppPkg, len(rows))
	byPackageName := make(map[string]*entity.AppPkg, len(rows))
	list := make([]*entity.AppPkg, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		byID[row.ID] = row
		byPackageName[row.PackageName] = row
		list = append(list, row)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID > list[j].ID
	})
	appPkgCache.Store(&appPkgSnapshot{
		byID:          byID,
		byPackageName: byPackageName,
		list:          list,
	})
}

func getAppPkgSnapshot() *appPkgSnapshot {
	v := appPkgCache.Load()
	if v == nil {
		return &appPkgSnapshot{
			byID:          make(map[uint64]*entity.AppPkg),
			byPackageName: make(map[string]*entity.AppPkg),
			list:          emptyAppPkgList,
		}
	}
	return v.(*appPkgSnapshot)
}

func getAppPkgByIDFromMemory(id uint64) *entity.AppPkg {
	return getAppPkgSnapshot().byID[id]
}

func findAppPkgByPackageNameFromMemory(packageName string, excludeID uint64) *entity.AppPkg {
	name := strings.TrimSpace(packageName)
	if name == "" {
		return nil
	}
	row := getAppPkgSnapshot().byPackageName[name]
	if row == nil || row.ID == excludeID {
		return nil
	}
	return row
}

func listAppPkgFromMemory(packageName string) []*entity.AppPkg {
	keyword := strings.ToLower(strings.TrimSpace(packageName))
	rows := getAppPkgSnapshot().list
	if keyword == "" {
		return append([]*entity.AppPkg(nil), rows...)
	}
	filtered := make([]*entity.AppPkg, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if strings.Contains(strings.ToLower(row.PackageName), keyword) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func paginateAppPkgList(rows []*entity.AppPkg, pageIndex, pageSize int) ([]*entity.AppPkg, int) {
	total := len(rows)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	start := (pageIndex - 1) * pageSize
	if start >= total {
		return []*entity.AppPkg{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return rows[start:end], total
}

func toAppPkgListRes(row *entity.AppPkg) *apppkgdto.AppPkgListRes {
	if row == nil {
		return nil
	}
	return &apppkgdto.AppPkgListRes{
		ID:          strconv.FormatUint(row.ID, 10),
		PackageName: row.PackageName,
		SecretKey:   row.SecretKey,
		Remark:      row.Remark,
		CreatedAt:   formatAppPkgTime(row.CreatedAt),
		UpdatedAt:   formatAppPkgTime(row.UpdatedAt),
	}
}

func formatAppPkgTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(appPkgTimeLayout)
}

func queryAppPkgListFromMemory(req *apppkgdto.AppPkgListReq) (int, []*apppkgdto.AppPkgListRes) {
	rows, total := paginateAppPkgList(listAppPkgFromMemory(req.PackageName), req.PageIndex, req.PageSize)
	list := make([]*apppkgdto.AppPkgListRes, 0, len(rows))
	for _, row := range rows {
		list = append(list, toAppPkgListRes(row))
	}
	return total, list
}

// GetAppPkgFromMemoryByPackageName 按包名从内存获取App包配置(供其它模块使用)
func GetAppPkgFromMemoryByPackageName(packageName string) *entity.AppPkg {
	return findAppPkgByPackageNameFromMemory(packageName, 0)
}

// GetAllAppPkgFromMemory 获取全部App包配置(供其它模块使用)
func GetAllAppPkgFromMemory() []*entity.AppPkg {
	return append([]*entity.AppPkg(nil), getAppPkgSnapshot().list...)
}
