package auth

func InitAuth() {
	initPhoneLoginGuard()
	initAppToken()
	initCmsToken()
}
