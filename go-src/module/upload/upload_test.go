package upload

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"xr-game-server/core/cfg"
)

type readCountingReader struct {
	reads int
}

func (r *readCountingReader) Read([]byte) (int, error) {
	r.reads++
	return 0, errors.New("reader must not be consumed")
}

func TestStoreUploadedContentRejectsKnownOversizeBeforeRead(t *testing.T) {
	src := &readCountingReader{}
	_, _, err := storeUploadedContentWithSize(
		context.Background(),
		src,
		11,
		StoreCatImages,
		".jpg",
		10,
		errImageFileTooLarge,
	)
	if !errors.Is(err, errImageFileTooLarge) {
		t.Fatalf("error = %v, want %v", err, errImageFileTooLarge)
	}
	if src.reads != 0 {
		t.Fatalf("reader consumed %d times before size validation", src.reads)
	}
}

func TestCopyUploadContentToFileRemovesOversizeFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "video.mp4")
	dst, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = copyUploadContentToFile(bytes.NewReader([]byte("123456")), dst, path, 5, errVideoFileTooLarge)
	if !errors.Is(err, errVideoFileTooLarge) {
		t.Fatalf("error = %v, want %v", err, errVideoFileTooLarge)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("oversize file was not removed: %v", statErr)
	}
}

func TestCopyUploadContentToFileRejectsEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "image.jpg")
	dst, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = copyUploadContentToFile(bytes.NewReader(nil), dst, path, 1024, errImageFileTooLarge)
	if !errors.Is(err, errUploadFileEmpty) {
		t.Fatalf("error = %v, want %v", err, errUploadFileEmpty)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("empty file was not removed: %v", statErr)
	}
}

func TestShouldRegisterResourceDomain(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		s3Enabled bool
		want      bool
	}{
		{name: "domain local storage", host: "cdn.example.com", want: true},
		{name: "domain cloud storage", host: "cdn.example.com", s3Enabled: true, want: false},
		{name: "ipv4", host: "203.0.113.10", want: false},
		{name: "ipv6", host: "2001:db8::10", want: false},
		{name: "bracketed ipv6", host: "[2001:db8::10]", want: false},
		{name: "ipv6 zone", host: "fe80::1%eth0", want: false},
		{name: "localhost", host: "localhost", want: false},
		{name: "empty", host: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRegisterResourceDomain(tt.host, tt.s3Enabled); got != tt.want {
				t.Fatalf("shouldRegisterResourceDomain(%q, %t) = %t, want %t", tt.host, tt.s3Enabled, got, tt.want)
			}
		})
	}
}

func TestResourceDomainHostRecognizesIPURLs(t *testing.T) {
	tests := map[string]string{
		"https://cdn.example.com/images":     "cdn.example.com",
		"http://203.0.113.10:9443/images":    "203.0.113.10",
		"https://[2001:db8::10]:9443/images": "2001:db8::10",
		"https://2001:db8::10":               "2001:db8::10",
	}
	for input, want := range tests {
		if got := resourceDomainHost(input); got != want {
			t.Fatalf("resourceDomainHost(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestIPv4ResourceDomainKeepsLocalAvatarPathAccessible(t *testing.T) {
	originalCfg := getResourceCfgCache()
	originalPrefix := cfg.GetImageStaticPrefix()
	root := t.TempDir()
	resourceCfgCache.Store(&resourceCfgSnapshot{
		ResourceDomain:           "http://192.167.2.99/iimage",
		ResourceDomainConfigured: true,
		StoragePath:              root,
		CmsExportTtlMinutes:      30,
		S3Enabled:                false,
	})
	cfg.ClearRuntimeStaticMappings()
	t.Cleanup(func() {
		cfg.ClearRuntimeStaticMappings()
		cfg.SetImageStaticPrefix(originalPrefix)
		resourceCfgCache.Store(originalCfg)
	})

	registerStaticMappings()

	if got := GetUrlByName("2.png"); got != "http://192.167.2.99/iimage/2.png" {
		t.Fatalf("avatar URL = %q", got)
	}
	mappedRoot, rel, ok := cfg.MatchStaticPath("/iimage/2.png")
	if !ok {
		t.Fatal("IPv4 avatar URL path was not registered")
	}
	if filepath.Clean(mappedRoot) != filepath.Clean(root) || rel != "2.png" {
		t.Fatalf("mapping = root %q rel %q, want root %q rel %q", mappedRoot, rel, root, "2.png")
	}
}
