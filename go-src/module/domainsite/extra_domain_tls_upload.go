package domainsite

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/dao/domainsitedao"
)

const (
	extraDomainTLSRoot = "/home/ec2-user/ssl/extra-domain"
	siteTLSRoot        = "/home/ec2-user/ssl/site-tls"
)

func extraDomainTLSDir(siteKey string) string {
	return filepath.Join(extraDomainTLSRoot, strings.TrimSpace(siteKey))
}

func siteTLSDir(siteKey string) string {
	return filepath.Join(siteTLSRoot, strings.TrimSpace(siteKey))
}

// storeExtraDomainTLSFiles 将上传的证书/私钥落到站点专属目录,返回绝对路径。
func storeExtraDomainTLSFiles(siteKey string, certPart, keyPart *ghttp.UploadFile) (certPath, keyPath string, err error) {
	return storeTLSFilesToDir(extraDomainTLSDir(siteKey), certPart, keyPart)
}

// StoreSiteTLSFiles 业务静态站(如 Sofie 官网)专用证书目录 /home/ec2-user/ssl/site-tls/{siteKey}/.
func StoreSiteTLSFiles(siteKey string, certPart, keyPart *ghttp.UploadFile) (certPath, keyPath string, err error) {
	return storeTLSFilesToDir(siteTLSDir(siteKey), certPart, keyPart)
}

// HasSiteTLS 该站点映射是否已绑定且磁盘上存在独立证书/私钥.
func HasSiteTLS(siteKey string) bool {
	row := domainsitedao.GetMappingBySiteKey(siteKey)
	if row == nil {
		return false
	}
	cert := strings.TrimSpace(row.CertFile)
	key := strings.TrimSpace(row.KeyFile)
	return cert != "" && key != "" && gfile.Exists(cert) && gfile.Exists(key)
}

func storeTLSFilesToDir(dir string, certPart, keyPart *ghttp.UploadFile) (certPath, keyPath string, err error) {
	if certPart == nil || keyPart == nil {
		return "", "", fmt.Errorf("请同时上传证书文件和私钥文件")
	}
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return "", "", fmt.Errorf("证书目录无效")
	}
	if err = gfile.Mkdir(dir); err != nil {
		return "", "", fmt.Errorf("创建证书目录失败: %w", err)
	}
	certPath = filepath.Join(dir, "fullchain.pem")
	keyPath = filepath.Join(dir, "privkey.pem")
	if err = writeUploadFile(certPart, certPath); err != nil {
		return "", "", fmt.Errorf("保存证书失败: %w", err)
	}
	if err = writeUploadFile(keyPart, keyPath); err != nil {
		_ = os.Remove(certPath)
		return "", "", fmt.Errorf("保存私钥失败: %w", err)
	}
	return certPath, keyPath, nil
}

func writeUploadFile(part *ghttp.UploadFile, dest string) error {
	src, err := part.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

func removeExtraDomainTLSDir(siteKey string) {
	dir := extraDomainTLSDir(siteKey)
	if dir == "" || !strings.HasPrefix(filepath.Clean(dir), filepath.Clean(extraDomainTLSRoot)) {
		return
	}
	_ = os.RemoveAll(dir)
}
