package entity

const (
	// VoiceChatHostSeatIndex 语聊房主播麦位(固定 0 号位)
	VoiceChatHostSeatIndex = 0
	// VoiceChatAudienceSeatCount 语聊房观众麦位数(1~8)
	VoiceChatAudienceSeatCount = 8
	// VoiceChatSeatCount 语聊房麦位总数(含主播位)
	VoiceChatSeatCount = VoiceChatAudienceSeatCount + 1
	// VoiceChatMicApplyListMax 单场申请上麦排队上限
	VoiceChatMicApplyListMax = 500
)
