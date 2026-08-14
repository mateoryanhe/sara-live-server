package game

import (
	"context"
	"math"
	"strconv"
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/gameplatformdto"
)

const vendorCallbackDefaultRTP = 0.9

// HandleVendorBalance 第三方获取余额回调(operator_player_session 为用户 ID).
func HandleVendorBalance(ctx context.Context, req *gameplatformdto.VendorBalanceReq) (*gameplatformdto.VendorBalanceRes, error) {
	if req == nil {
		return vendorBalanceFail(vendorCallbackCodeInvalidParam, "invalid request"), nil
	}
	bodyParams := collectVendorCallbackBodyParams(ctx)
	fillVendorBalanceReq(req, bodyParams)

	userIDStr := strings.TrimSpace(req.OperatorPlayerSession)
	vendorDetailLog().Infof(ctx, "vendor balance request user_id=%s player_name=%s operator_token=%s game_id=%s ip=%s",
		userIDStr, req.PlayerName, req.OperatorToken, req.GameID, req.IP)

	if fail := validateVendorTransferAuth(req.OperatorToken, req.SecretKey); fail != nil {
		resp := vendorBalanceFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor balance failed user_id=%s code=%d msg=%s", userIDStr, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}
	if fail := validateVendorBalanceReq(req); fail != nil {
		resp := vendorBalanceFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor balance failed user_id=%s code=%d msg=%s", userIDStr, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		resp := vendorBalanceFail(vendorCallbackCodeInvalidParam, "invalid operator_player_session")
		vendorDetailLog().Warningf(ctx, "vendor balance failed user_id=%s code=%d msg=%s", userIDStr, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}
	if accountdao.GetAccountById(userID) == nil {
		resp := vendorBalanceFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor balance failed user_id=%s code=%d msg=%s", userIDStr, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}

	userInfo := loadUserInfoForVendorVerify(userID)
	if userInfo == nil {
		resp := vendorBalanceFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor balance failed user_id=%s code=%d msg=%s", userIDStr, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}
	if !vendorTransferPlayerNameMatch(userInfo, userID, req.PlayerName) {
		resp := vendorBalanceFail(vendorCallbackCodeInvalidParam, "player_name mismatch")
		vendorDetailLog().Warningf(ctx, "vendor balance failed user_id=%s code=%d msg=%s", userIDStr, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}

	resp := vendorBalanceSuccess(userInfo.Gold)
	vendorDetailLog().Infof(ctx, "vendor balance success user_id=%s player_name=%s balance=%v currency_code=%s rtp=%v",
		userIDStr, req.PlayerName, resp.Data.Balance, resp.Data.CurrencyCode, resp.Data.RTP)
	return resp, nil
}

func validateVendorBalanceReq(req *gameplatformdto.VendorBalanceReq) *vendorCallbackFail {
	if req == nil {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid request"}
	}
	if strings.TrimSpace(req.OperatorPlayerSession) == "" {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid operator_player_session"}
	}
	if strings.TrimSpace(req.IP) == "" {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid ip"}
	}
	if strings.TrimSpace(req.GameID) == "" {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid game_id"}
	}
	if strings.TrimSpace(req.PlayerName) == "" {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid player_name"}
	}
	return nil
}

func fillVendorBalanceReq(req *gameplatformdto.VendorBalanceReq, params map[string]string) {
	if req == nil || len(params) == 0 {
		return
	}
	if req.OperatorToken == "" {
		req.OperatorToken = params["operator_token"]
	}
	if req.OperatorPlayerSession == "" {
		req.OperatorPlayerSession = params["operator_player_session"]
	}
	if req.SecretKey == "" {
		req.SecretKey = params["secret_key"]
	}
	if req.IP == "" {
		req.IP = params["ip"]
	}
	if req.GameID == "" {
		req.GameID = params["game_id"]
	}
	if req.PlayerName == "" {
		req.PlayerName = params["player_name"]
	}
}

func vendorBalanceSuccess(gold float64) *gameplatformdto.VendorBalanceRes {
	return &gameplatformdto.VendorBalanceRes{
		Data: &gameplatformdto.VendorBalanceData{
			Balance:      int64(math.Floor(gold)),
			CurrencyCode: vendorCallbackCurrency,
			RTP:          vendorCallbackDefaultRTP,
		},
		Error: nil,
	}
}

func vendorBalanceFail(code int, message string) *gameplatformdto.VendorBalanceRes {
	return &gameplatformdto.VendorBalanceRes{
		Data: nil,
		Error: &gameplatformdto.VendorBalanceError{
			Code:    code,
			Message: message,
		},
	}
}
