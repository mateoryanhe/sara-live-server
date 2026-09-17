package recharge

import (
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dto/haipaydto"
)

func normalizeHaiPayCollectionSelections(
	countryCode, currencyCode string,
	selections []*haipaydto.CollectionPaymentMethodSelection,
) ([]*haipaydto.CollectionPaymentMethodSelection, bool) {
	currencyCode = strings.ToUpper(strings.TrimSpace(currencyCode))
	if len(selections) > 1 {
		return nil, false
	}
	result := make([]*haipaydto.CollectionPaymentMethodSelection, 0, len(selections))
	seen := make(map[string]struct{}, len(selections))
	for _, selection := range selections {
		if selection == nil {
			continue
		}
		payType := strings.ToUpper(strings.TrimSpace(selection.PayType))
		inBankCode := strings.TrimSpace(selection.InBankCode)
		selectionCurrency := strings.ToUpper(strings.TrimSpace(selection.CurrencyCode))
		if selectionCurrency != currencyCode {
			return nil, false
		}
		if payType == "" && inBankCode == "" {
			continue
		}
		if payType == "" || inBankCode == "" {
			return nil, false
		}
		canonicalCode, ok := country.ResolveHaiPayCollectionPaymentMethodCode(
			countryCode, currencyCode, payType, inBankCode,
		)
		if !ok || canonicalCode == "" {
			return nil, false
		}
		key := payType + "\x00" + strings.ToUpper(canonicalCode)
		if _, ok = seen[key]; ok {
			continue
		}
		// 一个国家的同一代收币种只能选一个支付方式。USD 与本地币分别保存，互不覆盖。
		if len(result) > 0 {
			return nil, false
		}
		seen[key] = struct{}{}
		result = append(result, &haipaydto.CollectionPaymentMethodSelection{
			CurrencyCode: currencyCode, PayType: payType, InBankCode: canonicalCode,
		})
	}
	return result, true
}
