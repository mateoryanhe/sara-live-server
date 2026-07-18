package logquery

import (
	"runtime"

	"xr-game-server/errercode"
)

func ensureLinuxLogQuery() error {
	if runtime.GOOS != "linux" {
		return errercode.CreateCodeAndParam(errercode.InvalidParam, "日志查询仅支持 Linux 平台")
	}
	return nil
}
