package math

import (
	"math"
	"xr-game-server/constants/common"
)

func ChkUInt64(current uint64, add uint64) uint64 {
	if current == math.MaxUint64 {
		return common.Zero
	}
	diff := math.MaxUint64 - current
	if add >= diff {
		return diff
	} else {
		return add
	}
}

func ChkMax(current uint64, add uint64, max uint64) uint64 {
	if current == max {
		return common.Zero
	}
	diff := max - current
	if add >= diff {
		return diff
	} else {
		return add
	}
}
