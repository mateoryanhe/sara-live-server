package officialsitedeploy

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/dto/officialsitedeploydto"
)

func TestOfficialSiteFileChunkUploadAndAbort(t *testing.T) {
	fileSize := officialSiteUploadChunkSize + 3
	initRes, err := InitOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.InitOfficialSiteFileUploadReq{
		FileName: "app-release.apk",
		FileSize: fileSize,
	})
	if err != nil {
		t.Fatalf("init upload: %v", err)
	}
	t.Cleanup(func() {
		_, _ = AbortOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.AbortOfficialSiteFileUploadReq{UploadId: initRes.UploadId})
	})
	if initRes.TotalChunks != 2 {
		t.Fatalf("total chunks: got %d, want 2", initRes.TotalChunks)
	}

	chunks := [][]byte{
		bytes.Repeat([]byte{0x5a}, int(officialSiteUploadChunkSize)),
		{0x01, 0x02, 0x03},
	}
	for chunkIndex, chunk := range chunks {
		res, uploadErr := uploadChunkForTest(initRes.UploadId, int64(chunkIndex), chunk)
		if uploadErr != nil {
			t.Fatalf("upload chunk %d: %v", chunkIndex, uploadErr)
		}
		if res.FileSize != int64(len(chunk)) {
			t.Fatalf("chunk %d size: got %d, want %d", chunkIndex, res.FileSize, len(chunk))
		}
	}

	root, meta, err := loadUploadSession(initRes.UploadId, officialSiteTarget)
	if err != nil {
		t.Fatalf("load upload session: %v", err)
	}
	for chunkIndex, want := range chunks {
		got, readErr := os.ReadFile(uploadChunkPath(uploadSessionDir(root, meta.UploadId), int64(chunkIndex)))
		if readErr != nil {
			t.Fatalf("read chunk %d: %v", chunkIndex, readErr)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("chunk %d content mismatch", chunkIndex)
		}
	}

	if _, err := AbortOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.AbortOfficialSiteFileUploadReq{UploadId: initRes.UploadId}); err != nil {
		t.Fatalf("abort upload: %v", err)
	}
	if _, err := os.Stat(uploadSessionDir(root, initRes.UploadId)); !os.IsNotExist(err) {
		t.Fatalf("upload session still exists after abort: %v", err)
	}
}

func TestThirdPayOfficialSiteUploadUsesIndependentSessionAndAcceptsZip(t *testing.T) {
	initRes, err := InitThirdPayOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.InitThirdPayOfficialSiteFileUploadReq{
		FileName: "third-pay.apk",
		FileSize: 3,
	})
	if err != nil {
		t.Fatalf("init third-pay upload: %v", err)
	}
	t.Cleanup(func() {
		_, _ = AbortThirdPayOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.AbortThirdPayOfficialSiteFileUploadReq{UploadId: initRes.UploadId})
	})
	if _, _, err := loadUploadSession(initRes.UploadId, officialSiteTarget); err == nil {
		t.Fatal("expected official site endpoint to reject third-pay upload session")
	}
	if _, _, err := loadUploadSession(initRes.UploadId, thirdPayOfficialSiteTarget); err != nil {
		t.Fatalf("load third-pay upload session: %v", err)
	}
	zipRes, err := InitThirdPayOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.InitThirdPayOfficialSiteFileUploadReq{
		FileName: "third-pay.zip",
		FileSize: 3,
	})
	if err != nil {
		t.Fatalf("init third-pay zip upload: %v", err)
	}
	t.Cleanup(func() {
		_, _ = AbortThirdPayOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.AbortThirdPayOfficialSiteFileUploadReq{UploadId: zipRes.UploadId})
	})
}

func TestOfficialSiteFileChunkRejectsWrongSize(t *testing.T) {
	initRes, err := InitOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.InitOfficialSiteFileUploadReq{
		FileName: "site.zip",
		FileSize: 4,
	})
	if err != nil {
		t.Fatalf("init upload: %v", err)
	}
	t.Cleanup(func() {
		_, _ = AbortOfficialSiteFileUpload(context.Background(), &officialsitedeploydto.AbortOfficialSiteFileUploadReq{UploadId: initRes.UploadId})
	})
	if _, err := uploadChunkForTest(initRes.UploadId, 0, []byte{0x01, 0x02, 0x03}); err == nil {
		t.Fatal("expected wrong-size chunk to be rejected")
	}
}

func uploadChunkForTest(uploadId string, chunkIndex int64, chunk []byte) (*officialsitedeploydto.UploadOfficialSiteFileChunkRes, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "chunk.bin")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(chunk); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	httpReq := httptest.NewRequest("POST", "/officialSiteDeploy/uploadFileChunk", &body)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set(uploadIdHeader, uploadId)
	httpReq.Header.Set(chunkIndexHeader, strconv.FormatInt(chunkIndex, 10))
	return UploadOfficialSiteFileChunkFromRequest(&ghttp.Request{Request: httpReq})
}
