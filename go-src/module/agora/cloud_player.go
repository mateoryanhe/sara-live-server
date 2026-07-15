package agora

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/xrjson"
	"xr-game-server/dto/agoradto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

type cloudPlayerCreateReq struct {
	Player cloudPlayerCreateBody `json:"player"`
}

type cloudPlayerCreateBody struct {
	Account     string `json:"account"`
	ChannelName string `json:"channelName"`
	StreamUrl   string `json:"streamUrl"`
	Token       string `json:"token,omitempty"`
	Name        string `json:"name,omitempty"`
	RepeatTime  int    `json:"repeatTime,omitempty"`
}

type cloudPlayerCreateResp struct {
	Player struct {
		ID string `json:"id"`
	} `json:"player"`
}

// StartBotAnchorCloudPlayer 机器人主播开播时创建声网云播放器
func StartBotAnchorCloudPlayer(ctx context.Context, anchorId uint64, cloudPlayerVideo string) (string, error) {
	streamUrl := strings.TrimSpace(upload.ResolveCloudPlayerVideoUrl(cloudPlayerVideo))
	if streamUrl == "" {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}

	cfg := getAgoraCfgCache()
	if err := validateAgoraRestCfg(cfg); err != nil {
		return "", err
	}

	channelName := strconv.FormatUint(anchorId, 10)
	token, _, err := BuildChannelToken(anchorId, channelName, agoradto.RTCRolePublisher)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/%s/v1/projects/%s/cloud-player/players", cfg.CloudPlayerRegion, cfg.AppId)
	reqBody := cloudPlayerCreateReq{
		Player: cloudPlayerCreateBody{
			Account:     buildUserAccount(anchorId),
			ChannelName: channelName,
			StreamUrl:   streamUrl,
			Token:       token,
			Name:        fmt.Sprintf("bot_%d", anchorId),
			RepeatTime:  -1,
		},
	}

	resp, err := agoraRestPost(ctx, cfg, path, reqBody)
	if err != nil {
		g.Log().Errorf(ctx, "create cloud player failed anchorId=%d err=%v", anchorId, err)
		return "", errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
	}
	defer resp.Close()

	respBody := resp.ReadAll()
	playerId := strings.TrimSpace(resp.Header.Get("X-Resource-ID"))
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		if playerId == "" {
			var ret cloudPlayerCreateResp
			if unmarshalErr := xrjson.Unmarshal(respBody, &ret); unmarshalErr == nil {
				playerId = strings.TrimSpace(ret.Player.ID)
			}
		}
		if playerId != "" {
			return playerId, nil
		}
	}

	if resp.StatusCode == http.StatusConflict && playerId != "" {
		return playerId, nil
	}

	g.Log().Errorf(ctx, "create cloud player failed anchorId=%d status=%d playerId=%s body=%s", anchorId, resp.StatusCode, playerId, string(respBody))
	return "", errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
}

// StopBotAnchorCloudPlayer 机器人主播下播时销毁声网云播放器
func StopBotAnchorCloudPlayer(ctx context.Context, playerId string) error {
	playerId = strings.TrimSpace(playerId)
	if playerId == "" {
		return nil
	}

	cfg := getAgoraCfgCache()
	if err := validateAgoraRestCfg(cfg); err != nil {
		return err
	}

	path := fmt.Sprintf("/%s/v1/projects/%s/cloud-player/players/%s", cfg.CloudPlayerRegion, cfg.AppId, playerId)
	requestID := guid.S()
	client := newAgoraRestClient(cfg)
	client.SetHeader("X-Request-ID", requestID)

	resp, err := client.Delete(ctx, agoraRestBaseURL+path)
	if err != nil {
		g.Log().Errorf(ctx, "delete cloud player failed playerId=%s requestId=%s err=%v", playerId, requestID, err)
		return errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
	}
	defer resp.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	g.Log().Errorf(ctx, "delete cloud player failed playerId=%s requestId=%s status=%d", playerId, requestID, resp.StatusCode)
	return errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
}
