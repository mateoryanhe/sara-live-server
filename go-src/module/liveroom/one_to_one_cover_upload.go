package liveroom

import (
	"strings"

	"xr-game-server/dto/liveroomdto"
	"xr-game-server/module/upload"
)

func BuildOneToOneCoverUploadRes(storedName string) *liveroomdto.UploadOneToOneCoverRes {
	name := strings.TrimSpace(storedName)
	return &liveroomdto.UploadOneToOneCoverRes{
		Cover: name,
		Url:   upload.GetUrlByName(name),
	}
}