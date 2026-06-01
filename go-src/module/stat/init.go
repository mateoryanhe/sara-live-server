package stat

import "sync"

func Init() {
	initLoginEvent()
	initRegisterEvent()
}
