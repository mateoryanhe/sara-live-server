package gamevendordao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"strings"

	"xr-game-server/core/cache"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/gamewindao"
)

const vendorTransferIdempotencyCacheKeyPrefix = "vendor_transfer:done:"

var vendorTransferIdempotencyCache *cache.RowCache[bool]

func initVendorTransferIdempotencyCache() {
	vendorTransferIdempotencyCache = cache.NewRowCache[bool]()
}

func vendorTransferIdempotencyCacheKey(transactionID string) string {
	return vendorTransferIdempotencyCacheKeyPrefix + transactionID
}

func existsVendorTransferTransactionInDB(transactionID string) bool {
	return gamebetdao.ExistsGameBetLogByOrderId(transactionID) ||
		gamewindao.ExistsGameWinLogByOrderId(transactionID) ||
		currencylogdao.ExistsByTransactionId(transactionID)
}

// IsVendorTransferProcessed 判断 transaction_id 是否已处理; 缓存未命中时查 DB.
func IsVendorTransferProcessed(transactionID string) bool {
	transactionID = strings.TrimSpace(transactionID)
	if transactionID == "" || vendorTransferIdempotencyCache == nil {
		return false
	}
	v := vendorTransferIdempotencyCache.MustGetRow(gctx.New(), vendorTransferIdempotencyCacheKey(transactionID), func(ctx context.Context) (bool, error) {
		return existsVendorTransferTransactionInDB(transactionID), nil
	})
	return v
}

// MarkVendorTransferProcessed 处理成功后刷新缓存, 后续 GetData 直接命中.
func MarkVendorTransferProcessed(transactionID string) {
	transactionID = strings.TrimSpace(transactionID)
	if transactionID == "" || vendorTransferIdempotencyCache == nil {
		return
	}
	vendorTransferIdempotencyCache.PublishRow(gctx.New(), vendorTransferIdempotencyCacheKey(transactionID), true)
}
