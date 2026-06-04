package upload

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"io"
	"os"
	"path/filepath"
	"strings"
	"xr-game-server/core/snowflake"
)

const (
	// ImageSubDir 图片相对 serverRoot 的子目录
	ImageSubDir = "upload/images"
	// MaxImageSize 单张图片最大字节数(5MB)
	MaxImageSize int64 = 5 * 1024 * 1024
	// MaxCMSFileSize CMS后台单个文件最大字节数(100MB),用于短视频等资源
	MaxCMSFileSize int64 = 100 * 1024 * 1024
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
	// 图片
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".webp": {},
	".bmp":  {},
	".apng": {},
	// 动画 / 资源
	".svga":   {},
	".pag":    {},
	".json":   {}, // lottie
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

// getImageDir 计算图片保存的绝对目录,优先使用 server.serverRoot 配置
func getImageDir() string {
	return getUploadDir(ImageSubDir)
}

// getCMSDir 计算CMS上传资源保存的绝对目录
func getCMSDir() string {
	return getUploadDir(ImageSubDir)
}

// getUploadDir 计算上传保存的绝对目录,优先使用 server.serverRoot 配置
func getUploadDir(subDir string) string {
	ctx := gctx.New()
	root, _ := g.Cfg().Get(ctx, "server.serverRoot")
	base := root.String()
	if base == "" {
		base = "."
	}
	return filepath.Join(base, subDir)
}

func GetUrlByName(name string) string {
	return buildResourceUrl(fmt.Sprintf("/%s/%s", ImageSubDir, name))
}

// UploadCMSFile 保存CMS后台上传的图片或礼物动画资源到 <serverRoot>/upload/cms,返回保存后的文件名
func UploadCMSFile(file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", errors.New("upload file is empty")
	}
	if file.Size > MaxCMSFileSize {
		return "", fmt.Errorf("file too large, max=%d bytes", MaxCMSFileSize)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedCMSExt[ext]; !ok {
		return "", fmt.Errorf("file ext not allowed: %s", ext)
	}

	dir := getCMSDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	newName := fmt.Sprintf("%d%s", snowflake.GetId(), ext)
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

// UploadShortVideoFile 保存短视频文件,按传入大小限制校验
func UploadShortVideoFile(file *ghttp.UploadFile, maxSize uint64) (string, error) {
	if file == nil {
		return "", errors.New("upload file is empty")
	}
	if maxSize == 0 {
		maxSize = uint64(MaxCMSFileSize)
	}
	if file.Size > int64(maxSize) {
		return "", fmt.Errorf("file too large, max=%d bytes", maxSize)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedShortVideoExt[ext]; !ok {
		return "", fmt.Errorf("video ext not allowed: %s", ext)
	}

	dir := getImageDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	newName := fmt.Sprintf("%d%s", snowflake.GetId(), ext)
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
