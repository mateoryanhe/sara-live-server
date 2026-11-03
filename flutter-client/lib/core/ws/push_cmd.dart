/// 对齐 Go `constants/cmd`；负数为本机事件，非服务端下发。
abstract final class PushCmd {
  static const enter = 1;
  static const heart = 2;
  static const kick = 3;
  static const liveRoomGift = 8;
  static const liveRoomChat = 9;
  static const diamond = 10;
  static const gold = 11;

  /// VIP 等级变更：`{vipLevel, levelName}`
  static const vipLevel = 12;

  /// 充值到账（载荷含 `goldBalance` 余额）
  static const rechargeSuccess = 40;

  /// 视频通话请求（推送给接听方）
  static const callRequest = 25;
  static const callRejected = 26;
  static const callAccepted = 27;
  static const callEnded = 28;
  static const callStarted = 31;
  static const callTimeout = 32;

  /// 观众列表刷新（推送给房间内全体在线用户，含主播）
  static const liveRoomAudienceListRefresh = 24;
  static const liveRoomStopLive = 23;

  /// 付费弹幕（公屏）
  static const liveRoomPaidDanmaku = 22;

  /// 本机：WS 已连接（**已废弃**，请订阅 [WsConnectionBusEvent] `connected`）。
  @Deprecated('Subscribe to WsConnectionBusEvent(WsConnectionState.connected)')
  static const wsConnected = -100;
}
