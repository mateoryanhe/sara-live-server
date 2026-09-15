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
	baseCurrencyUSD        = "USD"
	fxRateCacheTTL         = time.Minute
	fxRateFailureTTL       = 5 * time.Minute
	fxRateStaleTTL         = 24 * time.Hour
	fxRateRequestTimeout   = 8 * time.Second
	frankfurterRateURL     = "https://api.frankfurter.dev/v2/rate/USD/%s"
	openExchangeRateAPIURL = "https://open.er-api.com/v6/latest/USD"
	sourceFrankfurter      = "frankfurter"
	sourceOpenERAPI        = "open.er-api"
	sourceStaleCache       = "stale-cache"
)

// ConversionResult 是 USD 金额转换为订单目标币种后的结果。
// 下单模块只依赖这个结构，外部汇率供应商的请求和响应变化由 fxrate 内部消化。
type ConversionResult struct {
	SourceCurrency string
	TargetCurrency string
	SourceAmount   float64
	TargetAmount   float64
	Rate           float64
	Source         string
	RateDate       string
	Cached         bool
	CacheExpiresAt int64
}

// quoteRate 表示 1 USD 可兑换的目标法币汇率，仅供 fxrate 内部使用。
type quoteRate struct {
	Base           string
	Quote          string
	Rate           float64
	Source         string
	RateDate       string
	Cached         bool
	CacheExpiresAt int64
}

type marketRateSnapshot struct {
	Rate      float64
	Source    string
	RateDate  string
	FetchedAt time.Time
}

type failureSnapshot struct {
	Message string
	RetryAt time.Time
}

type frankfurterResp struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}

type openERAPIResp struct {
	Result         string             `json:"result"`
	BaseCode       string             `json:"base_code"`
	TimeLastUpdate int64              `json:"time_last_update_unix"`
	Rates          map[string]float64 `json:"rates"`
}

var (
	fxRateCache = gcache.New()

	fxRateStaleMu   sync.RWMutex
	fxRateStaleData = make(map[string]marketRateSnapshot)

	fxRateFailureMu   sync.RWMutex
	fxRateFailureData = make(map[string]failureSnapshot)

	fxRateFetchLocks  sync.Map
	marketRateFetcher = fetchMarketRateFromProviders
)

// ConvertUSD 将 USD 金额换算为订单目标币种。调用方只需传入 USD 金额和订单币种，
// 币种规范化、汇率查询、缓存及金额换算全部由 fxrate 内部完成。
func ConvertUSD(ctx context.Context, usdAmount float64, orderCurrency string) (*ConversionResult, error) {
	if usdAmount <= 0 {
		return nil, fmt.Errorf("invalid usd amount=%v", usdAmount)
	}
	rate, err := getUSDToQuoteRate(ctx, orderCurrency)
	if err != nil {
		return nil, err
	}
	return &ConversionResult{
		SourceCurrency: rate.Base,
		TargetCurrency: rate.Quote,
		SourceAmount:   usdAmount,
		TargetAmount:   usdAmount * rate.Rate,
		Rate:           rate.Rate,
		Source:         rate.Source,
		RateDate:       rate.RateDate,
		Cached:         rate.Cached,
		CacheExpiresAt: rate.CacheExpiresAt,
	}, nil
}

// getUSDToQuoteRate 查询市场汇率。正常结果缓存 1 分钟；同币种并发请求只会触发一次外部查询。
// 外部接口临时失败时，24 小时内的最近成功值可继续使用；失败结果冷却 5 分钟，避免持续重试。
func getUSDToQuoteRate(ctx context.Context, currencyCode string) (*quoteRate, error) {
	quote, err := normalizeCurrencyCode(currencyCode)
	if err != nil {
		return nil, err
	}
	if quote == baseCurrencyUSD {
		return &quoteRate{
			Base:     baseCurrencyUSD,
			Quote:    baseCurrencyUSD,
			Rate:     1,
			Source:   "fixed",
			RateDate: time.Now().UTC().Format("2006-01-02"),
		}, nil
	}
	if ctx == nil {
		ctx = gctx.New()
	}
	key := cacheKey(quote)
	if cached, ok := getCachedQuoteRate(ctx, key); ok {
		return cached, nil
	}

	lock := fetchLock(key)
	lock.Lock()
	defer lock.Unlock()
	if cached, ok := getCachedQuoteRate(ctx, key); ok {
		return cached, nil
	}
	if failure, ok := getActiveFailure(quote); ok {
		if stale, staleOK := getStaleMarketRate(quote); staleOK {
			return staleQuoteRate(quote, stale), nil
		}
		return nil, fmt.Errorf("fx rate retry cooling down for %s: %s", quote, failure.Message)
	}

	snapshot, err := marketRateFetcher(ctx, quote)
	if err != nil {
		saveFailure(quote, err)
		if stale, ok := getStaleMarketRate(quote); ok {
			g.Log().Warningf(ctx, "fx rate fallback to stale cache, quote=%s err=%v", quote, err)
			return staleQuoteRate(quote, stale), nil
		}
		return nil, err
	}

	deleteFailure(quote)
	saveStaleMarketRate(quote, snapshot)
	rate := newQuoteRate(quote, snapshot, false, time.Now().Add(fxRateCacheTTL).Unix())
	setCachedQuoteRate(ctx, key, rate)
	return rate, nil
}

func normalizeCurrencyCode(code string) (string, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if len(code) != 3 {
		return "", fmt.Errorf("invalid currency code=%s", code)
	}
	for _, ch := range code {
		if ch < 'A' || ch > 'Z' {
			return "", fmt.Errorf("invalid currency code=%s", code)
		}
	}
	return code, nil
}

func cacheKey(currency string) string {
	return baseCurrencyUSD + ":" + currency
}

func fetchLock(key string) *sync.Mutex {
	value, _ := fxRateFetchLocks.LoadOrStore(key, &sync.Mutex{})
	return value.(*sync.Mutex)
}

func newQuoteRate(quote string, snapshot marketRateSnapshot, cached bool, expiresAt int64) *quoteRate {
	return &quoteRate{
		Base:           baseCurrencyUSD,
		Quote:          quote,
		Rate:           snapshot.Rate,
		Source:         snapshot.Source,
		RateDate:       snapshot.RateDate,
		Cached:         cached,
		CacheExpiresAt: expiresAt,
	}
}

func staleQuoteRate(quote string, snapshot marketRateSnapshot) *quoteRate {
	snapshot.Source = sourceStaleCache + "(" + snapshot.Source + ")"
	return newQuoteRate(quote, snapshot, true, 0)
}

func getCachedQuoteRate(ctx context.Context, key string) (*quoteRate, bool) {
	value, err := fxRateCache.Get(ctx, key)
	if err != nil || value.IsNil() {
		return nil, false
	}
	rate, ok := value.Val().(*quoteRate)
	if !ok || rate == nil || rate.Rate <= 0 {
		return nil, false
	}
	copyRate := *rate
	copyRate.Cached = true
	return &copyRate, true
}

func setCachedQuoteRate(ctx context.Context, key string, rate *quoteRate) {
	if rate == nil || rate.Rate <= 0 {
		return
	}
	copyRate := *rate
	copyRate.Cached = false
	_ = fxRateCache.Set(ctx, key, &copyRate, fxRateCacheTTL)
}

func saveStaleMarketRate(quote string, snapshot marketRateSnapshot) {
	fxRateStaleMu.Lock()
	defer fxRateStaleMu.Unlock()
	fxRateStaleData[quote] = snapshot
}

func getStaleMarketRate(quote string) (marketRateSnapshot, bool) {
	fxRateStaleMu.RLock()
	defer fxRateStaleMu.RUnlock()
	snapshot, ok := fxRateStaleData[quote]
	if !ok || snapshot.Rate <= 0 || snapshot.FetchedAt.IsZero() {
		return marketRateSnapshot{}, false
	}
	return snapshot, time.Since(snapshot.FetchedAt) <= fxRateStaleTTL
}

func saveFailure(quote string, err error) {
	fxRateFailureMu.Lock()
	defer fxRateFailureMu.Unlock()
	fxRateFailureData[quote] = failureSnapshot{Message: err.Error(), RetryAt: time.Now().Add(fxRateFailureTTL)}
}

func getActiveFailure(quote string) (failureSnapshot, bool) {
	fxRateFailureMu.RLock()
	defer fxRateFailureMu.RUnlock()
	failure, ok := fxRateFailureData[quote]
	return failure, ok && time.Now().Before(failure.RetryAt)
}

func deleteFailure(quote string) {
	fxRateFailureMu.Lock()
	defer fxRateFailureMu.Unlock()
	delete(fxRateFailureData, quote)
}

func fetchMarketRateFromProviders(ctx context.Context, quote string) (marketRateSnapshot, error) {
	frankfurterRate, frankfurterErr := fetchFrankfurterRate(ctx, quote)
	if frankfurterErr == nil {
		return frankfurterRate, nil
	}
	g.Log().Warningf(ctx, "frankfurter fx rate failed, quote=%s err=%v", quote, frankfurterErr)

	openERAPIRate, openERAPIErr := fetchOpenERAPIRate(ctx, quote)
	if openERAPIErr == nil {
		return openERAPIRate, nil
	}
	g.Log().Warningf(ctx, "open.er-api fx rate failed, quote=%s err=%v", quote, openERAPIErr)
	return marketRateSnapshot{}, fmt.Errorf(
		"fx rate unavailable for %s: frankfurter=%v; open.er-api=%v",
		quote, frankfurterErr, openERAPIErr,
	)
}

func fetchFrankfurterRate(ctx context.Context, quote string) (marketRateSnapshot, error) {
	client := gclient.New().SetTimeout(fxRateRequestTimeout)
	resp, err := client.Get(ctx, fmt.Sprintf(frankfurterRateURL, quote))
	if err != nil {
		return marketRateSnapshot{}, err
	}
	defer resp.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return marketRateSnapshot{}, fmt.Errorf("frankfurter http status=%d", resp.StatusCode)
	}

	var payload frankfurterResp
	if err = json.Unmarshal(resp.ReadAll(), &payload); err != nil {
		return marketRateSnapshot{}, err
	}
	if payload.Rate <= 0 {
		return marketRateSnapshot{}, fmt.Errorf("frankfurter missing rate for %s", quote)
	}
	if base := strings.ToUpper(strings.TrimSpace(payload.Base)); base != "" && base != baseCurrencyUSD {
		return marketRateSnapshot{}, fmt.Errorf("frankfurter base=%s", base)
	}
	if responseQuote := strings.ToUpper(strings.TrimSpace(payload.Quote)); responseQuote != "" && responseQuote != quote {
		return marketRateSnapshot{}, fmt.Errorf("frankfurter quote=%s want=%s", responseQuote, quote)
	}
	return marketRateSnapshot{
		Rate:      payload.Rate,
		Source:    sourceFrankfurter,
		RateDate:  payload.Date,
		FetchedAt: time.Now(),
	}, nil
}

func fetchOpenERAPIRate(ctx context.Context, quote string) (marketRateSnapshot, error) {
	client := gclient.New().SetTimeout(fxRateRequestTimeout)
	resp, err := client.Get(ctx, openExchangeRateAPIURL)
	if err != nil {
		return marketRateSnapshot{}, err
	}
	defer resp.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return marketRateSnapshot{}, fmt.Errorf("open.er-api http status=%d", resp.StatusCode)
	}

	var payload openERAPIResp
	if err = json.Unmarshal(resp.ReadAll(), &payload); err != nil {
		return marketRateSnapshot{}, err
	}
	if !strings.EqualFold(payload.Result, "success") {
		return marketRateSnapshot{}, fmt.Errorf("open.er-api result=%s", payload.Result)
	}
	if base := strings.ToUpper(strings.TrimSpace(payload.BaseCode)); base != "" && base != baseCurrencyUSD {
		return marketRateSnapshot{}, fmt.Errorf("open.er-api base=%s", base)
	}
	rate, ok := payload.Rates[quote]
	if !ok || rate <= 0 {
		return marketRateSnapshot{}, fmt.Errorf("open.er-api missing rate for %s", quote)
	}
	rateDate := time.Now().UTC().Format("2006-01-02")
	if payload.TimeLastUpdate > 0 {
		rateDate = time.Unix(payload.TimeLastUpdate, 0).UTC().Format("2006-01-02")
	}
	return marketRateSnapshot{
		Rate:      rate,
		Source:    sourceOpenERAPI,
		RateDate:  rateDate,
		FetchedAt: time.Now(),
	}, nil
}
