package auth

import (
	"time"

	"xr-game-server/dao/userinfodao"
	"xr-game-server/gameevent"
)

// setUserExtEmailFirstWrite 该 user_id 的 user_ext.email 从空首次写入时返回 true 并埋点 +1(按 user_id 去重; 新号注销再注册再计一次).
func setUserExtEmailFirstWrite(userId uint64, email string) bool {
	if userId == 0 {
		return false
	}
	email = normalizeEmailKey(email)
	if email == "" {
		return false
	}
	ext := userinfodao.GetUserExtByUserId(userId)
	if ext == nil || normalizeEmailKey(ext.Email) != "" {
		return false
	}
	ext.SetEmail(email)
	userinfodao.PublishUserExt(ext)
	gameevent.PubEmailAccountTracking(time.Now())
	return true
}