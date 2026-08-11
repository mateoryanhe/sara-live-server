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
	"xr-game-server/core/cfg"
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
	name, _, err := saveUploadedImageFile(file)
	return name, err
}

func getImageDir() string {
	if root := cfg.GetImageStaticRoot(); root != "" {
		return root
	}
	if prefix := cfg.GetImageStaticPrefix(); prefix != "" {
		return cfg.ResolvePhysicalDir(prefix)
	}
	return cfg.GetServerRoot()
}

func getCMSDir() string {
	return getImageDir()
}

func GetUrlByName(name string) string {
	segment := cfg.GetImageStaticPathSegment()
	if segment == "" {
		return buildResourceUrl("/" + name)
	}
	return buildResourceUrl(fmt.Sprintf("/%s/%s", segment, name))
}

func newStoredFileName(ext string) string {
	return guid.S() + ext
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
	dir := getCMSDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
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
		newName := newStoredFileName(ext)
		dstPath := filepath.Join(dir, newName)
		dst, err := os.Create(dstPath)
		if err != nil {
			part.Close()
			return "", err
		}
		_, copyErr := io.Copy(dst, part)
		part.Close()
		closeErr := dst.Close()
		if copyErr != nil {
			os.Remove(dstPath)
			return "", mapUploadReadErr(copyErr)
		}
		if closeErr != nil {
			os.Remove(dstPath)
			return "", closeErr
		}
		return newName, nil
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

	dir := getCMSDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	newName := newStoredFileName(ext)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.Create(filepath.Join(dir, newName))
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return newName, nil
}

// UploadShortVideoFile 保存短视频文件;大小上限由 server.clientMaxBodySize 控制
func UploadShortVideoFile(file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", errors.New("upload file is empty")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedShortVideoExt[ext]; !ok {
		return "", fmt.Errorf("video ext not allowed: %s", ext)
	}

	dir := getImageDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	newName := newStoredFileName(ext)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.Create(filepath.Join(dir, newName))
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return newName, nil
}

// StreamUploadShortVideoPart 流式保存短视频 multipart 文件字段
func StreamUploadShortVideoPart(part *multipart.Part, maxBytes int64) (string, error) {
	return streamUploadMultipartPart(part, allowedShortVideoExt, maxBytes)
}

// StreamUploadImagePart 流式保存图片 multipart 文件字段
func StreamUploadImagePart(part *multipart.Part, maxBytes int64) (string, error) {
	return streamUploadMultipartPart(part, allowedImageExt, maxBytes)
}

func streamUploadMultipartPart(part *multipart.Part, allowedExt map[string]struct{}, maxBytes int64) (string, error) {
	if part == nil || strings.TrimSpace(part.FileName()) == "" {
		return "", errors.New("upload file is empty")
	}
	ext := strings.ToLower(filepath.Ext(part.FileName()))
	if _, ok := allowedExt[ext]; !ok {
		return "", fmt.Errorf("file ext not allowed: %s", ext)
	}

	dir := getImageDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	newName := newStoredFileName(ext)
	dstPath := filepath.Join(dir, newName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}

	reader := io.Reader(part)
	if maxBytes > 0 {
		reader = io.LimitReader(part, maxBytes+1)
	}
	written, copyErr := io.Copy(dst, reader)
	closeErr := dst.Close()
	if copyErr != nil {
		os.Remove(dstPath)
		return "", mapUploadReadErr(copyErr)
	}
	if closeErr != nil {
		os.Remove(dstPath)
		return "", closeErr
	}
	if maxBytes > 0 && written > maxBytes {
		os.Remove(dstPath)
		return "", errFileTooLarge
	}
	return newName, nil
}

var errFileTooLarge = errors.New("upload file too large")

// IsUploadFileTooLarge 判断是否超过上传大小限制
func IsUploadFileTooLarge(err error) bool {
	return err == errFileTooLarge
}

func sanitizeStoredFileName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return ""
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return ""
	}
	return filepath.Base(name)
}

// ReadUploadedFileBytes 读取已上传资源文件内容
func ReadUploadedFileBytes(name string) ([]byte, error) {
	safeName := sanitizeStoredFileName(name)
	if safeName == "" {
		return nil, errors.New("invalid file name")
	}
	return os.ReadFile(filepath.Join(getImageDir(), safeName))
}

// SaveUploadedFileBytes 按原文件名写入资源文件(用于跨环境同步)
func SaveUploadedFileBytes(name string, data []byte) error {
	safeName := sanitizeStoredFileName(name)
	if safeName == "" {
		return errors.New("invalid file name")
	}
	dir := getImageDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, safeName), data, 0644)
}

// DeleteUploadedFile 删除 images 目录下的资源文件;无效文件名或文件不存在时忽略
func DeleteUploadedFile(name string) {
	if name == "" {
		return
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return
	}
	_ = os.Remove(filepath.Join(getImageDir(), name))
}
