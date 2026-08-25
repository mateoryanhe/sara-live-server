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

const excelSafeIntegerMax = uint64(9007199254740991) // 2^53-1, Excel/JS 安全整数上限

func formatCSVUint(value uint64) string {
	s := strconv.FormatUint(value, 10)
	if value > excelSafeIntegerMax {
		// Excel 打开 CSV 会将超长数字转为科学计数法并丢精度; \t 前缀强制按文本显示
		return "\t" + s
	}
	return s
}
