package userinfo

import (
	"context"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// SetUserType CMS 修改用户类型(仅允许普通用户/测试人员)
func SetUserType(_ context.Context, req *accountdto.SetUserTypeReq) (*accountdto.SetUserTypeRes, error) {
	if req.UserType != entity.UserTypeNormal && req.UserType != entity.UserTypeTester {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	user := userinfodao.GetUserInfoByUserId(req.AccountId)
	if entity.UserTypeIsAnchor(user.UserType) || user.UserType == entity.UserTypeCMSAuthor {
		return nil, errercode.CreateCode(errercode.UserAlreadyAnchor)
	}
	if user.UserType == req.UserType {
		return &accountdto.SetUserTypeRes{Success: true}, nil
	}
	user.SetUserType(req.UserType)
	return &accountdto.SetUserTypeRes{Success: true}, nil
}
