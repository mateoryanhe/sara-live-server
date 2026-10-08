package entity

const (
	// VoiceChatMicSeatStatusIdle 空闲(可上麦)
	VoiceChatMicSeatStatusIdle uint8 = 1
	// VoiceChatMicSeatStatusLocked 锁麦(空位,观众不可占)
	VoiceChatMicSeatStatusLocked uint8 = 2
	// VoiceChatMicSeatStatusOnMic 在麦(可说话)
	VoiceChatMicSeatStatusOnMic uint8 = 3
	// VoiceChatMicSeatStatusMicMuted 在麦禁音(不可说话)
	VoiceChatMicSeatStatusMicMuted uint8 = 4
)

// IsValidVoiceChatMicSeatStatus 麦位状态是否合法.
func IsValidVoiceChatMicSeatStatus(status uint8) bool {
	return status == VoiceChatMicSeatStatusIdle ||
		status == VoiceChatMicSeatStatusLocked ||
		status == VoiceChatMicSeatStatusOnMic ||
		status == VoiceChatMicSeatStatusMicMuted
}

// IsVoiceChatMicSeatLocked 空麦位是否已锁.
func IsVoiceChatMicSeatLocked(seat *VoiceChatMicSeat) bool {
	return seat != nil && seat.Status == VoiceChatMicSeatStatusLocked
}

// IsVoiceChatMicSeatMicMuted 该麦位是否禁音.
func IsVoiceChatMicSeatMicMuted(seat *VoiceChatMicSeat) bool {
	return seat != nil && seat.Status == VoiceChatMicSeatStatusMicMuted
}

// VoiceChatMicSeatStatusAllowsAudienceTake 观众是否可占用该麦位.
func VoiceChatMicSeatStatusAllowsAudienceTake(status uint8) bool {
	return status == VoiceChatMicSeatStatusIdle
}
