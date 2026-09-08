package upload

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const defaultS3Region = "auto"

const (
	StoreCatImages     = "images"
	StoreCatShortVideo = "shortvideo"
	StoreCatCMS        = "cms"
	StoreCatExport     = "export"
)

var (
	s3ClientMu  sync.Mutex
	s3Client    *s3.Client
	s3ClientKey string
	s3Bucket    string
)

func invalidateS3Client() {
	s3ClientMu.Lock()
	defer s3ClientMu.Unlock()
	s3Client = nil
	s3ClientKey = ""
	s3Bucket = ""
}

func s3Ready(snap *resourceCfgSnapshot) bool {
	if snap == nil || !snap.S3Enabled {
		return false
	}
	return snap.S3Endpoint != "" && snap.S3Bucket != "" && snap.S3AccessKeyId != "" && snap.S3SecretAccessKey != ""
}

func getS3Client() (*s3.Client, string, error) {
	snap := getResourceCfgCache()
	if !s3Ready(snap) {
		return nil, "", errors.New("s3 cfg incomplete")
	}
	// R2 / 多数 S3 兼容端点对 Region 无实质区分，固定 auto 即可（AWS SDK 仍要求非空）
	key := snap.S3Endpoint + "|" + defaultS3Region + "|" + snap.S3Bucket + "|" + snap.S3AccessKeyId + "|" + snap.S3SecretAccessKey
	s3ClientMu.Lock()
	defer s3ClientMu.Unlock()
	if s3Client != nil && s3ClientKey == key {
		return s3Client, s3Bucket, nil
	}
	endpoint := strings.TrimRight(snap.S3Endpoint, "/")
	// 若误把桶名写进 Endpoint(如 ...cloudflarestorage.com/1v1),只保留 scheme://host
	if u, err := url.Parse(endpoint); err == nil && u.Host != "" {
		endpoint = u.Scheme + "://" + u.Host
	}
	client := s3.NewFromConfig(aws.Config{
		Region:      defaultS3Region,
		Credentials: credentials.NewStaticCredentialsProvider(snap.S3AccessKeyId, snap.S3SecretAccessKey, ""),
		// R2 不兼容 SDK 默认强制 CRC32 校验,否则 PutObject 常返回 401 Unauthorized
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
	}, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})
	s3Client = client
	s3ClientKey = key
	s3Bucket = snap.S3Bucket
	return client, s3Bucket, nil
}

// objectKeyFromStored 业务相对路径 -> 桶内完整 Key(含环境前缀)
func objectKeyFromStored(storedName string) string {
	storedName = strings.Trim(strings.ReplaceAll(storedName, "\\", "/"), "/")
	if storedName == "" {
		return ""
	}
	prefix := strings.Trim(GetS3KeyPrefix(), "/")
	if prefix == "" {
		return storedName
	}
	if storedName == prefix || strings.HasPrefix(storedName, prefix+"/") {
		return storedName
	}
	return prefix + "/" + storedName
}

func joinStoreCategory(category, baseName string) string {
	baseName = strings.Trim(strings.ReplaceAll(baseName, "\\", "/"), "/")
	category = strings.Trim(strings.ReplaceAll(category, "\\", "/"), "/")
	if category == "" {
		return baseName
	}
	if baseName == "" {
		return category
	}
	return category + "/" + baseName
}

func contentTypeByExt(ext string) string {
	ext = strings.ToLower(ext)
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}
	return "application/octet-stream"
}

func putLocalFileToS3(storedName, localPath string) error {
	client, bucket, err := getS3Client()
	if err != nil {
		return err
	}
	key := objectKeyFromStored(storedName)
	if key == "" {
		return errors.New("empty s3 object key")
	}
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	ctx := gctx.New()
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          f,
		ContentLength: aws.Int64(st.Size()),
		ContentType:   aws.String(contentTypeByExt(filepath.Ext(storedName))),
	})
	if err != nil {
		g.Log().Warningf(ctx, "s3 PutObject failed key=%s err=%v", key, err)
		return err
	}
	return nil
}

func putBytesToS3(storedName string, data []byte) error {
	client, bucket, err := getS3Client()
	if err != nil {
		return err
	}
	key := objectKeyFromStored(storedName)
	if key == "" {
		return errors.New("empty s3 object key")
	}
	ctx := gctx.New()
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String(contentTypeByExt(filepath.Ext(storedName))),
	})
	if err != nil {
		g.Log().Warningf(ctx, "s3 PutObject(bytes) failed key=%s err=%v", key, err)
		return err
	}
	return nil
}

func getObjectBytesFromS3(storedName string) ([]byte, error) {
	client, bucket, err := getS3Client()
	if err != nil {
		return nil, err
	}
	key := objectKeyFromStored(storedName)
	if key == "" {
		return nil, errors.New("empty s3 object key")
	}
	ctx := gctx.New()
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}

func deleteObjectFromS3(storedName string) {
	if !IsS3Enabled() || storedName == "" {
		return
	}
	client, bucket, err := getS3Client()
	if err != nil {
		g.Log().Warningf(gctx.New(), "s3 delete skip cfg err=%v name=%s", err, storedName)
		return
	}
	key := objectKeyFromStored(storedName)
	if key == "" {
		return
	}
	ctx := context.Background()
	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		g.Log().Warningf(ctx, "s3 DeleteObject failed key=%s err=%v", key, err)
	}
}

// PublishLocalFileToS3 将本地已写好的文件上传到云桶(CMS 导出等)
func PublishLocalFileToS3(category, fileName, absPath string) error {
	if !IsS3Enabled() {
		return nil
	}
	stored := joinStoreCategory(category, fileName)
	return putLocalFileToS3(stored, absPath)
}

// BuildExportFileURL CMS 导出文件公开 URL
func BuildExportFileURL(fileName string) string {
	fileName = strings.Trim(strings.ReplaceAll(fileName, "\\", "/"), "/")
	if fileName == "" {
		return ""
	}
	if IsS3Enabled() {
		return GetUrlByName(joinStoreCategory(StoreCatExport, fileName))
	}
	return GetUrlByName(fileName)
}

// DeleteExportStored 删除导出对象(本地由 fileexport 删;此处删云桶)
func DeleteExportStored(fileName string) {
	fileName = strings.Trim(strings.ReplaceAll(fileName, "\\", "/"), "/")
	if fileName == "" {
		return
	}
	deleteObjectFromS3(joinStoreCategory(StoreCatExport, fileName))
}
