package agora

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"xr-game-server/core/xrlog"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/xrjson"
	"xr-game-server/dto/agoradto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const logSourceAgoraCloudPlayer = "AgoraCloudPlayer"

type cloudPlayerCreateReq struct {
	Player cloudPlayerBody `json:"player"`
}

type cloudPlayerBody struct {
	Account     string `json:"account"`
	ChannelName string `json:"channelName"`
	StreamUrl   string `json:"streamUrl"`
	Token       string `json:"token,omitempty"`
	Name        string `json:"name,omitempty"`
	RepeatTime  int    `json:"repeatTime,omitempty"`
	IdleTimeout int    `json:"idleTimeout,omitempty"`
}

type cloudPlayerUpdateReq struct {
	Player cloudPlayerBody `json:"player"`
}

func buildBotAnchorCloudPlayerBody(anchorId uint64, cloudPlayerVideo, token string) (cloudPlayerBody, error) {
	streamUrl := strings.TrimSpace(upload.ResolveCloudPlayerVideoUrl(cloudPlayerVideo))
	if streamUrl == "" {
		return cloudPlayerBody{}, errercode.CreateCode(errercode.InvalidParam)
	}
	return cloudPlayerBody{
		Account:     buildUserAccount(anchorId),
		ChannelName: strconv.FormatUint(anchorId, 10),
		StreamUrl:   streamUrl,
		Token:       token,
	}, nil
}

type cloudPlayerCreateResp struct {
	Player struct {
		ID string `json:"id"`
	} `json:"player"`
}

// StartBotAnchorCloudPlayer 机器人主播开播时创建声网云播放器
func StartBotAnchorCloudPlayer(ctx context.Context, anchorId uint64, cloudPlayerVideo string) (string, int64, error) {
	cfg := getAgoraCfgCache()
	if err := validateAgoraRestCfg(cfg); err != nil {
		return "", 0, err
	}

	channelName := strconv.FormatUint(anchorId, 10)
	token, expireAt, err := BuildChannelToken(anchorId, channelName, agoradto.RTCRolePublisher)
	if err != nil {
		return "", 0, err
	}

	playerBody, err := buildBotAnchorCloudPlayerBody(anchorId, cloudPlayerVideo, token)
	if err != nil {
		return "", 0, err
	}
	playerBody.Name = fmt.Sprintf("bot_%d", anchorId)
	playerBody.RepeatTime = -1
	playerBody.IdleTimeout = 600

	path := fmt.Sprintf("/%s/v1/projects/%s/cloud-player/players", cfg.CloudPlayerRegion, cfg.AppId)
	reqBody := cloudPlayerCreateReq{Player: playerBody}

	resp, err := agoraRestPost(ctx, cfg, path, reqBody)
	if err != nil {
		xrlog.ErrorWithErr(ctx, logSourceAgoraCloudPlayer,
			fmt.Sprintf("create request failed anchorId=%d channel=%s streamUrl=%s path=%s", anchorId, playerBody.ChannelName, playerBody.StreamUrl, path),
			err)
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

	xrlog.Error(ctx, logSourceAgoraCloudPlayer, fmt.Sprintf(
		"create failed anchorId=%d status=%d playerId=%s channel=%s streamUrl=%s region=%s appId=%s path=%s requestId=%s body=%s",
		anchorId, resp.StatusCode, playerId, playerBody.ChannelName, playerBody.StreamUrl, cfg.CloudPlayerRegion, cfg.AppId, path,
		strings.TrimSpace(resp.Header.Get("X-Request-ID")), string(respBody),
	))
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
		xrlog.ErrorWithErr(ctx, logSourceAgoraCloudPlayer,
			fmt.Sprintf("delete request failed playerId=%s requestId=%s path=%s", playerId, requestID, path),
			err)
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
	xrlog.Error(ctx, logSourceAgoraCloudPlayer, fmt.Sprintf(
		"delete failed playerId=%s requestId=%s status=%d path=%s body=%s",
		playerId, requestID, resp.StatusCode, path, string(respBody),
	))
	return errercode.CreateCode(errercode.AgoraCloudPlayerFailed)
}
