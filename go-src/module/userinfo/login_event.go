package userinfo

import (
	"time"
	"xr-game-server/core/event"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/dailyuserlogindao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity"
)

func initEvent() {
	event.Sub(event.AppToken, onLoginEvent)
}

func onLoginEvent(data any) {
	val := data.(*event.AppTokenData)
	userInfo := userinfodao.GetUserInfoByUserId(val.Id)
	now := time.Now()
	userInfo.SetLastLoginTime(&now)

	date := entity.FormatDailyLoginStatDate(now)
	if dailyuserlogindao.TryRecordLogin(date, val.Id) {
		stat := dailyloginstatdao.GetByDate(date)
		stat.AddLoginCount(1)
	}
}
