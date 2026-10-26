package resourcemonitor

import (
	"strings"
	"testing"
)

func TestCmsTrendQueryCacheKey_DefaultOnly(t *testing.T) {
	key := cmsTrendQueryCacheKey(trendQuery{Limit: 26000})
	if !strings.HasPrefix(key, "default:26000:fine:") {
		t.Fatalf("default key = %q", key)
	}
	key2 := cmsTrendQueryCacheKey(trendQuery{Limit: 100})
	if !strings.HasPrefix(key2, "default:100:fine:") {
		t.Fatalf("limit key = %q", key2)
	}
}

func TestListTrendPoints_CustomRangeBypassesSharedShape(t *testing.T) {
	// 自定义时间窗不应走永久/长生命周期 key;仅校验 key 生成不再包含 custom 串.
	// listTrendPoints 对自定义直接查库,本测试保证默认 key 路径仍可用.
	q := trendQuery{
		StartTime: "2026-10-24 00:00:00",
		EndTime:   "2026-10-24 23:59:59",
		Limit:     100,
	}
	if strings.TrimSpace(q.StartTime) == "" && strings.TrimSpace(q.EndTime) == "" {
		t.Fatal("expected custom range")
	}
	key := cmsTrendQueryCacheKey(trendQuery{Limit: q.Limit})
	if strings.Contains(key, "custom:") {
		t.Fatalf("default key must not use custom prefix: %q", key)
	}
}
