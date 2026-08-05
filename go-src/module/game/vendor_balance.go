package game

import (
	"context"
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/gameplatformdto"
)

// HandleVendorBalance 第三方获取余额回调(player_name 对应 verify 返回的玩家名).
func HandleVendorBalance(ctx context.Context, req *gameplatformdto.VendorBalanceReq) (*gameplatformdto.VendorBalanceRes, error) {
	if req == nil {
		return vendorBalanceFail(vendorCallbackCodeInvalidParam, "invalid request"), nil
	}
	signParams, signValue := collectVendorCallbackSignParams(ctx)
	fillVendorBalanceReq(req, signParams, signValue)

	vendorDetailLog().Infof(ctx, "vendor balance request player_name=%s operator_token=%s timestamp=%d",
		req.PlayerName, req.OperatorToken, req.Timestamp)

	if fail := validateVendorCallback(req.OperatorToken, signValue, signParams); fail != nil {
		resp := vendorBalanceFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor balance failed player_name=%s code=%d msg=%s", req.PlayerName, resp.Code, resp.Message)
		return resp, nil
	}

	playerName := strings.TrimSpace(req.PlayerName)
	if playerName == "" {
		resp := vendorBalanceFail(vendorCallbackCodeInvalidParam, "invalid player_name")
		vendorDetailLog().Warningf(ctx, "vendor balance failed player_name=%s code=%d msg=%s", req.PlayerName, resp.Code, resp.Message)
		return resp, nil
	}

	userInfo, userID := loadUserInfoByPlayerName(playerName)
	if userID == 0 || userInfo == nil || accountdao.GetAccountById(userID) == nil {
		resp := vendorBalanceFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor balance failed player_name=%s code=%d msg=%s", req.PlayerName, resp.Code, resp.Message)
		return resp, nil
	}

	resp := &gameplatformdto.VendorBalanceRes{
		Code: vendorCallbackCodeOK,
		Data: &gameplatformdto.VendorBalanceData{
			Balance:      userInfo.Gold,
			CurrencyCode: vendorCallbackCurrency,
		},
	}
	vendorDetailLog().Infof(ctx, "vendor balance success player_name=%s balance=%v currency_code=%s",
		req.PlayerName, resp.Data.Balance, resp.Data.CurrencyCode)
	return resp, nil
}

func fillVendorBalanceReq(req *gameplatformdto.VendorBalanceReq, params map[string]string, signValue string) {
	if req.OperatorToken == "" {
		req.OperatorToken = params["operator_token"]
	}
	if req.PlayerName == "" {
		req.PlayerName = params["player_name"]
	}
	if req.Sign == "" {
		req.Sign = signValue
	}
	if req.Timestamp == 0 {
		req.Timestamp = parseVendorCallbackTimestamp(params["timestamp"])
	}
}

func vendorBalanceFail(code int, message string) *gameplatformdto.VendorBalanceRes {
	return &gameplatformdto.VendorBalanceRes{
		Code:    code,
		Message: message,
	}
}
