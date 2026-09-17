package shortvideodto

// AppPublishShortVideoRes App 流式上传并发布短视频响应
type AppPublishShortVideoRes struct {
	ID           string `json:"id" dc:"短视频 ID"`
	Video        string `json:"video" dc:"视频完整URL"`
	PreviewVideo string `json:"previewVideo" dc:"试看视频完整URL,未上传时为空"`
	Cover        string `json:"cover" dc:"封面完整URL"`
}

// AppPublishShortVideoMultipart App 流式上传发布短视频(multipart/form-data, POST /shortVideo/appPublishShortVideo)
// 不走 GoFrame Bind,由服务端 MultipartReader 边读边写落盘.
// 字段:
//   - file: 短视频文件(必填,仅支持 mp4)
//   - previewFile: 试看视频文件(可选,仅支持 mp4)
//   - cover: 封面图片(可选)
//   - title, isPaid, payDiamond, categoryId, source, duration, freeWatchSeconds: 文本字段
