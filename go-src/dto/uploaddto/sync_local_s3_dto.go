package uploaddto

import "github.com/gogf/gf/v2/frame/g"

// SyncLocalStorageToS3Req 将本地 storagePath 下历史文件刷到云桶
type SyncLocalStorageToS3Req struct {
	g.Meta `path:"/syncLocalStorageToS3" method:"post" summary:"本地存储刷到云桶" tags:"上传配置"`
}

type SyncLocalStorageToS3Res struct {
	Started         bool   `json:"started" dc:"本次是否新启动任务"`
	AlreadyRunning  bool   `json:"alreadyRunning" dc:"已有任务在跑"`
	Message         string `json:"message"`
	*SyncLocalToS3Status
}

// GetSyncLocalStorageToS3StatusReq 查询刷桶进度
type GetSyncLocalStorageToS3StatusReq struct {
	g.Meta `path:"/getSyncLocalStorageToS3Status" method:"post" summary:"查询本地刷云桶进度" tags:"上传配置"`
}

type GetSyncLocalStorageToS3StatusRes struct {
	*SyncLocalToS3Status
}

// SyncLocalToS3Status 刷桶任务状态
type SyncLocalToS3Status struct {
	Running    bool   `json:"running"`
	Root       string `json:"root" dc:"本地扫描根目录=storagePath"`
	KeyPrefix  string `json:"keyPrefix" dc:"云桶环境前缀"`
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
	Total      int    `json:"total"`
	Done       int    `json:"done"`
	Success    int    `json:"success"`
	Failed     int    `json:"failed"`
	Skipped    int    `json:"skipped"`
	LastKey    string `json:"lastKey"`
	LastError  string `json:"lastError"`
}
