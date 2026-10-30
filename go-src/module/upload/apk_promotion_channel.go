package upload

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	apkVersionInBaseRE     = regexp.MustCompile(`(?:^|[_\-.\s])v?(\d+(?:\.\d+){1,3})(?:[_\-.\s]|$)`)
	apkVersionAnywhereRE   = regexp.MustCompile(`(\d+(?:\.\d+){1,3})`)
	apkVersionSegmentSanRE = regexp.MustCompile(`[^0-9a-zA-Z._-]+`)
)

const apkPromotionChannelStoredBaseMax = 180

// ApkPromotionChannelStoredName 旧版固定对象名(仅兼容历史数据理解,新上传勿用).
func ApkPromotionChannelStoredName(channelCode string) string {
	channelCode = strings.TrimSpace(channelCode)
	if channelCode == "" {
		return ""
	}
	return joinStoreCategory(StoreCatApkPromotionChannel, channelCode+".apk")
}

// BuildApkPromotionChannelStoredName 每次上传生成带版本标识的唯一对象名(CDN 不命中旧包).
// originalUploadFileName 为 CMS 上传时的原始文件名,用于解析版本号(如 app-1.2.3.apk).
func BuildApkPromotionChannelStoredName(channelCode, originalUploadFileName string) string {
	channelCode = strings.TrimSpace(channelCode)
	if channelCode == "" {
		return ""
	}
	version := extractApkVersionLabel(originalUploadFileName)
	if version == "" {
		version = "build"
	}
	now := time.Now().UTC()
	ts := now.Format("20060102150405")
	ms := now.Nanosecond() / 1e6
	base := fmt.Sprintf("%s_v%s_%s%03d", channelCode, version, ts, ms)
	base = truncateApkPromotionStoredBase(base)
	return joinStoreCategory(StoreCatApkPromotionChannel, base+".apk")
}

func extractApkVersionLabel(originalUploadFileName string) string {
	base := strings.TrimSpace(originalUploadFileName)
	if base == "" {
		return ""
	}
	base = strings.TrimSuffix(filepath.Base(base), filepath.Ext(base))
	if m := apkVersionInBaseRE.FindStringSubmatch(base); len(m) > 1 {
		return sanitizeApkVersionSegment(m[1])
	}
	if m := apkVersionAnywhereRE.FindStringSubmatch(base); len(m) > 1 {
		return sanitizeApkVersionSegment(m[1])
	}
	return ""
}

func sanitizeApkVersionSegment(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = apkVersionSegmentSanRE.ReplaceAllString(raw, "_")
	raw = strings.Trim(raw, "._-")
	if raw == "" {
		return ""
	}
	if len(raw) > 32 {
		raw = raw[:32]
	}
	return raw
}

func truncateApkPromotionStoredBase(base string) string {
	base = strings.TrimSpace(base)
	if len(base) <= apkPromotionChannelStoredBaseMax {
		return base
	}
	return base[:apkPromotionChannelStoredBaseMax]
}

// PublishApkPromotionChannelFile 上传渠道 APK 到云桶(复用 CMS 上传资源配置).
func PublishApkPromotionChannelFile(channelCode, localPath, originalUploadFileName string) (publicURL string, storedName string, err error) {
	if !IsS3Enabled() {
		return "", "", errors.New("s3 is required for apk promotion channel upload")
	}
	storedName = BuildApkPromotionChannelStoredName(channelCode, originalUploadFileName)
	if storedName == "" {
		return "", "", errors.New("invalid channel code")
	}
	if err := putLocalFileToS3(storedName, localPath); err != nil {
		return "", "", err
	}
	return GetUrlByName(storedName), storedName, nil
}

// PublicURLForApkPromotionStoredName 由入库的 storedName 拼公开 URL(启服预热/改 S3 域后重建).
func PublicURLForApkPromotionStoredName(storedName string) string {
	storedName = strings.Trim(strings.ReplaceAll(storedName, "\\", "/"), "/")
	if storedName == "" {
		return ""
	}
	return GetUrlByName(storedName)
}

// DeleteApkPromotionChannelStored 删除云桶对象(渠道删/换包时).
func DeleteApkPromotionChannelStored(storedName string) {
	storedName = strings.Trim(strings.ReplaceAll(storedName, "\\", "/"), "/")
	if storedName == "" {
		return
	}
	deleteObjectFromS3(storedName)
}
