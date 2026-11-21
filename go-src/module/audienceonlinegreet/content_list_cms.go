package audienceonlinegreet

import (
	"strings"
	"unicode/utf8"

	messageentity "xr-game-server/entity/message"
)

func normalizeGreetContentList(list []string) []string {
	if len(list) == 0 {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, s := range list {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if utf8.RuneCountInString(s) > messageentity.AudienceOnlineAnchorGreetContentItemMaxLen {
			s = truncateRunes(s, messageentity.AudienceOnlineAnchorGreetContentItemMaxLen)
		}
		out = append(out, s)
		if len(out) >= messageentity.AudienceOnlineAnchorGreetContentMaxItems {
			break
		}
	}
	return out
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max])
}
