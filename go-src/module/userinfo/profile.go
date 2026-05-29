package userinfo

import (
	"context"
	"time"
	"xr-game-server/constants/gender"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

const birthdayLayout = "2006-01-02"

// UpdateGender 修改性别
func UpdateGender(ctx context.Context, req *userinfodto.UpdateGenderReq) (*userinfodto.UpdateGenderRes, error) {
	if !gender.IsValid(req.Gender) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	if data == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	data.SetGender(req.Gender)
	return &userinfodto.UpdateGenderRes{Gender: data.Gender}, nil
}

// UpdateBirthday 修改出生日期
func UpdateBirthday(ctx context.Context, req *userinfodto.UpdateBirthdayReq) (*userinfodto.UpdateBirthdayRes, error) {
	birthday, err := time.ParseInLocation(birthdayLayout, req.Birthday, time.Local)
	if err != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if birthday.After(today) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	if data == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	data.SetBirthday(&birthday)
	return &userinfodto.UpdateBirthdayRes{Birthday: formatBirthday(data.Birthday)}, nil
}

func formatBirthday(val *time.Time) string {
	if val == nil || val.IsZero() {
		return ""
	}
	return val.Format(birthdayLayout)
}
