package appversioncfg

import (
	"context"

	"xr-game-server/dto/appversioncfgdto"
)

// AppVersionQuery App端查询版本(读内存缓存;开关值仅透传给App,服务端不做拦截)
func AppVersionQuery(_ context.Context, _ *appversioncfgdto.AppVersionQueryReq) (*appversioncfgdto.AppVersionQueryRes, error) {
	enabled, version, buildVersion, downloadUrl, updateDetails := GetVersionQuerySnapshot()
	return &appversioncfgdto.AppVersionQueryRes{
		Enabled:       enabled,
		Version:       version,
		BuildVersion:  buildVersion,
		DownloadUrl:   downloadUrl,
		UpdateDetails: updateDetails,
	}, nil
}
