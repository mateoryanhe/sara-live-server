package gift

import "math/rand"

func NormalizeLuckyGiftRtp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// RollLuckyGiftWin 按礼物 RTP(%) 判定本次是否中奖；未中奖则不抽倍数、不返奖
func RollLuckyGiftWin(rtpPercent float64) bool {
	rtp := NormalizeLuckyGiftRtp(rtpPercent)
	if rtp <= 0 {
		return false
	}
	if rtp >= 100 {
		return true
	}
	return rand.Float64()*100 < rtp
}
