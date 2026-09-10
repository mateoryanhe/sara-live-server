package apppkg

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/apppkgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *apppkgdto.AppPkgListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, list := queryAppPkgListFromMemory(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func applyAttributionFields(row *entity.AppPkg, enabled bool, provider, devKey, appId string) {
	row.AttributionEnabled = enabled
	row.AttributionProvider = strings.TrimSpace(provider)
	row.AppsFlyerDevKey = strings.TrimSpace(devKey)
	row.AppsFlyerAppId = strings.TrimSpace(appId)
}

func Create(_ context.Context, req *apppkgdto.CreateAppPkgReq) (*apppkgdto.CreateAppPkgRes, error) {
	packageName := strings.TrimSpace(req.PackageName)
	if packageName == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if findAppPkgByPackageNameFromMemory(packageName, 0) != nil {
		return nil, errercode.CreateCode(errercode.AppPkgExist)
	}
	row := &entity.AppPkg{
		PackageName: packageName,
		Remark:      strings.TrimSpace(req.Remark),
	}
	applyAttributionFields(row, req.AttributionEnabled, req.AttributionProvider, req.AppsFlyerDevKey, req.AppsFlyerAppId)
	if err := cfgdao.CreateAppPkg(row); err != nil {
		return nil, err
	}
	reloadAppPkgMemory()
	return &apppkgdto.CreateAppPkgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *apppkgdto.UpdateAppPkgReq) (*apppkgdto.UpdateAppPkgRes, error) {
	row := getAppPkgByIDFromMemory(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.AppPkgNonExist)
	}
	packageName := strings.TrimSpace(req.PackageName)
	if packageName == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := findAppPkgByPackageNameFromMemory(packageName, req.ID); existing != nil {
		return nil, errercode.CreateCode(errercode.AppPkgExist)
	}
	updated := *row
	updated.PackageName = packageName
	updated.Remark = strings.TrimSpace(req.Remark)
	applyAttributionFields(&updated, req.AttributionEnabled, req.AttributionProvider, req.AppsFlyerDevKey, req.AppsFlyerAppId)
	if err := cfgdao.UpdateAppPkg(&updated); err != nil {
		return nil, err
	}
	reloadAppPkgMemory()
	return &apppkgdto.UpdateAppPkgRes{Success: true}, nil
}

func Delete(_ context.Context, req *apppkgdto.DeleteAppPkgReq) (*apppkgdto.DeleteAppPkgRes, error) {
	if getAppPkgByIDFromMemory(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.AppPkgNonExist)
	}
	if err := cfgdao.DeleteAppPkg(req.ID); err != nil {
		return nil, err
	}
	reloadAppPkgMemory()
	return &apppkgdto.DeleteAppPkgRes{Success: true}, nil
}
