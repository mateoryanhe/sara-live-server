package luckygiftpoolcfg

import (
	"time"

	"xr-game-server/core/syndb"
	"xr-game-server/errercode"
)

const luckyGiftPoolSyndbFlushTimeout = 3 * time.Second

func flushLuckyGiftPoolSyndb() error {
	if !syndb.FlushUntilIdle(luckyGiftPoolSyndbFlushTimeout) {
		return errercode.CreateCode(errercode.SysError)
	}
	return nil
}
