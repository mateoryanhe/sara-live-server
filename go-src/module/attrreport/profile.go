package attrreport

import (
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
)

const emailAccountChannel uint = 7 // 与 auth.EmailChannel 一致

// ResolveHashedEmail 优先 user_ext.email,其次邮箱渠道 open_id; 返回 SHA-256.
func ResolveHashedEmail(userId uint64) string {
	if userId == 0 {
		return ""
	}
	if ext := userinfodao.GetUserExtByUserId(userId); ext != nil {
		if h := HashEmail(ext.Email); h != "" {
			return h
		}
	}
	acc := accountdao.GetAccountById(userId)
	if acc != nil && acc.Channel == emailAccountChannel {
		return HashEmail(acc.OpenId)
	}
	return ""
}

// ResolveHashedPhone 用 user_infos.phone + accounts.phone_area_code; 返回 SHA-256.
func ResolveHashedPhone(userId uint64) string {
	if userId == 0 {
		return ""
	}
	info := userinfodao.GetUserInfoByUserId(userId)
	if info == nil {
		return ""
	}
	phone := strings.TrimSpace(info.Phone)
	if phone == "" {
		return ""
	}
	area := ""
	if acc := accountdao.GetAccountById(userId); acc != nil {
		area = strings.TrimSpace(acc.PhoneAreaCode)
	}
	return HashPhone(area, phone)
}
