package game

import (
	entity "xr-game-server/entity/game"
)

// VendorGame CMS/上架使用的第三方游戏视图.
type VendorGame struct {
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
	NameEn   string `json:"nameEn"`
	Category string `json:"category"`
	Cover    string `json:"cover"`
	Platform string `json:"platform"`
}

func toVendorGame(row *entity.VendorGameLib) *VendorGame {
	if row == nil {
		return nil
	}
	return &VendorGame{
		GameCode: row.GameCode,
		Name:     row.Name,
		NameEn:   row.NameEn,
		Category: row.Category,
		Cover:    row.Cover,
		Platform: row.Platform,
	}
}

func toVendorGameLibEntity(v *VendorGame) *entity.VendorGameLib {
	if v == nil {
		return nil
	}
	return &entity.VendorGameLib{
		GameCode: v.GameCode,
		Name:     v.Name,
		NameEn:   v.NameEn,
		Category: v.Category,
		Cover:    v.Cover,
		Platform: v.Platform,
	}
}
