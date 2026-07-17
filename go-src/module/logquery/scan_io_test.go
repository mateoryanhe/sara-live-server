package logquery

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestReadLogLineTruncatesWithoutStopping(t *testing.T) {
	var builder strings.Builder
	builder.WriteString("head,")
	builder.WriteString(strings.Repeat("x", maxLogLineBytes+1024))
	builder.WriteByte('\n')
	builder.WriteString("tail-line\n")

	reader := bufio.NewReader(strings.NewReader(builder.String()))
	first, err := readLogLine(reader)
	if err != nil {
		t.Fatalf("read first line failed: %v", err)
	}
	if !strings.Contains(first, "...[line truncated]") {
		t.Fatalf("expected truncated marker, got len=%d", len(first))
	}

	second, err := readLogLine(reader)
	if err != nil {
		t.Fatalf("read second line failed: %v", err)
	}
	if second != "tail-line" {
		t.Fatalf("unexpected second line: %q", second)
	}

	_, err = readLogLine(reader)
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}
