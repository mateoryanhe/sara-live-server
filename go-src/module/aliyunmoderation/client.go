package aliyunmoderation

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green "github.com/alibabacloud-go/green-20220302/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gogf/gf/v2/frame/g"
)

type reasonPayload struct {
	RiskLevel string `json:"riskLevel"`
}

var clientMu sync.Mutex
var greenClient *green.Client
var greenClientKey string

func getGreenClient(cfg *cfgSnapshot) (*green.Client, error) {
	key := cfg.AccessKeyId + "|" + cfg.AccessKeySecret + "|" + cfg.RegionId + "|" + cfg.Endpoint
	clientMu.Lock()
	defer clientMu.Unlock()
	if greenClient != nil && greenClientKey == key {
		return greenClient, nil
	}
	conf := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyId),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		RegionId:        tea.String(cfg.RegionId),
		Endpoint:        tea.String(cfg.Endpoint),
		ConnectTimeout:  tea.Int(3000),
		ReadTimeout:     tea.Int(6000),
	}
	c, err := green.NewClient(conf)
	if err != nil {
		return nil, err
	}
	greenClient = c
	greenClientKey = key
	return c, nil
}

func invalidateGreenClient() {
	clientMu.Lock()
	defer clientMu.Unlock()
	greenClient = nil
	greenClientKey = ""
}

func moderateText(ctx context.Context, cfg *cfgSnapshot, service, content string) (bool, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return true, nil
	}
	if service == "" {
		service = cfg.ChatService
	}
	client, err := getGreenClient(cfg)
	if err != nil {
		return false, err
	}
	params, err := json.Marshal(map[string]string{"content": content})
	if err != nil {
		return false, err
	}
	req := &green.TextModerationRequest{
		Service:           tea.String(service),
		ServiceParameters: tea.String(string(params)),
	}
	runtime := &util.RuntimeOptions{}
	runtime.SetConnectTimeout(3000)
	runtime.SetReadTimeout(8000)

	resp, err := client.TextModerationWithOptions(req, runtime)
	if err != nil {
		g.Log().Warningf(ctx, "aliyun TextModeration err service=%s: %v", service, err)
		return false, err
	}
	if resp == nil || resp.Body == nil {
		return false, nil
	}
	if resp.StatusCode != nil && *resp.StatusCode != http.StatusOK {
		g.Log().Warningf(ctx, "aliyun TextModeration http status=%d", *resp.StatusCode)
		return false, nil
	}
	body := resp.Body
	if body.Code != nil && *body.Code != 200 {
		g.Log().Warningf(ctx, "aliyun TextModeration code=%d msg=%s", *body.Code, tea.StringValue(body.Message))
		return false, nil
	}
	if body.Data == nil {
		return true, nil
	}
	return !hasRiskLabels(body.Data.Labels, body.Data.Reason), nil
}

func hasRiskLabels(labels, reason *string) bool {
	labelStr := strings.TrimSpace(tea.StringValue(labels))
	if labelStr != "" {
		lower := strings.ToLower(labelStr)
		if lower != "nonlabel" && lower != "normal" && lower != "pass" {
			return true
		}
	}
	reasonStr := strings.TrimSpace(tea.StringValue(reason))
	if reasonStr == "" {
		return false
	}
	var rp reasonPayload
	if err := json.Unmarshal([]byte(reasonStr), &rp); err != nil {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(rp.RiskLevel)) {
	case "high", "medium":
		return true
	default:
		return false
	}
}
