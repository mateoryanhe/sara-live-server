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

// QueryAnchorListByGuildIds CMS分页查询指定多个工会名下主播列表
func QueryAnchorListByGuildIds(_ context.Context, guildIds []uint64, req *accountdto.QueryAnchorListReq) (*httpserver.CMSQueryResp, error) {
	total, data := queryAnchorListByGuildIdsFromMemory(guildIds, req)
	return httpserver.NewCMSQueryResp(total, data), nil
}
