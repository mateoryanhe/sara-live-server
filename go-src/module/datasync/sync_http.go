package datasync

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/errercode"
)

const syncHTTPTimeout = 120 * time.Second

type syncAPIEnvelope struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

func loadSyncTargetConfig() (targetBase, syncToken string, err error) {
	cfg := cfgdao.GetDataSyncCfg()
	if cfg == nil {
		return "", "", errInvalidParam()
	}
	targetBase = strings.TrimRight(strings.TrimSpace(cfg.TargetApiBase), "/")
	syncToken = strings.TrimSpace(cfg.Token)
	if targetBase == "" || syncToken == "" {
		return "", "", errInvalidParam()
	}
	return targetBase, syncToken, nil
}

func postSyncReceive(path string, payload any, result any) error {
	targetBase, syncToken, err := loadSyncTargetConfig()
	if err != nil {
		return err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reqCtx, cancel := context.WithTimeout(context.Background(), syncHTTPTimeout)
	defer cancel()

	client := gclient.New().SetTimeout(syncHTTPTimeout)
	client.SetHeader(HeaderDataSyncToken, syncToken)
	client.SetHeader("Content-Type", "application/json")

	url := targetBase + path
	resp, err := client.Post(reqCtx, url, body)
	if err != nil {
		return fmt.Errorf("sync request failed: %w", err)
	}
	defer resp.Close()

	raw := resp.ReadAll()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sync http status=%d body=%s", resp.StatusCode, string(raw))
	}

	var envelope syncAPIEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return fmt.Errorf("parse sync response: %w", err)
	}
	if envelope.Code != errercode.Success {
		return fmt.Errorf("sync rejected code=%d", envelope.Code)
	}
	if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return fmt.Errorf("sync response data is empty")
	}
	if err := json.Unmarshal(envelope.Data, result); err != nil {
		return fmt.Errorf("parse sync response data: %w", err)
	}
	return nil
}
