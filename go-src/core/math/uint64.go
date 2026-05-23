package math

import "math"

// Add 两个 uint64 安全相加,发生溢出时返回 math.MaxUint64
func Add(a, b uint64) uint64 {
	if b > math.MaxUint64-a {
		return math.MaxUint64
	}
	return a + b
}

// AddMax 计算 current 在不超过 max 时最多还能加多少
func AddMax(current, add, max uint64) uint64 {
	ret := Add(current, add)
	if ret >= max {
		return max
	}
	return ret
}
