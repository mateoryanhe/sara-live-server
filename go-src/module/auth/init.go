package auth

func InitAuth() {
	initPhoneLoginGuard()
	initAppCancelGuard()
	initAppToken()
	initCmsToken()
	initUserMaxId()
}
