package gender

const (
	Unknown uint8 = 0 // 未知
	Male    uint8 = 1 // 男
	Female  uint8 = 2 // 女
)

// IsValid 是否为合法性别值
func IsValid(v uint8) bool {
	return v == Unknown || v == Male || v == Female
}
