package upload

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type readCountingReader struct {
	reads int
}

func (r *readCountingReader) Read([]byte) (int, error) {
	r.reads++
	return 0, errors.New("reader must not be consumed")
}

func TestStoreUploadedContentRejectsKnownOversizeBeforeRead(t *testing.T) {
	src := &readCountingReader{}
	_, _, err := storeUploadedContentWithSize(
		context.Background(),
		src,
		11,
		StoreCatImages,
		".jpg",
		10,
		errImageFileTooLarge,
	)
	if !errors.Is(err, errImageFileTooLarge) {
		t.Fatalf("error = %v, want %v", err, errImageFileTooLarge)
	}
	if src.reads != 0 {
		t.Fatalf("reader consumed %d times before size validation", src.reads)
	}
}

func TestCopyUploadContentToFileRemovesOversizeFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "video.mp4")
	dst, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = copyUploadContentToFile(bytes.NewReader([]byte("123456")), dst, path, 5, errVideoFileTooLarge)
	if !errors.Is(err, errVideoFileTooLarge) {
		t.Fatalf("error = %v, want %v", err, errVideoFileTooLarge)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("oversize file was not removed: %v", statErr)
	}
}

func TestCopyUploadContentToFileRejectsEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "image.jpg")
	dst, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = copyUploadContentToFile(bytes.NewReader(nil), dst, path, 1024, errImageFileTooLarge)
	if !errors.Is(err, errUploadFileEmpty) {
		t.Fatalf("error = %v, want %v", err, errUploadFileEmpty)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("empty file was not removed: %v", statErr)
	}
}
