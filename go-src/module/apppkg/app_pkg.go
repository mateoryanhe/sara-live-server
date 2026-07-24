package apppkg

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/apppkgdao"
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

func Create(_ context.Context, req *apppkgdto.CreateAppPkgReq) (*apppkgdto.CreateAppPkgRes, error) {
	packageName := strings.TrimSpace(req.PackageName)
	secretKey := strings.TrimSpace(req.SecretKey)
	if packageName == "" || secretKey == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if findAppPkgByPackageNameFromMemory(packageName, 0) != nil {
		return nil, errercode.CreateCode(errercode.AppPkgExist)
	}
	row := &entity.AppPkg{
		PackageName:       packageName,
		SecretKey:         secretKey,
		PrivacyPolicyUrl:  strings.TrimSpace(req.PrivacyPolicyUrl),
		TermsOfServiceUrl: strings.TrimSpace(req.TermsOfServiceUrl),
		Remark:            strings.TrimSpace(req.Remark),
	}
	if err := apppkgdao.Create(row); err != nil {
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
	secretKey := strings.TrimSpace(req.SecretKey)
	if packageName == "" || secretKey == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := findAppPkgByPackageNameFromMemory(packageName, req.ID); existing != nil {
		return nil, errercode.CreateCode(errercode.AppPkgExist)
	}
	updated := *row
	updated.PackageName = packageName
	updated.SecretKey = secretKey
	updated.PrivacyPolicyUrl = strings.TrimSpace(req.PrivacyPolicyUrl)
	updated.TermsOfServiceUrl = strings.TrimSpace(req.TermsOfServiceUrl)
	updated.Remark = strings.TrimSpace(req.Remark)
	if err := apppkgdao.Update(&updated); err != nil {
		return nil, err
	}
	reloadAppPkgMemory()
	return &apppkgdto.UpdateAppPkgRes{Success: true}, nil
}

func Delete(_ context.Context, req *apppkgdto.DeleteAppPkgReq) (*apppkgdto.DeleteAppPkgRes, error) {
	if getAppPkgByIDFromMemory(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.AppPkgNonExist)
	}
	if err := apppkgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadAppPkgMemory()
	return &apppkgdto.DeleteAppPkgRes{Success: true}, nil
}
