package cmsexport

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type csvWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func newCSVWriter(path string, headers []string) (*csvWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	writer := &csvWriter{
		file:   file,
		writer: bufio.NewWriterSize(file, 64*1024),
	}
	if _, err := writer.writer.WriteString("\uFEFF"); err != nil {
		_ = file.Close()
		return nil, err
	}
	if err := writer.writeRow(headers); err != nil {
		_ = file.Close()
		return nil, err
	}
	return writer, nil
}

func (w *csvWriter) writeRow(cells []string) error {
	parts := make([]string, len(cells))
	for i, cell := range cells {
		parts[i] = escapeCSVCell(cell)
	}
	_, err := w.writer.WriteString(strings.Join(parts, ",") + "\r\n")
	return err
}

func (w *csvWriter) close() error {
	if w == nil {
		return nil
	}
	if err := w.writer.Flush(); err != nil {
		_ = w.file.Close()
		return err
	}
	return w.file.Close()
}

func escapeCSVCell(value string) string {
	if value == "" {
		return ""
	}
	if strings.ContainsAny(value, "\",\r\n") {
		return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	return value
}

func formatCSVFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatCSVUint(value uint64) string {
	return strconv.FormatUint(value, 10)
}
