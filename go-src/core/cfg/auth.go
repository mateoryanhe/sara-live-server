package cfg

var authCfg = &AuthCfg{}

type AuthCfg struct {
	LoginOff bool
}

func GetAuthCfg() *AuthCfg {
	return authCfg
}
