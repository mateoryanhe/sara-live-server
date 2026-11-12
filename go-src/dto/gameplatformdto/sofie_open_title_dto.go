package gameplatformdto

// OpenTitleReq Sofie 打开标题(对应旧 AppGameStart;字段名避免 game).
type OpenTitleReq struct {
	TitleCode string `json:"titleCode" v:"required#标题编码不能为空" dc:"标题编码"`
	RoomId    uint64 `json:"roomId" dc:"直播间ID:0=大厅(清空关联);>0=该房间流水写入live_room_id"`
	UrlParams string `json:"urlParams" dc:"透传给厂商的 urlParams（原样上报）"`
}

// OpenTitleRes Sofie 打开标题应答.
type OpenTitleRes struct {
	Link string `json:"link" dc:"启动链接(有效期约2分钟,仅可使用一次)"`
}
