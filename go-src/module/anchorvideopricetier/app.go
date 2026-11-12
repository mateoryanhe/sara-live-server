package anchorvideopricetier

import (
	"context"

	"xr-game-server/dto/anchorvideopricetierdto"
)

// GetAppPriceTierList App 查询已上架挡位:只读内存缓存,不查库.
func GetAppPriceTierList(_ context.Context, _ *anchorvideopricetierdto.AppPriceTierListReq) (*anchorvideopricetierdto.AppPriceTierListRes, error) {
	cached := GetOnShelfPriceTiers()
	list := make([]*anchorvideopricetierdto.AppPriceTierItem, 0, len(cached))
	for _, item := range cached {
		if item == nil {
			continue
		}
		list = append(list, &anchorvideopricetierdto.AppPriceTierItem{
			ID:    item.ID,
			Price: item.Price,
		})
	}
	return &anchorvideopricetierdto.AppPriceTierListRes{List: list}, nil
}
