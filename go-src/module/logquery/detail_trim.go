package logquery

import "strings"

const (
	maxDetailLogRespContentBytes = 1024
	detailLogPayloadTruncated    = "...[truncated]"
	detailLogResponseWriteMarker = "应答序列化,写入框架缓冲区"
)

func trimDetailLogEntryForQuery(entry *DetailLogEntry) *DetailLogEntry {
	if entry == nil || !isResponseWriteDetailLog(entry) {
		return entry
	}
	entry.Message = truncateLogFieldValue(entry.Message, "respContent=", maxDetailLogRespContentBytes)
	entry.Raw = truncateLogFieldValue(entry.Raw, "respContent=", maxDetailLogRespContentBytes)
	return entry
}

func isResponseWriteDetailLog(entry *DetailLogEntry) bool {
	return strings.Contains(entry.Message, detailLogResponseWriteMarker) ||
		strings.Contains(entry.Raw, detailLogResponseWriteMarker)
}

func truncateLogFieldValue(text, prefix string, maxBytes int) string {
	start := strings.Index(text, prefix)
	if start < 0 {
		return text
	}
	valueStart := start + len(prefix)
	if valueStart >= len(text) {
		return text
	}
	valueLen := len(text) - valueStart
	if valueLen <= maxBytes {
		return text
	}
	return text[:valueStart] + text[valueStart:valueStart+maxBytes] + detailLogPayloadTruncated
}
