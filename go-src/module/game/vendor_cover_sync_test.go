package game

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

var testPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func testCoverHTTPClient(statusCode int, contentLength int64, body []byte) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    statusCode,
			ContentLength: contentLength,
			Body:          io.NopCloser(bytes.NewReader(body)),
			Header:        make(http.Header),
		}, nil
	})}
}

func TestResolveVendorGameCoverSourceURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		raw      string
		base     string
		expected string
		wantErr  bool
	}{
		{
			name:     "absolute URL",
			raw:      "https://vendor.example/cover/a.png",
			expected: "https://vendor.example/cover/a.png",
		},
		{
			name:     "protocol relative URL",
			raw:      "//vendor.example/cover/a.png",
			expected: "https://vendor.example/cover/a.png",
		},
		{
			name:     "vendor relative path",
			raw:      "/uploads/file/game/slot/a.png?version=2",
			base:     "https://cdn.example/game",
			expected: "https://cdn.example/game/slot/a.png?version=2",
		},
		{
			name:    "relative path without base",
			raw:     "slot/a.png",
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			raw:     "file:///tmp/a.png",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := resolveVendorGameCoverSourceURL(tt.raw, tt.base)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got URL %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolve URL: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDownloadVendorGameCover(t *testing.T) {
	t.Parallel()

	client := testCoverHTTPClient(http.StatusOK, int64(len(testPNG)), testPNG)
	data, ext, err := downloadVendorGameCover(context.Background(), client, "https://vendor.example/cover.png")
	if err != nil {
		t.Fatalf("download cover: %v", err)
	}
	if string(data) != string(testPNG) {
		t.Fatalf("downloaded data mismatch")
	}
	if ext != ".png" {
		t.Fatalf("got extension %q, want .png", ext)
	}
}

func TestNewVendorGameCoverHTTPClientIsIndependentPerImage(t *testing.T) {
	t.Parallel()

	first := newVendorGameCoverHTTPClient()
	second := newVendorGameCoverHTTPClient()
	defer first.CloseIdleConnections()
	defer second.CloseIdleConnections()

	if first == second || first.Transport == second.Transport {
		t.Fatal("each image must use an independent HTTP client and transport")
	}
	if first.Timeout != vendorGameCoverTimeout || second.Timeout != vendorGameCoverTimeout {
		t.Fatalf("unexpected client timeout: %v, %v", first.Timeout, second.Timeout)
	}
	transport, ok := first.Transport.(*http.Transport)
	if !ok || !transport.DisableKeepAlives {
		t.Fatal("image transport must close its connection after the request")
	}
}

func TestDownloadVendorGameCoverRejectsInvalidContent(t *testing.T) {
	t.Parallel()

	body := []byte("not an image")
	client := testCoverHTTPClient(http.StatusOK, int64(len(body)), body)
	if _, _, err := downloadVendorGameCover(context.Background(), client, "https://vendor.example/cover.txt"); err == nil {
		t.Fatal("expected invalid image error")
	}
}

func TestDownloadVendorGameCoverRejectsOversizedContentLength(t *testing.T) {
	t.Parallel()

	client := testCoverHTTPClient(http.StatusOK, vendorGameCoverMaxBytes+1, nil)
	if _, _, err := downloadVendorGameCover(context.Background(), client, "https://vendor.example/large.png"); err == nil {
		t.Fatal("expected oversized image error")
	}
}

func TestDownloadVendorGameCoverClassifiesUnavailableStatus(t *testing.T) {
	t.Parallel()

	for _, statusCode := range []int{http.StatusNotFound, http.StatusGone} {
		client := testCoverHTTPClient(statusCode, 0, nil)
		_, _, err := downloadVendorGameCover(context.Background(), client, "https://vendor.example/missing.webp")
		if !isVendorGameCoverUnavailable(err) {
			t.Fatalf("status %d should be unavailable, got %v", statusCode, err)
		}
	}

	client := testCoverHTTPClient(http.StatusInternalServerError, 0, nil)
	_, _, err := downloadVendorGameCover(context.Background(), client, "https://vendor.example/error.webp")
	if isVendorGameCoverUnavailable(err) {
		t.Fatalf("status 500 must remain retryable, got %v", err)
	}
}

func TestMirrorVendorGameCoversCommitsOnlyAfterAllStored(t *testing.T) {
	t.Parallel()

	games := []*VendorGame{
		{GameCode: "100", Platform: "PG", Cover: "source-a"},
		{GameCode: "200", Platform: "PP", Cover: ""},
		{GameCode: "300", Platform: "PG", Cover: "source-c"},
	}
	seen := make(map[string]string)
	order := make([]string, 0, 2)
	err := mirrorVendorGameCoversWithStore(context.Background(), games, func(_ context.Context, game *VendorGame, raw string) (string, error) {
		seen[game.GameCode] = raw
		order = append(order, game.GameCode)
		return "game-covers/" + game.GameCode + ".png", nil
	})
	if err != nil {
		t.Fatalf("mirror covers: %v", err)
	}
	if games[0].Cover != "game-covers/100.png" || games[1].Cover != "" || games[2].Cover != "game-covers/300.png" {
		t.Fatalf("unexpected mirrored covers: %#v", games)
	}
	if len(seen) != 2 || seen["100"] != "source-a" || seen["300"] != "source-c" {
		t.Fatalf("unexpected store calls: %#v", seen)
	}
	if strings.Join(order, ",") != "100,300" {
		t.Fatalf("covers were not stored sequentially: %#v", order)
	}
}

func TestMirrorVendorGameCoversDoesNotMutateOnFailure(t *testing.T) {
	t.Parallel()

	games := []*VendorGame{
		{GameCode: "100", Platform: "PG", Cover: "source-a"},
		{GameCode: "200", Platform: "PP", Cover: "source-b"},
	}
	err := mirrorVendorGameCoversWithStore(context.Background(), games, func(_ context.Context, game *VendorGame, _ string) (string, error) {
		if game.GameCode == "200" {
			return "", errors.New("download failed")
		}
		return "game-covers/100.png", nil
	})
	if err == nil || !strings.Contains(err.Error(), "gameCode=200") {
		t.Fatalf("expected game-specific error, got %v", err)
	}
	if games[0].Cover != "source-a" || games[1].Cover != "source-b" {
		t.Fatalf("games mutated after failure: %#v", games)
	}
}

func TestMirrorVendorGameCoversSkipsUnavailableCover(t *testing.T) {
	t.Parallel()

	games := []*VendorGame{
		{GameCode: "100", Platform: "PG", Cover: "source-a"},
		{GameCode: "85", Platform: "JILI", Cover: "missing-source"},
		{GameCode: "300", Platform: "PG", Cover: "source-c"},
	}
	err := mirrorVendorGameCoversWithStore(context.Background(), games, func(_ context.Context, game *VendorGame, _ string) (string, error) {
		if game.GameCode == "85" {
			return "", fmt.Errorf("%w: test 404", errVendorGameCoverUnavailable)
		}
		return "game-covers/" + game.GameCode + ".webp", nil
	})
	if err != nil {
		t.Fatalf("mirror covers: %v", err)
	}
	if games[0].Cover != "game-covers/100.webp" || games[1].Cover != "" || games[2].Cover != "game-covers/300.webp" {
		t.Fatalf("unexpected mirrored covers: %#v", games)
	}
}

func TestVendorGameCoverStoredNameIsStableAndScoped(t *testing.T) {
	t.Parallel()

	first := vendorGameCoverStoredName("100", "PG", ".PNG")
	second := vendorGameCoverStoredName("100", "PG", ".png")
	otherPlatform := vendorGameCoverStoredName("100", "PP", ".png")
	if first != second {
		t.Fatalf("stored name is not stable: %q != %q", first, second)
	}
	if first == otherPlatform {
		t.Fatalf("stored name must include platform")
	}
	if !strings.HasPrefix(first, "game-covers/") || !strings.HasSuffix(first, ".png") {
		t.Fatalf("unexpected stored name %q", first)
	}
}

func TestDedupeVendorGames(t *testing.T) {
	t.Parallel()

	first := &VendorGame{GameCode: "100", Platform: "PG", Cover: "first"}
	duplicate := &VendorGame{GameCode: "100", Platform: "PG", Cover: "second"}
	otherPlatform := &VendorGame{GameCode: "100", Platform: "PP", Cover: "third"}
	got := dedupeVendorGames([]*VendorGame{nil, first, duplicate, otherPlatform, &VendorGame{GameCode: " "}})
	if len(got) != 2 || got[0] != first || got[1] != otherPlatform {
		t.Fatalf("unexpected deduplicated games: %#v", got)
	}
}
