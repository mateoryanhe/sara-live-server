package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
	"xr-game-server/dao/cfgdao"
)

type firebaseServiceAccount struct {
	Type        string `json:"type"`
	ProjectId   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

var (
	firebaseAuthClientMu  sync.Mutex
	firebaseAuthClient    *firebaseauth.Client
	firebaseAuthClientErr error
	firebaseAuthClientKey string
)

// ResetFirebaseAuthClient CMS 保存配置后重建 Firebase Admin 客户端.
func ResetFirebaseAuthClient() {
	firebaseAuthClientMu.Lock()
	defer firebaseAuthClientMu.Unlock()
	firebaseAuthClient = nil
	firebaseAuthClientErr = nil
	firebaseAuthClientKey = ""
}

func getFirebaseAuthClient(ctx context.Context) (*firebaseauth.Client, error) {
	cfg := cfgdao.GetFirebaseCfgCached()
	if cfg == nil {
		return nil, fmt.Errorf("firebase config is empty")
	}
	projectId := strings.TrimSpace(cfg.ProjectId)
	serviceAccountJson := strings.TrimSpace(cfg.ServiceAccountJson)
	if projectId == "" || serviceAccountJson == "" {
		return nil, fmt.Errorf("firebase config is incomplete")
	}

	clientKey := firebaseClientKey(projectId, serviceAccountJson)
	firebaseAuthClientMu.Lock()
	defer firebaseAuthClientMu.Unlock()
	if firebaseAuthClientKey == clientKey {
		return firebaseAuthClient, firebaseAuthClientErr
	}

	firebaseAuthClient, firebaseAuthClientErr = newFirebaseAuthClient(ctx, projectId, serviceAccountJson)
	firebaseAuthClientKey = clientKey
	return firebaseAuthClient, firebaseAuthClientErr
}

func verifyFirebaseUID(ctx context.Context, idToken string) (string, error) {
	idToken = strings.TrimSpace(idToken)
	if idToken == "" || len(idToken) > 16384 {
		return "", fmt.Errorf("firebase id token is invalid")
	}
	client, err := getFirebaseAuthClient(ctx)
	if err != nil {
		return "", err
	}
	verifiedToken, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("verify firebase id token: %w", err)
	}
	firebaseUID := strings.TrimSpace(verifiedToken.UID)
	if firebaseUID == "" {
		return "", fmt.Errorf("firebase uid is empty")
	}
	return firebaseUID, nil
}

func newFirebaseAuthClient(ctx context.Context, projectId, serviceAccountJson string) (*firebaseauth.Client, error) {
	app, err := firebase.NewApp(
		ctx,
		&firebase.Config{ProjectID: projectId},
		option.WithCredentialsJSON([]byte(serviceAccountJson)),
	)
	if err != nil {
		return nil, fmt.Errorf("create firebase app: %w", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("create firebase auth client: %w", err)
	}
	return client, nil
}

func validateFirebaseConfig(ctx context.Context, projectId, serviceAccountJson string) error {
	projectId = strings.TrimSpace(projectId)
	serviceAccountJson = strings.TrimSpace(serviceAccountJson)
	if projectId == "" || serviceAccountJson == "" {
		return fmt.Errorf("firebase project id and service account json are required")
	}

	var serviceAccount firebaseServiceAccount
	if err := json.Unmarshal([]byte(serviceAccountJson), &serviceAccount); err != nil {
		return fmt.Errorf("parse firebase service account json: %w", err)
	}
	if serviceAccount.Type != "service_account" ||
		strings.TrimSpace(serviceAccount.ProjectId) == "" ||
		strings.TrimSpace(serviceAccount.ClientEmail) == "" ||
		strings.TrimSpace(serviceAccount.PrivateKey) == "" {
		return fmt.Errorf("firebase service account json is incomplete")
	}
	if serviceAccount.ProjectId != projectId {
		return fmt.Errorf("firebase project id does not match service account")
	}
	_, err := newFirebaseAuthClient(ctx, projectId, serviceAccountJson)
	return err
}

func firebaseClientKey(projectId, serviceAccountJson string) string {
	sum := sha256.Sum256([]byte(projectId + "\x00" + serviceAccountJson))
	return hex.EncodeToString(sum[:])
}
