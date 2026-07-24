package randomnick

import (
	"context"

	"xr-game-server/dao/randomnickdao"
	"xr-game-server/errercode"
)

// GetCMSCfg 查询昵称库概览(读内存)
func GetCMSCfg(_ context.Context) *CMSCfgRes {
	langs := make([]*CMSLangStatItem, 0, len(SupportedLangs()))
	counts := CountAllLangs()
	for _, lang := range SupportedLangs() {
		langs = append(langs, &CMSLangStatItem{
			Lang:      lang,
			LangCode:  LangCode(lang),
			LangLabel: LangLabel(lang),
			Count:     counts[lang],
			Samples:   SampleNicknames(lang, 5),
		})
	}
	return &CMSCfgRes{
		UseDB: UseDB(),
		Langs: langs,
	}
}

// ImportNicknames CMS 批量导入昵称
func ImportNicknames(_ context.Context, req *ImportNicknamesReq) (*ImportNicknamesRes, error) {
	lang := NormalizeLang(req.Lang)
	names := ParseImportText(req.Content)
	if len(names) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	imported := names
	var err error
	if req.Replace {
		err = randomnickdao.ReplaceByLang(lang, names)
	} else {
		existing, loadErr := randomnickdao.ListNicknamesByLang(lang)
		if loadErr != nil {
			return nil, loadErr
		}
		seen := make(map[string]struct{}, len(existing))
		for _, n := range existing {
			seen[n] = struct{}{}
		}
		appendNames := make([]string, 0, len(names))
		for _, n := range names {
			if _, ok := seen[n]; ok {
				continue
			}
			seen[n] = struct{}{}
			appendNames = append(appendNames, n)
		}
		imported = appendNames
		if len(appendNames) > 0 {
			err = randomnickdao.BatchInsert(lang, appendNames)
		}
	}
	if err != nil {
		return nil, err
	}
	if err := reloadMemory(); err != nil {
		return nil, err
	}
	return &ImportNicknamesRes{
		Imported: len(imported),
		Total:    CountByLang(lang),
	}, nil
}

// ClearNicknames 清空指定语言昵称库
func ClearNicknames(_ context.Context, req *ClearNicknamesReq) (*ClearNicknamesRes, error) {
	lang := NormalizeLang(req.Lang)
	if err := randomnickdao.DeleteByLang(lang); err != nil {
		return nil, err
	}
	if err := reloadMemory(); err != nil {
		return nil, err
	}
	return &ClearNicknamesRes{Success: true, Total: CountByLang(lang)}, nil
}

type CMSCfgRes struct {
	UseDB bool               `json:"useDB"`
	Langs []*CMSLangStatItem `json:"langs"`
}

type CMSLangStatItem struct {
	Lang      uint8    `json:"lang"`
	LangCode  string   `json:"langCode"`
	LangLabel string   `json:"langLabel"`
	Count     int      `json:"count"`
	Samples   []string `json:"samples"`
}

type ImportNicknamesReq struct {
	Lang    uint8  `json:"lang"`
	Content string `json:"content"`
	Replace bool   `json:"replace"`
}

type ImportNicknamesRes struct {
	Imported int `json:"imported"`
	Total    int `json:"total"`
}

type ClearNicknamesReq struct {
	Lang uint8 `json:"lang"`
}

type ClearNicknamesRes struct {
	Success bool `json:"success"`
	Total   int  `json:"total"`
}
