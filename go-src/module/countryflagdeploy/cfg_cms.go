package countryflagdeploy

import (
	"context"
	"strconv"

	"xr-game-server/core/cfg"
	"xr-game-server/dto/countryflagdeploydto"
)

func GetCountryFlagDeployInfo(_ context.Context, _ *countryflagdeploydto.GetCountryFlagDeployInfoReq) (*countryflagdeploydto.GetCountryFlagDeployInfoRes, error) {
	root, err := getFlagsRoot()
	if err != nil {
		return nil, err
	}
	snap := getCfgCache()
	return &countryflagdeploydto.GetCountryFlagDeployInfoRes{
		Info: &countryflagdeploydto.CountryFlagDeployInfoItem{
			ID:         formatUintID(snap.ID),
			Version:    snap.Version,
			UrlPrefix:  buildURLPrefix(snap.Version),
			DeployPath: root,
			AcceptExt:  ".zip",
			UpdatedAt:  snap.UpdatedAt,
		},
	}, nil
}

func getFlagsRoot() (string, error) {
	imageRoot := cfg.GetImageStaticRoot()
	if imageRoot == "" {
		return "", errImageRootNotConfigured
	}
	return joinFlagsRoot(imageRoot), nil
}

func buildURLPrefix(version string) string {
	seg := cfg.GetImageStaticPathSegment()
	if seg == "" {
		seg = "images"
	}
	prefix := "/" + seg + "/" + countryflagdeploydto.CountryFlagStaticDir
	if version == "" {
		return prefix
	}
	return prefix + "/" + version
}

func formatUintID(id uint64) string {
	if id == 0 {
		return "0"
	}
	return strconv.FormatUint(id, 10)
}
