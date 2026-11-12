package attrreportdto

import "github.com/gogf/gf/v2/frame/g"

type GetAttrReportCfgReq struct {
	g.Meta `path:"/getAttrReportCfg" method:"post" summary:"查询归因数据上报配置" tags:"归因数据上报"`
}

type AttrReportCfgItem struct {
	ID        string `json:"id"`
	BaseApi   string `json:"baseApi" dc:"归因上报主机根地址(自动拼/tj/api/accept.php)"`
	SecretKey string `json:"secretKey" dc:"归因上报密钥(api_token)"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GetAttrReportCfgRes struct {
	Cfg *AttrReportCfgItem `json:"cfg"`
}

type SaveAttrReportCfgReq struct {
	g.Meta    `path:"/saveAttrReportCfg" method:"post" summary:"保存归因数据上报配置" tags:"归因数据上报"`
	ID        uint64 `json:"id" dc:"配置ID,新建传0"`
	BaseApi   string `json:"baseApi" v:"required|max-length:512#Base API不能为空|Base API最长512字符" dc:"主机根地址,如 https://cdntk.showsme.click"`
	SecretKey string `json:"secretKey" v:"required|max-length:255#密钥不能为空|密钥最长255字符" dc:"归因上报密钥(api_token)"`
}

type SaveAttrReportCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
