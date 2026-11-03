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
	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/dto/officialsitedeploydto"
	"xr-game-server/dto/promotionchannelurldto"
	"xr-game-server/module/promotionchannelurl"
	"xr-game-server/module/upload"
)

var promotionChannelUrlUploadTarget = deployTarget{
	prefix:      "/_upload/promotion-channel-url",
	defaultPath: "",
}

func InitPromotionChannelUrlApkUpload(_ context.Context, req *promotionchannelurldto.InitPromotionChannelUrlApkUploadReq) (*promotionchannelurldto.InitPromotionChannelUrlApkUploadRes, error) {
	if req == nil || req.ChannelID == 0 {
		return nil, errors.New("invalid channel")
	}
	row := promotionchannelurldao.GetById(req.ChannelID)
	if row == nil {
		return nil, errors.New("channel not found")
	}
	if row.ChannelCode <= 0 {
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
		TargetPrefix:            promotionChannelUrlUploadTarget.prefix,
		FileName:                fileName,
		FileSize:                req.FileSize,
		FileType:                fileType,
		ChunkSize:               officialSiteUploadChunkSize,
		TotalChunk:              totalChunk,
		CreatedAt:               time.Now().Unix(),
		PromotionChannelURLCode: row.ChannelCode,
	}
	if err := saveUploadSession(root, meta); err != nil {
		return nil, err
	}
	return &promotionchannelurldto.InitPromotionChannelUrlApkUploadRes{
		UploadId:   uploadId,
		ChunkSize:  officialSiteUploadChunkSize,
		TotalChunk: totalChunk,
	}, nil
}

func UploadPromotionChannelUrlApkChunkFromRequest(r *ghttp.Request) (*officialsitedeploydto.UploadOfficialSiteFileChunkRes, error) {
	return uploadSiteFileChunkFromRequest(r, promotionChannelUrlUploadTarget)
}

func CompletePromotionChannelUrlApkUpload(ctx context.Context, req *promotionchannelurldto.CompletePromotionChannelUrlApkUploadReq) (*promotionchannelurldto.CompletePromotionChannelUrlApkUploadRes, error) {
	if req == nil || req.ChannelID == 0 {
		return nil, errors.New("invalid request")
	}
	row := promotionchannelurldao.GetById(req.ChannelID)
	if row == nil {
		return nil, errors.New("channel not found")
	}
	if row.ChannelCode <= 0 {
		return nil, errors.New("invalid channel code")
	}
	assembledPath, metaCode, originalFileName, err := mergePromotionChannelUrlUpload(req.UploadId)
	if err != nil {
		return nil, err
	}
	defer os.Remove(assembledPath)
	if metaCode != row.ChannelCode {
		return nil, errors.New("channel code mismatch")
	}
	oldStored := strings.TrimSpace(row.ApkFileName)
	publicURL, storedName, err := upload.PublishPromotionChannelUrlApkFile(row.ChannelCode, assembledPath, originalFileName)
	if err != nil {
		return nil, err
	}
	if oldStored != "" && oldStored != storedName {
		upload.DeletePromotionChannelUrlApkStored(oldStored)
	}
	row.ApkFileName = storedName
	row.ApkDownloadURL = publicURL
	if err := promotionchannelurldao.Update(row); err != nil {
		return nil, err
	}
	promotionchannelurl.SetCachedApkURL(row.ChannelCode, publicURL)
	return &promotionchannelurldto.CompletePromotionChannelUrlApkUploadRes{
		ApkDownloadURL: publicURL,
		ApkFileName:    storedName,
	}, nil
}

func AbortPromotionChannelUrlApkUpload(_ context.Context, req *promotionchannelurldto.AbortPromotionChannelUrlApkUploadReq) (*promotionchannelurldto.AbortPromotionChannelUrlApkUploadRes, error) {
	if req == nil {
		return nil, errors.New("invalid upload id")
	}
	if err := abortSiteFileUpload(req.UploadId, promotionChannelUrlUploadTarget); err != nil {
		return nil, err
	}
	return &promotionchannelurldto.AbortPromotionChannelUrlApkUploadRes{}, nil
}

func mergePromotionChannelUrlUpload(uploadId string) (assembledPath string, channelCode int, originalFileName string, err error) {
	root, meta, err := loadUploadSession(uploadId, promotionChannelUrlUploadTarget)
	if err != nil {
		return "", 0, "", err
	}
	channelCode = meta.PromotionChannelURLCode
	if channelCode <= 0 {
		return "", 0, "", errors.New("invalid channel code in session")
	}
	originalFileName = strings.TrimSpace(meta.FileName)
	sessionDir := uploadSessionDir(root, meta.UploadId)
	unlock, err := lockUploadSession(sessionDir)
	if err != nil {
		return "", 0, "", err
	}
	defer unlock()

	assembled, err := os.CreateTemp(root, "promotion-channel-url-*.apk")
	if err != nil {
		return "", 0, "", err
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
				return "", 0, "", fmt.Errorf("chunk %d has not been uploaded", chunkIndex)
			}
			return "", 0, "", statErr
		}
		expectedSize := expectedChunkSize(meta, chunkIndex)
		if info.Size() != expectedSize {
			return "", 0, "", fmt.Errorf("invalid chunk %d size", chunkIndex)
		}
		chunkFile, openErr := os.Open(chunkPath)
		if openErr != nil {
			return "", 0, "", openErr
		}
		written, copyErr := io.Copy(assembled, chunkFile)
		_ = chunkFile.Close()
		if copyErr != nil {
			return "", 0, "", copyErr
		}
		mergedSize += written
	}
	if mergedSize != meta.FileSize {
		return "", 0, "", fmt.Errorf("invalid merged file size: expected %d, got %d", meta.FileSize, mergedSize)
	}
	if err = assembled.Close(); err != nil {
		return "", 0, "", err
	}
	if err = os.RemoveAll(sessionDir); err != nil {
		return "", 0, "", err
	}
	return assembledPath, channelCode, originalFileName, nil
}
