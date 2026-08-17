package simulatorcpukeyword

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/simulatorcpukeyworddao"
	"xr-game-server/entity/cms"
)

// 首次空表时写入的默认模拟器 CPU 关键词
var defaultKeywords = []string{
	"goldfish",
	"ranchu",
	"vbox86",
	"ttvm",
	"genymotion",
	"qemu",
	"virtualbox",
	"bluestacks",
	"nox",
	"memu",
	"ldplayer",
	"mumuvm",
	"android_x86",
}

var keywordCache atomic.Value // []string

func Init() {
	seedDefaultsIfEmpty()
	reloadKeywordMemory()
}

func seedDefaultsIfEmpty() {
	if simulatorcpukeyworddao.CountAll() > 0 {
		return
	}
	for _, kw := range defaultKeywords {
		_ = simulatorcpukeyworddao.Create(&entity.SimulatorCpuKeyword{
			Keyword: kw,
			Remark:  "default",
		})
	}
}

func reloadKeywordMemory() {
	keywordCache.Store(simulatorcpukeyworddao.ListAllKeywords())
}

func getCachedKeywords() []string {
	v := keywordCache.Load()
	if v == nil {
		return nil
	}
	list, ok := v.([]string)
	if !ok {
		return nil
	}
	return list
}

// MatchCPU 上报 cpuModel 是否命中任一关键词(小写模糊匹配)
func MatchCPU(cpuModel string) bool {
	cpu := strings.ToLower(strings.TrimSpace(cpuModel))
	if cpu == "" {
		return false
	}
	for _, kw := range getCachedKeywords() {
		if kw != "" && strings.Contains(cpu, kw) {
			return true
		}
	}
	return false
}
