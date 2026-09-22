package officialsitedeploy

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/cfg"
	"xr-game-server/dto/officialsitedeploydto"
	"xr-game-server/module/domainsite"
)

const (
	officialSiteUploadChunkSize = int64(8 * 1024 * 1024)
	uploadSessionMaxAge         = 24 * time.Hour
	uploadIdHeader              = "X-Upload-Id"
	chunkIndexHeader            = "X-Chunk-Index"
)

var uploadIdPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

type uploadSessionMeta struct {
	UploadId     string `json:"uploadId"`
	TargetPrefix string `json:"targetPrefix"`
	FileName     string `json:"fileName"`
	FileSize     int64  `json:"fileSize"`
	FileType     string `json:"fileType"`
	ChunkSize    int64  `json:"chunkSize"`
	TotalChunk   int64  `json:"totalChunk"`
	CreatedAt    int64  `json:"createdAt"`
}

func InitOfficialSiteFileUpload(_ context.Context, req *officialsitedeploydto.InitOfficialSiteFileUploadReq) (*officialsitedeploydto.InitOfficialSiteFileUploadRes, error) {
	if req == nil {
		return nil, errors.New("upload request is empty")
	}
	return initSiteFileUpload(req.FileName, req.FileSize, officialSiteTarget)
}

func InitThirdPayOfficialSiteFileUpload(_ context.Context, req *officialsitedeploydto.InitThirdPayOfficialSiteFileUploadReq) (*officialsitedeploydto.InitThirdPayOfficialSiteFileUploadRes, error) {
	if req == nil {
		return nil, errors.New("upload request is empty")
	}
	return initSiteFileUpload(req.FileName, req.FileSize, thirdPayOfficialSiteTarget)
}

func initSiteFileUpload(rawFileName string, fileSize int64, target deployTarget) (*officialsitedeploydto.InitOfficialSiteFileUploadRes, error) {
	fileName, fileType, err := validateOfficialSiteUploadFile(rawFileName, fileSize)
	if err != nil {
		return nil, err
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
	totalChunk := 1 + (fileSize-1)/officialSiteUploadChunkSize
	meta := &uploadSessionMeta{
		UploadId:     uploadId,
		TargetPrefix: target.prefix,
		FileName:     fileName,
		FileSize:     fileSize,
		FileType:     fileType,
		ChunkSize:    officialSiteUploadChunkSize,
		TotalChunk:   totalChunk,
		CreatedAt:    time.Now().Unix(),
	}
	if err := saveUploadSession(root, meta); err != nil {
		return nil, err
	}
	return &officialsitedeploydto.InitOfficialSiteFileUploadRes{
		UploadId:    uploadId,
		ChunkSize:   officialSiteUploadChunkSize,
		TotalChunks: totalChunk,
	}, nil
}

func UploadOfficialSiteFileChunkFromRequest(r *ghttp.Request) (*officialsitedeploydto.UploadOfficialSiteFileChunkRes, error) {
	return uploadSiteFileChunkFromRequest(r, officialSiteTarget)
}

func UploadThirdPayOfficialSiteFileChunkFromRequest(r *ghttp.Request) (*officialsitedeploydto.UploadOfficialSiteFileChunkRes, error) {
	return uploadSiteFileChunkFromRequest(r, thirdPayOfficialSiteTarget)
}

func uploadSiteFileChunkFromRequest(r *ghttp.Request, target deployTarget) (*officialsitedeploydto.UploadOfficialSiteFileChunkRes, error) {
	if r == nil || r.Request == nil {
		return nil, errors.New("upload chunk is empty")
	}
	uploadId := strings.TrimSpace(r.GetHeader(uploadIdHeader))
	chunkIndex, err := strconv.ParseInt(strings.TrimSpace(r.GetHeader(chunkIndexHeader)), 10, 64)
	if err != nil {
		return nil, errors.New("invalid chunk index")
	}
	root, meta, err := loadUploadSession(uploadId, target)
	if err != nil {
		return nil, err
	}
	if chunkIndex < 0 || chunkIndex >= meta.TotalChunk {
		return nil, errors.New("chunk index out of range")
	}
	expectedSize := expectedChunkSize(meta, chunkIndex)
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return nil, mapUploadReadErr(err)
	}

	sessionDir := uploadSessionDir(root, uploadId)
	tmpFile, err := os.CreateTemp(sessionDir, fmt.Sprintf("chunk-%06d-*.tmp", chunkIndex))
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()
	keepTmp := false
	defer func() {
		_ = tmpFile.Close()
		if !keepTmp {
			_ = os.Remove(tmpPath)
		}
	}()

	var written int64
	found := false
	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return nil, mapUploadReadErr(nextErr)
		}
		if part.FormName() != "file" {
			_ = part.Close()
			continue
		}
		found = true
		written, err = io.Copy(tmpFile, io.LimitReader(part, expectedSize+1))
		_ = part.Close()
		if err != nil {
			return nil, mapUploadReadErr(err)
		}
		break
	}
	if !found {
		return nil, errors.New("upload chunk is empty")
	}
	if written != expectedSize {
		return nil, fmt.Errorf("invalid chunk size: expected %d, got %d", expectedSize, written)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	chunkPath := uploadChunkPath(sessionDir, chunkIndex)
	if err := replaceFile(tmpPath, chunkPath); err != nil {
		return nil, err
	}
	keepTmp = true
	return &officialsitedeploydto.UploadOfficialSiteFileChunkRes{
		ChunkIndex: chunkIndex,
		FileSize:   written,
	}, nil
}

func CompleteOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.CompleteOfficialSiteFileUploadReq) (*officialsitedeploydto.CompleteOfficialSiteFileUploadRes, error) {
	if req == nil {
		return nil, errors.New("upload request is empty")
	}
	return completeSiteFileUpload(ctx, req.UploadId, officialSiteTarget)
}

func CompleteThirdPayOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.CompleteThirdPayOfficialSiteFileUploadReq) (*officialsitedeploydto.CompleteThirdPayOfficialSiteFileUploadRes, error) {
	if req == nil {
		return nil, errors.New("upload request is empty")
	}
	return completeSiteFileUpload(ctx, req.UploadId, thirdPayOfficialSiteTarget)
}

func completeSiteFileUpload(ctx context.Context, uploadId string, target deployTarget) (*officialsitedeploydto.CompleteOfficialSiteFileUploadRes, error) {
	root, meta, err := loadUploadSession(uploadId, target)
	if err != nil {
		return nil, err
	}
	sessionDir := uploadSessionDir(root, meta.UploadId)
	unlock, err := lockUploadSession(sessionDir)
	if err != nil {
		return nil, err
	}
	defer unlock()

	deployDir, err := getDeployDir(target)
	if err != nil {
		return nil, err
	}
	assembled, err := os.CreateTemp(deployDir, ".official-site-upload-*.tmp")
	if err != nil {
		return nil, err
	}
	assembledPath := assembled.Name()
	defer os.Remove(assembledPath)

	var mergedSize int64
	for chunkIndex := int64(0); chunkIndex < meta.TotalChunk; chunkIndex++ {
		chunkPath := uploadChunkPath(sessionDir, chunkIndex)
		info, statErr := os.Stat(chunkPath)
		if statErr != nil {
			_ = assembled.Close()
			if os.IsNotExist(statErr) {
				return nil, fmt.Errorf("chunk %d has not been uploaded", chunkIndex)
			}
			return nil, statErr
		}
		expectedSize := expectedChunkSize(meta, chunkIndex)
		if info.Size() != expectedSize {
			_ = assembled.Close()
			return nil, fmt.Errorf("invalid chunk %d size", chunkIndex)
		}
		chunkFile, openErr := os.Open(chunkPath)
		if openErr != nil {
			_ = assembled.Close()
			return nil, openErr
		}
		written, copyErr := io.Copy(assembled, chunkFile)
		_ = chunkFile.Close()
		if copyErr != nil {
			_ = assembled.Close()
			return nil, copyErr
		}
		mergedSize += written
	}
	if mergedSize != meta.FileSize {
		_ = assembled.Close()
		return nil, fmt.Errorf("invalid merged file size: expected %d, got %d", meta.FileSize, mergedSize)
	}
	if err := assembled.Close(); err != nil {
		return nil, err
	}

	res := &officialsitedeploydto.CompleteOfficialSiteFileUploadRes{
		FileType:  meta.FileType,
		FileName:  meta.FileName,
		FileSize:  meta.FileSize,
		UrlPrefix: target.prefix,
	}
	if meta.FileType == "zip" {
		fileCount, dirCount, extractErr := extractZip(assembledPath, deployDir)
		if extractErr != nil {
			return nil, extractErr
		}
		res.FileCount = fileCount
		res.DirCount = dirCount
		res.DeployPath = deployDir
		if _, err := domainsite.RecordDeploySuccess(ctx, target.siteKey); err != nil {
			return nil, fmt.Errorf("files deployed but record upload time failed: %w", err)
		}
	} else {
		apkDir := filepath.Join(deployDir, target.apkSubDir)
		if err := os.MkdirAll(apkDir, 0755); err != nil {
			return nil, err
		}
		destination := filepath.Join(apkDir, meta.FileName)
		if err := os.Chmod(assembledPath, 0644); err != nil {
			return nil, err
		}
		if err := replaceFile(assembledPath, destination); err != nil {
			return nil, err
		}
		res.DeployPath = destination
		res.DownloadUrl = staticSiteFileURL(target, filepath.ToSlash(filepath.Join(target.apkSubDir, meta.FileName)))
	}
	if err := os.RemoveAll(sessionDir); err != nil {
		return nil, err
	}
	return res, nil
}

func AbortOfficialSiteFileUpload(_ context.Context, req *officialsitedeploydto.AbortOfficialSiteFileUploadReq) (*officialsitedeploydto.AbortOfficialSiteFileUploadRes, error) {
	if req == nil {
		return nil, errors.New("invalid upload id")
	}
	if err := abortSiteFileUpload(req.UploadId, officialSiteTarget); err != nil {
		return nil, err
	}
	return &officialsitedeploydto.AbortOfficialSiteFileUploadRes{}, nil
}

func AbortThirdPayOfficialSiteFileUpload(_ context.Context, req *officialsitedeploydto.AbortThirdPayOfficialSiteFileUploadReq) (*officialsitedeploydto.AbortThirdPayOfficialSiteFileUploadRes, error) {
	if req == nil {
		return nil, errors.New("invalid upload id")
	}
	if err := abortSiteFileUpload(req.UploadId, thirdPayOfficialSiteTarget); err != nil {
		return nil, err
	}
	return &officialsitedeploydto.AbortThirdPayOfficialSiteFileUploadRes{}, nil
}

func abortSiteFileUpload(rawUploadId string, target deployTarget) error {
	uploadId := strings.TrimSpace(rawUploadId)
	if !validUploadId(uploadId) {
		return errors.New("invalid upload id")
	}
	root, _, err := loadUploadSession(uploadId, target)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(uploadSessionDir(root, uploadId)); err != nil {
		return err
	}
	return nil
}

func validateOfficialSiteUploadFile(rawFileName string, fileSize int64) (fileName, fileType string, err error) {
	fileName = strings.TrimSpace(rawFileName)
	if fileName == "" || fileSize <= 0 {
		return "", "", errors.New("invalid upload file")
	}
	if strings.ContainsAny(fileName, "/\\\x00") || filepath.Base(fileName) != fileName {
		return "", "", errors.New("invalid upload file name")
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".zip":
		return fileName, "zip", nil
	case ".apk":
		return fileName, "apk", nil
	default:
		return "", "", fmt.Errorf("file ext not allowed: %s", ext)
	}
}

func uploadSessionRoot() (string, error) {
	root := filepath.Join(os.TempDir(), "sara-official-site-file-upload")
	if err := os.MkdirAll(root, 0700); err != nil {
		return "", err
	}
	return root, nil
}

func saveUploadSession(root string, meta *uploadSessionMeta) error {
	sessionDir := uploadSessionDir(root, meta.UploadId)
	if err := os.Mkdir(sessionDir, 0700); err != nil {
		return err
	}
	data, err := json.Marshal(meta)
	if err != nil {
		_ = os.RemoveAll(sessionDir)
		return err
	}
	if err := os.WriteFile(filepath.Join(sessionDir, "metadata.json"), data, 0600); err != nil {
		_ = os.RemoveAll(sessionDir)
		return err
	}
	return nil
}

func loadUploadSession(uploadId string, target deployTarget) (root string, meta *uploadSessionMeta, err error) {
	uploadId = strings.TrimSpace(uploadId)
	if !validUploadId(uploadId) {
		return "", nil, errors.New("invalid upload id")
	}
	root, err = uploadSessionRoot()
	if err != nil {
		return "", nil, err
	}
	data, err := os.ReadFile(filepath.Join(uploadSessionDir(root, uploadId), "metadata.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, errors.New("upload session does not exist")
		}
		return "", nil, err
	}
	meta = &uploadSessionMeta{}
	if err := json.Unmarshal(data, meta); err != nil {
		return "", nil, errors.New("invalid upload session")
	}
	fileName, fileType, validateErr := validateOfficialSiteUploadFile(meta.FileName, meta.FileSize)
	expectedTotalChunk := int64(0)
	if meta.FileSize > 0 {
		expectedTotalChunk = 1 + (meta.FileSize-1)/officialSiteUploadChunkSize
	}
	if validateErr != nil || meta.UploadId != uploadId || meta.TargetPrefix != target.prefix || meta.FileName != fileName || meta.FileType != fileType || meta.ChunkSize != officialSiteUploadChunkSize || meta.TotalChunk != expectedTotalChunk {
		return "", nil, errors.New("invalid upload session")
	}
	return root, meta, nil
}

func newUploadId() (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

func validUploadId(uploadId string) bool {
	return uploadIdPattern.MatchString(strings.TrimSpace(uploadId))
}

func uploadSessionDir(root, uploadId string) string {
	return filepath.Join(root, uploadId)
}

func uploadChunkPath(sessionDir string, chunkIndex int64) string {
	return filepath.Join(sessionDir, fmt.Sprintf("chunk-%06d", chunkIndex))
}

func expectedChunkSize(meta *uploadSessionMeta, chunkIndex int64) int64 {
	start := chunkIndex * meta.ChunkSize
	remaining := meta.FileSize - start
	if remaining < meta.ChunkSize {
		return remaining
	}
	return meta.ChunkSize
}

func lockUploadSession(sessionDir string) (func(), error) {
	lockPath := filepath.Join(sessionDir, "complete.lock")
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsExist(err) {
			return nil, errors.New("upload is already being completed")
		}
		return nil, err
	}
	_ = lockFile.Close()
	return func() { _ = os.Remove(lockPath) }, nil
}

func replaceFile(source, destination string) error {
	renameErr := os.Rename(source, destination)
	if renameErr == nil {
		return nil
	}
	if _, statErr := os.Stat(destination); statErr != nil {
		return renameErr
	}
	backup := fmt.Sprintf("%s.replace-backup-%d", destination, time.Now().UnixNano())
	if err := os.Rename(destination, backup); err != nil {
		return renameErr
	}
	if err := os.Rename(source, destination); err != nil {
		_ = os.Rename(backup, destination)
		return err
	}
	_ = os.Remove(backup)
	return nil
}

func staticSiteFileURL(target deployTarget, relativePath string) string {
	domain := strings.TrimSpace(cfg.GetStaticSiteDomain(target.prefix))
	if domain == "" {
		return ""
	}
	return (&url.URL{Scheme: "https", Host: domain, Path: "/" + strings.TrimPrefix(filepath.ToSlash(relativePath), "/")}).String()
}

func cleanupExpiredUploadSessions(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	expireBefore := time.Now().Add(-uploadSessionMaxAge)
	for _, entry := range entries {
		if !entry.IsDir() || !validUploadId(entry.Name()) {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr == nil && info.ModTime().Before(expireBefore) {
			_ = os.RemoveAll(uploadSessionDir(root, entry.Name()))
		}
	}
}
