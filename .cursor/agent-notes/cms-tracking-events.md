# CMS「埋点」一级菜单与事件页约定（2026-10-10）

## 产品形态

- CMS **新开一级目录**：菜单名 **埋点**（与「用户管理 / 日志 / 系统配置」同级，非挂在系统配置下）。
- **一个业务事件 = 一个 CMS 页面**（独立路由、独立权限页 `pageName`），不做「一个大页里下拉选事件」。
- 每个事件页统一布局：
  1. **总数**：当前统计口径下的汇总（与下方图表维度一致，或单独说明口径）。
  2. **日 / 周 / 月** 三个维度（Tab 或子 Tab，对齐 Dashboard 用户趋势）。
  3. **折线图**：该维度下按时间 bucket 的 **次数**（count），X 轴为日期/周/月，Y 轴为次数。

示例首事件：**发起视频通话次数**（用户/App 侧「发起呼叫」成功或创建 `call_order` 时 +1，口径实现时再定）。

## CMS 工程约定

| 项 | 约定 |
|----|------|
| 路由分组 | 新增 `cms/src/router/routes/tracking.ts`，`path: '/tracking'`，挂到 `layoutRouteGroups` |
| 视图目录 | `cms/src/views/tracking/<event-slug>.vue` |
| 权限树 | `permission-menu-tree.ts` 增加 `id: 'tracking'` 分组，`children` 下每事件一个 `page('TrackingXxx')` |
| 菜单 | `layout/index.vue` 顶栏 `el-sub-menu index="/tracking"` |
| 按钮 | 默认 **查看 + 刷新**（`DEFAULT_VIEW_BUTTONS` 或 `view` + 无 save） |
| 页面复用 | 抽 `TrackingEventTrendPage`  composable/组件：props `eventKey` + i18n 标题；图表复用 Dashboard 折线风格（如 `UserStatChart` 同类） |
| 时间 | 展示与筛选统一 `server-datetime.ts`（UTC+0） |

**pageName 命名**：`Tracking` + 事件 PascalCase，如 `TrackingVideoCallInitiated`（菜单文案用 i18n「发起视频通话次数」）。

**URL 示例**：`/tracking/video-call-initiated`

## Go 后端约定（与 `module/stat` 对齐）

- **不写** Firebase/第三方 SDK 配置页；埋点 = **服务端计数 + CMS 只读报表**。
- 事件标识：稳定字符串 `event_key`（如 `video_call_initiated`），代码常量 + 文档表。
- 写入：业务成功路径 **同步或 `xrpool.AddWithRecover` 入队**（高 QPS 用队列 + 定时刷盘，参考 `module/stat` 登录/观众统计）。
- 存储：按事件 **日 / 周 / 月** 聚合表或通用表 + `event_key` + `period_bucket`（与 `statdao` 日/周/月 login 类似）；热数据可 RowCache，CMS 读趋势接口可短 TTL 缓存。
- API（CMS）：`POST /trackingEvent/getVideoCallInitiatedTrend` 或通用 `POST /trackingEvent/getEventTrend` + `eventKey`（**每事件单独接口** 更易做权限映射时，优先 per-event）。
- **不进 Sofie**（纯 CMS 运营数据）；App 只负责触发业务，不上报埋点到 CMS API。

## 已落地事件

| 页面 | event | 口径 | CMS API |
|------|-------|------|---------|
| `TrackingHotLiveRoomJoinManagement`（秀场直播间） | `hot_live_room_join` | 观众 `joinRoom` 成功且 `live_room_cfgs.category=1`(秀场/Hot)，不含主播本人；每次 +1；`JoinRoom` Pub `ShowcaseLiveRoomJoinTrackingEvent`，`module/tracking` 队列异步消费 | `POST /trackingEvent/getHotLiveRoomJoinTrend` |
| 同上页 **Tab「直播首帧画面渲染完成」** | `live_first_frame_rendered` | 秀场(category=1)、直播中；App `POST /liveRoom/reportLiveFirstFrameRendered`（Sofie `POST /sofie/studio/notifyLiveFirstFrameReady`）；不含主播；每次成功 +1 | `POST /trackingEvent/getLiveFirstFrameRenderTrend` |
| 同上页 **Tab「退出直播间」** | `hot_live_room_leave` | 秀场(category=1)；`leaveRoom` 成功 +1，不含主播；仅主动退房 API（非踢出/心跳清理） | `POST /trackingEvent/getHotLiveRoomLeaveTrend` |
| 同上页 **Tab「礼物面板点击」** | `hot_live_room_gift_click` | `reportHotLiveRoomGiftClick` 登录成功 +1（暂不要求 body） | `POST /trackingEvent/getHotLiveRoomGiftClickTrend` |

App 秀场礼物面板：`POST /liveRoom/reportHotLiveRoomGiftClick`（Sofie `POST /sofie/studio/notifyHotLiveRoomGiftClick`）。

表：`daily_hot_live_room_join_stats` / `weekly_*` / `monthly_*`（仅 `count`）；首帧：`daily_live_first_frame_render_stats` / `weekly_*` / `monthly_*`；退房：`daily_hot_live_room_leave_stats` / `weekly_*` / `monthly_*`；游戏进房：`daily_game_live_room_join_stats` / `weekly_*` / `monthly_*`；秀场礼物点击：`daily_hot_live_room_gift_click_stats` / `weekly_*` / `monthly_*`。

| 页面 | event | 口径 | CMS API |
|------|-------|------|---------|
| `TrackingCall1v1InitiateManagement`（1v1视频通话） | `call_1v1_initiate` | `liveRoomCall` 成功且 `call_order.source=1`（直播间来源）+1，**不区分**直播间 category | `POST /trackingEvent/getCall1v1InitiateTrend` |

表：`daily_call1v1_initiate_stats` / `weekly_*` / `monthly_*`（GORM 表名无 `call_` 与 `1v1` 间下划线；`TbName`/syndb/dao 须与此一致，2026-10-19 正式服曾因此全页 0）。

| 同上页 **Tab「1v1房间视频通话」** | `call_1v1_room` | `oneToOneRoomCall` 成功且 `call_order.source=3` +1 | `POST /trackingEvent/getCall1v1RoomCallTrend` |

表：`daily_call1v1_room_stats` / `weekly_*` / `monthly_*`。

| 同上页 **Tab「双方成功接通」** | `call_1v1_connect_success` | 视频通话；双方 `confirmCall` 且首次扣费成功进入通话中 +1；source=1 限秀场 category=1，source=3 含 1v1 房间 | `POST /trackingEvent/getCall1v1ConnectSuccessTrend` |

表：`daily_call1v1_connect_success_stats` / `weekly_*` / `monthly_*`。

| 同上页 **Tab「1v1房间列表」** | `one_to_one_room_list` | App `oneToOneRoomList`（Sofie `browseDirectBooths`）成功返回 +1（含翻页） | `POST /trackingEvent/getOneToOneRoomListTrend` |

表：`daily_one_to_one_room_list_stats` / `weekly_*` / `monthly_*`。

| 页面 | event | 口径 | CMS API |
|------|-------|------|---------|
| `TrackingMiniGameRoundStartManagement`（游戏直播间） | — | **2026-10-21 起整页停用**：CMS 无 Tab/图表；`gameevent.GameLiveRoomTrackingPersistDisabled=true`，App 上报仍 200/success **不入库** | CMS 趋势接口仍注册，仅读历史 |
| ~~游戏直播间~~ | `game_live_room_join` | 游戏直播(category=2)；观众 `joinRoom` 成功 +1，不含主播 | `POST /trackingEvent/getGameLiveRoomJoinTrend` |
| 同上 **Tab「直播首帧画面渲染完成」** | `game_live_watch_stream` | 游戏直播(category=2)；`reportLiveFirstFrameRendered` 成功 +1，不含主播 | `POST /trackingEvent/getGameLiveWatchStreamTrend` |
| 同上 **Tab「点击开始游戏」** | `mini_game_round_start` | App `appGameStart` 成功返回非空启动 link +1 | `POST /trackingEvent/getMiniGameRoundStartTrend` |

表：`daily_mini_game_round_start_stats` / `weekly_*` / `monthly_*`。

| 同上页 **Tab「单局游戏结算」** | `mini_game_round_result` | 第三方 `vendor transfer` 回调钱包更新成功且 `transaction_id` 首次处理 +1 | `POST /trackingEvent/getMiniGameRoundResultTrend` |

表：`daily_mini_game_round_result_stats` / `weekly_*` / `monthly_*`。

| 同上页 **Tab「半屏游戏曝光」** | `mini_game_exposure` | 登录用户 `reportMiniGameExposure` 成功 +1（暂不要求 body 参数、不做房间/游戏校验） | `POST /trackingEvent/getMiniGameExposureTrend` |

App：`POST /liveRoom/reportMiniGameExposure`（Sofie `POST /sofie/studio/notifyMiniGameExposure`）。

表：`daily_mini_game_exposure_stats` / `weekly_*` / `monthly_*`。

| 同上页 **Tab「WebView 加载成功」** | `mini_game_webview_load_success` | 登录用户 `reportMiniGameWebViewLoadSuccess` 成功 +1（暂不要求 body） | `POST /trackingEvent/getMiniGameWebViewLoadTrend` |

App：`POST /liveRoom/reportMiniGameWebViewLoadSuccess`（Sofie `POST /sofie/studio/notifyMiniGameWebViewLoadSuccess`）。

表：`daily_mini_game_webview_load_stats` / `weekly_*` / `monthly_*`。

| 同上页 **Tab「点击游戏手柄」** | `mini_game_controller_click` | 登录用户 `reportMiniGameControllerClick` 成功 +1（暂不要求 body） | `POST /trackingEvent/getMiniGameControllerClickTrend` |

App：`POST /liveRoom/reportMiniGameControllerClick`（Sofie `POST /sofie/studio/notifyMiniGameControllerClick`）。

表：`daily_mini_game_controller_click_stats` / `weekly_*` / `monthly_*`。

| 同上页 **Tab「进入离线房间玩游戏」** | `offline_game_room_play` | 游戏房 category=2、未开播(`liveRecordId=0`)；`appGameStart` 传 `roomId` 且返回有效 link +1，不含主播 | `POST /trackingEvent/getOfflineGameRoomPlayTrend` |
| **游戏直播间页（2026-10）Tab「Spin点击」** | `mini_game_spin_click` | 登录用户 `reportMiniGameSpinClick` 成功 +1 | `POST /trackingEvent/getMiniGameSpinClickTrend` |
| 同上 **Tab「游戏+」** | `mini_game_plus_click` | 登录用户 `reportMiniGamePlusClick` 成功 +1 | `POST /trackingEvent/getMiniGamePlusClickTrend` |

App：`POST /liveRoom/reportMiniGameSpinClick`、`POST /liveRoom/reportMiniGamePlusClick`（Sofie `notifyMiniGameSpinClick` / `notifyMiniGamePlusClick`）。

表：`daily_mini_game_spin_click_stats` / `weekly_*` / `monthly_*`；`daily_mini_game_plus_click_stats` / `weekly_*` / `monthly_*`。

表：`daily_offline_game_room_play_stats` / `weekly_*` / `monthly_*`。

| 页面 | event | 口径 | CMS API |
|------|-------|------|---------|
| `TrackingShortVideoManagement`（短视频） | `short_video_play_start` | `watchShortVideoStart` 成功且视频已上架 +1 | `POST /trackingEvent/getShortVideoPlayStartTrend` |
| 同上 **Tab「付费解锁点击」** | `short_video_paywall_click` | `payShortVideo` 未购时每次确认解锁请求 +1（已购早退不计；扣款失败仍计点击） | `POST /trackingEvent/getShortVideoPaywallClickTrend` |
| 同上 **Tab「付费成功」** | `short_video_pay_success` | `payShortVideo` 首次扣款成功 +1 | `POST /trackingEvent/getShortVideoPaySuccessTrend` |
| 同上 **Tab「播放结束」** | `short_video_watch_end` | `watchShortVideoEnd` 登录成功 +1 | `POST /trackingEvent/getShortVideoWatchEndTrend` |

App 确认付费解锁：`POST /shortVideo/payShortVideo`（Sofie `POST /sofie/clip/unlockMedia`），无额外上报接口。

表：各事件 `daily_short_video_*_stats` / `weekly_*` / `monthly_*`。

| 页面 | event | 口径 | CMS API |
|------|-------|------|---------|
| `TrackingVoiceChatRoomManagement`（语聊房） | `voice_room_enter` | 语聊 category=5；观众 `joinRoom` 成功 +1，不含主播 | `POST /trackingEvent/getVoiceRoomEnterTrend` |
| 同上 **Tab「申请上麦」** | `voice_room_mic_request` | `applyVoiceChatMic` 成功 +1 | `POST /trackingEvent/getVoiceRoomMicRequestTrend` |
| 同上 **Tab「礼物面板点击」** | `voice_room_gift_click` | `reportVoiceRoomGiftClick` 登录成功 +1（暂不要求 body） | `POST /trackingEvent/getVoiceRoomGiftClickTrend` |

App 礼物面板：`POST /liveRoom/reportVoiceRoomGiftClick`（Sofie `POST /sofie/studio/notifyVoiceRoomGiftClick`）。

表：各事件 `daily_voice_room_*_stats` / `weekly_*` / `monthly_*`（进房表名为 `*_voice_room_enter_*`）。

| 页面 | event | 口径 | CMS API |
|------|-------|------|---------|
| `TrackingRechargeManagement`（充值） | `charge_panel_enter` | 登录用户 `reportChargePanelEnter` 成功 +1（暂不要求 body） | `POST /trackingEvent/getChargePanelEnterTrend` |
| 同上 **Tab「充值档位点击」** | `charge_click` | App/渠道/币商建单成功（`createRechargeOrder` 等）+1 | `POST /trackingEvent/getChargeClickTrend` |
| 同上 **Tab「充值到账成功」** | `charge_result_server` | 支付回调 `completeOrder` 首次发放金币成功 +1 | `POST /trackingEvent/getChargeResultServerTrend` |

App 面板曝光：`POST /rechargeOrder/reportChargePanelEnter`（Sofie `POST /sofie/topup/notifyChargePanelEnter`）。

表：各事件 `daily_charge_*_stats` / `weekly_*` / `monthly_*`（如 `daily_charge_panel_enter_stats`）。

## 新增一个事件的 checklist

1. Go：常量 `event_key` + 业务点 `Incr` +（如需）stat 表/DAO
2. Go：CMS controller + DTO trend 响应（daily/weekly/monthly points + total）
3. CMS：`tracking/<slug>.vue`（薄包装）+ i18n + api module 方法
4. CMS：`permission-menu-tree` + `permission-api-paths` + 路由 + 侧栏
5. 角色权限里勾选新页（或超管默认可见）

## 与已移除功能的关系

- 已去掉 **Firebase Analytics CMS 配置页** 及 Go `firebaseanalytics` 模块；后续埋点走 **本约定**，不恢复 GA4 配置 UI。
