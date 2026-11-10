package gift

import "math/rand"

// NormalizeLuckyGiftRtp 将 RTP 归一到 [0,100]；RTP 表示单次送礼的中奖概率(%)，不是期望返还比。
func NormalizeLuckyGiftRtp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// RollLuckyGiftWin 按 RTP(%) 判定本次是否中奖。
// 流程第 1 步：未中奖则不进入挡位抽取、不返奖；中奖后再按挡位权重抽倍率。
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
