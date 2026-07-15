package agora

import (
	"context"
	"encoding/base64"

	"github.com/gogf/gf/v2/net/gclient"
	"xr-game-server/errercode"
)

const agoraRestBaseURL = "https://api.sd-rtn.com"

func buildRestAuthHeader(customerId, customerSecret string) string {
	plain := customerId + ":" + customerSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(plain))
}

func newAgoraRestClient(cfg *agoraCfgSnapshot) *gclient.Client {
	client := gclient.New()
	client.SetHeader("Accept", "application/json")
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", buildRestAuthHeader(cfg.RestCustomerId, cfg.RestCustomerSecret))
	return client
}

func validateAgoraRestCfg(cfg *agoraCfgSnapshot) error {
	if err := validateAgoraCfg(cfg); err != nil {
		return err
	}
	if cfg.RestCustomerId == "" || cfg.RestCustomerSecret == "" {
		return errercode.CreateCode(errercode.AgoraCfgInvalid)
	}
	if cfg.CloudPlayerRegion == "" {
		return errercode.CreateCode(errercode.AgoraCfgInvalid)
	}
	return nil
}

func agoraRestPost(ctx context.Context, cfg *agoraCfgSnapshot, path string, body any) (*gclient.Response, error) {
	client := newAgoraRestClient(cfg)
	resp, err := client.Post(ctx, agoraRestBaseURL+path, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
