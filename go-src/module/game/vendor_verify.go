package game

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/entity"
)

// HandleVendorVerify 第三方身份验证回调(ops=user_id).
func HandleVendorVerify(ctx context.Context, req *gameplatformdto.VendorVerifyReq) (*gameplatformdto.VendorVerifyRes, error) {
	if req == nil {
		return vendorVerifyFail(vendorCallbackCodeInvalidParam, "invalid request"), nil
	}
	signParams, signValue := collectVendorCallbackSignParams(ctx)
	fillVendorVerifyReq(req, signParams, signValue)

	vendorDetailLog().Infof(ctx, "vendor verify request ops=%s operator_token=%s timestamp=%d",
		req.Ops, req.OperatorToken, req.Timestamp)

	if fail := validateVendorCallback(req.OperatorToken, signValue, signParams); fail != nil {
		resp := vendorVerifyFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor verify failed ops=%s code=%d msg=%s", req.Ops, resp.Code, resp.Message)
		return resp, nil
	}

	userID, err := strconv.ParseUint(strings.TrimSpace(req.Ops), 10, 64)
	if err != nil || userID == 0 {
		resp := vendorVerifyFail(vendorCallbackCodeInvalidParam, "invalid ops")
		vendorDetailLog().Warningf(ctx, "vendor verify failed ops=%s code=%d msg=%s", req.Ops, resp.Code, resp.Message)
		return resp, nil
	}
	if accountdao.GetAccountById(userID) == nil {
		resp := vendorVerifyFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor verify failed ops=%s code=%d msg=%s", req.Ops, resp.Code, resp.Message)
		return resp, nil
	}

	userInfo := loadUserInfoForVendorVerify(userID)
	if userInfo == nil {
		resp := vendorVerifyFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor verify failed ops=%s code=%d msg=%s", req.Ops, resp.Code, resp.Message)
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
	vendorDetailLog().Infof(ctx, "vendor verify success ops=%s player_name=%s balance=%v currency=%s",
		req.Ops, resp.Data.PlayerName, resp.Data.Balance, resp.Data.Currency)
	return resp, nil
}

func loadUserInfoForVendorVerify(userID uint64) *entity.UserInfo {
	if userInfo := userinfodao.GetUserInfoFromMemory(userID); userInfo != nil {
		return userInfo
	}
	var row entity.UserInfo
	if err := g.DB().Model(string(entity.TbUserInfo)).Unscoped().
		Where("id = ?", userID).
		Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func fillVendorVerifyReq(req *gameplatformdto.VendorVerifyReq, params map[string]string, signValue string) {
	if req.OperatorToken == "" {
		req.OperatorToken = params["operator_token"]
	}
	if req.Ops == "" {
		req.Ops = params["ops"]
	}
	if req.Sign == "" {
		req.Sign = signValue
	}
	if req.Timestamp == 0 {
		req.Timestamp = parseVendorCallbackTimestamp(params["timestamp"])
	}
}

func vendorVerifyFail(code int, message string) *gameplatformdto.VendorVerifyRes {
	return &gameplatformdto.VendorVerifyRes{
		Code:    code,
		Message: message,
	}
}

func collectVendorCallbackSignParams(ctx context.Context) (map[string]string, string) {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return map[string]string{}, ""
	}
	return CollectSignParamsFromRequest(r)
}

func parseVendorCallbackTimestamp(raw string) int64 {
	parsed, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}
