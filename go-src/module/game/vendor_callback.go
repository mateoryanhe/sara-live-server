package game

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity"
)

const vendorCallbackCurrency = "gold"

const (
	vendorCallbackCodeOK                  = 0
	vendorCallbackCodeInvalidParam        = 1
	vendorCallbackCodeInvalidSign         = 2
	vendorCallbackCodePlayerNotFound      = 3
	vendorCallbackCodePlatformNotReady    = 4
	vendorCallbackCodeInsufficientBalance = 5
)

// VendorCallbackCodeOK 第三方回调成功码.
const VendorCallbackCodeOK = vendorCallbackCodeOK

type vendorCallbackFail struct {
	Code    int
	Message string
}

func validateVendorCallback(operatorToken, signValue string, signParams map[string]string) *vendorCallbackFail {
	cfg := cfgdao.GetGamePlatformCfgFromMemory()
	if cfg == nil || !cfgdao.GamePlatformCfgReady() {
		return &vendorCallbackFail{Code: vendorCallbackCodePlatformNotReady, Message: "platform cfg not ready"}
	}
	operatorToken = strings.TrimSpace(operatorToken)
	if operatorToken == "" || operatorToken != strings.TrimSpace(cfg.Token) {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidSign, Message: "invalid operator_token"}
	}
	if !VerifySign(signParams, signValue, strings.TrimSpace(cfg.SecretKey)) {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidSign, Message: "invalid sign"}
	}
	return nil
}

func validateVendorTransferAuth(operatorToken, secretKey string) *vendorCallbackFail {
	cfg := cfgdao.GetGamePlatformCfgFromMemory()
	if cfg == nil || !cfgdao.GamePlatformCfgReady() {
		return &vendorCallbackFail{Code: vendorCallbackCodePlatformNotReady, Message: "platform cfg not ready"}
	}
	operatorToken = strings.TrimSpace(operatorToken)
	if operatorToken == "" || operatorToken != strings.TrimSpace(cfg.Token) {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidSign, Message: "invalid operator_token"}
	}
	secretKey = strings.TrimSpace(secretKey)
	if secretKey == "" || secretKey != strings.TrimSpace(cfg.SecretKey) {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidSign, Message: "invalid secret_key"}
	}
	return nil
}

func resolvePlayerName(userInfo *entity.UserInfo, userID uint64) string {
	if userInfo == nil {
		return strconv.FormatUint(userID, 10)
	}
	playerName := strings.TrimSpace(userInfo.Nickname)
	if playerName == "" {
		return strconv.FormatUint(userID, 10)
	}
	return playerName
}

func loadUserInfoByPlayerName(playerName string) (*entity.UserInfo, uint64) {
	playerName = strings.TrimSpace(playerName)
	if playerName == "" {
		return nil, 0
	}
	if userID, err := strconv.ParseUint(playerName, 10, 64); err == nil && userID > 0 {
		if userInfo := loadUserInfoForVendorVerify(userID); userInfo != nil {
			return userInfo, userID
		}
	}
	var row entity.UserInfo
	if err := g.DB().Model(string(entity.TbUserInfo)).Unscoped().
		Where(string(entity.UserInfoNickname)+" = ?", playerName).
		Limit(1).
		Scan(&row); err != nil || row.ID == 0 {
		return nil, 0
	}
	return loadUserInfoForVendorVerify(row.ID), row.ID
}
