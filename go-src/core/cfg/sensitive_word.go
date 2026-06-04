package cfg

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const sensitiveWordCfgKey = "sensitiveWord"

// DefaultEnglishSensitiveWordURL 默认英文铭感词库(每行一词, MIT)
const DefaultEnglishSensitiveWordURL = "https://raw.githubusercontent.com/readme-SVG/Banned-words/main/Banned-words-list/en.txt"

// SensitiveWordCfg 铭感词库配置
type SensitiveWordCfg struct {
	// LocalDir 本地词库目录(相对进程工作目录或绝对路径)
	LocalDir string `json:"localDir"`
	// DownloadUrls 各语言词库下载地址,key 为语言代码 zh-CN / zh-TW / en
	DownloadUrls map[string]string `json:"downloadUrls"`
}

var sensitiveWordCfg *SensitiveWordCfg

func GetSensitiveWordCfg() *SensitiveWordCfg {
	return sensitiveWordCfg
}

func initSensitiveWordCfg() {
	sensitiveWordCfg = &SensitiveWordCfg{
		LocalDir: "data/sensitive_words",
		DownloadUrls: map[string]string{
			"en": DefaultEnglishSensitiveWordURL,
		},
	}
	data, err := g.Cfg().Get(gctx.New(), sensitiveWordCfgKey)
	if err != nil || data.IsNil() {
		return
	}
	_ = data.Scan(sensitiveWordCfg)
	if sensitiveWordCfg.LocalDir == "" {
		sensitiveWordCfg.LocalDir = "data/sensitive_words"
	}
	if sensitiveWordCfg.DownloadUrls == nil {
		sensitiveWordCfg.DownloadUrls = map[string]string{}
	}
	if strings.TrimSpace(sensitiveWordCfg.DownloadUrls["en"]) == "" {
		sensitiveWordCfg.DownloadUrls["en"] = DefaultEnglishSensitiveWordURL
	}
}
