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

// UploadImage 保存单张图片到 <serverRoot>/upload/images,返回保存后的文件名
func UploadImage(file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", errors.New("upload file is empty")
	}
	if file.Size > MaxImageSize {
		return "", fmt.Errorf("image too large, max=%d bytes", MaxImageSize)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedImageExt[ext]; !ok {
		return "", fmt.Errorf("image ext not allowed: %s", ext)
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

// getImageDir 计算图片保存的绝对目录,优先使用 server.serverRoot 配置
func getImageDir() string {
	ctx := gctx.New()
	root, _ := g.Cfg().Get(ctx, "server.serverRoot")
	base := root.String()
	if base == "" {
		base = "."
	}
	return filepath.Join(base, ImageSubDir)
}
