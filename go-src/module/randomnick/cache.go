package randomnick

import (
	"crypto/rand"
	"math/big"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/randomnickdao"
)

type nicknamePoolSnapshot struct {
	useDB    bool
	byLang   map[uint8][]string
	fallback []string
}

var poolCache atomic.Value // *nicknamePoolSnapshot

func Init() {
	if err := reloadMemory(); err != nil {
		g.Log().Errorf(gctx.New(), "randomnick reload failed: %v", err)
	}
}

// ReloadMemory CMS 导入后刷新内存
func ReloadMemory() error {
	return reloadMemory()
}

func reloadMemory() error {
	count, err := randomnickdao.CountAll()
	if err != nil {
		return err
	}
	snap := &nicknamePoolSnapshot{
		byLang:   make(map[uint8][]string),
		fallback: buildBuiltinEnglishNicknames(),
	}
	if count == 0 {
		snap.useDB = false
		snap.byLang[LangEN] = snap.fallback
		poolCache.Store(snap)
		return nil
	}
	rows, err := randomnickdao.LoadAll()
	if err != nil {
		return err
	}
	snap.useDB = true
	for _, row := range rows {
		if row == nil || row.Nickname == "" {
			continue
		}
		lang := NormalizeLang(row.Lang)
		snap.byLang[lang] = append(snap.byLang[lang], row.Nickname)
	}
	poolCache.Store(snap)
	return nil
}

func getPoolSnapshot() *nicknamePoolSnapshot {
	v := poolCache.Load()
	if v == nil {
		fb := buildBuiltinEnglishNicknames()
		return &nicknamePoolSnapshot{
			byLang:   map[uint8][]string{LangEN: fb},
			fallback: fb,
		}
	}
	snap, ok := v.(*nicknamePoolSnapshot)
	if !ok || snap == nil {
		fb := buildBuiltinEnglishNicknames()
		return &nicknamePoolSnapshot{
			byLang:   map[uint8][]string{LangEN: fb},
			fallback: fb,
		}
	}
	return snap
}

func poolForLang(lang uint8) []string {
	lang = NormalizeLang(lang)
	snap := getPoolSnapshot()
	if list, ok := snap.byLang[lang]; ok && len(list) > 0 {
		return list
	}
	if list, ok := snap.byLang[LangEN]; ok && len(list) > 0 {
		return list
	}
	return snap.fallback
}

// UseDB 当前是否使用数据库昵称库
func UseDB() bool {
	return getPoolSnapshot().useDB
}

// Count 指定语言昵称数量(内存)
func CountByLang(lang uint8) int {
	return len(poolForLang(lang))
}

// CountAllLangs 各语言数量(内存)
func CountAllLangs() map[uint8]int {
	snap := getPoolSnapshot()
	out := make(map[uint8]int, len(SupportedLangs()))
	for _, lang := range SupportedLangs() {
		if list, ok := snap.byLang[lang]; ok {
			out[lang] = len(list)
		}
	}
	return out
}

// PickRandom 按语言随机昵称(只读内存)
func PickRandom(lang uint8) string {
	list := poolForLang(lang)
	if len(list) == 0 {
		return "LuckyStar"
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(list))))
	if err != nil {
		return list[0]
	}
	return list[n.Int64()]
}

// PickRandomDefault 默认英文
func PickRandomDefault() string {
	return PickRandom(DefaultLang)
}

// SampleNicknames 预览若干昵称(内存)
func SampleNicknames(lang uint8, limit int) []string {
	list := poolForLang(lang)
	if limit <= 0 || len(list) == 0 {
		return nil
	}
	if limit > len(list) {
		limit = len(list)
	}
	out := make([]string, limit)
	copy(out, list[:limit])
	return out
}
