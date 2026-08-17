package shortvideo

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/user"
	"xr-game-server/module/auth"
	"xr-game-server/module/randomnick"
)

const maxCMSAuthorNicknameRunes = 32

// createCMSAuthorUser 为 CMS 上传短视频自动创建作者账号
func createCMSAuthorUser(nickname, avatar string) (uint64, error) {
	nickname = normalizeCMSAuthorNickname(nickname)
	if nickname == "" {
		return 0, fmt.Errorf("empty cms author nickname")
	}
	openId := fmt.Sprintf("sv_%s", guid.S())
	account := accountdao.RegisterAccount(openId, auth.ShortVideoAuthorChannel)
	if account == nil || account.ID == 0 {
		return 0, fmt.Errorf("register cms author account failed")
	}
	user := userinfodao.GetUserInfoByUserId(account.ID)
	user.SetNickname(nickname)
	if avatar = strings.TrimSpace(avatar); avatar != "" {
		user.SetAvatar(avatar)
	}
	user.SetUserType(entity.UserTypeCMSAuthor)
	userinfodao.GetUserCumulativeStatByUserId(account.ID)
	return account.ID, nil
}

func normalizeCMSAuthorNickname(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxCMSAuthorNicknameRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxCMSAuthorNicknameRunes])
}

func resolveCMSAuthorNickname(authorNickname string) string {
	if nickname := normalizeCMSAuthorNickname(authorNickname); nickname != "" {
		return nickname
	}
	return normalizeCMSAuthorNickname(randomnick.PickRandomDefault())
}
