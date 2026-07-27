package liveroom

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/accountdto"
)

// QueryAnchorList CMS分页查询主播列表(基于 roomListCache)
func QueryAnchorList(_ context.Context, req *accountdto.QueryAnchorListReq) (*httpserver.CMSQueryResp, error) {
	total, data := queryAnchorListFromMemory(req)
	return httpserver.NewCMSQueryResp(total, data), nil
}
