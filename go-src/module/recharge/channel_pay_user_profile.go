package recharge

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"xr-game-server/dao/channelpaydao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/errercode"
)

// requireExistingAppUser 校验 user_infos 中确有该用户；禁止用 GetUserInfoByUserId（不存在会新建）
func requireExistingAppUser(userId uint64) error {
	if userId == 0 {
		return errercode.CreateCode(errercode.EmptyUserId)
	}
	if userinfodao.GetUserInfoFromDB(userId) == nil {
		return errercode.CreateCode(errercode.EmptyUserId)
	}
	return nil
}

// resolveChannelPayPayer 优先用已存付款资料，其次请求入参，再回落昵称/绑定邮箱
func resolveChannelPayPayer(userId uint64, reqName, reqEmail string) (name, email string) {
	reqName = strings.TrimSpace(reqName)
	reqEmail = strings.ToLower(strings.TrimSpace(reqEmail))

	if profile := channelpaydao.GetExisting(userId); profile != nil {
		if strings.TrimSpace(profile.Name) != "" {
			name = strings.TrimSpace(profile.Name)
		}
		if strings.TrimSpace(profile.Email) != "" {
			email = strings.ToLower(strings.TrimSpace(profile.Email))
		}
	}
	if reqName != "" {
		name = reqName
	}
	if reqEmail != "" && !strings.HasSuffix(reqEmail, "@noreply.local") {
		email = reqEmail
	}
	if name == "" {
		if user := userinfodao.GetUserInfoFromMemory(userId); user != nil && strings.TrimSpace(user.Nickname) != "" {
			name = strings.TrimSpace(user.Nickname)
		} else if user := userinfodao.GetUserInfoFromDB(userId); user != nil && strings.TrimSpace(user.Nickname) != "" {
			name = strings.TrimSpace(user.Nickname)
		}
	}
	if email == "" {
		if ext := userinfodao.GetUserExtFromMemory(userId); ext != nil && strings.TrimSpace(ext.Email) != "" {
			email = strings.ToLower(strings.TrimSpace(ext.Email))
		}
	}
	if name == "" {
		name = fmt.Sprintf("User %d", userId)
	}
	if email == "" {
		email = fmt.Sprintf("user%d@noreply.local", userId)
	}
	// 仅真实入参才落库（且调用方已校验用户存在）
	if reqName != "" || (reqEmail != "" && !strings.HasSuffix(reqEmail, "@noreply.local")) {
		channelpaydao.UpsertPayerInfo(userId, reqName, reqEmail)
	}
	return name, email
}

func normalizePayDisplayName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return name
	}
	if !strings.Contains(name, " ") {
		return name + " User"
	}
	return name
}

func isLikelyRealPayEmail(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return false
	}
	if strings.HasSuffix(email, "@noreply.local") || strings.HasSuffix(email, "@noreply.sara.live") {
		return false
	}
	return true
}

func isLikelyRealPayName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "user ") {
		rest := strings.TrimSpace(name[5:])
		allDigit := rest != ""
		for _, r := range rest {
			if !unicode.IsDigit(r) {
				allDigit = false
				break
			}
		}
		if allDigit {
			return false
		}
	}
	if _, err := strconv.ParseUint(name, 10, 64); err == nil {
		return false
	}
	return true
}

// GetChannelPayUserProfile App 查询已存付款人资料(无需鉴权,userId 由 App 上报)
func GetChannelPayUserProfile(_ context.Context, req *rechargeorderdto.AppGetChannelPayUserProfileReq) (*rechargeorderdto.AppGetChannelPayUserProfileRes, error) {
	userId, err := strconv.ParseUint(strings.TrimSpace(req.UserId), 10, 64)
	if err != nil || userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := requireExistingAppUser(userId); err != nil {
		return nil, err
	}
	res := &rechargeorderdto.AppGetChannelPayUserProfileRes{}
	if row := channelpaydao.GetExisting(userId); row != nil {
		res.Name = strings.TrimSpace(row.Name)
		res.Email = strings.TrimSpace(row.Email)
	}
	return res, nil
}

// SaveChannelPayUserProfile App 保存付款人资料(无需鉴权；填好后再下单)
func SaveChannelPayUserProfile(_ context.Context, req *rechargeorderdto.AppSaveChannelPayUserProfileReq) (*rechargeorderdto.AppSaveChannelPayUserProfileRes, error) {
	userId, err := strconv.ParseUint(strings.TrimSpace(req.UserId), 10, 64)
	if err != nil || userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := requireExistingAppUser(userId); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if name == "" || email == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !isLikelyRealPayEmail(email) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := channelpaydao.UpsertPayerInfo(userId, name, email)
	res := &rechargeorderdto.AppSaveChannelPayUserProfileRes{Success: true}
	if row != nil {
		res.Name = strings.TrimSpace(row.Name)
		res.Email = strings.TrimSpace(row.Email)
	}
	return res, nil
}
