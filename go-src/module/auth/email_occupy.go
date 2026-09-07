package auth

import (
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// isEmailInUseByActiveAccount 邮箱是否被未注销账号占用(EmailChannel open_id 或 user_exts.email).
// excludeUserId>0 时忽略该用户自身.
func isEmailInUseByActiveAccount(email string, excludeUserId uint64) bool {
	email = normalizeEmailKey(email)
	if email == "" {
		return false
	}
	if acc := accountdao.FindActiveAccount(email, EmailChannel); acc != nil && acc.ID != excludeUserId {
		return true
	}
	userId := userinfodao.GetActiveUserIdByEmail(email)
	return userId != 0 && userId != excludeUserId
}

func ensureEmailAvailable(email string, excludeUserId uint64) error {
	if isEmailInUseByActiveAccount(email, excludeUserId) {
		return errercode.CreateCode(errercode.EmailAlreadyInUse)
	}
	return nil
}

// findActiveAccountByBoundEmail 通过 user_exts.email 找到未注销账号(走 email→userId 缓存)
func findActiveAccountByBoundEmail(email string) *entity.Account {
	email = normalizeEmailKey(email)
	if email == "" {
		return nil
	}
	userId := userinfodao.GetActiveUserIdByEmail(email)
	if userId == 0 {
		return nil
	}
	dbAcc := accountdao.GetAccountById(userId)
	if dbAcc == nil || dbAcc.ID == 0 || dbAcc.Cancel {
		userinfodao.InvalidateEmailUserIdCache(email)
		return nil
	}
	if cached := accountdao.GetAccountFromCache(dbAcc.OpenId, dbAcc.Channel, userId); cached != nil {
		return cached
	}
	return dbAcc
}

func accountBoundEmail(accountId uint64) string {
	if accountId == 0 {
		return ""
	}
	ext := userinfodao.GetUserExtByUserId(accountId)
	if ext == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(ext.Email))
}
