package message

import (
	"encoding/json"
	"strings"

	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/message"
	"xr-game-server/module/upload"
)

const followAnchorNicknamePlaceholder = "{nickname}"

var followAnchorPersonalMessageContent = personalSystemMessageDisplay{
	TitleEn:   "New follower",
	TitleEs:   "Nuevo seguidor",
	TitlePt:   "Novo seguidor",
	TitleHi:   "नया अनुयायी",
	TitleId:   "Pengikut baru",
	ContentEn: followAnchorNicknamePlaceholder + " started following you. Say hello when you go live!",
	ContentEs: followAnchorNicknamePlaceholder + " comenzó a seguirte. ¡Saluda cuando estés en vivo!",
	ContentPt: followAnchorNicknamePlaceholder + " começou a seguir você. Dê um oi quando entrar ao vivo!",
	ContentHi: followAnchorNicknamePlaceholder + " ने आपको फ़ॉलो करना शुरू किया। लाइव आने पर नमस्ते कहें!",
	ContentId: followAnchorNicknamePlaceholder + " mulai mengikuti Anda. Sapa mereka saat Anda live!",
}

type followAnchorPersonalMessageParams struct {
	FollowerId uint64 `json:"followerId,string"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar" dc:"头像资源名"`
}

// AddFollowAnchorPersonalMessage 用户新关注主播后,向主播下发个人系统消息
func AddFollowAnchorPersonalMessage(anchorId, followerId uint64) {
	if anchorId == 0 || followerId == 0 {
		return
	}
	nickname, avatar := resolveFollowAnchorFollowerProfile(followerId, "", "")
	payload := followAnchorPersonalMessageParams{
		FollowerId: followerId,
		Nickname:   nickname,
		Avatar:     avatar,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return
	}
	addPersonalSystemMessage(anchorId, entity.PersonalSystemMessageTypeFollowAnchor, string(raw))
}

func followAnchorPersonalMessageDisplay(row *entity.UserPersonalSystemMessage) personalSystemMessageDisplay {
	display := followAnchorPersonalMessageContent
	if row == nil {
		return display
	}
	var params followAnchorPersonalMessageParams
	if row.Params != "" {
		_ = json.Unmarshal([]byte(row.Params), &params)
	}
	nickname, avatarName := resolveFollowAnchorFollowerProfile(params.FollowerId, params.Nickname, params.Avatar)
	if nickname == "" {
		nickname = "Someone"
	}
	display.ContentEn = strings.ReplaceAll(display.ContentEn, followAnchorNicknamePlaceholder, nickname)
	display.ContentEs = strings.ReplaceAll(display.ContentEs, followAnchorNicknamePlaceholder, nickname)
	display.ContentPt = strings.ReplaceAll(display.ContentPt, followAnchorNicknamePlaceholder, nickname)
	display.ContentHi = strings.ReplaceAll(display.ContentHi, followAnchorNicknamePlaceholder, nickname)
	display.ContentId = strings.ReplaceAll(display.ContentId, followAnchorNicknamePlaceholder, nickname)
	display.IconEn = avatarName
	display.IconEs = avatarName
	display.IconPt = avatarName
	display.IconHi = avatarName
	display.IconId = avatarName
	return display
}

func followAnchorFollowerAvatarURL(row *entity.UserPersonalSystemMessage) string {
	if row == nil {
		return ""
	}
	var params followAnchorPersonalMessageParams
	if row.Params != "" {
		_ = json.Unmarshal([]byte(row.Params), &params)
	}
	_, avatarName := resolveFollowAnchorFollowerProfile(params.FollowerId, params.Nickname, params.Avatar)
	return upload.ResolveAvatarUrlForUser(params.FollowerId, avatarName)
}

func resolveFollowAnchorFollowerProfile(followerId uint64, storedNickname, storedAvatar string) (nickname, avatarName string) {
	nickname = storedNickname
	avatarName = storedAvatar
	if followerId == 0 {
		return nickname, avatarName
	}
	if nickname != "" && avatarName != "" {
		return nickname, avatarName
	}
	u := userinfodao.FindUserInfoByUserId(followerId)
	if u == nil {
		return nickname, avatarName
	}
	if nickname == "" {
		nickname = u.Nickname
	}
	if avatarName == "" {
		avatarName = u.Avatar
	}
	return nickname, avatarName
}