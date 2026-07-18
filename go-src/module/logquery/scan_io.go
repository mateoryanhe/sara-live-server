package logquery

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const maxLogLineBytes = 4 * 1024 * 1024

func openLogFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return file, nil
}

func readLogLine(reader *bufio.Reader) (string, error) {
	return readLogLineLimited(reader, maxLogLineBytes)
}

func readLogLineLimited(reader *bufio.Reader, maxBytes int) (string, error) {
	if maxBytes <= 0 {
		maxBytes = maxLogLineBytes
	}
	var builder strings.Builder
	builder.Grow(4096)
	truncated := false

	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			if builder.Len() == 0 {
				return "", io.EOF
			}
			return builder.String(), nil
		}
		if err != nil {
			return builder.String(), err
		}
		if b == '\n' {
			break
		}
		if builder.Len() < maxBytes {
			builder.WriteByte(b)
			continue
		}
		truncated = true
	}

	if truncated {
		for {
			b, err := reader.ReadByte()
			if err != nil {
				break
			}
			if b == '\n' {
				break
			}
		}
		builder.WriteString("...[line truncated]")
	}

	return builder.String(), nil
}

func scanLogFile(filePath string, fn func(line string) bool) error {
	return scanLogFileWithMaxLineBytes(filePath, maxLogLineBytes, fn)
}

func scanLogFileWithMaxLineBytes(filePath string, maxLineBytes int, fn func(line string) bool) error {
	file, err := openLogFile(filePath)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 256*1024)
	for {
		line, err := readLogLineLimited(reader, maxLineBytes)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		if !fn(line) {
			return nil
		}
	}
}
