package datasync

import "xr-game-server/errercode"

func errInvalidParam() error {
	return errercode.CreateCode(errercode.InvalidParam)
}
