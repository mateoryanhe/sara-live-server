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
	IdleTimeout int    `json:"idleTimeout,omitempty"`
}

type cloudPlayerUpdateReq struct {
	Player cloudPlayerUpdateBody `json:"player"`
}

type cloudPlayerUpdateBody struct {
	Token string `json:"token,omitempty"`
}

type cloudPlayerCreateResp struct {
	Player struct {
		ID string `json:"id"`
	} `json:"player"`
}

// StartBotAnchorCloudPlayer 机器人主播开播时创建声网云播放器
func StartBotAnchorCloudPlayer(ctx context.Context, anchorId uint64, cloudPlayerVideo string) (string, int64, error) {
	streamUrl := strings.TrimSpace(upload.ResolveCloudPlayerVideoUrl(cloudPlayerVideo))
	if streamUrl == "" {
		return "", 0, errercode.CreateCode(errercode.InvalidParam)
	}

	cfg := getAgoraCfgCache()
	if err := validateAgoraRestCfg(cfg); err != nil {
		return "", 0, err
	}

	channelName := strconv.FormatUint(anchorId, 10)
	token, expireAt, err := BuildChannelToken(anchorId, channelName, agoradto.RTCRolePublisher)
	if err != nil {
		return "", 0, err
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
			IdleTimeout: 600,
		},
	}

	resp, err := agoraRestPost(ctx, cfg, path, reqBody)
	if err != nil {
		g.Log().Errorf(ctx, "create cloud player failed anchorId=%d err=%v", anchorId, err)
		return "", 0, errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
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
			registerCloudPlayerSequence(playerId)
			return playerId, expireAt, nil
		}
	}

	if resp.StatusCode == http.StatusConflict && playerId != "" {
		registerCloudPlayerSequence(playerId)
		return playerId, expireAt, nil
	}

	g.Log().Errorf(ctx, "create cloud player failed anchorId=%d status=%d playerId=%s body=%s", anchorId, resp.StatusCode, playerId, string(respBody))
	return "", 0, errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
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
		unregisterCloudPlayerSequence(playerId)
		return nil
	}
	if resp.StatusCode == http.StatusNotFound {
		g.Log().Warningf(ctx, "delete cloud player skipped, player not found playerId=%s requestId=%s", playerId, requestID)
		unregisterCloudPlayerSequence(playerId)
		return nil
	}

	respBody := resp.ReadAll()
	g.Log().Errorf(ctx, "delete cloud player failed playerId=%s requestId=%s status=%d body=%s", playerId, requestID, resp.StatusCode, string(respBody))
	return errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
}
