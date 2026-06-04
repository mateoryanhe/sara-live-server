package sensitiveword

import (
	"sync"

	"xr-game-server/constants/lang"
)

var (
	storeMu sync.RWMutex
	// wordsByLang 语言 -> 铭感词列表(小写归一化后的匹配用副本见 filter)
	wordsByLang = map[lang.Lang][]string{}
)

func setWords(l lang.Lang, words []string) {
	storeMu.Lock()
	defer storeMu.Unlock()
	wordsByLang[l] = words
}

func getWords(l lang.Lang) []string {
	storeMu.RLock()
	defer storeMu.RUnlock()
	if w, ok := wordsByLang[l]; ok && len(w) > 0 {
		return w
	}
	if w, ok := wordsByLang[lang.DefaultLang]; ok {
		return w
	}
	return nil
}
