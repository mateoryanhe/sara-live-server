package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"
)

func TestValidateFirebaseConfig(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate private key: %v", err)
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal private key: %v", err)
	}
	credentials, err := json.Marshal(map[string]string{
		"type":         "service_account",
		"project_id":   "firebase-project",
		"private_key":  string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})),
		"client_email": "firebase-adminsdk@example.iam.gserviceaccount.com",
		"token_uri":    "https://oauth2.googleapis.com/token",
	})
	if err != nil {
		t.Fatalf("marshal credentials: %v", err)
	}

	if err := validateFirebaseConfig(context.Background(), "firebase-project", string(credentials)); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	if err := validateFirebaseConfig(context.Background(), "other-project", string(credentials)); err == nil {
		t.Fatal("expected project mismatch error")
	}
	if err := validateFirebaseConfig(context.Background(), "firebase-project", `{}`); err == nil {
		t.Fatal("expected incomplete service account error")
	}
}

func TestParseFirebaseClientConfig(t *testing.T) {
	raw := `{
		"apiKey":"client-api-key",
		"authDomain":"firebase-project.firebaseapp.com",
		"projectId":"firebase-project",
		"storageBucket":"firebase-project.firebasestorage.app",
		"messagingSenderId":"123456",
		"appId":"1:123456:web:abcdef"
	}`
	cfg, err := parseFirebaseClientConfig("firebase-project", raw)
	if err != nil {
		t.Fatalf("valid client config rejected: %v", err)
	}
	if cfg.ProjectId != "firebase-project" || cfg.ApiKey != "client-api-key" || cfg.AppId == "" {
		t.Fatalf("unexpected parsed client config: %+v", cfg)
	}
	if _, err := parseFirebaseClientConfig("other-project", raw); err == nil {
		t.Fatal("expected project mismatch error")
	}
	if _, err := parseFirebaseClientConfig("firebase-project", `{}`); err == nil {
		t.Fatal("expected incomplete client config error")
	}
}
