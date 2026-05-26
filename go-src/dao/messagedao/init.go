package messagedao

func Init() {
	initPrivateMessageDao()
	initSystemMessageDao()
	initMessageUnreadDao()
	initMessageUnreadDetailDao()
}
