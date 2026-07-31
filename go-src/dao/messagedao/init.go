package messagedao

func Init() {
	initPrivateMessageDao()
	initMessageUnreadDao()
	initMessageUnreadDetailDao()
	initSystemMessageUnreadDao()
	initActivityMessageDao()
	initUserActivityMessageDao()
}
