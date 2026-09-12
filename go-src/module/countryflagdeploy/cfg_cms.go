package countryflagdeploy

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dto/countryflagdeploydto"
	"xr-game-server/module/upload"
)

func GetCountryFlagDeployInfo(_ context.Context, _ *countryflagdeploydto.GetCountryFlagDeployInfoReq) (*countryflagdeploydto.GetCountryFlagDeployInfoRes, error) {
	root, err := getFlagsRoot()
	if err != nil {
		return nil, err
	}
	snap := getCfgCache()
	flags := listVersionFlags(snap.Version)
	return &countryflagdeploydto.GetCountryFlagDeployInfoRes{
		Info: &countryflagdeploydto.CountryFlagDeployInfoItem{
			ID:         formatUintID(snap.ID),
			Version:    snap.Version,
			UrlPrefix:  buildURLPrefix(snap.Version),
			DeployPath: root,
			AcceptExt:  ".zip",
			UpdatedAt:  snap.UpdatedAt,
			FileCount:  len(flags),
			Flags:      flags,
		},
	}, nil
}

// getFlagsRoot 与头像同一 storagePath 下的 country-flags 根目录
func getFlagsRoot() (string, error) {
	storage := strings.TrimSpace(upload.GetStoragePath())
	if storage == "" {
		return "", errImageRootNotConfigured
	}
	return filepath.Join(storage, country.AssetDir), nil
}

// buildURLPrefix 与 App flagIcon 相同:upload.GetUrlByName(资源相对路径)
func buildURLPrefix(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return upload.GetUrlByName(country.AssetDir)
	}
	return upload.GetUrlByName(country.AssetDir + "/" + version)
}

func listVersionFlags(version string) []*countryflagdeploydto.CountryFlagPreviewItem {
	version = strings.TrimSpace(version)
	if version == "" {
		return []*countryflagdeploydto.CountryFlagPreviewItem{}
	}
	codes := listVersionFlagCodes(version)
	if len(codes) == 0 {
		return []*countryflagdeploydto.CountryFlagPreviewItem{}
	}
	sort.Strings(codes)
	out := make([]*countryflagdeploydto.CountryFlagPreviewItem, 0, len(codes))
	for _, code := range codes {
		c := country.MustGet(code)
		rel := country.RelPath(code, version)
		out = append(out, &countryflagdeploydto.CountryFlagPreviewItem{
			Code:   strings.ToUpper(code),
			NameEn: c.NameEn,
			NameZh: c.NameZh,
			File:   country.FlagFileName(code),
			Icon:   upload.GetUrlByName(rel),
		})
	}
	return out
}

func listVersionFlagCodes(version string) []string {
	seen := map[string]struct{}{}
	// 优先本地目录(未开 S3 或本地仍有缓存)
	if root, err := getFlagsRoot(); err == nil {
		dir := filepath.Join(root, version)
		if entries, err := os.ReadDir(dir); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				if code, ok := flagCodeFromFileName(e.Name()); ok {
					seen[code] = struct{}{}
				}
			}
		}
	}
	// 开云桶时再扫对象前缀(与头像同一存储)
	if upload.IsS3Enabled() {
		prefix := country.AssetDir + "/" + version
		if objs, err := upload.ListStoredObjectsByPrefix(prefix); err == nil {
			for _, o := range objs {
				base := filepath.Base(strings.ReplaceAll(o.StoredName, "\\", "/"))
				if code, ok := flagCodeFromFileName(base); ok {
					seen[code] = struct{}{}
				}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for code := range seen {
		out = append(out, code)
	}
	return out
}

func flagCodeFromFileName(name string) (string, bool) {
	base := strings.ToLower(strings.TrimSpace(filepath.Base(name)))
	if !strings.HasSuffix(base, country.FlagIconExt) {
		return "", false
	}
	code := strings.TrimSuffix(base, country.FlagIconExt)
	if len(code) != 2 {
		return "", false
	}
	for _, r := range code {
		if r < 'a' || r > 'z' {
			return "", false
		}
	}
	return code, true
}

func formatUintID(id uint64) string {
	if id == 0 {
		return "0"
	}
	return strconv.FormatUint(id, 10)
}
