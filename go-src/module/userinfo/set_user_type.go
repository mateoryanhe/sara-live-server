package userinfo

import (
	"context"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// SetUserType CMS 修改用户类型(仅允许普通用户/测试人员/币商)
// 设为币商时:仅允许从普通用户(0)变更
func SetUserType(_ context.Context, req *accountdto.SetUserTypeReq) (*accountdto.SetUserTypeRes, error) {
	if req.UserType != entity.UserTypeNormal &&
		req.UserType != entity.UserTypeTester &&
		req.UserType != entity.UserTypeCoinMerchant {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	user := userinfodao.GetUserInfoByUserId(req.AccountId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	if entity.UserTypeIsAnchor(user.UserType) || user.UserType == entity.UserTypeCMSAuthor {
		return nil, errercode.CreateCode(errercode.UserAlreadyAnchor)
	}
	if req.UserType == entity.UserTypeCoinMerchant && user.UserType != entity.UserTypeNormal {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if user.UserType == req.UserType {
		return &accountdto.SetUserTypeRes{Success: true}, nil
	}
	user.SetUserType(req.UserType)
	userinfodao.PublishUserInfo(user)
	return &accountdto.SetUserTypeRes{Success: true}, nil
}
