package userwithdrawdto

import "github.com/gogf/gf/v2/frame/g"

type CMSListWithdrawLogsReq struct {
	g.Meta   `path:"/listWithdrawLogs" method:"post" summary:"用户提现记录" tags:"用户提现CMS"`
	UserId   string `json:"userId"`
	Status   *uint8 `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type CMSWithdrawLogItem struct {
	Id                  string  `json:"id"`
	UserId              string  `json:"userId"`
	Nickname            string  `json:"nickname"`
	DiamondAmount       float64 `json:"diamondAmount"`
	GoldAmount          float64 `json:"goldAmount"`
	WithdrawUsd         float64 `json:"withdrawUsd"`
	Status              uint8   `json:"status"`
	TransferOrderId     string  `json:"transferOrderId"`
	TransferPlatformNo  string  `json:"transferPlatformNo"`
	TransferCurrency    string  `json:"transferCurrency"`
	TransferLocalAmount float64 `json:"transferLocalAmount"`
	TransferFailMsg     string  `json:"transferFailMsg"`
	RejectReason        string  `json:"rejectReason"`
	CreatedAt           string  `json:"createdAt"`
	TransferAt          string  `json:"transferAt"`
}

type CMSListWithdrawLogsRes struct {
	Total int                   `json:"total"`
	List  []*CMSWithdrawLogItem `json:"list"`
}

type CMSBatchApproveWithdrawReq struct {
	g.Meta `path:"/batchApproveWithdraw" method:"post" summary:"批量审核通过" tags:"用户提现CMS"`
	Ids    []string `json:"ids"`
}

type CMSBatchApproveWithdrawRes struct {
	SuccessCount int `json:"successCount"`
	FailCount    int `json:"failCount"`
}

type CMSBatchRejectWithdrawReq struct {
	g.Meta `path:"/batchRejectWithdraw" method:"post" summary:"批量审核拒绝" tags:"用户提现CMS"`
	Ids    []string `json:"ids"`
	Reason string `json:"reason"`
}

type CMSBatchRejectWithdrawRes struct {
	SuccessCount int `json:"successCount"`
	FailCount    int `json:"failCount"`
}

type CMSBatchRevokeWithdrawApprovalReq struct {
	g.Meta `path:"/batchRevokeWithdrawApproval" method:"post" summary:"批量驳回到待审核" tags:"用户提现CMS"`
	Ids    []string `json:"ids"`
}

type CMSBatchRevokeWithdrawApprovalRes struct {
	SuccessCount int `json:"successCount"`
	FailCount    int `json:"failCount"`
}

type CMSBatchSubmitPayoutReq struct {
	g.Meta `path:"/batchSubmitPayout" method:"post" summary:"批量提交HaiPay代付" tags:"用户提现CMS"`
	Ids    []string `json:"ids"`
}

type CMSBatchSubmitPayoutRes struct {
	SuccessCount int    `json:"successCount"`
	FailCount    int    `json:"failCount"`
	Message      string `json:"message"`
}
