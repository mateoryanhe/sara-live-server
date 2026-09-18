package golddto

import "github.com/gogf/gf/v2/frame/g"

// AppTransferGoldReq App端转赠金币(当前登录用户扣款,目标用户到账)
type AppTransferGoldReq struct {
	g.Meta       `path:"/transferGold" method:"post" summary:"转赠金币给指定用户" tags:"金币"`
	TargetUserId string  `json:"targetUserId" v:"required#目标用户ID不能为空" dc:"接收金币的用户ID"`
	Amount       float64 `json:"amount" v:"required|min:0.01#金币数量不能为空|金币数量至少0.01" dc:"转赠金币数量(最多2位小数)"`
}

type AppTransferGoldRes struct {
	Amount         float64 `json:"amount" dc:"实际转赠数量(已按2位小数规范化)"`
	SenderGold     float64 `json:"senderGold" dc:"转出后当前用户金币余额"`
	TargetUserId   string  `json:"targetUserId" dc:"接收用户ID"`
	TargetUserGold float64 `json:"targetUserGold" dc:"接收后目标用户金币余额"`
}

// AppCoinMerchantGoldTransferRecordListReq 币商分页查询自己的金币转账记录。
type AppCoinMerchantGoldTransferRecordListReq struct {
	g.Meta       `path:"/getCoinMerchantTransferRecordList" method:"post" summary:"分页查询币商金币转账记录" tags:"金币"`
	TargetUserId string `json:"targetUserId" dc:"目标用户ID查询关键字(可选,精确匹配)"`
	StartTime    int64  `json:"startTime" dc:"开始时间(可选,秒时间戳,包含)"`
	EndTime      int64  `json:"endTime" dc:"结束时间(可选,秒时间戳,包含)"`
	PageIndex    int    `json:"pageIndex" dc:"页码(从1开始,默认1)"`
	PageSize     int    `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

type AppCoinMerchantGoldTransferRecordItem struct {
	Id             string  `json:"id" dc:"转账记录ID"`
	TargetUserId   string  `json:"targetUserId" dc:"目标用户ID"`
	TargetNickname string  `json:"targetNickname" dc:"目标用户昵称"`
	TargetAvatar   string  `json:"targetAvatar" dc:"目标用户头像完整URL"`
	Amount         float64 `json:"amount" dc:"转赠金币数量"`
	CreatedAt      int64   `json:"createdAt" dc:"转账时间(秒时间戳)"`
	CreatedAtText  string  `json:"createdAtText" dc:"服务器格式化的转账时间(YYYY-MM-DD HH:mm:ss)"`
}

type AppCoinMerchantGoldTransferRecordListRes struct {
	Total     int                                      `json:"total" dc:"记录总数"`
	PageIndex int                                      `json:"pageIndex" dc:"当前页码"`
	PageSize  int                                      `json:"pageSize" dc:"每页数量"`
	List      []*AppCoinMerchantGoldTransferRecordItem `json:"list" dc:"转账记录列表"`
}
