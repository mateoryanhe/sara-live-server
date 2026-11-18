package auth

import (
	"context"

	"xr-game-server/gameevent"
)

func publishAppLoginEvent(ctx context.Context, userId uint64) {
	gameevent.PublishLoginEvent(ctx, userId)
}