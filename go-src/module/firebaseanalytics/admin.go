package firebaseanalytics

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"xr-game-server/dao/cfgdao"
)

var (
	firebaseAdminMu  sync.Mutex
	firebaseAdminApp *firebase.App
	firebaseAdminErr error
	firebaseAdminKey string
)

// ResetFirebaseAdmin CMS 保存或关闭埋点后重置 Admin SDK 实例.
func ResetFirebaseAdmin() {
	firebaseAdminMu.Lock()
	defer firebaseAdminMu.Unlock()
	firebaseAdminApp = nil
	firebaseAdminErr = nil
	firebaseAdminKey = ""
}

func warmInitAdmin(ctx context.Context) {
	if !IsEnabled() {
		ResetFirebaseAdmin()
		return
	}
	_, _ = GetFirebaseAdminApp(ctx)
}

// GetFirebaseAdminApp 获取 Analytics 模块独立的 Firebase Admin App(服务端上报等).
func GetFirebaseAdminApp(ctx context.Context) (*firebase.App, error) {
	cfg := cfgdao.ResolveFirebaseAnalyticsCfg()
	if cfg == nil || !cfg.IsActive() {
		return nil, fmt.Errorf("firebase analytics is disabled or config incomplete")
	}
	projectId := strings.TrimSpace(cfg.ProjectId)
	serviceAccountJson := strings.TrimSpace(cfg.ServiceAccountJson)
	if projectId == "" || serviceAccountJson == "" {
		return nil, fmt.Errorf("firebase analytics server config is incomplete")
	}

	clientKey := firebaseAdminClientKey(projectId, serviceAccountJson)
	firebaseAdminMu.Lock()
	defer firebaseAdminMu.Unlock()
	if firebaseAdminKey == clientKey && firebaseAdminApp != nil {
		return firebaseAdminApp, firebaseAdminErr
	}
	firebaseAdminApp, firebaseAdminErr = newFirebaseAdminApp(ctx, projectId, serviceAccountJson)
	firebaseAdminKey = clientKey
	return firebaseAdminApp, firebaseAdminErr
}

func newFirebaseAdminApp(ctx context.Context, projectId, serviceAccountJson string) (*firebase.App, error) {
	app, err := firebase.NewApp(
		ctx,
		&firebase.Config{ProjectID: projectId},
		option.WithCredentialsJSON([]byte(serviceAccountJson)),
	)
	if err != nil {
		return nil, fmt.Errorf("create firebase admin app: %w", err)
	}
	return app, nil
}

func firebaseAdminClientKey(projectId, serviceAccountJson string) string {
	sum := sha256.Sum256([]byte(projectId + "\x00" + serviceAccountJson))
	return hex.EncodeToString(sum[:])
}
