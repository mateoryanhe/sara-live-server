package incomesettlement

import (
	"testing"

	"xr-game-server/dto/incomesettlementdto"
	"xr-game-server/entity/live"
)

func TestGuildTypeFilterForCMSList(t *testing.T) {
	normal := entity.LiveGuildTypeNormal
	coin := entity.LiveGuildTypeCoinMerchant
	tests := []struct {
		name string
		req  *incomesettlementdto.CMSGuildIncomeSettlementLogListReq
		want *uint8
	}{
		{
			name: "explicit guild type 0",
			req:  &incomesettlementdto.CMSGuildIncomeSettlementLogListReq{GuildType: &normal},
			want: &normal,
		},
		{
			name: "normal guild only flag",
			req:  &incomesettlementdto.CMSGuildIncomeSettlementLogListReq{NormalGuildOnly: true},
			want: &normal,
		},
		{
			name: "coin merchant only flag",
			req:  &incomesettlementdto.CMSGuildIncomeSettlementLogListReq{CoinMerchantGuildOnly: true},
			want: &coin,
		},
		{
			name: "no filter",
			req:  &incomesettlementdto.CMSGuildIncomeSettlementLogListReq{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := guildTypeFilterForCMSList(tt.req)
			if (got == nil) != (tt.want == nil) {
				t.Fatalf("got %v want %v", got, tt.want)
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Fatalf("got %d want %d", *got, *tt.want)
			}
		})
	}
}
