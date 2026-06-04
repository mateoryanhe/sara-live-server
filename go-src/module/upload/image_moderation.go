package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"

	"xr-game-server/core/snowflake"
	"xr-game-server/errercode"
)

const (
	defaultImageModerationRegion   = "cn-shanghai"
	defaultImageModerationEndpoint = "green-cip.cn-shanghai.aliyuncs.com"
	defaultImageModerationService  = "profilePhotoCheck"
	imageModerationAPIVersion      = "2022-03-02"
	imageModerationAPIAction       = "ImageModeration"
)

var (
	imageModerationClientMu  sync.Mutex
	imageModerationClient    *openapi.Client
	imageModerationClientKey string
)

type imageModerationCfg struct {
	enabled   bool
	accessKey string
	secret    string
	regionId  string
	endpoint  string
	service   string
}

type imageModerationAPIResp struct {
	Code    int                  `json:"Code"`
	Message string               `json:"Message"`
	Data    *imageModerationData `json:"Data"`
}

type imageModerationData struct {
	RiskLevel string                  `json:"RiskLevel"`
	Result    []imageModerationResult `json:"Result"`
}

type imageModerationResult struct {
	RiskLevel   string  `json:"RiskLevel"`
	Label       string  `json:"Label"`
	Description string  `json:"Description"`
	Confidence  float32 `json:"Confidence"`
}

func imageModerationFromSnapshot(s *resourceCfgSnapshot) imageModerationCfg {
	if s == nil {
		return imageModerationCfg{}
	}
	return imageModerationCfg{
		enabled:   s.ImageModerationEnabled,
		accessKey: s.ImageModerationAccessKeyId,
		secret:    s.ImageModerationAccessKeySecret,
		regionId:  s.ImageModerationRegionId,
		endpoint:  s.ImageModerationEndpoint,
		service:   s.ImageModerationService,
	}
}

func (c imageModerationCfg) ready() bool {
	return c.enabled && c.accessKey != "" && c.secret != "" && c.endpoint != ""
}

// ImageModerationEnabled 是否开启 App 图片审核
func ImageModerationEnabled() bool {
	return getResourceCfgCache().ImageModerationEnabled
}

func getImageModerationOpenAPIClient(cfg imageModerationCfg) (*openapi.Client, error) {
	key := cfg.accessKey + "|" + cfg.secret + "|" + cfg.regionId + "|" + cfg.endpoint
	imageModerationClientMu.Lock()
	defer imageModerationClientMu.Unlock()
	if imageModerationClient != nil && imageModerationClientKey == key {
		return imageModerationClient, nil
	}
	conf := &openapi.Config{
		AccessKeyId:     tea.String(cfg.accessKey),
		AccessKeySecret: tea.String(cfg.secret),
		RegionId:        tea.String(cfg.regionId),
		Endpoint:        tea.String(cfg.endpoint),
		ConnectTimeout:  tea.Int(3000),
		ReadTimeout:     tea.Int(10000),
	}
	client, err := openapi.NewClient(conf)
	if err != nil {
		return nil, err
	}
	imageModerationClient = client
	imageModerationClientKey = key
	return client, nil
}

func invalidateImageGreenClient() {
	imageModerationClientMu.Lock()
	defer imageModerationClientMu.Unlock()
	imageModerationClient = nil
	imageModerationClientKey = ""
}

func buildImageURL(fileName string) string {
	return buildResourceUrl(fmt.Sprintf("/%s/%s", ImageSubDir, fileName))
}

func isPublicResourceURL(rawURL string) bool {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Host == "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if host == "127.0.0.1" || host == "localhost" {
		return false
	}
	if strings.HasPrefix(host, "192.168.") || strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "172.") {
		return false
	}
	return strings.HasPrefix(strings.ToLower(u.Scheme), "http")
}

// callImageModeration 调用 Green ImageModeration(2022-03-02),通过 imageUrl 检测
func callImageModeration(ctx context.Context, cfg imageModerationCfg, imageURL string) (*imageModerationAPIResp, error) {
	client, err := getImageModerationOpenAPIClient(cfg)
	if err != nil {
		return nil, err
	}
	service := cfg.service
	if service == "" {
		service = defaultImageModerationService
	}
	serviceParams, err := json.Marshal(map[string]string{
		"imageUrl": imageURL,
		"dataId":   uuid.New().String(),
	})
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"Service":           service,
		"ServiceParameters": string(serviceParams),
	}
	params := &openapiutil.Params{
		Action:      tea.String(imageModerationAPIAction),
		Version:     tea.String(imageModerationAPIVersion),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	runtime := &util.RuntimeOptions{}
	runtime.SetConnectTimeout(3000)
	runtime.SetReadTimeout(10000)

	result, err := client.CallApi(params, &openapiutil.OpenApiRequest{Body: body}, runtime)
	if err != nil {
		g.Log().Warningf(ctx, "ImageModeration CallApi err: %v", err)
		return nil, err
	}
	statusCode := 0
	if v, ok := result["statusCode"]; ok {
		switch n := v.(type) {
		case json.Number:
			code, _ := n.Int64()
			statusCode = int(code)
		case float64:
			statusCode = int(n)
		case int:
			statusCode = n
		}
	}
	if statusCode != 0 && statusCode != http.StatusOK {
		g.Log().Warningf(ctx, "ImageModeration http status=%d", statusCode)
		return nil, fmt.Errorf("http status %d", statusCode)
	}
	bodyMap, err := util.AssertAsMap(result["body"])
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}
	var parsed imageModerationAPIResp
	if err = json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func hasImageRisk(data *imageModerationData) bool {
	if data == nil {
		return false
	}
	if isRiskLevel(data.RiskLevel) {
		return true
	}
	for _, r := range data.Result {
		if isRiskLevel(r.RiskLevel) {
			return true
		}
		label := strings.ToLower(strings.TrimSpace(r.Label))
		if label != "" && label != "nonlabel" && label != "normal" {
			desc := strings.ToLower(r.Description)
			if !strings.Contains(desc, "no risk") {
				return true
			}
		}
	}
	return false
}

func isRiskLevel(level string) bool {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "high", "medium":
		return true
	default:
		return false
	}
}

func moderateAppImage(ctx context.Context, cfg imageModerationCfg, fileName string) error {
	imageURL := buildImageURL(fileName)
	if !isPublicResourceURL(imageURL) {
		g.Log().Warningf(ctx, "image moderation requires public resource domain, url=%s", imageURL)
		return errercode.CreateCode(errercode.ImageModerationCfgInvalid)
	}
	resp, err := callImageModeration(ctx, cfg, imageURL)
	if err != nil {
		return errercode.CreateCode(errercode.ImageModerationFailed)
	}
	if resp.Code != 0 && resp.Code != 200 {
		g.Log().Warningf(ctx, "ImageModeration code=%d msg=%s", resp.Code, resp.Message)
		return errercode.CreateCode(errercode.ImageModerationFailed)
	}
	if hasImageRisk(resp.Data) {
		return errercode.CreateCode(errercode.ImageSensitiveContent)
	}
	return nil
}

// RequireAppImageCompliant 审核 App 已保存到本地的图片文件名;未开启时跳过
func RequireAppImageCompliant(ctx context.Context, fileName string) error {
	cfg := imageModerationFromSnapshot(getResourceCfgCache())
	if !cfg.enabled {
		return nil
	}
	if !cfg.ready() {
		return errercode.CreateCode(errercode.ImageModerationCfgInvalid)
	}
	return moderateAppImage(ctx, cfg, fileName)
}

// UploadImageForApp App 端上传:先落本地,再 ImageModeration 审核,违规则删文件并返回错误码
func UploadImageForApp(ctx context.Context, file *ghttp.UploadFile) (string, error) {
	name, fullPath, err := saveUploadedImageFile(file)
	if err != nil {
		return "", err
	}
	if err := RequireAppImageCompliant(ctx, name); err != nil {
		_ = os.Remove(fullPath)
		return "", err
	}
	return name, nil
}

// saveUploadedImageFile 校验并保存图片到 upload/images,返回文件名与绝对路径
func saveUploadedImageFile(file *ghttp.UploadFile) (name, fullPath string, err error) {
	if file == nil {
		return "", "", fmt.Errorf("upload file is empty")
	}
	if file.Size > MaxImageSize {
		return "", "", fmt.Errorf("image too large, max=%d bytes", MaxImageSize)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedImageExt[ext]; !ok {
		return "", "", fmt.Errorf("image ext not allowed: %s", ext)
	}
	dir := getImageDir()
	if err = os.MkdirAll(dir, 0755); err != nil {
		return "", "", err
	}
	name = fmt.Sprintf("%d%s", snowflake.GetId(), ext)
	fullPath = filepath.Join(dir, name)
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", "", err
	}
	if _, err = io.Copy(dst, src); err != nil {
		_ = dst.Close()
		return "", "", err
	}
	if err = dst.Close(); err != nil {
		return "", "", err
	}
	return name, fullPath, nil
}
