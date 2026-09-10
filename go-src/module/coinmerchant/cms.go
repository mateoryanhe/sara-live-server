package coinmerchant

import (
	"context"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/str"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/coinmerchantdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/auth"
)

// ListCoinMerchants CMS分页查询币商(channel=CoinMerchantChannel)
func ListCoinMerchants(_ context.Context, req *coinmerchantdto.CoinMerchantListReq) (*httpserver.CMSQueryResp, error) {
	sql := `select a.id, a.open_id, a.cancel, a.channel, a.created_at, IFNULL(u.gold, 0) as gold
		from accounts a
		left join user_infos u on u.id = a.id
		where a.channel = ?`
	param := []any{auth.CoinMerchantChannel}
	if key := strings.TrimSpace(req.Key); key != "" {
		sql += ` and (CAST(a.id AS CHAR) LIKE ? or a.open_id = ?)`
		param = append(param, "%"+key+"%", key)
	}
	if req.Cancel != nil {
		if *req.Cancel == 1 {
			sql += ` and a.cancel = 1`
		} else {
			sql += ` and a.cancel = 0`
		}
	}
	sql += ` order by a.id desc`
	ctx := gctx.New()
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())

	type row struct {
		ID        uint64     `json:"id"`
		OpenId    string     `json:"open_id"`
		Cancel    bool       `json:"cancel"`
		Channel   uint       `json:"channel"`
		CreatedAt *time.Time `json:"created_at"`
		Gold      float64    `json:"gold"`
	}
	rows := make([]*row, 0)
	_ = g.DB().GetScan(ctx, &rows, sql, param)

	list := make([]*coinmerchantdto.CoinMerchantItem, 0, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		username := r.OpenId
		cancel := r.Cancel
		gold := r.Gold
		channel := r.Channel
		if acc := accountdao.GetAccountFromCache(r.OpenId, auth.CoinMerchantChannel, r.ID); acc != nil {
			username = acc.OpenId
			cancel = acc.Cancel
			channel = acc.Channel
		}
		if ui := userinfodao.GetUserInfoFromMemory(r.ID); ui != nil {
			gold = ui.Gold
		}
		list = append(list, &coinmerchantdto.CoinMerchantItem{
			ID:        strconv.FormatUint(r.ID, 10),
			Username:  username,
			Gold:      gold,
			Cancel:    cancel,
			Channel:   channel,
			CreatedAt: r.CreatedAt,
		})
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

// CreateCoinMerchant CMS新建币商账号
func CreateCoinMerchant(_ context.Context, req *coinmerchantdto.CreateCoinMerchantReq) (*coinmerchantdto.CreateCoinMerchantRes, error) {
	username := normalizeUsername(req.Username)
	if username == "" || !isValidUsername(username) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	password := strings.TrimSpace(req.Password)
	if len(password) != 32 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if accountdao.FindActiveAccount(username, auth.CoinMerchantChannel) != nil {
		return nil, errercode.CreateCode(errercode.AccountAlreadyExists)
	}

	account := accountdao.RegisterAccount(username, auth.CoinMerchantChannel)
	if account == nil || account.ID == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	account.SetPassword(strings.ToLower(password))
	accountdao.PublishAccountList(account.OpenId, account.Channel)

	user := userinfodao.GetUserInfoByUserId(account.ID)
	user.SetNickname(username)
	user.SetUserType(entity.UserTypeCoinMerchant)
	userinfodao.PublishUserInfo(user)
	userinfodao.GetUserCumulativeStatByUserId(account.ID)

	return &coinmerchantdto.CreateCoinMerchantRes{
		ID:       strconv.FormatUint(account.ID, 10),
		Username: account.OpenId,
	}, nil
}

// ResetCoinMerchantPassword CMS重置币商密码
func ResetCoinMerchantPassword(_ context.Context, req *coinmerchantdto.ResetCoinMerchantPasswordReq) (*coinmerchantdto.ResetCoinMerchantPasswordRes, error) {
	password := strings.TrimSpace(req.Password)
	if req.AccountId == 0 || len(password) != 32 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	account, err := getCoinMerchantAccount(req.AccountId)
	if err != nil {
		return nil, err
	}
	account.SetPassword(strings.ToLower(password))
	accountdao.PublishAccountList(account.OpenId, account.Channel)
	return &coinmerchantdto.ResetCoinMerchantPasswordRes{Success: true}, nil
}

// CancelCoinMerchant CMS注销币商
func CancelCoinMerchant(ctx context.Context, req *coinmerchantdto.CancelCoinMerchantReq) (*coinmerchantdto.CancelCoinMerchantRes, error) {
	account, err := getCoinMerchantAccount(req.AccountId)
	if err != nil {
		return nil, err
	}
	if account.Cancel {
		return &coinmerchantdto.CancelCoinMerchantRes{Success: true}, nil
	}
	ok, err := auth.CancelUser(ctx, &accountdto.CancelReq{
		AccountId: account.ID,
		OpenId:    account.OpenId,
		Channel:   account.Channel,
	})
	if err != nil {
		return nil, err
	}
	return &coinmerchantdto.CancelCoinMerchantRes{Success: ok}, nil
}

func getCoinMerchantAccount(accountId uint64) (*entity.Account, error) {
	dbAcc := accountdao.GetAccountById(accountId)
	if dbAcc == nil || dbAcc.Channel != auth.CoinMerchantChannel {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	account := accountdao.GetAccountFromCache(dbAcc.OpenId, dbAcc.Channel, accountId)
	if account == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return account, nil
}

func normalizeUsername(raw string) string {
	return strings.TrimSpace(raw)
}

func isValidUsername(username string) bool {
	if len(username) < 2 || len(username) > 32 {
		return false
	}
	for _, r := range username {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' || r == '@' {
			continue
		}
		return false
	}
	return true
}
