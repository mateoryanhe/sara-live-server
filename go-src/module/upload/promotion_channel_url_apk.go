package upload

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// PublishPromotionChannelUrlApkFile 推广渠道 URL(H5 ch) APK 上传至云桶.
func PublishPromotionChannelUrlApkFile(channelCode int, localPath, originalUploadFileName string) (publicURL string, storedName string, err error) {
	if channelCode <= 0 {
		return "", "", errors.New("invalid channel code")
	}
	if !IsS3Enabled() {
		return "", "", errors.New("s3 is required for promotion channel apk upload")
	}
	codeKey := strconv.Itoa(channelCode)
	storedName = buildPromotionChannelUrlApkStoredName(codeKey, originalUploadFileName)
	if storedName == "" {
		return "", "", errors.New("invalid channel code")
	}
	if err := putLocalFileToS3(storedName, localPath); err != nil {
		return "", "", err
	}
	return GetUrlByName(storedName), storedName, nil
}

func buildPromotionChannelUrlApkStoredName(channelCodeKey, originalUploadFileName string) string {
	channelCodeKey = strings.TrimSpace(channelCodeKey)
	if channelCodeKey == "" {
		return ""
	}
	version := extractApkVersionLabel(originalUploadFileName)
	if version == "" {
		version = "build"
	}
	now := time.Now().UTC()
	ts := now.Format("20060102150405")
	ms := now.Nanosecond() / 1e6
	base := fmt.Sprintf("pc%s_v%s_%s%03d", channelCodeKey, version, ts, ms)
	base = truncateApkPromotionStoredBase(base)
	return joinStoreCategory(StoreCatPromotionChannelUrlApk, base+".apk")
}

// PublicURLForPromotionChannelUrlApkStoredName 由 storedName 拼公开 URL.
func PublicURLForPromotionChannelUrlApkStoredName(storedName string) string {
	storedName = strings.Trim(strings.ReplaceAll(storedName, "\\", "/"), "/")
	if storedName == "" {
		return ""
	}
	return GetUrlByName(storedName)
}

// DeletePromotionChannelUrlApkStored 删除云桶对象.
func DeletePromotionChannelUrlApkStored(storedName string) {
	storedName = strings.Trim(strings.ReplaceAll(storedName, "\\", "/"), "/")
	if storedName == "" {
		return
	}
	deleteObjectFromS3(storedName)
}
