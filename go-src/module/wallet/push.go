package wallet

import (
	"fmt"
	"math"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dto/diamonddto"
	"xr-game-server/dto/golddto"
)

// pushDiamondToApp 推送用户最新钻石余额到 App 端
func pushDiamondToApp(userId uint64, diamond float64) {
	push.Data(userId, cmd.DiamondPush, &diamonddto.DiamondPushItem{
		Diamond: formatCurrencyBalance(diamond),
	})
}

// pushGoldToApp 推送用户最新金币余额到 App 端
func pushGoldToApp(userId uint64, goldBalance float64) {
	push.Data(userId, cmd.GoldPush, &golddto.GoldPushItem{
		Gold: formatCurrencyBalance(goldBalance),
	})
}

func formatCurrencyBalance(v float64) string {
	scaled := int64(math.Round(v * currencyScale))
	if scaled < 0 {
		scaled = 0
	}
	whole := scaled / currencyScale
	frac := scaled % currencyScale
	return fmt.Sprintf("%d.%04d", whole, frac)
}
