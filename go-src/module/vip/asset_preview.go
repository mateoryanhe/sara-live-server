package vip

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	liveentity "xr-game-server/entity/live"
	"xr-game-server/module/upload"
)

type vipAssetField struct {
	Key       string
	Label     string
	FileName  string
	URL       string
	MediaKind string
	Link      string
}

type vipAssetPreviewItem struct {
	ID         uint64
	Level      uint32
	LevelName  string
	LevelIcon  string
	IconURL    string
	DetailLink string
	Fields     []vipAssetField
}

type vipAssetPreviewPageData struct {
	Title string
	Items []vipAssetPreviewItem
}

type vipAssetPreviewResourceData struct {
	Title      string
	CfgID      uint64
	Level      uint32
	LevelName  string
	FieldLabel string
	FileName   string
	URL        string
	MediaKind  string
}

func RenderVipAssetPreviewListHTML() (string, error) {
	data := vipAssetPreviewPageData{
		Title: "VIP 资源预览",
		Items: buildVipAssetPreviewItems(),
	}
	return renderVipAssetPreviewTemplate("vipAssetPreviewList", vipAssetPreviewListHTML, data)
}

func RenderVipAssetPreviewResourceHTML(cfgID uint64, fieldKey string) (string, error) {
	data, err := buildVipAssetPreviewResourcePage(cfgID, fieldKey)
	if err != nil {
		return "", err
	}
	return renderVipAssetPreviewTemplate("vipAssetPreviewResource", vipAssetPreviewResourceHTML, data)
}

func buildVipAssetPreviewItems() []vipAssetPreviewItem {
	rows := listVipCfgFromMemory("")
	items := make([]vipAssetPreviewItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		fields := collectVipAssetFields(row)
		iconName := strings.TrimSpace(row.LevelIcon)
		items = append(items, vipAssetPreviewItem{
			ID:         row.ID,
			Level:      row.Level,
			LevelName:  row.LevelName,
			LevelIcon:  iconName,
			IconURL:    resolveVipAssetURL(iconName),
			DetailLink: fmt.Sprintf("/vipCfg/assetPreview/resource?id=%d&field=levelIcon", row.ID),
			Fields:     fields,
		})
	}
	return items
}

func buildVipAssetPreviewResourcePage(cfgID uint64, fieldKey string) (*vipAssetPreviewResourceData, error) {
	row := getVipCfgByIDFromMemory(cfgID)
	if row == nil {
		return nil, fmt.Errorf("vip cfg not found")
	}
	fieldKey = strings.TrimSpace(fieldKey)
	if fieldKey == "" {
		fieldKey = "levelIcon"
	}
	var selected *vipAssetField
	for _, field := range collectVipAssetFields(row) {
		if field.Key == fieldKey {
			copyField := field
			selected = &copyField
			break
		}
	}
	if selected == nil {
		return nil, fmt.Errorf("vip asset field not found")
	}
	return &vipAssetPreviewResourceData{
		Title:      fmt.Sprintf("VIP 资源预览 #%d", row.ID),
		CfgID:      row.ID,
		Level:      row.Level,
		LevelName:  row.LevelName,
		FieldLabel: selected.Label,
		FileName:   selected.FileName,
		URL:        selected.URL,
		MediaKind:  selected.MediaKind,
	}, nil
}

func collectVipAssetFields(row *liveentity.VipCfg) []vipAssetField {
	defs := []struct {
		key   string
		label string
		name  string
	}{
		{"levelIcon", "等级图标", row.LevelIcon},
		{"animationIcon", "进场特效图标", row.AnimationIcon},
		{"animation", "进场特效动画", row.Animation},
		{"commentEffectIcon", "公屏特效图标", row.CommentEffectIcon},
		{"commentEffect", "公屏特效动画", row.CommentEffect},
		{"customerServiceIcon", "客服优先图标", row.CustomerServiceIcon},
	}
	fields := make([]vipAssetField, 0, len(defs))
	for _, def := range defs {
		name := strings.TrimSpace(def.name)
		fields = append(fields, vipAssetField{
			Key:       def.key,
			Label:     def.label,
			FileName:  name,
			URL:       resolveVipAssetURL(name),
			MediaKind: detectVipAssetMediaKind(name),
			Link:      fmt.Sprintf("/vipCfg/assetPreview/resource?id=%d&field=%s", row.ID, def.key),
		})
	}
	return fields
}

func renderVipAssetPreviewTemplate(name, tmplText string, data any) (string, error) {
	tmpl, err := template.New(name).Parse(tmplText)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func resolveVipAssetURL(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return ""
	}
	if url := upload.GetUrlByName(fileName); url != "" {
		return url
	}
	return "/images/" + strings.TrimLeft(strings.ReplaceAll(fileName, "\\", "/"), "/")
}

func detectVipAssetMediaKind(fileName string) string {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(fileName)))
	switch ext {
	case ".mp4", ".webm", ".mov":
		return "video"
	case ".gif", ".png", ".jpg", ".jpeg", ".webp", ".bmp", ".apng":
		return "image"
	default:
		if fileName == "" {
			return "empty"
		}
		return "other"
	}
}

const vipAssetPreviewListHTML = `<!DOCTYPE html>
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
    main { padding: 24px; display: grid; gap: 16px; }
    .card { background: #fff; border: 1px solid #e5e7eb; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 2px rgba(0,0,0,.04); display: grid; grid-template-columns: 140px 1fr; }
    .icon-link { display: block; background: #fafafa; min-height: 140px; }
    .icon-link img { width: 100%; height: 140px; object-fit: contain; display: block; }
    .icon-placeholder { width: 100%; height: 140px; display: flex; align-items: center; justify-content: center; color: #9ca3af; font-size: 13px; }
    .meta { padding: 14px 16px; }
    .name { font-weight: 700; font-size: 16px; margin-bottom: 6px; }
    .sub { font-size: 12px; color: #6b7280; line-height: 1.5; word-break: break-all; margin-bottom: 10px; }
    .fields { display: flex; flex-wrap: wrap; gap: 8px; }
    .field { display: inline-flex; align-items: center; gap: 8px; padding: 6px 10px; border: 1px solid #e5e7eb; border-radius: 999px; text-decoration: none; color: #111827; background: #f9fafb; font-size: 12px; }
    .field img { width: 28px; height: 28px; object-fit: contain; border-radius: 4px; background: #fff; }
    .field.empty { color: #9ca3af; }
    .empty { padding: 48px; text-align: center; color: #6b7280; background: #fff; border-radius: 12px; }
    @media (max-width: 720px) {
      .card { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <header>
    <h1>{{.Title}}</h1>
    <p class="hint">点击等级图标或下方资源标签，打开新页面核对 VIP 图标/动画是否正确。</p>
  </header>
  <main>
    {{if .Items}}
      {{range .Items}}
      <article class="card">
        <a class="icon-link" href="{{.DetailLink}}" target="_blank" rel="noopener noreferrer" title="查看 {{.LevelName}} 等级图标">
          {{if .IconURL}}
          <img src="{{.IconURL}}" alt="{{.LevelName}}">
          {{else}}
          <div class="icon-placeholder">无等级图标</div>
          {{end}}
        </a>
        <div class="meta">
          <div class="name">#{{.ID}} LV{{.Level}} {{.LevelName}}</div>
          <div class="sub">levelIcon: {{if .LevelIcon}}{{.LevelIcon}}{{else}}-{{end}}</div>
          <div class="fields">
            {{range .Fields}}
            <a class="field {{if not .FileName}}empty{{end}}" href="{{.Link}}" target="_blank" rel="noopener noreferrer">
              {{if and .URL (eq .MediaKind "image")}}
              <img src="{{.URL}}" alt="{{.Label}}">
              {{end}}
              <span>{{.Label}}{{if .FileName}}: {{.FileName}}{{else}}: 未配置{{end}}</span>
            </a>
            {{end}}
          </div>
        </div>
      </article>
      {{end}}
    {{else}}
    <div class="empty">暂无 VIP 配置</div>
    {{end}}
  </main>
</body>
</html>`

const vipAssetPreviewResourceHTML = `<!DOCTYPE html>
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
    main { min-height: calc(100vh - 96px); display: flex; align-items: center; justify-content: center; padding: 24px; }
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
    <h1>#{{.CfgID}} LV{{.Level}} {{.LevelName}} · {{.FieldLabel}}</h1>
    <p class="sub">file: {{if .FileName}}{{.FileName}}{{else}}未配置{{end}}</p>
  </header>
  <main>
    <div class="panel">
      <div class="preview">
        {{if eq .MediaKind "video"}}
        <video src="{{.URL}}" controls autoplay loop playsinline></video>
        {{else if eq .MediaKind "image"}}
        <img src="{{.URL}}" alt="{{.FieldLabel}}">
        {{else if .URL}}
        <div class="fallback">
          <p>当前资源格式暂不支持内嵌预览。</p>
          <p><a href="{{.URL}}" target="_blank" rel="noopener noreferrer">点击下载/打开资源</a></p>
        </div>
        {{else}}
        <div class="fallback">未配置该资源</div>
        {{end}}
      </div>
      <a class="back" href="/vipCfg/assetPreview">返回 VIP 列表</a>
    </div>
  </main>
</body>
</html>`
