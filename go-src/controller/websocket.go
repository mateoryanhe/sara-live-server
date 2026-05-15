package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/websocketdto"
)

type WebSocket struct {
}

func initWebSocket() {
	httpserver.RegCMS("/ws", new(WebSocket))
}

func (s *WebSocket) BasePush(ctx context.Context, req *websocketdto.WebSocketReq) (res []*httpserver.PushResp, err error) {
	return nil, err
}
func (s *WebSocket) Hear(ctx context.Context, req *websocketdto.HeartReq) (res *websocketdto.HeartRes, err error) {
	return &websocketdto.HeartRes{}, nil
}

func (s *WebSocket) Auth(ctx context.Context, req *websocketdto.WebSocketAuthReq) (res *websocketdto.WebSocketAuthRes, err error) {
	return &websocketdto.WebSocketAuthRes{}, nil
}
