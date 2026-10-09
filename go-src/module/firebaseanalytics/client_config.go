package firebaseanalytics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ClientConfig 客户端 Firebase/GA4 初始化字段(可公开给 App/H5).
type ClientConfig struct {
	ApiKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain,omitempty"`
	ProjectId         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket,omitempty"`
	MessagingSenderId string `json:"messagingSenderId,omitempty"`
	AppId             string `json:"appId"`
	MeasurementId     string `json:"measurementId,omitempty"`
	DatabaseUrl       string `json:"databaseURL,omitempty"`
}

func parseClientConfig(projectId, raw string) (*ClientConfig, error) {
	projectId = strings.TrimSpace(projectId)
	raw = strings.TrimSpace(raw)
	if projectId == "" || raw == "" {
		return nil, fmt.Errorf("project id and client config json are required")
	}
	var cfg ClientConfig
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, fmt.Errorf("parse client config json: %w", err)
	}
	cfg.ApiKey = strings.TrimSpace(cfg.ApiKey)
	cfg.ProjectId = strings.TrimSpace(cfg.ProjectId)
	cfg.AppId = strings.TrimSpace(cfg.AppId)
	cfg.MeasurementId = strings.TrimSpace(cfg.MeasurementId)
	if cfg.ApiKey == "" || cfg.ProjectId == "" || cfg.AppId == "" {
		return nil, fmt.Errorf("client config requires apiKey, projectId and appId")
	}
	if cfg.ProjectId != projectId {
		return nil, fmt.Errorf("client config projectId does not match")
	}
	if cfg.MeasurementId == "" {
		return nil, fmt.Errorf("client config requires measurementId for GA4 analytics")
	}
	return &cfg, nil
}
