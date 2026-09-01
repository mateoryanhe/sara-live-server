package gift

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/live"
	"xr-game-server/module/upload"
)

type assetPreviewItem struct {
	ID            uint64
	Name          string
	Status        uint8
	StatusLabel   string
	IconURL       string
	IconName      string
	AnimationURL  string
	AnimationName string
	AnimationLink string
}

type assetPreviewPageData struct {
	Title string
	Items []assetPreviewItem
}

type assetPreviewAnimationData struct {
	Title         string
	GiftID        uint64
	Name          string
	IconURL       string
	AnimationURL  string
	AnimationName string
	MediaKind     string
}

func RenderAssetPreviewListHTML() (string, error) {
	data := assetPreviewPageData{
		Title: "礼物资源预览",
		Items: buildAssetPreviewItems(),
	}
	return renderAssetPreviewTemplate(assetPreviewListHTML, data)
}

func RenderAssetPreviewAnimationHTML(giftID uint64) (string, error) {
	data, err := buildAssetPreviewAnimationPage(giftID)
	if err != nil {
		return "", err
	}
	return renderAssetPreviewTemplate(assetPreviewAnimationHTML, data)
}

func buildAssetPreviewItems() []assetPreviewItem {
	rows := cfgdao.GetAllGiftsForAssetPreview()
	items := make([]assetPreviewItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		iconName := strings.TrimSpace(row.Icon)
		animationName := strings.TrimSpace(row.Animation)
		items = append(items, assetPreviewItem{
			ID:            row.ID,
			Name:          row.Name,
			Status:        row.Status,
			StatusLabel:   giftStatusLabel(row.Status),
			IconURL:       resolveGiftAssetURL(iconName),
			IconName:      iconName,
			AnimationURL:  resolveGiftAssetURL(animationName),
			AnimationName: animationName,
			AnimationLink: fmt.Sprintf("/gift/assetPreview/animation?id=%d", row.ID),
		})
	}
	return items
}

func buildAssetPreviewAnimationPage(giftID uint64) (*assetPreviewAnimationData, error) {
	row := cfgdao.GetGiftById(giftID)
	if row == nil {
		return nil, fmt.Errorf("gift not found")
	}
	iconName := strings.TrimSpace(row.Icon)
	animationName := strings.TrimSpace(row.Animation)
	return &assetPreviewAnimationData{
		Title:         fmt.Sprintf("礼物动画预览 #%d", row.ID),
		GiftID:        row.ID,
		Name:          row.Name,
		IconURL:       resolveGiftAssetURL(iconName),
		AnimationURL:  resolveGiftAssetURL(animationName),
		AnimationName: animationName,
		MediaKind:     detectGiftAnimationMediaKind(animationName),
	}, nil
}

func renderAssetPreviewTemplate(tmplText string, data any) (string, error) {
	tmpl, err := template.New("giftAssetPreview").Parse(tmplText)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func resolveGiftAssetURL(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return ""
	}
	if url := upload.GetUrlByName(fileName); url != "" {
		return url
	}
	return "/images/" + strings.TrimLeft(strings.ReplaceAll(fileName, "\\", "/"), "/")
}

func giftStatusLabel(status uint8) string {
	if status == entity.LiveGiftStatusOnShelf {
		return "上架"
	}
	return "下架"
}

func detectGiftAnimationMediaKind(fileName string) string {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(fileName)))
	switch ext {
	case ".mp4", ".webm", ".mov":
		return "video"
	case ".gif", ".png", ".jpg", ".jpeg", ".webp", ".bmp", ".apng":
		return "image"
	default:
		return "other"
	}
}

const assetPreviewListHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}}</title>
  <style>
    * { box-sizing: border-box; }
    body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f5f7fb; color: #1f2937; }
    header { padding: 20px 24px; background: #fff; border-bottom: 1px solid #e5e7eb; }
    h1 { margin: 0 0 8px; font-size: 22px; }
    .hint { margin: 0; color: #6b7280; font-size: 14px; }
    main { padding: 24px; }
    .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
    .card { background: #fff; border: 1px solid #e5e7eb; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 2px rgba(0,0,0,.04); }
    .icon-link { display: block; aspect-ratio: 1; background: #fafafa; }
    .icon-link img { width: 100%; height: 100%; object-fit: contain; display: block; }
    .icon-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #9ca3af; font-size: 13px; }
    .meta { padding: 12px; }
    .name { font-weight: 600; font-size: 14px; margin-bottom: 6px; word-break: break-all; }
    .sub { font-size: 12px; color: #6b7280; line-height: 1.5; word-break: break-all; }
    .tag { display: inline-block; padding: 2px 8px; border-radius: 999px; font-size: 12px; margin-top: 6px; }
    .tag.on { background: #ecfdf5; color: #047857; }
    .tag.off { background: #f3f4f6; color: #6b7280; }
    .empty { padding: 48px; text-align: center; color: #6b7280; background: #fff; border-radius: 12px; }
  </style>
</head>
<body>
  <header>
    <h1>{{.Title}}</h1>
    <p class="hint">点击图标将在新页面打开并播放动画资源，用于核对 icon / animation 是否正确。</p>
  </header>
  <main>
    {{if .Items}}
    <div class="grid">
      {{range .Items}}
      <article class="card">
        <a class="icon-link" href="{{.AnimationLink}}" target="_blank" rel="noopener noreferrer" title="查看 {{.Name}} 动画">
          {{if .IconURL}}
          <img src="{{.IconURL}}" alt="{{.Name}}">
          {{else}}
          <div class="icon-placeholder">无图标</div>
          {{end}}
        </a>
        <div class="meta">
          <div class="name">#{{.ID}} {{.Name}}</div>
          <div class="sub">icon: {{if .IconName}}{{.IconName}}{{else}}-{{end}}</div>
          <div class="sub">animation: {{if .AnimationName}}{{.AnimationName}}{{else}}-{{end}}</div>
          <span class="tag {{if eq .StatusLabel "上架"}}on{{else}}off{{end}}">{{.StatusLabel}}</span>
        </div>
      </article>
      {{end}}
    </div>
    {{else}}
    <div class="empty">暂无礼物数据</div>
    {{end}}
  </main>
</body>
</html>`

const assetPreviewAnimationHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}}</title>
  <style>
    * { box-sizing: border-box; }
    body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #111827; color: #f9fafb; }
    header { padding: 16px 20px; background: #1f2937; border-bottom: 1px solid #374151; }
    h1 { margin: 0 0 6px; font-size: 20px; }
    .sub { margin: 0; color: #9ca3af; font-size: 13px; word-break: break-all; }
    main { min-height: calc(100vh - 88px); display: flex; align-items: center; justify-content: center; padding: 24px; }
    .panel { width: min(960px, 100%); background: #1f2937; border: 1px solid #374151; border-radius: 16px; padding: 20px; }
    .preview { background: #000; border-radius: 12px; min-height: 360px; display: flex; align-items: center; justify-content: center; overflow: hidden; }
    video, img { max-width: 100%; max-height: 70vh; display: block; }
    .fallback { padding: 24px; text-align: center; color: #d1d5db; line-height: 1.6; }
    .fallback a { color: #93c5fd; }
    .back { display: inline-block; margin-top: 16px; color: #93c5fd; text-decoration: none; }
  </style>
</head>
<body>
  <header>
    <h1>#{{.GiftID}} {{.Name}}</h1>
    <p class="sub">animation: {{if .AnimationName}}{{.AnimationName}}{{else}}未配置{{end}}</p>
  </header>
  <main>
    <div class="panel">
      <div class="preview">
        {{if eq .MediaKind "video"}}
        <video src="{{.AnimationURL}}" controls autoplay loop playsinline></video>
        {{else if eq .MediaKind "image"}}
        <img src="{{.AnimationURL}}" alt="{{.Name}}">
        {{else if .AnimationURL}}
        <div class="fallback">
          <p>当前动画格式暂不支持内嵌预览。</p>
          <p><a href="{{.AnimationURL}}" target="_blank" rel="noopener noreferrer">点击下载/打开资源</a></p>
        </div>
        {{else}}
        <div class="fallback">未配置动画资源</div>
        {{end}}
      </div>
      <a class="back" href="/gift/assetPreview">返回礼物列表</a>
    </div>
  </main>
</body>
</html>`
