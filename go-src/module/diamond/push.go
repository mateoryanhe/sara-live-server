package diamond

import (
	"fmt"
	"math"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dto/diamonddto"
)

// pushDiamondToApp 推送用户最新钻石余额到 App 端
func pushDiamondToApp(userId uint64, diamond float64) {
	push.Data(userId, cmd.DiamondPush, &diamonddto.DiamondPushItem{
		Diamond: formatCurrencyBalance(diamond),
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
