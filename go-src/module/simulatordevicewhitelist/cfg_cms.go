package simulatordevicewhitelist

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/simulatordevicewhitelistdao"
	"xr-game-server/dto/simulatordevicewhitelistdto"
	"xr-game-server/entity/cms"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *simulatordevicewhitelistdto.SimulatorDeviceWhitelistListReq) (*httpserver.CMSQueryResp, error) {
	total, list := simulatordevicewhitelistdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *simulatordevicewhitelistdto.CreateSimulatorDeviceWhitelistReq) (*simulatordevicewhitelistdto.CreateSimulatorDeviceWhitelistRes, error) {
	deviceId := strings.TrimSpace(req.DeviceId)
	if deviceId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if simulatordevicewhitelistdao.GetByDeviceId(deviceId) != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.SimulatorDeviceWhitelist{DeviceId: deviceId}
	if err := simulatordevicewhitelistdao.Create(row); err != nil {
		return nil, err
	}
	reloadDeviceMemory()
	return &simulatordevicewhitelistdto.CreateSimulatorDeviceWhitelistRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *simulatordevicewhitelistdto.UpdateSimulatorDeviceWhitelistReq) (*simulatordevicewhitelistdto.UpdateSimulatorDeviceWhitelistRes, error) {
	row := simulatordevicewhitelistdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	deviceId := strings.TrimSpace(req.DeviceId)
	if deviceId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := simulatordevicewhitelistdao.GetByDeviceId(deviceId); existing != nil && existing.ID != row.ID {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.DeviceId = deviceId
	if err := simulatordevicewhitelistdao.Update(row); err != nil {
		return nil, err
	}
	reloadDeviceMemory()
	return &simulatordevicewhitelistdto.UpdateSimulatorDeviceWhitelistRes{Success: true}, nil
}

func Delete(_ context.Context, req *simulatordevicewhitelistdto.DeleteSimulatorDeviceWhitelistReq) (*simulatordevicewhitelistdto.DeleteSimulatorDeviceWhitelistRes, error) {
	if simulatordevicewhitelistdao.GetById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := simulatordevicewhitelistdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadDeviceMemory()
	return &simulatordevicewhitelistdto.DeleteSimulatorDeviceWhitelistRes{Success: true}, nil
}
