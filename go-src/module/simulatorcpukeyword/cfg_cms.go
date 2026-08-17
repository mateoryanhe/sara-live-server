package simulatorcpukeyword

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/simulatorcpukeyworddao"
	"xr-game-server/dto/simulatorcpukeyworddto"
	"xr-game-server/entity/cms"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *simulatorcpukeyworddto.SimulatorCpuKeywordListReq) (*httpserver.CMSQueryResp, error) {
	total, list := simulatorcpukeyworddao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *simulatorcpukeyworddto.CreateSimulatorCpuKeywordReq) (*simulatorcpukeyworddto.CreateSimulatorCpuKeywordRes, error) {
	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	if keyword == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if simulatorcpukeyworddao.GetByKeyword(keyword) != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.SimulatorCpuKeyword{
		Keyword: keyword,
		Remark:  strings.TrimSpace(req.Remark),
	}
	if err := simulatorcpukeyworddao.Create(row); err != nil {
		return nil, err
	}
	reloadKeywordMemory()
	return &simulatorcpukeyworddto.CreateSimulatorCpuKeywordRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *simulatorcpukeyworddto.UpdateSimulatorCpuKeywordReq) (*simulatorcpukeyworddto.UpdateSimulatorCpuKeywordRes, error) {
	row := simulatorcpukeyworddao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	if keyword == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := simulatorcpukeyworddao.GetByKeyword(keyword); existing != nil && existing.ID != row.ID {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.Keyword = keyword
	row.Remark = strings.TrimSpace(req.Remark)
	if err := simulatorcpukeyworddao.Update(row); err != nil {
		return nil, err
	}
	reloadKeywordMemory()
	return &simulatorcpukeyworddto.UpdateSimulatorCpuKeywordRes{Success: true}, nil
}

func Delete(_ context.Context, req *simulatorcpukeyworddto.DeleteSimulatorCpuKeywordReq) (*simulatorcpukeyworddto.DeleteSimulatorCpuKeywordRes, error) {
	if simulatorcpukeyworddao.GetById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := simulatorcpukeyworddao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadKeywordMemory()
	return &simulatorcpukeyworddto.DeleteSimulatorCpuKeywordRes{Success: true}, nil
}
