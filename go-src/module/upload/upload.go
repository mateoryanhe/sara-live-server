package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

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
	".mp4":  {},
	".webm": {},
	".mov":  {},
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

// storeUploadedContent 先落本地,开启 S3 后再 PutObject;返回业务侧存储名
func storeUploadedContent(src io.Reader, category, ext string, maxBytes int64, tooLargeErr error) (storedName, localPath string, err error) {
	baseName := newStoredFileName(ext)
	storedName = storedNameForCategory(category, baseName)
	localPath = localPathForStored(storedName)
	if localPath == "" {
		return "", "", errors.New("invalid stored name")
	}
	if err = os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return "", "", err
	}
	if _, err = copyUploadContent(src, localPath, maxBytes, tooLargeErr); err != nil {
		return "", "", err
	}
	if IsS3Enabled() {
		if err = putLocalFileToS3(storedName, localPath); err != nil {
			_ = os.Remove(localPath)
			return "", "", err
		}
	}
	return storedName, localPath, nil
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
		name, _, storeErr := storeUploadedContent(part, StoreCatCMS, ext, 0, nil)
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
	name, _, err := storeUploadedContent(src, StoreCatCMS, ext, 0, nil)
	return name, err
}

// UploadShortVideoFile 保存短视频文件;maxBytes 为业务大小上限(字节),0 表示不限制
func UploadShortVideoFile(file *ghttp.UploadFile, maxBytes int64) (string, error) {
	if file == nil {
		return "", errors.New("upload file is empty")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedShortVideoExt[ext]; !ok {
		return "", fmt.Errorf("video ext not allowed: %s", ext)
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	name, _, err := storeUploadedContent(src, StoreCatShortVideo, ext, maxBytes, errVideoFileTooLarge)
	return name, err
}

// StreamUploadShortVideoPart 流式保存短视频 multipart 文件字段
func StreamUploadShortVideoPart(part *multipart.Part, maxBytes int64) (string, error) {
	return streamUploadMultipartPart(part, StoreCatShortVideo, allowedShortVideoExt, maxBytes, errVideoFileTooLarge)
}

// StreamUploadImagePart 流式保存图片 multipart 文件字段
func StreamUploadImagePart(part *multipart.Part, maxBytes int64) (string, error) {
	return streamUploadMultipartPart(part, StoreCatImages, allowedImageExt, maxBytes, errImageFileTooLarge)
}

func streamUploadMultipartPart(part *multipart.Part, category string, allowedExt map[string]struct{}, maxBytes int64, tooLargeErr error) (string, error) {
	if part == nil || strings.TrimSpace(part.FileName()) == "" {
		return "", errors.New("upload file is empty")
	}
	ext := strings.ToLower(filepath.Ext(part.FileName()))
	if _, ok := allowedExt[ext]; !ok {
		return "", fmt.Errorf("file ext not allowed: %s", ext)
	}
	name, _, err := storeUploadedContent(part, category, ext, maxBytes, tooLargeErr)
	return name, err
}

var (
	errFileTooLarge      = errors.New("upload file too large")
	errImageFileTooLarge = errors.New("upload image file too large")
	errVideoFileTooLarge = errors.New("upload video file too large")
)

func copyUploadContent(src io.Reader, dstPath string, maxBytes int64, tooLargeErr error) (int64, error) {
	dst, err := os.Create(dstPath)
	if err != nil {
		return 0, err
	}

	reader := io.Reader(src)
	if maxBytes > 0 {
		reader = io.LimitReader(src, maxBytes+1)
	}
	written, copyErr := io.Copy(dst, reader)
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
		if tooLargeErr != nil {
			return 0, tooLargeErr
		}
		return 0, errFileTooLarge
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

// ReadUploadedFileBytes 读取已上传资源文件内容
func ReadUploadedFileBytes(name string) ([]byte, error) {
	safeName := sanitizeStoredRelativePath(name)
	if safeName == "" {
		return nil, errors.New("invalid file name")
	}
	localPath := filepath.Join(getImageDir(), filepath.FromSlash(safeName))
	data, err := os.ReadFile(localPath)
	if err == nil {
		return data, nil
	}
	if IsS3Enabled() {
		return getObjectBytesFromS3(safeName)
	}
	return nil, err
}

// SaveUploadedFileBytes 按原文件名写入资源文件(用于跨环境同步)
func SaveUploadedFileBytes(name string, data []byte) error {
	safeName := sanitizeStoredRelativePath(name)
	if safeName == "" {
		return errors.New("invalid file name")
	}
	localPath := filepath.Join(getImageDir(), filepath.FromSlash(safeName))
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return err
	}
	if IsS3Enabled() {
		return putBytesToS3(safeName, data)
	}
	return nil
}

// DeleteUploadedFile 删除本地与云桶中的资源;无效文件名忽略
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
