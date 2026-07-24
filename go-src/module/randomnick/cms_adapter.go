package randomnick

import (
	"context"

	"xr-game-server/dto/randomnickdto"
)

func GetCMSCfgDTO(_ context.Context) (*randomnickdto.GetRandomNicknameCfgRes, error) {
	cfg := GetCMSCfg(context.Background())
	langs := make([]*randomnickdto.RandomNicknameLangItem, 0, len(cfg.Langs))
	for _, item := range cfg.Langs {
		if item == nil {
			continue
		}
		langs = append(langs, &randomnickdto.RandomNicknameLangItem{
			Lang:      item.Lang,
			LangCode:  item.LangCode,
			LangLabel: item.LangLabel,
			Count:     item.Count,
			Samples:   item.Samples,
		})
	}
	return &randomnickdto.GetRandomNicknameCfgRes{
		UseDB: cfg.UseDB,
		Langs: langs,
	}, nil
}

func ImportNicknamesDTO(_ context.Context, req *randomnickdto.ImportRandomNicknamesReq) (*randomnickdto.ImportRandomNicknamesRes, error) {
	res, err := ImportNicknames(context.Background(), &ImportNicknamesReq{
		Lang:    req.Lang,
		Content: req.Content,
		Replace: req.Replace,
	})
	if err != nil {
		return nil, err
	}
	return &randomnickdto.ImportRandomNicknamesRes{
		Imported: res.Imported,
		Total:    res.Total,
	}, nil
}

func ClearNicknamesDTO(_ context.Context, req *randomnickdto.ClearRandomNicknamesReq) (*randomnickdto.ClearRandomNicknamesRes, error) {
	res, err := ClearNicknames(context.Background(), &ClearNicknamesReq{Lang: req.Lang})
	if err != nil {
		return nil, err
	}
	return &randomnickdto.ClearRandomNicknamesRes{
		Success: res.Success,
		Total:   res.Total,
	}, nil
}
