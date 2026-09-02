package fxrate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	baseCurrencyUSD       = "USD"
	fxRateCacheTTL        = 15 * time.Minute
	fxRateRequestTimeout  = 8 * time.Second
	frankfurterLatestURL  = "https://api.frankfurter.app/latest?from=USD&to=%s"
	openExchangeRateAPIURL = "https://open.er-api.com/v6/latest/USD"
)

const (
	sourceFrankfurter = "frankfurter"
	sourceOpenERAPI   = "open.er-api"
	sourceStaleCache  = "stale-cache"
)

// QuoteRate USD 兑目标法币汇率快照
type QuoteRate struct {
	Base           string
	Quote          string
	MarketRate     float64
	Rate           float64
	AdjustPercent  float64
	InverseRate    float64
	Source         string
	RateDate       string
	Cached         bool
	CacheExpiresAt int64
}

var (
	fxRateCache     = gcache.New()
	fxRateStaleMu   sync.RWMutex
	fxRateStaleData = make(map[string]marketRateSnapshot)
)

type marketRateSnapshot struct {
	Rate     float64
	Source   string
	RateDate string
}

type frankfurterResp struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

type openERAPIResp struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	TimeLastUpdate  int64              `json:"time_last_update_unix"`
	Rates           map[string]float64 `json:"rates"`
}

func normalizeCurrencyCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func cacheKey(currency string) string {
	return baseCurrencyUSD + ":" + currency
}

// FetchMarketUsdToQuoteRate 拉取市场参考汇率(1 USD = X Quote),带 15 分钟缓存与失败降级
func FetchMarketUsdToQuoteRate(ctx context.Context, currencyCode string) (*QuoteRate, error) {
	quote := normalizeCurrencyCode(currencyCode)
	if quote == "" {
		return nil, fmt.Errorf("empty currency code")
	}
	if quote == baseCurrencyUSD {
		now := time.Now()
		return &QuoteRate{
			Base:       baseCurrencyUSD,
			Quote:      quote,
			MarketRate: 1,
			Rate:       1,
			InverseRate: 1,
			Source:     sourceFrankfurter,
			RateDate:   now.Format("2006-01-02"),
			Cached:     false,
		}, nil
	}

	if ctx == nil {
		ctx = gctx.New()
	}
	key := cacheKey(quote)
	if cached, ok := getCachedQuoteRate(ctx, key); ok {
		return cached, nil
	}

	snap, source, err := fetchMarketRateFromProviders(ctx, quote)
	if err != nil {
		if stale, ok := getStaleMarketRate(quote); ok {
			g.Log().Warningf(ctx, "fx rate fallback to stale cache, quote=%s err=%v", quote, err)
			rate := buildQuoteRate(quote, stale.Rate, 0, stale.Source, stale.RateDate, true, 0)
			rate.Source = sourceStaleCache + "(" + stale.Source + ")"
			return rate, nil
		}
		return nil, err
	}

	saveStaleMarketRate(quote, snap)
	expiresAt := time.Now().Add(fxRateCacheTTL).Unix()
	rate := buildQuoteRate(quote, snap.Rate, 0, source, snap.RateDate, false, expiresAt)
	setCachedQuoteRate(ctx, key, rate)
	return rate, nil
}

// ApplyAdjustPercent 在 market rate 上应用加点比例(%)
func ApplyAdjustPercent(marketRate, adjustPercent float64) float64 {
	if marketRate <= 0 {
		return 0
	}
	return marketRate * (1 + adjustPercent/100)
}

func buildQuoteRate(quote string, marketRate, adjustPercent float64, source, rateDate string, cached bool, cacheExpiresAt int64) *QuoteRate {
	finalRate := ApplyAdjustPercent(marketRate, adjustPercent)
	inverse := 0.0
	if finalRate > 0 {
		inverse = 1 / finalRate
	}
	return &QuoteRate{
		Base:           baseCurrencyUSD,
		Quote:          quote,
		MarketRate:     marketRate,
		AdjustPercent:  adjustPercent,
		Rate:           finalRate,
		InverseRate:    inverse,
		Source:         source,
		RateDate:       rateDate,
		Cached:         cached,
		CacheExpiresAt: cacheExpiresAt,
	}
}

func getCachedQuoteRate(ctx context.Context, key string) (*QuoteRate, bool) {
	v, err := fxRateCache.Get(ctx, key)
	if err != nil || v.IsNil() {
		return nil, false
	}
	rate, ok := v.Val().(*QuoteRate)
	if !ok || rate == nil {
		return nil, false
	}
	copyRate := *rate
	copyRate.Cached = true
	return &copyRate, true
}

func setCachedQuoteRate(ctx context.Context, key string, rate *QuoteRate) {
	if rate == nil {
		return
	}
	copyRate := *rate
	copyRate.Cached = false
	_ = fxRateCache.Set(ctx, key, &copyRate, fxRateCacheTTL)
}

func saveStaleMarketRate(quote string, snap marketRateSnapshot) {
	fxRateStaleMu.Lock()
	defer fxRateStaleMu.Unlock()
	fxRateStaleData[quote] = snap
}

func getStaleMarketRate(quote string) (marketRateSnapshot, bool) {
	fxRateStaleMu.RLock()
	defer fxRateStaleMu.RUnlock()
	snap, ok := fxRateStaleData[quote]
	return snap, ok && snap.Rate > 0
}

// ReloadExchangeRateCache 清除汇率缓存; currencyCode 为空则清除全部
func ReloadExchangeRateCache(currencyCode string) {
	ctx := gctx.New()
	code := normalizeCurrencyCode(currencyCode)
	if code == "" {
		keys, err := fxRateCache.Keys(ctx)
		if err == nil {
			_, _ = fxRateCache.Remove(ctx, keys...)
		}
		return
	}
	_, _ = fxRateCache.Remove(ctx, cacheKey(code))
}

func fetchMarketRateFromProviders(ctx context.Context, quote string) (marketRateSnapshot, string, error) {
	if snap, err := fetchFrankfurterRate(ctx, quote); err == nil {
		return snap, sourceFrankfurter, nil
	} else {
		g.Log().Warningf(ctx, "frankfurter fx rate failed, quote=%s err=%v", quote, err)
	}
	if snap, err := fetchOpenERAPIRate(ctx, quote); err == nil {
		return snap, sourceOpenERAPI, nil
	} else {
		g.Log().Warningf(ctx, "open.er-api fx rate failed, quote=%s err=%v", quote, err)
	}
	return marketRateSnapshot{}, "", fmt.Errorf("fx rate unavailable for %s", quote)
}

func fetchFrankfurterRate(ctx context.Context, quote string) (marketRateSnapshot, error) {
	client := gclient.New().SetTimeout(fxRateRequestTimeout)
	resp, err := client.Get(ctx, fmt.Sprintf(frankfurterLatestURL, quote))
	if err != nil {
		return marketRateSnapshot{}, err
	}
	defer resp.Close()

	body := resp.ReadAll()
	var payload frankfurterResp
	if err := json.Unmarshal(body, &payload); err != nil {
		return marketRateSnapshot{}, err
	}
	rate, ok := payload.Rates[quote]
	if !ok || rate <= 0 {
		return marketRateSnapshot{}, fmt.Errorf("frankfurter missing rate for %s", quote)
	}
	return marketRateSnapshot{Rate: rate, Source: sourceFrankfurter, RateDate: payload.Date}, nil
}

func fetchOpenERAPIRate(ctx context.Context, quote string) (marketRateSnapshot, error) {
	client := gclient.New().SetTimeout(fxRateRequestTimeout)
	resp, err := client.Get(ctx, openExchangeRateAPIURL)
	if err != nil {
		return marketRateSnapshot{}, err
	}
	defer resp.Close()

	body := resp.ReadAll()
	var payload openERAPIResp
	if err := json.Unmarshal(body, &payload); err != nil {
		return marketRateSnapshot{}, err
	}
	if payload.Result != "success" {
		return marketRateSnapshot{}, fmt.Errorf("open.er-api result=%s", payload.Result)
	}
	rate, ok := payload.Rates[quote]
	if !ok || rate <= 0 {
		return marketRateSnapshot{}, fmt.Errorf("open.er-api missing rate for %s", quote)
	}
	rateDate := time.Unix(payload.TimeLastUpdate, 0).UTC().Format("2006-01-02")
	return marketRateSnapshot{Rate: rate, Source: sourceOpenERAPI, RateDate: rateDate}, nil
}

// GetUsdToQuoteRate 服务端查询最终汇率(含 CMS 加点),要求币种已在 CMS 启用
func GetUsdToQuoteRate(ctx context.Context, currencyCode string, adjustPercent float64) (*QuoteRate, error) {
	market, err := FetchMarketUsdToQuoteRate(ctx, currencyCode)
	if err != nil {
		return nil, err
	}
	quote := normalizeCurrencyCode(currencyCode)
	final := buildQuoteRate(quote, market.MarketRate, adjustPercent, market.Source, market.RateDate, market.Cached, market.CacheExpiresAt)
	return final, nil
}

// ConvertUsdToQuote 将 USD 金额换算为目标法币(使用最终汇率)
func ConvertUsdToQuote(usdAmount float64, rate *QuoteRate) float64 {
	if usdAmount <= 0 || rate == nil || rate.Rate <= 0 {
		return 0
	}
	return usdAmount * rate.Rate
}
