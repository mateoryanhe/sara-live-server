package apptokendto

import (
	"time"

	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/entity"
)

type AppTokenDto struct {
	Id       string     `json:"id"`
	Token    string     `json:"token"`
	ExpireAt *time.Time `json:"expireAt"`
	Expired  bool       `json:"expired"`
}

func NewAppTokenDto(token *entity.AppToken) *AppTokenDto {
	if token == nil {
		return nil
	}
	ret := &AppTokenDto{
		Id:       gconv.String(token.ID),
		Token:    token.Token,
		ExpireAt: &token.ExpireAt,
		Expired:  token.ExpireAt.Before(time.Now()),
	}
	return ret
}

func NewAppTokenDtoFromCache(userId uint64, token string, expireAt time.Time) *AppTokenDto {
	ret := &AppTokenDto{
		Id:       gconv.String(userId),
		Token:    token,
		ExpireAt: &expireAt,
		Expired:  expireAt.Before(time.Now()),
	}
	return ret
}
