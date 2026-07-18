package logquery

import (
	"os"
	"testing"

	"xr-game-server/dto/logquerydto"
)

func TestSanitizeGrepPatternRejectsUnsafeInput(t *testing.T) {
	if _, ok := sanitizeGrepPattern("abc\n"); ok {
		t.Fatal("expected newline pattern to be rejected")
	}
	if pattern, ok := sanitizeGrepPattern(" trace "); !ok || pattern != "trace" {
		t.Fatalf("unexpected sanitize result: %q %v", pattern, ok)
	}
}

func TestBuildDetailGrepPatternsPrefersTraceId(t *testing.T) {
	primary, secondary := buildDetailGrepPatterns(&logquerydto.CMSQueryDetailLogsReq{
		TraceId: "abc123",
		Url:     "/api/test",
	})
	if primary != "{abc123}" {
		t.Fatalf("unexpected primary pattern: %s", primary)
	}
	if len(secondary) != 1 || secondary[0] != "/api/test" {
		t.Fatalf("unexpected secondary patterns: %#v", secondary)
	}
}

func TestSplitOutputLinesRespectsLimit(t *testing.T) {
	out := []byte("a\nb\nc\n")
	lines := splitOutputLines(out, 2)
	if len(lines) != 2 || lines[0] != "a" || lines[1] != "b" {
		t.Fatalf("unexpected lines: %#v", lines)
	}
}

func TestFilterAllowedExistingLogFiles(t *testing.T) {
	dir := t.TempDir()
	filePath := dir + "/2026-07-17.log"
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file failed: %v", err)
	}
	files := filterAllowedExistingLogFiles([]string{filePath, dir + "/missing.log", "/etc/passwd"}, dir)
	if len(files) != 1 || files[0] != filePath {
		t.Fatalf("unexpected files: %#v", files)
	}
}
