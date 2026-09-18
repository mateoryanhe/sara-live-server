package auth

import (
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func normalizeFirebaseUID(firebaseUID string) string {
	return strings.TrimSpace(firebaseUID)
}

// isFirebaseUIDInUseByActiveAccount 检查 UID 是否已被 Firebase 渠道账号或绑定关系占用。
func isFirebaseUIDInUseByActiveAccount(firebaseUID string, excludeUserId uint64) bool {
	firebaseUID = normalizeFirebaseUID(firebaseUID)
	if firebaseUID == "" {
		return false
	}
	if acc := accountdao.FindActiveAccount(firebaseUID, FirebaseChannel); acc != nil && acc.ID != excludeUserId {
		return true
	}
	userId := userinfodao.GetActiveUserIdByFirebaseUID(firebaseUID)
	return userId != 0 && userId != excludeUserId
}

func ensureFirebaseUIDAvailable(firebaseUID string, excludeUserId uint64) error {
	if isFirebaseUIDInUseByActiveAccount(firebaseUID, excludeUserId) {
		return errercode.CreateCode(errercode.FirebaseAlreadyInUse)
	}
	return nil
}

// findActiveAccountByBoundFirebaseUID 通过 user_exts.firebase_uid 找到未注销账号。
func findActiveAccountByBoundFirebaseUID(firebaseUID string) *entity.Account {
	firebaseUID = normalizeFirebaseUID(firebaseUID)
	if firebaseUID == "" {
		return nil
	}
	userId := userinfodao.GetActiveUserIdByFirebaseUID(firebaseUID)
	if userId == 0 {
		return nil
	}
	dbAcc := accountdao.GetAccountById(userId)
	if dbAcc == nil || dbAcc.ID == 0 || dbAcc.Cancel {
		userinfodao.InvalidateFirebaseUserIdCache(firebaseUID)
		return nil
	}
	if cached := accountdao.GetAccountFromCache(dbAcc.OpenId, dbAcc.Channel, userId); cached != nil {
		return cached
	}
	return dbAcc
}

func accountBoundFirebaseUID(accountId uint64) string {
	if accountId == 0 {
		return ""
	}
	ext := userinfodao.GetUserExtByUserId(accountId)
	if ext == nil {
		return ""
	}
	return normalizeFirebaseUID(ext.FirebaseUID)
}
