package wallet

import (
	"testing"
	"time"

	"xr-game-server/dao/currencylogdao"
)

func TestNormalizeCoinMerchantTransferRecordPage(t *testing.T) {
	if page, size := normalizeCoinMerchantTransferRecordPage(0, 0); page != 1 || size != 20 {
		t.Fatalf("default page=%d size=%d", page, size)
	}
	if page, size := normalizeCoinMerchantTransferRecordPage(3, 101); page != 3 || size != 100 {
		t.Fatalf("limited page=%d size=%d", page, size)
	}
}

func TestBuildCoinMerchantTransferRecordItems(t *testing.T) {
	createdAt := time.Unix(1_700_000_000, 0)
	rows := []*currencylogdao.CoinMerchantGoldTransferListRow{
		{
			ID:             987654321,
			TargetUserId:   200,
			TargetNickname: "target-user",
			Amount:         12.34,
			CreatedAt:      createdAt,
		},
	}

	items := buildCoinMerchantTransferRecordItems(rows)
	if len(items) != 1 {
		t.Fatalf("items=%d", len(items))
	}
	item := items[0]
	if item.Id != "987654321" || item.TargetUserId != "200" || item.TargetNickname != "target-user" ||
		item.Amount != 12.34 || item.CreatedAt != createdAt.Unix() || item.TargetAvatar == "" ||
		item.CreatedAtText != createdAt.In(time.Local).Format(coinMerchantTransferRecordTimeLayout) {
		t.Fatalf("unexpected item=%+v", item)
	}
}

func TestParseCoinMerchantTransferTargetUserId(t *testing.T) {
	if id, err := parseCoinMerchantTransferTargetUserId(" 123 "); err != nil || id != 123 {
		t.Fatalf("id=%d err=%v", id, err)
	}
	if id, err := parseCoinMerchantTransferTargetUserId(""); err != nil || id != 0 {
		t.Fatalf("empty id=%d err=%v", id, err)
	}
	if _, err := parseCoinMerchantTransferTargetUserId("abc"); err == nil {
		t.Fatal("invalid target user id should fail")
	}
}

func TestValidateCoinMerchantTransferRecordTimeRange(t *testing.T) {
	if err := validateCoinMerchantTransferRecordTimeRange(100, 200); err != nil {
		t.Fatalf("valid range err=%v", err)
	}
	if err := validateCoinMerchantTransferRecordTimeRange(200, 100); err == nil {
		t.Fatal("reversed range should fail")
	}
	if err := validateCoinMerchantTransferRecordTimeRange(-1, 0); err == nil {
		t.Fatal("negative range should fail")
	}
}
