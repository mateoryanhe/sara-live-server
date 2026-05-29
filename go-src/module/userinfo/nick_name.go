package userinfo

import (
	"crypto/rand"
	"math/big"

	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/gameevent"
)

const (
	randomNicknameLength  = 8
	randomNicknameCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func initNicknameEvent() {
	event.Sub(gameevent.RegisterEvent, onRegisterAssignNickname)
}

func onRegisterAssignNickname(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	user := userinfodao.GetUserInfoByUserId(val.UserId)
	if user == nil || user.Nickname != "" {
		return
	}
	user.SetNickname(genRandomNickname())
}

func genRandomNickname() string {
	charsetLen := big.NewInt(int64(len(randomNicknameCharset)))
	b := make([]byte, randomNicknameLength)
	for i := 0; i < randomNicknameLength; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			b[i] = randomNicknameCharset[i%len(randomNicknameCharset)]
			continue
		}
		b[i] = randomNicknameCharset[n.Int64()]
	}
	return string(b)
}
