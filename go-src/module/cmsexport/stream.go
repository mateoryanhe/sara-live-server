package cmsexport

import (
	"context"
	"fmt"
)

const defaultExportPageSize = 500

type pageRowsFunc func(pageIndex, pageSize int) (total int, rows [][]string)

func streamCSVExport(
	ctx context.Context,
	headers []string,
	pageSize int,
	fetch pageRowsFunc,
	onProgress func(exportedRows, totalRows int),
) (*exportResult, error) {
	if len(headers) == 0 {
		return nil, errExportEmpty
	}
	if pageSize <= 0 {
		pageSize = defaultExportPageSize
	}

	record, filePath, err := createExportFile(".csv")
	if err != nil {
		return nil, err
	}

	writer, err := newCSVWriter(filePath, headers)
	if err != nil {
		removeFile(filePath)
		exportRecords.Delete(record.exportID)
		return nil, err
	}

	pageIndex := 1
	totalRows := 0
	exportedRows := 0

	abort := func() (*exportResult, error) {
		_ = writer.close()
		removeFile(filePath)
		exportRecords.Delete(record.exportID)
		return nil, errExportEmpty
	}

	for {
		select {
		case <-ctx.Done():
			_ = writer.close()
			removeFile(filePath)
			exportRecords.Delete(record.exportID)
			return nil, ctx.Err()
		default:
		}

		total, rows := fetch(pageIndex, pageSize)
		if pageIndex == 1 {
			totalRows = total
			if totalRows == 0 {
				return abort()
			}
			if onProgress != nil {
				onProgress(0, totalRows)
			}
		}
		if len(rows) == 0 {
			if pageIndex == 1 && totalRows == 0 {
				return abort()
			}
			if exportedRows < totalRows {
				_ = writer.close()
				removeFile(filePath)
				exportRecords.Delete(record.exportID)
				return nil, fmt.Errorf("export incomplete: exported %d of %d rows", exportedRows, totalRows)
			}
			break
		}
		for _, row := range rows {
			if err := writer.writeRow(row); err != nil {
				_ = writer.close()
				removeFile(filePath)
				exportRecords.Delete(record.exportID)
				return nil, err
			}
			exportedRows++
		}
		if onProgress != nil {
			onProgress(exportedRows, totalRows)
		}
		if exportedRows >= totalRows {
			break
		}
		pageIndex++
	}

	if err := writer.close(); err != nil {
		removeFile(filePath)
		exportRecords.Delete(record.exportID)
		return nil, err
	}
	return finalizeExportRecord(record, exportedRows), nil
}
