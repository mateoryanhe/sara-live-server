package game

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/module/upload"
)

const (
	vendorGameCoverStoredPrefix = "game-covers"
	vendorGameCoverMaxBytes     = int64(10 * 1024 * 1024)
	vendorGameCoverTimeout      = 20 * time.Second
	vendorGameCoverMaxAttempts  = 3
	vendorGameCoverRetryDelay   = 500 * time.Millisecond
)

type vendorGameCoverStoreFunc func(context.Context, *VendorGame, string) (string, error)

var errVendorGameCoverUnavailable = errors.New("vendor game cover unavailable")

type vendorGameCoverHTTPStatusError struct {
	StatusCode int
}

func (e *vendorGameCoverHTTPStatusError) Error() string {
	return fmt.Sprintf("download cover HTTP status=%d", e.StatusCode)
}

// mirrorVendorGameCovers 下载第三方封面，并按当前云桶开关写入 R2 或本地 storagePath.
func mirrorVendorGameCovers(ctx context.Context, games []*VendorGame) error {
	if ctx == nil {
		ctx = context.Background()
	}
	iconBaseURL := cfgdao.GetGameIconBaseUrlFromMemory()
	return mirrorVendorGameCoversWithStore(ctx, games, func(ctx context.Context, game *VendorGame, rawCover string) (string, error) {
		sourceURL, err := resolveVendorGameCoverSourceURL(rawCover, iconBaseURL)
		if err != nil {
			return "", err
		}
		// 每张封面使用独立 HTTP 连接，完成后立即释放；每次尝试单独限制超时时间。
		client := newVendorGameCoverHTTPClient()
		defer client.CloseIdleConnections()
		var lastErr error
		for attempt := 1; attempt <= vendorGameCoverMaxAttempts; attempt++ {
			downloadCtx, cancel := context.WithTimeout(ctx, vendorGameCoverTimeout)
			data, ext, downloadErr := downloadVendorGameCover(downloadCtx, client, sourceURL)
			cancel()
			if downloadErr == nil {
				storedName := vendorGameCoverStoredName(game.GameCode, game.Platform, ext)
				if saveErr := upload.SaveUploadedFileBytes(storedName, data); saveErr == nil {
					return storedName, nil
				} else {
					lastErr = fmt.Errorf("save cover: %w", saveErr)
				}
			} else if isVendorGameCoverUnavailable(downloadErr) {
				return "", fmt.Errorf("%w: sourceURL=%s: %v", errVendorGameCoverUnavailable, sourceURL, downloadErr)
			} else {
				lastErr = downloadErr
			}
			if attempt >= vendorGameCoverMaxAttempts {
				break
			}
			vendorDetailLog().Warningf(ctx, "sync vendor game cover retry gameCode=%s platform=%s attempt=%d/%d err=%v",
				game.GameCode, game.Platform, attempt, vendorGameCoverMaxAttempts, lastErr)
			if err := sleepWithContext(ctx, vendorGameCoverRetryDelay*time.Duration(attempt)); err != nil {
				return "", err
			}
		}
		return "", fmt.Errorf("store cover after %d attempts: %w", vendorGameCoverMaxAttempts, lastErr)
	})
}

func newVendorGameCoverHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	transport.MaxConnsPerHost = 1
	return &http.Client{
		Transport: transport,
		Timeout:   vendorGameCoverTimeout,
	}
}

// mirrorVendorGameCoversWithStore 严格逐张处理封面，仅在全部成功后回填 games.
func mirrorVendorGameCoversWithStore(ctx context.Context, games []*VendorGame, store vendorGameCoverStoreFunc) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if store == nil {
		return errors.New("vendor game cover store is nil")
	}

	paths := make([]string, len(games))
	for index, game := range games {
		if game == nil {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		rawCover := strings.TrimSpace(game.Cover)
		if rawCover == "" {
			continue
		}
		storedName, err := store(ctx, game, rawCover)
		if err != nil {
			if errors.Is(err, errVendorGameCoverUnavailable) {
				vendorDetailLog().Warningf(ctx, "skip unavailable vendor game cover gameCode=%s platform=%s err=%v",
					game.GameCode, game.Platform, err)
				continue
			}
			return fmt.Errorf("gameCode=%s platform=%s: %w", game.GameCode, game.Platform, err)
		}
		storedName = strings.TrimSpace(storedName)
		if storedName == "" {
			return fmt.Errorf("gameCode=%s platform=%s: stored cover path is empty", game.GameCode, game.Platform)
		}
		paths[index] = storedName
	}
	for index, game := range games {
		if game == nil {
			continue
		}
		game.Cover = paths[index]
	}
	return nil
}

func resolveVendorGameCoverSourceURL(rawCover, iconBaseURL string) (string, error) {
	rawCover = strings.TrimSpace(rawCover)
	if rawCover == "" {
		return "", errors.New("empty vendor cover")
	}
	if strings.HasPrefix(rawCover, "//") {
		rawCover = "https:" + rawCover
	}
	if parsed, err := url.Parse(rawCover); err == nil && parsed.IsAbs() {
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return "", fmt.Errorf("unsupported cover URL scheme %q", parsed.Scheme)
		}
		if parsed.Host == "" {
			return "", errors.New("vendor cover URL has no host")
		}
		return parsed.String(), nil
	}

	iconBaseURL = strings.TrimRight(strings.TrimSpace(iconBaseURL), "/")
	if iconBaseURL == "" {
		return "", errors.New("game icon base URL is empty")
	}
	base, err := url.Parse(iconBaseURL + "/")
	if err != nil || (base.Scheme != "http" && base.Scheme != "https") || base.Host == "" {
		return "", errors.New("game icon base URL is invalid")
	}
	relativePath := normalizeVendorGameCoverPath(rawCover)
	if relativePath == "" {
		return "", errors.New("vendor cover path is empty")
	}
	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return "", errors.New("vendor cover path is invalid")
	}
	return base.ResolveReference(relativeURL).String(), nil
}

func downloadVendorGameCover(ctx context.Context, client *http.Client, sourceURL string) ([]byte, string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		client = &http.Client{Timeout: vendorGameCoverTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("download cover: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, "", &vendorGameCoverHTTPStatusError{StatusCode: resp.StatusCode}
	}
	if resp.ContentLength > vendorGameCoverMaxBytes {
		return nil, "", fmt.Errorf("cover exceeds %d bytes", vendorGameCoverMaxBytes)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, vendorGameCoverMaxBytes+1))
	if err != nil {
		return nil, "", fmt.Errorf("read cover: %w", err)
	}
	if len(data) == 0 {
		return nil, "", errors.New("cover is empty")
	}
	if int64(len(data)) > vendorGameCoverMaxBytes {
		return nil, "", fmt.Errorf("cover exceeds %d bytes", vendorGameCoverMaxBytes)
	}
	contentType := http.DetectContentType(data)
	contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	extension, ok := vendorGameCoverExtension(contentType)
	if !ok {
		return nil, "", fmt.Errorf("unsupported cover content type %q", contentType)
	}
	return data, extension, nil
}

func isVendorGameCoverUnavailable(err error) bool {
	var statusErr *vendorGameCoverHTTPStatusError
	if !errors.As(err, &statusErr) {
		return false
	}
	return statusErr.StatusCode == http.StatusNotFound || statusErr.StatusCode == http.StatusGone
}

func vendorGameCoverExtension(contentType string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/gif":
		return ".gif", true
	case "image/webp":
		return ".webp", true
	case "image/bmp", "image/x-ms-bmp":
		return ".bmp", true
	default:
		return "", false
	}
}

func vendorGameCoverStoredName(gameCode, platform, extension string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(platform) + "\x00" + strings.TrimSpace(gameCode)))
	extension = strings.ToLower(strings.TrimSpace(extension))
	if extension == "" || !strings.HasPrefix(extension, ".") {
		extension = ".jpg"
	}
	return vendorGameCoverStoredPrefix + "/" + hex.EncodeToString(sum[:]) + extension
}
