package agora

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gogf/gf/v2/net/gclient"
)

const agoraRestBaseURL = "https://api.agora.io"

type agoraRestResp struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

type agoraUserStatusData struct {
	InChannel bool  `json:"in_channel"`
	Join      int64 `json:"join"`
	Role      int   `json:"role"`
}

type agoraChannelUserListData struct {
	ChannelExist bool          `json:"channel_exist"`
	Broadcasters []interface{} `json:"broadcasters"`
	Audience     []interface{} `json:"audience"`
	Users        []interface{} `json:"users"`
}

func buildRestAuthHeader(customerId, customerSecret string) string {
	plain := customerId + ":" + customerSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(plain))
}

func newAgoraRestClient(cfg *agoraCfgSnapshot) *gclient.Client {
	client := gclient.New()
	client.SetHeader("Accept", "application/json")
	client.SetHeader("Authorization", buildRestAuthHeader(cfg.RestCustomerId, cfg.RestCustomerSecret))
	return client
}

func agoraRestGet(ctx context.Context, cfg *agoraCfgSnapshot, path string) (*agoraRestResp, error) {
	client := newAgoraRestClient(cfg)
	resp, err := client.Get(ctx, agoraRestBaseURL+path)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	body := resp.ReadAll()
	if len(body) == 0 {
		return nil, fmt.Errorf("agora rest empty response")
	}

	var ret agoraRestResp
	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// queryUserStatus 调用声网官方 REST API 查询指定用户频道状态
// GET /dev/v1/channel/user/property/{appid}/{uid}/{channelName}
func queryUserStatus(ctx context.Context, cfg *agoraCfgSnapshot, channelName string, uid uint64) (*agoraUserStatusData, error) {
	path := fmt.Sprintf(
		"/dev/v1/channel/user/property/%s/%d/%s",
		url.PathEscape(cfg.AppId),
		uid,
		url.PathEscape(channelName),
	)
	resp, err := agoraRestGet(ctx, cfg, path)
	if err != nil {
		return nil, err
	}
	if !resp.Success || len(resp.Data) == 0 {
		return &agoraUserStatusData{InChannel: false}, nil
	}
	var data agoraUserStatusData
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// queryChannelUserList 调用声网官方 REST API 查询频道用户列表
// GET /dev/v1/channel/user/{appid}/{channelName}
func queryChannelUserList(ctx context.Context, cfg *agoraCfgSnapshot, channelName string) (*agoraChannelUserListData, error) {
	path := fmt.Sprintf(
		"/dev/v1/channel/user/%s/%s",
		url.PathEscape(cfg.AppId),
		url.PathEscape(channelName),
	)
	resp, err := agoraRestGet(ctx, cfg, path)
	if err != nil {
		return nil, err
	}
	if !resp.Success || len(resp.Data) == 0 {
		return &agoraChannelUserListData{ChannelExist: false}, nil
	}
	var data agoraChannelUserListData
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func userAccountInChannelLists(userAccount string, data *agoraChannelUserListData) bool {
	if data == nil || !data.ChannelExist {
		return false
	}
	return matchUserAccount(data.Broadcasters, userAccount) ||
		matchUserAccount(data.Audience, userAccount) ||
		matchUserAccount(data.Users, userAccount)
}

func matchUserAccount(list []interface{}, target string) bool {
	for _, item := range list {
		if normalizeAgoraUserID(item) == target {
			return true
		}
	}
	return false
}

func normalizeAgoraUserID(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return val.String()
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case uint64:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprint(val)
	}
}
