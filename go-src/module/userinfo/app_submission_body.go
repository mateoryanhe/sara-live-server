package userinfo

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

const appSubmissionMaxContentBytes = 256 * 1024

func readAppSubmissionBody(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	raw := r.GetBodyString()
	if len(raw) > appSubmissionMaxContentBytes {
		raw = raw[:appSubmissionMaxContentBytes]
	}
	return raw
}
