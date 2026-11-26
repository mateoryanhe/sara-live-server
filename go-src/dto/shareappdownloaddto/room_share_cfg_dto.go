package shareappdownloaddto

import "github.com/gogf/gf/v2/frame/g"

type GetLiveRoomShareCfgReq struct {
	g.Meta `path:"/getLiveRoomShareCfg" method:"post" summary:"查询直播间分享配置" tags:"直播间分享"`
}

type LiveRoomShareCfgItem struct {
	ID            string `json:"id"`
	ShareTextZhCN string `json:"shareTextZhCN"`
	ShareTextZhTW string `json:"shareTextZhTW"`
	ShareTextEn   string `json:"shareTextEn"`
	ShareTextEs   string `json:"shareTextEs"`
	ShareTextPt   string `json:"shareTextPt"`
	ShareTextHi   string `json:"shareTextHi"`
	ShareTextId   string `json:"shareTextId"`
	AfLink        string `json:"afLink"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	// 内置默认文案(只读, CMS 空栏预览与 App 回落一致)
	BuiltinShareTextEn string `json:"builtinShareTextEn"`
	BuiltinShareTextEs string `json:"builtinShareTextEs"`
	BuiltinShareTextPt string `json:"builtinShareTextPt"`
	BuiltinShareTextHi string `json:"builtinShareTextHi"`
	BuiltinShareTextId string `json:"builtinShareTextId"`
}

type GetLiveRoomShareCfgRes struct {
	Cfg *LiveRoomShareCfgItem `json:"cfg"`
}

type SaveLiveRoomShareCfgReq struct {
	g.Meta        `path:"/saveLiveRoomShareCfg" method:"post" summary:"保存直播间分享配置" tags:"直播间分享"`
	ID            uint64 `json:"id"`
	ShareTextZhCN string `json:"shareTextZhCN"`
	ShareTextZhTW string `json:"shareTextZhTW"`
	ShareTextEn   string `json:"shareTextEn"`
	ShareTextEs   string `json:"shareTextEs"`
	ShareTextPt   string `json:"shareTextPt"`
	ShareTextHi   string `json:"shareTextHi"`
	ShareTextId   string `json:"shareTextId"`
	AfLink        string `json:"afLink" v:"required|max-length:1024#请填写 AF 链接|链接过长"`
}

type SaveLiveRoomShareCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}