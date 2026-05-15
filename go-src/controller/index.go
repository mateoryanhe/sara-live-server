package controller

import (
	"xr-game-server/core/httpserver"
)

type IndexController struct {
}

func initIndex() {
	httpserver.RegAPI("/index", new(IndexController))
}
