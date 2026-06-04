package aliyunmoderation

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/aliyuntextmoderationcfgdao"
	"xr-game-server/entity"
)

type cfgSnapshot struct {
	Enabled         bool
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
	Endpoint        string
	ChatService     string
	NicknameService string
	CommentService  string
}

var (
	cfgCache         atomic.Value // *cfgSnapshot
	emptyCfgSnapshot = &cfgSnapshot{
		RegionId:        "cn-shanghai",
		Endpoint:        "green-cip.cn-shanghai.aliyuncs.com",
		ChatService:     "chat_detection",
		NicknameService: "nickname_detection",
		CommentService:  "comment_detection",
	}
)

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(aliyuntextmoderationcfgdao.Load()))
}

func getCfgCache() *cfgSnapshot {
	v := cfgCache.Load()
	if v == nil {
		return emptyCfgSnapshot
	}
	cfg, ok := v.(*cfgSnapshot)
	if !ok || cfg == nil {
		return emptyCfgSnapshot
	}
	return cfg
}

func toCfgSnapshot(row *entity.AliyunTextModerationCfg) *cfgSnapshot {
	if row == nil {
		return emptyCfgSnapshot
	}
	s := &cfgSnapshot{
		Enabled:         row.Enabled,
		AccessKeyId:     strings.TrimSpace(row.AccessKeyId),
		AccessKeySecret: strings.TrimSpace(row.AccessKeySecret),
		RegionId:        strings.TrimSpace(row.RegionId),
		Endpoint:        strings.TrimSpace(row.Endpoint),
		ChatService:     strings.TrimSpace(row.ChatService),
		NicknameService: strings.TrimSpace(row.NicknameService),
		CommentService:  strings.TrimSpace(row.CommentService),
	}
	if s.RegionId == "" {
		s.RegionId = emptyCfgSnapshot.RegionId
	}
	if s.Endpoint == "" {
		s.Endpoint = emptyCfgSnapshot.Endpoint
	}
	if s.ChatService == "" {
		s.ChatService = emptyCfgSnapshot.ChatService
	}
	if s.NicknameService == "" {
		s.NicknameService = emptyCfgSnapshot.NicknameService
	}
	if s.CommentService == "" {
		s.CommentService = emptyCfgSnapshot.CommentService
	}
	return s
}

func (s *cfgSnapshot) ready() bool {
	return s.Enabled && s.AccessKeyId != "" && s.AccessKeySecret != "" && s.Endpoint != ""
}
