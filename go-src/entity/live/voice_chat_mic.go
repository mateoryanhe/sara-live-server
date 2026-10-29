package entity

const (
	// VoiceChatHostSeatIndex 语聊房主播麦位(固定 0 号位)
	VoiceChatHostSeatIndex = 0
	// VoiceChatAudienceSeatCount 语聊房观众麦位数上限(1~24,与 Max 配合)
	VoiceChatAudienceSeatCount = 24
	// MaxVoiceChatMicSeatCount 语聊房麦位总数上限(含主播位 0 号)
	MaxVoiceChatMicSeatCount = VoiceChatAudienceSeatCount + 1
	// DefaultVoiceChatMicSeatCount 未配置 live_room_cfgs.voice_chat_mic_seat_count 时的默认麦位总数
	DefaultVoiceChatMicSeatCount = MaxVoiceChatMicSeatCount
	// MinVoiceChatMicSeatCount 最少麦位(主播位 + 至少 1 个观众位)
	MinVoiceChatMicSeatCount = 2
	// VoiceChatSeatCount 兼容旧代码:等于默认/上限麦位总数
	VoiceChatSeatCount = MaxVoiceChatMicSeatCount
	// VoiceChatMicApplyListMax 单场申请上麦排队上限
	VoiceChatMicApplyListMax = 500
)
