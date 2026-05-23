package math

import stdmath "math"

const float64DecimalScale = 10000 // 10^4,保留4位小数

var maxFloat64With4Decimals = stdmath.Round(float64(stdmath.MaxInt64)) / float64DecimalScale

// AddFloat64 两个 float64 安全相加,结果四舍五入保留4位小数,溢出时返回可表示的最大值
func AddFloat64(a, b float64) float64 {
	if !isValidFloat64(a) || !isValidFloat64(b) {
		return 0
	}

	aScaled := toScaledInt64(a)
	bScaled := toScaledInt64(b)

	if bScaled > 0 && aScaled > stdmath.MaxInt64-bScaled {
		return maxFloat64With4Decimals
	}
	if bScaled < 0 && aScaled < stdmath.MinInt64-bScaled {
		return fromScaledInt64(stdmath.MinInt64)
	}

	return fromScaledInt64(aScaled + bScaled)
}

// SubFloat64 两个 float64 安全相减(a-b),结果四舍五入保留4位小数,溢出时返回边界值
func SubFloat64(a, b float64) float64 {
	if !isValidFloat64(a) || !isValidFloat64(b) {
		return 0
	}

	aScaled := toScaledInt64(a)
	bScaled := toScaledInt64(b)

	if bScaled < 0 && aScaled > stdmath.MaxInt64+bScaled {
		return maxFloat64With4Decimals
	}
	if bScaled > 0 && aScaled < stdmath.MinInt64+bScaled {
		return fromScaledInt64(stdmath.MinInt64)
	}

	return fromScaledInt64(aScaled - bScaled)
}

func isValidFloat64(v float64) bool {
	return !stdmath.IsNaN(v) && !stdmath.IsInf(v, 0)
}

func toScaledInt64(v float64) int64 {
	if v >= float64(stdmath.MaxInt64)/float64DecimalScale {
		return stdmath.MaxInt64
	}
	if v <= float64(stdmath.MinInt64)/float64DecimalScale {
		return stdmath.MinInt64
	}
	return int64(stdmath.Round(v * float64DecimalScale))
}

func fromScaledInt64(scaled int64) float64 {
	return stdmath.Round(float64(scaled)) / float64DecimalScale
}
