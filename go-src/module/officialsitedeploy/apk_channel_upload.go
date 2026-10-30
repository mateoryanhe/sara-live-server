package officialsitedeploy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/dto/apkpromotionchanneldto"
	"xr-game-server/dto/officialsitedeploydto"
	"xr-game-server/module/apkpromotionchannel"
	"xr-game-server/module/upload"
)

// 分片上传会话标识(仅内存/临时目录,不写官网静态目录).
var apkPromotionChannelUploadTarget = deployTarget{
	prefix:     "/_upload/apk-promotion-channel",
	defaultPath: "",
}

// InitApkPromotionChannelApkUpload 渠道 APK 分片上传(合并后写入云桶 apk-promotion-channel/).
func InitApkPromotionChannelApkUpload(_ context.Context, req *apkpromotionchanneldto.InitApkPromotionChannelApkUploadReq) (*apkpromotionchanneldto.InitApkPromotionChannelApkUploadRes, error) {
	if req == nil || req.ChannelID == 0 {
		return nil, errors.New("invalid channel")
	}
	row := apkpromotionchanneldao.GetById(req.ChannelID)
	if row == nil {
		return nil, errors.New("channel not found")
	}
	channelCode := strings.TrimSpace(row.ChannelCode)
	if channelCode == "" {
		return nil, errors.New("invalid channel code")
	}
	fileName, fileType, err := validateOfficialSiteUploadFile(req.FileName, req.FileSize)
	if err != nil {
		return nil, err
	}
	if fileType != "apk" {
		return nil, errors.New("apk only")
	}
	root, err := uploadSessionRoot()
	if err != nil {
		return nil, err
	}
	cleanupExpiredUploadSessions(root)
	uploadId, err := newUploadId()
	if err != nil {
		return nil, err
	}
	totalChunk := 1 + (req.FileSize-1)/officialSiteUploadChunkSize
	meta := &uploadSessionMeta{
		UploadId:                uploadId,
		TargetPrefix:            apkPromotionChannelUploadTarget.prefix,
		FileName:                fileName,
		FileSize:                req.FileSize,
		FileType:                fileType,
		ChunkSize:               officialSiteUploadChunkSize,
		TotalChunk:              totalChunk,
		CreatedAt:               time.Now().Unix(),
		ApkPromotionChannelCode: channelCode,
	}
	if err := saveUploadSession(root, meta); err != nil {
		return nil, err
	}
	return &apkpromotionchanneldto.InitApkPromotionChannelApkUploadRes{
		UploadId:   uploadId,
		ChunkSize:  officialSiteUploadChunkSize,
		TotalChunk: totalChunk,
	}, nil
}

func UploadApkPromotionChannelApkChunkFromRequest(r *ghttp.Request) (*officialsitedeploydto.UploadOfficialSiteFileChunkRes, error) {
	return uploadSiteFileChunkFromRequest(r, apkPromotionChannelUploadTarget)
}

func CompleteApkPromotionChannelApkUpload(ctx context.Context, req *apkpromotionchanneldto.CompleteApkPromotionChannelApkUploadReq) (*apkpromotionchanneldto.CompleteApkPromotionChannelApkUploadRes, error) {
	if req == nil || req.ChannelID == 0 {
		return nil, errors.New("invalid request")
	}
	row := apkpromotionchanneldao.GetById(req.ChannelID)
	if row == nil {
		return nil, errors.New("channel not found")
	}
	channelCode := strings.TrimSpace(row.ChannelCode)
	if channelCode == "" {
		return nil, errors.New("invalid channel code")
	}
	assembledPath, metaCode, originalFileName, err := mergeApkPromotionChannelUpload(req.UploadId)
	if err != nil {
		return nil, err
	}
	defer os.Remove(assembledPath)
	if metaCode != channelCode {
		return nil, errors.New("channel code mismatch")
	}
	oldStored := strings.TrimSpace(row.ApkFileName)
	publicURL, storedName, err := upload.PublishApkPromotionChannelFile(channelCode, assembledPath, originalFileName)
	if err != nil {
		return nil, err
	}
	if oldStored != "" && oldStored != storedName {
		upload.DeleteApkPromotionChannelStored(oldStored)
	}
	row.ApkFileName = storedName
	row.ApkDownloadURL = publicURL
	if err := apkpromotionchanneldao.Update(row); err != nil {
		return nil, err
	}
	apkpromotionchannel.SetCachedApkURL(channelCode, publicURL)
	return &apkpromotionchanneldto.CompleteApkPromotionChannelApkUploadRes{
		ApkDownloadURL: publicURL,
		ApkFileName:    storedName,
	}, nil
}

func AbortApkPromotionChannelApkUpload(_ context.Context, req *apkpromotionchanneldto.AbortApkPromotionChannelApkUploadReq) (*apkpromotionchanneldto.AbortApkPromotionChannelApkUploadRes, error) {
	if req == nil {
		return nil, errors.New("invalid upload id")
	}
	if err := abortSiteFileUpload(req.UploadId, apkPromotionChannelUploadTarget); err != nil {
		return nil, err
	}
	return &apkpromotionchanneldto.AbortApkPromotionChannelApkUploadRes{}, nil
}

func mergeApkPromotionChannelUpload(uploadId string) (assembledPath, channelCode, originalFileName string, err error) {
	root, meta, err := loadUploadSession(uploadId, apkPromotionChannelUploadTarget)
	if err != nil {
		return "", "", "", err
	}
	channelCode = strings.TrimSpace(meta.ApkPromotionChannelCode)
	if channelCode == "" {
		return "", "", "", errors.New("invalid channel code in session")
	}
	originalFileName = strings.TrimSpace(meta.FileName)
	sessionDir := uploadSessionDir(root, meta.UploadId)
	unlock, err := lockUploadSession(sessionDir)
	if err != nil {
		return "", "", "", err
	}
	defer unlock()

	// 合并文件不能放在 sessionDir 内,否则末尾 RemoveAll(sessionDir) 会删掉待上传 S3 的临时 APK.
	assembled, err := os.CreateTemp(root, "apk-promotion-channel-*.apk")
	if err != nil {
		return "", "", "", err
	}
	assembledPath = assembled.Name()
	defer func() {
		if err != nil {
			_ = assembled.Close()
			os.Remove(assembledPath)
		}
	}()

	var mergedSize int64
	for chunkIndex := int64(0); chunkIndex < meta.TotalChunk; chunkIndex++ {
		chunkPath := uploadChunkPath(sessionDir, chunkIndex)
		info, statErr := os.Stat(chunkPath)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				return "", "", "", fmt.Errorf("chunk %d has not been uploaded", chunkIndex)
			}
			return "", "", "", statErr
		}
		expectedSize := expectedChunkSize(meta, chunkIndex)
		if info.Size() != expectedSize {
			return "", "", "", fmt.Errorf("invalid chunk %d size", chunkIndex)
		}
		chunkFile, openErr := os.Open(chunkPath)
		if openErr != nil {
			return "", "", "", openErr
		}
		written, copyErr := io.Copy(assembled, chunkFile)
		_ = chunkFile.Close()
		if copyErr != nil {
			return "", "", "", copyErr
		}
		mergedSize += written
	}
	if mergedSize != meta.FileSize {
		return "", "", "", fmt.Errorf("invalid merged file size: expected %d, got %d", meta.FileSize, mergedSize)
	}
	if err = assembled.Close(); err != nil {
		return "", "", "", err
	}
	if err = os.RemoveAll(sessionDir); err != nil {
		return "", "", "", err
	}
	return assembledPath, channelCode, originalFileName, nil
}
