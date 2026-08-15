package recharge

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/recharge"
)

var (
	googlePlayClientMu   sync.Mutex
	googlePlayService    *androidpublisher.Service
	googlePlayServiceErr error
	googlePlayCredKey    string
)

// ResetGooglePlayClient CMS 保存配置后重置 Google API 客户端
func ResetGooglePlayClient() {
	googlePlayClientMu.Lock()
	defer googlePlayClientMu.Unlock()
	googlePlayService = nil
	googlePlayServiceErr = nil
	googlePlayCredKey = ""
}

func getGooglePlayService() (*androidpublisher.Service, error) {
	cfg := cfgdao.GetGooglePlayCfgCached()
	if cfg == nil || cfg.ServiceAccountJson == "" {
		return nil, fmt.Errorf("google play service account json is empty")
	}
	credKey := googlePlayCredKeyOf(cfg)
	googlePlayClientMu.Lock()
	defer googlePlayClientMu.Unlock()
	if googlePlayService != nil && googlePlayCredKey == credKey {
		return googlePlayService, googlePlayServiceErr
	}
	ctx := gctx.New()
	googlePlayService, googlePlayServiceErr = androidpublisher.NewService(ctx, option.WithCredentialsJSON([]byte(cfg.ServiceAccountJson)))
	googlePlayCredKey = credKey
	return googlePlayService, googlePlayServiceErr
}

func googlePlayCredKeyOf(cfg *entity.GooglePlayCfg) string {
	if cfg == nil {
		return ""
	}
	return fmt.Sprintf("%d:%d", cfg.ID, len(cfg.ServiceAccountJson))
}

func getGoogleProductPurchase(ctx context.Context, packageName, productId, purchaseToken string) (*androidpublisher.ProductPurchase, error) {
	svc, err := getGooglePlayService()
	if err != nil {
		return nil, err
	}
	return svc.Purchases.Products.Get(packageName, productId, purchaseToken).Context(ctx).Do()
}

func consumeGoogleProductPurchase(ctx context.Context, packageName, productId, purchaseToken string) error {
	svc, err := getGooglePlayService()
	if err != nil {
		return err
	}
	return svc.Purchases.Products.Consume(packageName, productId, purchaseToken).Context(ctx).Do()
}

func logGooglePlayInfo(ctx context.Context, format string, args ...any) {
	g.Log().Infof(ctx, "google play "+format, args...)
}

func logGooglePlayError(ctx context.Context, format string, args ...any) {
	g.Log().Errorf(ctx, "google play "+format, args...)
}

func getActiveGooglePlayCfg() *entity.GooglePlayCfg {
	return cfgdao.GetGooglePlayCfgFromMemory()
}

func googlePlayEnabled() bool {
	return cfgdao.GooglePlayEnabled()
}
