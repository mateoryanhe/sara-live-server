package fxrate

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConvertUSDCoalescesAndCaches(t *testing.T) {
	const quote = "INR"
	clearRateTestState(quote)
	originalFetcher := marketRateFetcher
	defer func() {
		marketRateFetcher = originalFetcher
		clearRateTestState(quote)
	}()

	var calls atomic.Int32
	marketRateFetcher = func(_ context.Context, gotQuote string) (marketRateSnapshot, error) {
		calls.Add(1)
		if gotQuote != quote {
			return marketRateSnapshot{}, fmt.Errorf("quote=%s want=%s", gotQuote, quote)
		}
		time.Sleep(20 * time.Millisecond)
		return marketRateSnapshot{
			Rate:      83.25,
			Source:    sourceFrankfurter,
			RateDate:  "2026-09-15",
			FetchedAt: time.Now(),
		}, nil
	}

	const workers = 20
	var waitGroup sync.WaitGroup
	waitGroup.Add(workers)
	errors := make(chan error, workers)
	for range workers {
		go func() {
			defer waitGroup.Done()
			conversion, err := ConvertUSD(context.Background(), 2, quote)
			if err != nil {
				errors <- err
				return
			}
			if conversion.Rate != 83.25 {
				t.Errorf("rate=%v want=83.25", conversion.Rate)
			}
			if conversion.TargetAmount != 166.5 {
				t.Errorf("amount=%v want=166.5", conversion.TargetAmount)
			}
			if conversion.TargetCurrency != quote {
				t.Errorf("quote=%s want=%s", conversion.TargetCurrency, quote)
			}
		}()
	}
	waitGroup.Wait()
	close(errors)
	for err := range errors {
		t.Fatal(err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("external fetches=%d want=1", got)
	}

	cached, err := ConvertUSD(context.Background(), 3, quote)
	if err != nil {
		t.Fatal(err)
	}
	if !cached.Cached {
		t.Fatal("second lookup should use cache")
	}
	if cached.TargetAmount != 249.75 {
		t.Fatalf("cached amount=%v want=249.75", cached.TargetAmount)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("external fetches after cache hit=%d want=1", got)
	}
}

func TestConvertUSDNormalizesOrderCurrencyAndSkipsProviderForUSD(t *testing.T) {
	originalFetcher := marketRateFetcher
	defer func() { marketRateFetcher = originalFetcher }()

	var calls atomic.Int32
	marketRateFetcher = func(_ context.Context, quote string) (marketRateSnapshot, error) {
		calls.Add(1)
		return marketRateSnapshot{}, fmt.Errorf("unexpected provider call for %s", quote)
	}

	conversion, err := ConvertUSD(context.Background(), 12.34, " usd ")
	if err != nil {
		t.Fatal(err)
	}
	if conversion.SourceCurrency != "USD" || conversion.TargetCurrency != "USD" {
		t.Fatalf("currencies=%s/%s want=USD/USD", conversion.SourceCurrency, conversion.TargetCurrency)
	}
	if conversion.SourceAmount != 12.34 || conversion.TargetAmount != 12.34 || conversion.Rate != 1 {
		t.Fatalf("conversion=%+v", conversion)
	}
	if calls.Load() != 0 {
		t.Fatalf("external fetches=%d want=0", calls.Load())
	}
}

func clearRateTestState(quote string) {
	ctx := context.Background()
	key := cacheKey(quote)
	_, _ = fxRateCache.Remove(ctx, key)
	fxRateFetchLocks.Delete(key)

	fxRateStaleMu.Lock()
	delete(fxRateStaleData, quote)
	fxRateStaleMu.Unlock()

	fxRateFailureMu.Lock()
	delete(fxRateFailureData, quote)
	fxRateFailureMu.Unlock()
}
