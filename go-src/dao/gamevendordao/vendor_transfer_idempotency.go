package gamevendordao

import (
	"context"
	"strings"

	"xr-game-server/core/cache"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/gamewindao"
)

const vendorTransferIdempotencyCacheKeyPrefix = "vendor_transfer:done:"

var vendorTransferIdempotencyCache *cache.CacheMgr

func initVendorTransferIdempotencyCache() {
	vendorTransferIdempotencyCache = cache.NewCacheMgr()
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
	v := vendorTransferIdempotencyCache.GetData(vendorTransferIdempotencyCacheKey(transactionID), func(ctx context.Context) (interface{}, error) {
		return existsVendorTransferTransactionInDB(transactionID), nil
	})
	processed, _ := v.(bool)
	return processed
}

// MarkVendorTransferProcessed 处理成功后刷新缓存, 后续 GetData 直接命中.
func MarkVendorTransferProcessed(transactionID string) {
	transactionID = strings.TrimSpace(transactionID)
	if transactionID == "" || vendorTransferIdempotencyCache == nil {
		return
	}
	vendorTransferIdempotencyCache.FlushCache(vendorTransferIdempotencyCacheKey(transactionID), true)
}
