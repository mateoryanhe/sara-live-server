package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"
)

// HostTLSRegistration 纯 Host→证书(无静态目录);Domain 支持精确 Host 或单级泛域名 *.example.com.
type HostTLSRegistration struct {
	Domain   string
	CertFile string
	KeyFile  string
}

type hostTLSRegistrySnapshot struct {
	certificates map[string]*tls.Certificate // domain or *.base -> cert
}

var (
	hostTLSRegistry   atomic.Pointer[hostTLSRegistrySnapshot]
	hostTLSRegistryMu sync.Mutex
)

// RefreshHostTLSRegistry 原子替换纯 Host 证书表(可传空清空).
func RefreshHostTLSRegistry(ctx context.Context, registrations []HostTLSRegistration) error {
	hostTLSRegistryMu.Lock()
	defer hostTLSRegistryMu.Unlock()

	snapshot := &hostTLSRegistrySnapshot{
		certificates: make(map[string]*tls.Certificate, len(registrations)),
	}
	for index, item := range registrations {
		domain, err := normalizeHostTLSName(item.Domain)
		if err != nil {
			return fmt.Errorf("第 %d 项 Host 无效: %w", index+1, err)
		}
		certFile := strings.TrimSpace(item.CertFile)
		keyFile := strings.TrimSpace(item.KeyFile)
		if certFile == "" || keyFile == "" {
			return fmt.Errorf("Host %s 必须同时配置证书和私钥", domain)
		}
		loaded, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("加载 Host %s 证书失败: %w", domain, err)
		}
		if _, exists := snapshot.certificates[domain]; exists {
			return fmt.Errorf("Host 重复注册: %s", domain)
		}
		cert := loaded
		snapshot.certificates[domain] = &cert
	}
	hostTLSRegistry.Store(snapshot)
	g.Log().Infof(ctx, "纯Host证书注册表刷新成功,hosts=%d", len(snapshot.certificates))
	return nil
}

func lookupHostTLSCertificate(serverName string) *tls.Certificate {
	name := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(serverName), "."))
	if name == "" {
		return nil
	}
	snapshot := hostTLSRegistry.Load()
	if snapshot == nil {
		return nil
	}
	if cert := snapshot.certificates[name]; cert != nil {
		return cert
	}
	// 单级泛域名: api.sofielive.net → *.sofielive.net(不匹配 apex 与多级子域)
	if i := strings.IndexByte(name, '.'); i > 0 && i+1 < len(name) {
		if cert := snapshot.certificates["*."+name[i+1:]]; cert != nil {
			return cert
		}
	}
	return nil
}

// NormalizeHostTLSName 精确 Host,或单级泛域名 *.example.com(仅最左标签可为 *).
func NormalizeHostTLSName(domain string) (string, error) {
	return normalizeHostTLSName(domain)
}

func normalizeHostTLSName(domain string) (string, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	domain = strings.TrimSuffix(domain, ".")
	if domain == "" {
		return "", fmt.Errorf("域名不能为空")
	}
	if strings.Contains(domain, "://") || strings.ContainsAny(domain, `/\\?#@`) {
		return "", fmt.Errorf("只允许填写 Host,不能包含协议、路径或查询参数: %s", domain)
	}
	if strings.ContainsAny(domain, " \t\r\n") {
		return "", fmt.Errorf("域名不能包含空白字符: %s", domain)
	}
	if strings.Contains(domain, ":") {
		return "", fmt.Errorf("域名不能包含端口: %s", domain)
	}
	if strings.Contains(domain, ",") {
		return "", fmt.Errorf("只允许配置一个 Host,不能使用逗号")
	}
	if net.ParseIP(domain) != nil || domain == "localhost" {
		return domain, nil
	}
	if len(domain) > 253 {
		return "", fmt.Errorf("域名长度超过253字符: %s", domain)
	}

	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return "", fmt.Errorf("Host 至少需要两级: %s", domain)
	}
	for i, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return "", fmt.Errorf("域名标签长度无效: %s", domain)
		}
		if label == "*" {
			if i != 0 {
				return "", fmt.Errorf("泛域名 * 只能出现在最左侧,如 *.example.com: %s", domain)
			}
			if len(labels) < 3 {
				return "", fmt.Errorf("泛域名至少形如 *.example.com: %s", domain)
			}
			continue
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return "", fmt.Errorf("域名标签不能以连字符开头或结尾: %s", domain)
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return "", fmt.Errorf("域名包含非法字符: %s", domain)
			}
		}
	}
	return domain, nil
}
