package entity

// HaiPayBizType identifies HaiPay business flows. Storage may be scenario-specific;
// coin-merchant collection uses HaiPayCoinMerchantCollectionCfg and its own cache.
type HaiPayBizType uint8

const (
	HaiPayBizTypeNormalCollection HaiPayBizType = iota + 1
	HaiPayBizTypeCoinMerchantCollection
	HaiPayBizTypeNormalPayout
	HaiPayBizTypeGuildPayout
)
