package currency

// 商业类型
const (
	BusinessTypeSocial uint8 = 1 // 社交
	BusinessTypeGame   uint8 = 2 // 游戏
)

// BusinessTypeText 返回商业类型中文文案
func BusinessTypeText(v uint8) string {
	switch v {
	case BusinessTypeGame:
		return "游戏"
	case BusinessTypeSocial:
		return "社交"
	default:
		return "社交"
	}
}
