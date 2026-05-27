package auth

import (
	"time"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/entity"
)

func initAppToken() {
	event.Sub(event.AppToken, onAppToken)
	tokens := accountdao.ListValidAppTokens()
	for _, token := range tokens {
		xrtoken.InitAppToken(token.ID, token.Token, token.ExpireAt)
	}
}

func onAppToken(val any) {
	data := val.(*event.AppTokenData)
	expireAt := time.Now().Add(xrtoken.Time)
	entity.NewAppToken(data.Id, data.Token, expireAt)
}
