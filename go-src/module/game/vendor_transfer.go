package game

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/constants/gameplatform"
	"xr-game-server/core/event"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/gamevendordao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/entity/game"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/wallet"
)

const vendorTransferAmountEpsilon = 0.0001

// HandleVendorTransfer 第三方下注转账回调(扣款/派彩).
func HandleVendorTransfer(ctx context.Context, req *gameplatformdto.VendorTransferReq) (*gameplatformdto.VendorTransferRes, error) {
	if req == nil {
		return vendorTransferFail(vendorCallbackCodeInvalidParam, "invalid request"), nil
	}
	bodyParams := collectVendorTransferBodyParams(ctx)
	fillVendorTransferReq(req, bodyParams)
	req.IP = resolveVendorCallbackIP(ctx, req.IP, bodyParams)

	vendorDetailLog().Infof(ctx,
		"vendor transfer request transaction_id=%s player_name=%s bet_id=%s bet_amount=%v win_amount=%v real_transfer_amount=%v",
		req.TransactionID, req.PlayerName, req.BetID, req.BetAmount, req.WinAmount, req.RealTransferAmount,
	)

	if fail := validateVendorTransferAuth(req.OperatorToken, req.SecretKey); fail != nil {
		resp := vendorTransferFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor transfer failed transaction_id=%s code=%d msg=%s", req.TransactionID, fail.Code, fail.Message)
		return resp, nil
	}
	if fail := validateVendorTransferReq(req); fail != nil {
		resp := vendorTransferFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor transfer failed transaction_id=%s code=%d msg=%s", req.TransactionID, fail.Code, fail.Message)
		return resp, nil
	}

	userID, err := strconv.ParseUint(strings.TrimSpace(req.OperatorPlayerSession), 10, 64)
	if err != nil || userID == 0 {
		resp := vendorTransferFail(vendorCallbackCodeInvalidParam, "invalid operator_player_session")
		vendorDetailLog().Warningf(ctx, "vendor transfer failed transaction_id=%s code=%d msg=%s", req.TransactionID, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}
	if accountdao.GetAccountById(userID) == nil {
		resp := vendorTransferFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor transfer failed transaction_id=%s code=%d msg=%s", req.TransactionID, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}

	userInfo := loadUserInfoForVendorVerify(userID)
	if userInfo == nil {
		resp := vendorTransferFail(vendorCallbackCodePlayerNotFound, "player not found")
		vendorDetailLog().Warningf(ctx, "vendor transfer failed transaction_id=%s code=%d msg=%s", req.TransactionID, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}
	if !vendorTransferPlayerNameMatch(userInfo, userID, req.PlayerName) {
		resp := vendorTransferFail(vendorCallbackCodeInvalidParam, "player_name mismatch")
		vendorDetailLog().Warningf(ctx, "vendor transfer failed transaction_id=%s code=%d msg=%s", req.TransactionID, resp.Error.Code, resp.Error.Message)
		return resp, nil
	}

	transactionID := strings.TrimSpace(req.TransactionID)
	gmlock.Lock(vendorTransferLockKey(transactionID))
	defer gmlock.Unlock(vendorTransferLockKey(transactionID))
	if gamevendordao.IsVendorTransferProcessed(transactionID) {
		balance := userinfodao.GetUserInfoByUserId(userID).Gold
		resp := vendorTransferSuccess(balance, resolveVendorTransferUpdatedTime(req))
		vendorDetailLog().Infof(ctx, "vendor transfer duplicate transaction_id=%s balance=%v", transactionID, balance)
		return resp, nil
	}

	platformType, _ := gameplatform.ParsePlatform(req.Platform)
	gameID := strings.TrimSpace(req.GameID)
	gameName, gameCategory := resolveVendorTransferGameMeta(gameID, req.Platform)
	nameEn, cover := resolveVendorTransferGameInfo(gameID)

	balance, fail := applyVendorTransferWallet(userID, req.BetAmount, req.WinAmount, gameName, gameCategory, transactionID)
	if fail != nil {
		resp := vendorTransferFail(fail.Code, fail.Message)
		vendorDetailLog().Warningf(ctx, "vendor transfer failed transaction_id=%s code=%d msg=%s", transactionID, fail.Code, fail.Message)
		return resp, nil
	}

	if req.BetAmount > 0 {
		recordVendorTransferBetLog(userID, gameID, nameEn, cover, platformType, transactionID, req.BetAmount)
	}
	if req.WinAmount > 0 {
		recordVendorTransferWinLog(userID, gameID, nameEn, cover, platformType, transactionID, req.WinAmount)
	}

	gamevendordao.MarkVendorTransferProcessed(transactionID)

	resp := vendorTransferSuccess(balance, resolveVendorTransferUpdatedTime(req))
	vendorDetailLog().Infof(ctx, "vendor transfer success transaction_id=%s user_id=%d balance=%v bet_amount=%v win_amount=%v",
		transactionID, userID, balance, req.BetAmount, req.WinAmount)
	return resp, nil
}

func validateVendorTransferReq(req *gameplatformdto.VendorTransferReq) *vendorCallbackFail {
	requiredStrings := map[string]string{
		"operator_token":          req.OperatorToken,
		"operator_player_session": req.OperatorPlayerSession,
		"secret_key":              req.SecretKey,
		"game_id":                 req.GameID,
		"player_name":             req.PlayerName,
		"parent_bet_id":           req.ParentBetID,
		"bet_id":                  req.BetID,
		"transaction_id":          req.TransactionID,
		"currency_code":           req.CurrencyCode,
		"platform":                req.Platform,
		"wallet_type":             req.WalletType,
		"is_feature":              req.IsFeature,
		"is_feature_buy":          req.IsFeatureBuy,
	}
	for field, value := range requiredStrings {
		if strings.TrimSpace(value) == "" {
			return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "missing " + field}
		}
	}
	if req.CreateTime <= 0 || req.UpdatedTime <= 0 {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid time"}
	}
	if req.BetType != 1 {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid bet_type"}
	}
	if strings.TrimSpace(req.WalletType) != "C" {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid wallet_type"}
	}
	if _, ok := gameplatform.ParsePlatform(req.Platform); !ok {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid platform"}
	}
	if req.BetAmount < 0 || req.WinAmount < 0 {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid amount"}
	}
	expectedTransfer := req.WinAmount - req.BetAmount
	if math.Abs(req.TransferAmount-expectedTransfer) > vendorTransferAmountEpsilon {
		return &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "invalid transfer_amount"}
	}
	return nil
}

// applyVendorTransferWallet 下注扣金币、派彩加金币.
func applyVendorTransferWallet(userID uint64, betAmount, winAmount float64, gameName, gameCategory, transactionID string) (float64, *vendorCallbackFail) {
	balance := userinfodao.GetUserInfoByUserId(userID).Gold
	gameMeta := &gameevent.CurrencyChangeMeta{
		GameName:      gameName,
		GameCategory:  gameCategory,
		BusinessType:  currency.BusinessTypeGame,
		TransactionId: strings.TrimSpace(transactionID),
	}

	if betAmount > 0 {
		var err error
		balance, err = wallet.GoldSub(userID, betAmount, currency.ReasonGameBet, gameMeta)
		if err != nil {
			if isGoldNotEnoughErr(err) {
				return 0, &vendorCallbackFail{Code: vendorCallbackCodeInsufficientBalance, Message: "insufficient balance"}
			}
			return 0, &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "wallet update failed"}
		}
	}

	if winAmount > 0 {
		var err error
		balance, err = wallet.GoldAdd(userID, winAmount, currency.ReasonGameBetWin, gameMeta)
		if err != nil {
			return 0, &vendorCallbackFail{Code: vendorCallbackCodeInvalidParam, Message: "wallet update failed"}
		}
	}

	return balance, nil
}

func isGoldNotEnoughErr(err error) bool {
	if err == nil {
		return false
	}
	var gErr *errercode.XError
	if errors.As(err, &gErr) {
		return gErr.Code() == errercode.GoldNotEnough
	}
	return false
}

func resolveVendorTransferGameInfo(gameID string) (nameEn, cover string) {
	gameCfg := cfgdao.GetGameCfgByGameCode(gameID)
	if gameCfg == nil {
		return "", ""
	}
	return gameCfg.NameEn, gameCfg.Cover
}

func resolveVendorTransferGameMeta(gameID, platform string) (gameName, gameCategory string) {
	gameID = strings.TrimSpace(gameID)
	platform = strings.TrimSpace(platform)
	if lib := cfgdao.GetVendorGameLib(gameID, platform); lib != nil {
		gameName = strings.TrimSpace(lib.NameEn)
		if gameName == "" {
			gameName = strings.TrimSpace(lib.Name)
		}
		gameCategory = strings.TrimSpace(lib.Category)
	}
	if gameName == "" {
		if gameCfg := cfgdao.GetGameCfgByGameCode(gameID); gameCfg != nil {
			gameName = strings.TrimSpace(gameCfg.NameEn)
			if gameName == "" {
				gameName = strings.TrimSpace(gameCfg.LiveGameName)
			}
		}
	}
	return gameName, gameCategory
}

func recordVendorTransferBetLog(userID uint64, gameCode, nameEn, cover string, platformType gameplatform.Platform, orderID string, amount float64) {
	if userID == 0 || !gameplatform.IsValid(platformType) || strings.TrimSpace(orderID) == "" || amount <= 0 {
		return
	}
	liveRoomID, liveRecordID := gamebetdao.ResolveAudienceLiveContext(userID)
	row := entity.NewGameBetLog(
		userID,
		strings.TrimSpace(gameCode),
		strings.TrimSpace(nameEn),
		strings.TrimSpace(cover),
		platformType,
		strings.TrimSpace(orderID),
		amount,
		liveRoomID,
		liveRecordID,
	)
	gamebetdao.PrependGameBetToAppListCache(userID, row)
	if liveRecordID > 0 {
		event.Pub(gameevent.GameBetCreatedEvent, gameevent.NewGameBetCreatedEventData(userID, liveRecordID, amount))
	}
}

func recordVendorTransferWinLog(userID uint64, gameCode, nameEn, cover string, platformType gameplatform.Platform, orderID string, amount float64) {
	if userID == 0 || !gameplatform.IsValid(platformType) || strings.TrimSpace(orderID) == "" || amount <= 0 {
		return
	}
	entity.NewGameWinLog(
		userID,
		strings.TrimSpace(gameCode),
		strings.TrimSpace(nameEn),
		strings.TrimSpace(cover),
		platformType,
		strings.TrimSpace(orderID),
		amount,
	)
}

func vendorTransferLockKey(transactionID string) string {
	return "vendor_transfer:" + strings.TrimSpace(transactionID)
}

func resolveVendorTransferUpdatedTime(req *gameplatformdto.VendorTransferReq) int64 {
	if req != nil && req.UpdatedTime > 0 {
		return req.UpdatedTime
	}
	return 0
}

func vendorTransferPlayerNameMatch(userInfo *userentity.UserInfo, userID uint64, playerName string) bool {
	expected := resolvePlayerName(userInfo, userID)
	return strings.EqualFold(strings.TrimSpace(expected), strings.TrimSpace(playerName))
}

func vendorTransferSuccess(balance float64, updatedTime int64) *gameplatformdto.VendorTransferRes {
	return &gameplatformdto.VendorTransferRes{
		Data: &gameplatformdto.VendorTransferData{
			Balance:      balance,
			CurrencyCode: vendorCallbackCurrency,
			UpdatedTime:  updatedTime,
		},
		Error: nil,
	}
}

func vendorTransferFail(code int, message string) *gameplatformdto.VendorTransferRes {
	return &gameplatformdto.VendorTransferRes{
		Data: nil,
		Error: &gameplatformdto.VendorTransferError{
			Code:    code,
			Message: message,
		},
	}
}

func collectVendorTransferBodyParams(ctx context.Context) map[string]string {
	params, _ := collectVendorCallbackSignParams(ctx)
	return params
}

func fillVendorTransferReq(req *gameplatformdto.VendorTransferReq, params map[string]string) {
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
	if req.ParentBetID == "" {
		req.ParentBetID = params["parent_bet_id"]
	}
	if req.BetID == "" {
		req.BetID = params["bet_id"]
	}
	if req.TransactionID == "" {
		req.TransactionID = params["transaction_id"]
	}
	if req.CurrencyCode == "" {
		req.CurrencyCode = params["currency_code"]
	}
	if req.Platform == "" {
		req.Platform = params["platform"]
	}
	if req.WalletType == "" {
		req.WalletType = params["wallet_type"]
	}
	if req.IsFeature == "" {
		req.IsFeature = params["is_feature"]
	}
	if req.IsFeatureBuy == "" {
		req.IsFeatureBuy = params["is_feature_buy"]
	}
	if req.InventoryPoolID == "" {
		req.InventoryPoolID = params["inventory_pool_id"]
	}
	if req.BetAmount == 0 {
		req.BetAmount = parseVendorTransferFloat(params["bet_amount"])
	}
	if req.WinAmount == 0 {
		req.WinAmount = parseVendorTransferFloat(params["win_amount"])
	}
	if req.TransferAmount == 0 {
		req.TransferAmount = parseVendorTransferFloat(params["transfer_amount"])
	}
	if req.RealTransferAmount == 0 {
		req.RealTransferAmount = parseVendorTransferFloat(params["real_transfer_amount"])
	}
	if req.InventoryAmount == 0 {
		req.InventoryAmount = parseVendorTransferFloat(params["inventory_amount"])
	}
	if req.CreateTime == 0 {
		req.CreateTime = parseVendorCallbackTimestamp(params["create_time"])
	}
	if req.UpdatedTime == 0 {
		req.UpdatedTime = parseVendorCallbackTimestamp(params["updated_time"])
	}
	if req.BetType == 0 {
		req.BetType = int(parseVendorCallbackTimestamp(params["bet_type"]))
	}
}

func parseVendorTransferFloat(raw string) float64 {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0
	}
	return parsed
}
