package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gamebetdto"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/module/game"
)

const GameAppUrl = "/game"

type GameAppController struct{}

func initGameAppController() {
	httpserver.RegAPI(GameAppUrl, &GameAppController{})
}

// AppGameList App 分页查询已上架游戏列表
func (c *GameAppController) AppGameList(ctx context.Context, req *gameplatformdto.AppGameListReq) (*gameplatformdto.AppGameListRes, error) {
	return game.GetAppGameList(ctx, req)
}

// BrowseTitles Sofie 标题列表(同 AppGameList 数据;应答字段不用 game 字眼).
func (c *GameAppController) BrowseTitles(ctx context.Context, req *gameplatformdto.BrowseTitlesReq) (*gameplatformdto.BrowseTitlesRes, error) {
	pageIndex, pageSize := 0, 0
	if req != nil {
		pageIndex, pageSize = req.PageIndex, req.PageSize
	}
	old, err := game.GetAppGameList(ctx, &gameplatformdto.AppGameListReq{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, err
	}
	if old == nil {
		return &gameplatformdto.BrowseTitlesRes{List: []*gameplatformdto.BrowseTitlesItem{}}, nil
	}
	list := make([]*gameplatformdto.BrowseTitlesItem, 0, len(old.List))
	for _, row := range old.List {
		if row == nil {
			continue
		}
		list = append(list, &gameplatformdto.BrowseTitlesItem{
			TitleCode: row.GameCode,
			NameEn:    row.NameEn,
			Cover:     row.Cover,
			Category:  row.Category,
			Platform:  row.Platform,
		})
	}
	return &gameplatformdto.BrowseTitlesRes{
		Total:     old.Total,
		PageIndex: old.PageIndex,
		PageSize:  old.PageSize,
		List:      list,
	}, nil
}

// AppGameStart App 获取游戏启动链接
func (c *GameAppController) AppGameStart(ctx context.Context, req *gameplatformdto.AppGameStartReq) (*gameplatformdto.AppGameStartRes, error) {
	return game.GetAppGameStartLink(ctx, req)
}

// OpenTitle Sofie 打开标题(同 AppGameStart;请求字段用 titleCode).
func (c *GameAppController) OpenTitle(ctx context.Context, req *gameplatformdto.OpenTitleReq) (*gameplatformdto.OpenTitleRes, error) {
	titleCode, roomId, urlParams := "", uint64(0), ""
	if req != nil {
		titleCode = req.TitleCode
		roomId = req.RoomId
		urlParams = req.UrlParams
	}
	old, err := game.GetAppGameStartLink(ctx, &gameplatformdto.AppGameStartReq{
		GameCode:  titleCode,
		RoomId:    roomId,
		UrlParams: urlParams,
	})
	if err != nil {
		return nil, err
	}
	if old == nil {
		return &gameplatformdto.OpenTitleRes{}, nil
	}
	return &gameplatformdto.OpenTitleRes{Link: old.Link}, nil
}

// AppGameBetList App 分页查询游戏下注记录
func (c *GameAppController) AppGameBetList(ctx context.Context, req *gamebetdto.AppGameBetListReq) (*gamebetdto.AppGameBetListRes, error) {
	return game.GetAppBetList(ctx, req)
}

// AppGameConsumeRank App 分页查询单场直播游戏消费榜
func (c *GameAppController) AppGameConsumeRank(ctx context.Context, req *gamebetdto.AppGameConsumeRankReq) (*gamebetdto.AppGameConsumeRankRes, error) {
	return game.GetAppGameConsumeRank(ctx, req)
}
