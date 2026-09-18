package upload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
)

// allowedImageExt 允许的图片扩展名
var allowedImageExt = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".webp": {},
	".bmp":  {},
}

// allowedCMSExt CMS后台允许的扩展名(图片 + 礼物动画资源)
var allowedCMSExt = map[string]struct{}{
	".jpg":    {},
	".jpeg":   {},
	".png":    {},
	".gif":    {},
	".webp":   {},
	".bmp":    {},
	".apng":   {},
	".svga":   {},
	".pag":    {},
	".json":   {},
	".lottie": {},
	".mp4":    {},
	".webm":   {},
	".mov":    {},
	".zip":    {},
}

var allowedShortVideoExt = map[string]struct{}{
	".mp4": {},
}

// IsShortVideoFileNameAllowed 短视频主文件和试看文件统一只允许 MP4.
func IsShortVideoFileNameAllowed(name string) bool {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(name)))
	_, ok := allowedShortVideoExt[ext]
	return ok
}

// UploadImage 保存单张图片(CMS 等后台使用,不做内容审核)
func UploadImage(file *ghttp.UploadFile) (string, error) {
	name, _, err := saveUploadedImageFile(file, 0)
	return name, err
}

func getImageDir() string {
	return GetStoragePath()
}

func getCMSDir() string {
	return getImageDir()
}

func GetUrlByName(name string) string {
	return buildImageResourceUrl(name)
}

func newStoredFileName(ext string) string {
	return guid.S() + ext
}

func storedNameForCategory(category, baseName string) string {
	if IsS3Enabled() {
		return joinStoreCategory(category, baseName)
	}
	return baseName
}

func localPathForStored(storedName string) string {
	safe := sanitizeStoredRelativePath(storedName)
	if safe == "" {
		return ""
	}
	return filepath.Join(getImageDir(), filepath.FromSlash(safe))
}

// storeUploadedContent 保存长度未知的流。
func storeUploadedContent(src io.Reader, category, ext string, maxBytes int64, tooLargeErr error) (storedName, localPath string, err error) {
	return storeUploadedContentContext(context.Background(), src, -1, category, ext, maxBytes, tooLargeErr)
}

// storeUploadedContentWithSize 保存长度已知的流。云存储可直接 PutObject，避免 Uploader 的 5MB 分片缓冲。
func storeUploadedContentWithSize(ctx context.Context, src io.Reader, size int64, category, ext string, maxBytes int64, tooLargeErr error) (storedName, localPath string, err error) {
	return storeUploadedContentContext(ctx, src, size, category, ext, maxBytes, tooLargeErr)
}

// storeUploadedContentContext 按云桶开关写入：开则流式写云，关则流式写本地。
// size < 0 表示未知长度；已知长度会在上传前校验大小并直接使用 PutObject。
func storeUploadedContentContext(ctx context.Context, src io.Reader, size int64, category, ext string, maxBytes int64, tooLargeErr error) (storedName, localPath string, err error) {
	if src == nil {
		return "", "", errUploadFileEmpty
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if size == 0 {
		return "", "", errUploadFileEmpty
	}
	if maxBytes > 0 && size >= 0 && size > maxBytes {
		return "", "", uploadTooLargeError(tooLargeErr)
	}
	baseName := newStoredFileName(ext)
	storedName = storedNameForCategory(category, baseName)
	if IsS3Enabled() {
		if size >= 0 {
			if err = putSizedStreamToS3(ctx, storedName, src, size); err != nil {
				deleteObjectFromS3(storedName)
				return "", "", mapUploadReadErr(err)
			}
			return storedName, "", nil
		}
		counter := &countReader{r: src}
		body := io.Reader(counter)
		if maxBytes > 0 {
			body = io.LimitReader(counter, maxBytes+1)
		}
		if err = putStreamToS3(ctx, storedName, body); err != nil {
			deleteObjectFromS3(storedName)
			return "", "", mapUploadReadErr(err)
		}
		if maxBytes > 0 && counter.n > maxBytes {
			deleteObjectFromS3(storedName)
			return "", "", uploadTooLargeError(tooLargeErr)
		}
		if counter.n == 0 {
			deleteObjectFromS3(storedName)
			return "", "", errUploadFileEmpty
		}
		return storedName, "", nil
	}
	localPath = localPathForStored(storedName)
	if localPath == "" {
		return "", "", errors.New("invalid stored name")
	}
	if err = os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return "", "", err
	}
	dst, createErr := os.Create(localPath)
	if createErr != nil {
		return "", "", createErr
	}
	if _, err = copyUploadContentToFile(src, dst, localPath, maxBytes, tooLargeErr); err != nil {
		return "", "", err
	}
	return storedName, localPath, nil
}

func uploadTooLargeError(tooLargeErr error) error {
	if tooLargeErr != nil {
		return tooLargeErr
	}
	return errFileTooLarge
}

// countReader 统计已读字节,供云直传后校验大小上限
type countReader struct {
	r io.Reader
	n int64
}

func (c *countReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}

// UploadCMSFileFromRequest 流式读取 multipart 的 file 字段,不触发 ParseMultipartForm
func UploadCMSFileFromRequest(r *ghttp.Request) (string, error) {
	if r == nil || r.Request == nil {
		return "", errors.New("upload file is empty")
	}
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return "", mapUploadReadErr(err)
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", mapUploadReadErr(err)
		}
		if part.FormName() != "file" {
			part.Close()
			continue
		}
		ext := strings.ToLower(filepath.Ext(part.FileName()))
		if _, ok := allowedCMSExt[ext]; !ok {
			part.Close()
			return "", fmt.Errorf("file ext not allowed: %s", ext)
		}
		name, _, storeErr := storeUploadedContentContext(r.Context(), part, -1, StoreCatCMS, ext, 0, nil)
		part.Close()
		if storeErr != nil {
			return "", mapUploadReadErr(storeErr)
		}
		return name, nil
	}
	return "", errors.New("upload file is empty")
}

func mapUploadReadErr(err error) error {
	if err == nil {
		return nil
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return errors.New("upload request incomplete")
	}
	if strings.Contains(strings.ToLower(err.Error()), "unexpected eof") {
		return errors.New("upload request incomplete")
	}
	return err
}

// UploadCMSFile 保存CMS后台上传文件(兼容旧绑定方式)
func UploadCMSFile(file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", errors.New("upload file is empty")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedCMSExt[ext]; !ok {
		return "", fmt.Errorf("file ext not allowed: %s", ext)
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	name, _, err := storeUploadedContentWithSize(context.Background(), src, file.Size, StoreCatCMS, ext, 0, nil)
	return name, err
}

// UploadShortVideoFile 保存短视频文件;maxBytes 为业务大小上限(字节),0 表示不限制
func UploadShortVideoFile(file *ghttp.UploadFile, maxBytes int64) (string, error) {
	if file == nil {
		return "", errors.New("upload file is empty")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !IsShortVideoFileNameAllowed(file.Filename) {
		return "", fmt.Errorf("video ext not allowed: %s", ext)
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	name, _, err := storeUploadedContentWithSize(context.Background(), src, file.Size, StoreCatShortVideo, ext, maxBytes, errVideoFileTooLarge)
	return name, err
}

// StreamUploadShortVideoPart 流式保存短视频 multipart 文件字段
func StreamUploadShortVideoPart(part *multipart.Part, maxBytes int64) (string, error) {
	return StreamUploadShortVideoPartContext(context.Background(), part, maxBytes)
}

// StreamUploadShortVideoPartContext 流式保存短视频，并将请求取消传递给云存储。
func StreamUploadShortVideoPartContext(ctx context.Context, part *multipart.Part, maxBytes int64) (string, error) {
	return streamUploadMultipartPart(ctx, part, StoreCatShortVideo, allowedShortVideoExt, maxBytes, errVideoFileTooLarge)
}

// StreamUploadImagePart 流式保存图片 multipart 文件字段
func StreamUploadImagePart(part *multipart.Part, maxBytes int64) (string, error) {
	return StreamUploadImagePartContext(context.Background(), part, maxBytes)
}

// StreamUploadImagePartContext 流式保存图片，并将请求取消传递给云存储。
func StreamUploadImagePartContext(ctx context.Context, part *multipart.Part, maxBytes int64) (string, error) {
	return streamUploadMultipartPart(ctx, part, StoreCatImages, allowedImageExt, maxBytes, errImageFileTooLarge)
}

func streamUploadMultipartPart(ctx context.Context, part *multipart.Part, category string, allowedExt map[string]struct{}, maxBytes int64, tooLargeErr error) (string, error) {
	if part == nil || strings.TrimSpace(part.FileName()) == "" {
		return "", errors.New("upload file is empty")
	}
	ext := strings.ToLower(filepath.Ext(part.FileName()))
	if _, ok := allowedExt[ext]; !ok {
		return "", fmt.Errorf("file ext not allowed: %s", ext)
	}
	name, _, err := storeUploadedContentContext(ctx, part, -1, category, ext, maxBytes, tooLargeErr)
	return name, err
}

var (
	errUploadFileEmpty   = errors.New("upload file is empty")
	errFileTooLarge      = errors.New("upload file too large")
	errImageFileTooLarge = errors.New("upload image file too large")
	errVideoFileTooLarge = errors.New("upload video file too large")
)

var uploadCopyBufferPool = sync.Pool{
	New: func() any {
		buffer := make([]byte, 32*1024)
		return &buffer
	},
}

func copyUploadContentToFile(src io.Reader, dst *os.File, dstPath string, maxBytes int64, tooLargeErr error) (int64, error) {
	reader := io.Reader(src)
	if maxBytes > 0 {
		reader = io.LimitReader(src, maxBytes+1)
	}
	buffer := uploadCopyBufferPool.Get().(*[]byte)
	defer uploadCopyBufferPool.Put(buffer)
	written, copyErr := io.CopyBuffer(dst, reader, *buffer)
	closeErr := dst.Close()
	if copyErr != nil {
		os.Remove(dstPath)
		return 0, mapUploadReadErr(copyErr)
	}
	if closeErr != nil {
		os.Remove(dstPath)
		return 0, closeErr
	}
	if maxBytes > 0 && written > maxBytes {
		os.Remove(dstPath)
		return 0, uploadTooLargeError(tooLargeErr)
	}
	if written == 0 {
		os.Remove(dstPath)
		return 0, errUploadFileEmpty
	}
	return written, nil
}

// IsUploadFileTooLarge 判断是否超过上传大小限制
func IsUploadFileTooLarge(err error) bool {
	return err == errFileTooLarge || err == errImageFileTooLarge || err == errVideoFileTooLarge
}

// IsUploadImageFileTooLarge 判断 App 图片是否超过大小限制
func IsUploadImageFileTooLarge(err error) bool {
	return err == errImageFileTooLarge
}

// IsUploadVideoFileTooLarge 判断短视频是否超过大小限制
func IsUploadVideoFileTooLarge(err error) bool {
	return err == errVideoFileTooLarge
}

// sanitizeStoredRelativePath 允许 images/uuid.ext 一类相对路径,禁止 .. 与绝对 URL
func sanitizeStoredRelativePath(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return ""
	}
	name = strings.Trim(strings.ReplaceAll(name, "\\", "/"), "/")
	if name == "" || strings.Contains(name, "..") {
		return ""
	}
	parts := strings.Split(name, "/")
	for _, p := range parts {
		if p == "" || p == "." || p == ".." {
			return ""
		}
	}
	return name
}

func sanitizeStoredFileName(name string) string {
	return sanitizeStoredRelativePath(name)
}

// ReadUploadedFileBytes 读取已上传资源;开云桶优先读云,本地仅作历史兼容回落
func ReadUploadedFileBytes(name string) ([]byte, error) {
	safeName := sanitizeStoredRelativePath(name)
	if safeName == "" {
		return nil, errors.New("invalid file name")
	}
	if IsS3Enabled() {
		data, err := getObjectBytesFromS3(safeName)
		if err == nil {
			return data, nil
		}
	}
	localPath := filepath.Join(getImageDir(), filepath.FromSlash(safeName))
	return os.ReadFile(localPath)
}

// SaveUploadedFileBytes 按原文件名写入资源(跨环境同步等);开云桶只写云
func SaveUploadedFileBytes(name string, data []byte) error {
	safeName := sanitizeStoredRelativePath(name)
	if safeName == "" {
		return errors.New("invalid file name")
	}
	if IsS3Enabled() {
		return putBytesToS3(safeName, data)
	}
	localPath := filepath.Join(getImageDir(), filepath.FromSlash(safeName))
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(localPath, data, 0644)
}

// DeleteUploadedFile 删除资源;开云桶删云对象,并顺带清可能残留的本地文件
func DeleteUploadedFile(name string) {
	if name == "" {
		return
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return
	}
	safeName := sanitizeStoredRelativePath(name)
	if safeName == "" {
		return
	}
	_ = os.Remove(filepath.Join(getImageDir(), filepath.FromSlash(safeName)))
	deleteObjectFromS3(safeName)
}
