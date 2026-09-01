package cmsexport

import (
	"errors"

	"xr-game-server/module/fileexport"
)

type exportRecord = fileexport.Record

type exportResult struct {
	ExportID string
	FileName string
	FileUrl  string
	Total    int
}

var (
	errExportUnavailable = fileexport.ErrUnavailable
	errExportEmpty       = errors.New("没有可导出的数据")
)

func ensureExportReady() error {
	if !fileexport.Ready() {
		return errExportUnavailable
	}
	return nil
}

func createExportFile(suffix string) (*exportRecord, string, error) {
	if err := ensureExportReady(); err != nil {
		return nil, "", err
	}
	record, err := fileexport.Create(suffix)
	if err != nil {
		return nil, "", err
	}
	return record, record.AbsPath, nil
}

func finalizeExportRecord(record *exportRecord, total int) *exportResult {
	if record == nil {
		return nil
	}
	return &exportResult{
		ExportID: record.ExportID,
		FileName: record.FileName,
		FileUrl:  record.FileURL,
		Total:    total,
	}
}

func deleteExport(exportID string) error {
	return fileexport.Delete(exportID)
}
