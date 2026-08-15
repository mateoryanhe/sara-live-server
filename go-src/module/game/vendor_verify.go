package game

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/gameplatformdto"
	userentity "xr-game-server/entity/user"
)

// HandleVendorVerify 第三方身份验证回调(operator_player_session 或 ops 为用户 ID).
func HandleVendorVerify(ctx context.Context, req *gameplatformdto.VendorVerifyReq) (*gameplatformdto.VendorVerifyRes, error) {
	if req == nil {
		return vendorVerifyFail(vendorCallbackCodeInvalidParam, "invalid request"), nil
	}
	bodyParams := collectVendorCallbackBodyParams(ctx)
	fillVendorVerifyReq(req, bodyParams)

	userIDStr := resolveVendorVerifyUserID(req)
	vendorDetailLog().Infof(ctx, "vendor verify request user_id=%s operator_token=%s game_id=%s bet_type=%d",
		userIDStr, req.OperatorToken, req.GameID, req.BetType)

	if fail := validateVendorTransferAuth(req.OperatorToken, req.SecretKey); fail != nil {
		resp := vendorVerifyFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor verify failed user_id=%s code=%d msg=%s", userIDStr, resp.Code, resp.Message)
		return resp, nil
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		resp := vendorVerifyFail(vendorCallbackCodeInvalidParam, "invalid operator_player_session")
		vendorDetailLog().Warningf(ctx, "vendor verify failed user_id=%s code=%d msg=%s", userIDStr, resp.Code, resp.Message)
		return resp, nil
	}
	if accountdao.GetAccountById(userID) == nil {
		resp := vendorVerifyFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor verify failed user_id=%s code=%d msg=%s", userIDStr, resp.Code, resp.Message)
		return resp, nil
	}

	userInfo := loadUserInfoForVendorVerify(userID)
	if userInfo == nil {
		resp := vendorVerifyFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor verify failed user_id=%s code=%d msg=%s", userIDStr, resp.Code, resp.Message)
		return resp, nil
	}

	resp := &gameplatformdto.VendorVerifyRes{
		Code: vendorCallbackCodeOK,
		Data: &gameplatformdto.VendorVerifyData{
			PlayerName: resolvePlayerName(userInfo, userID),
			Currency:   vendorCallbackCurrency,
			Balance:    userInfo.Gold,
		},
	}
	vendorDetailLog().Infof(ctx, "vendor verify success user_id=%s player_name=%s balance=%v currency=%s",
		userIDStr, resp.Data.PlayerName, resp.Data.Balance, resp.Data.Currency)
	return resp, nil
}

func loadUserInfoForVendorVerify(userID uint64) *userentity.UserInfo {
	if userInfo := userinfodao.GetUserInfoFromMemory(userID); userInfo != nil {
		return userInfo
	}
	var row userentity.UserInfo
	if err := g.DB().Model(string(userentity.TbUserInfo)).Unscoped().
		Where("id = ?", userID).
		Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func fillVendorVerifyReq(req *gameplatformdto.VendorVerifyReq, params map[string]string) {
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
	if req.Ops == "" {
		req.Ops = params["ops"]
	}
	if req.BetType == 0 {
		req.BetType = int(parseVendorCallbackTimestamp(params["bet_type"]))
	}
}

func resolveVendorVerifyUserID(req *gameplatformdto.VendorVerifyReq) string {
	if req == nil {
		return ""
	}
	if userID := strings.TrimSpace(req.OperatorPlayerSession); userID != "" {
		return userID
	}
	return strings.TrimSpace(req.Ops)
}

func vendorVerifyFail(code int, message string) *gameplatformdto.VendorVerifyRes {
	return &gameplatformdto.VendorVerifyRes{
		Code:    code,
		Message: message,
	}
}
