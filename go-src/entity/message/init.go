package entity

// Init 用户消息相关表迁移与 syndb 注册
func Init() {
	initActivityMessage()
	initUserActivityMessage()
	initUserMessage()
	initUserMessageSession()
	initUserMessageUnread()
	initUserMessageUnreadDetail()
	initUserSystemMessageUnread()
	initUserPersonalSystemMessage()
}
