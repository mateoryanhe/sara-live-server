package entity

import (
	"encoding/json"
	"strings"
)

const (
	AudienceOnlineAnchorGreetContentMaxItems    = 30
	AudienceOnlineAnchorGreetContentItemMaxLen  = 1024
)

// ParseAudienceGreetContentListJSON 解析 CMS 存的招呼列表 JSON；非法或非数组时回落为单条明文.
func ParseAudienceGreetContentListJSON(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if !strings.HasPrefix(raw, "[") {
		if t := strings.TrimSpace(raw); t != "" {
			return []string{t}
		}
		return nil
	}
	var list []string
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, s := range list {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

// EncodeAudienceGreetContentListJSON 序列化招呼列表.
func EncodeAudienceGreetContentListJSON(list []string) string {
	if len(list) == 0 {
		return ""
	}
	b, err := json.Marshal(list)
	if err != nil {
		return ""
	}
	return string(b)
}

func mergeAudienceGreetContentList(jsonField, legacySingle string, legacyExtra ...string) []string {
	fromJSON := ParseAudienceGreetContentListJSON(jsonField)
	if len(fromJSON) > 0 {
		return fromJSON
	}
	for _, leg := range append([]string{legacySingle}, legacyExtra...) {
		if t := strings.TrimSpace(leg); t != "" {
			return []string{t}
		}
	}
	return nil
}

// ResolvedContentsEn 英文招呼列表(含历史单字段回落).
func (c *AudienceOnlineAnchorGreetCfg) ResolvedContentsEn() []string {
	if c == nil {
		return nil
	}
	return mergeAudienceGreetContentList(c.ContentsEn, c.ContentEn, c.Content)
}

func (c *AudienceOnlineAnchorGreetCfg) ResolvedContentsEs() []string {
	if c == nil {
		return nil
	}
	return mergeAudienceGreetContentList(c.ContentsEs, c.ContentEs)
}

func (c *AudienceOnlineAnchorGreetCfg) ResolvedContentsPt() []string {
	if c == nil {
		return nil
	}
	return mergeAudienceGreetContentList(c.ContentsPt, c.ContentPt)
}

func (c *AudienceOnlineAnchorGreetCfg) ResolvedContentsHi() []string {
	if c == nil {
		return nil
	}
	return mergeAudienceGreetContentList(c.ContentsHi, c.ContentHi)
}

func (c *AudienceOnlineAnchorGreetCfg) ResolvedContentsId() []string {
	if c == nil {
		return nil
	}
	return mergeAudienceGreetContentList(c.ContentsId, c.ContentId)
}