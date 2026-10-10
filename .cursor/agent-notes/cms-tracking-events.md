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
| `TrackingHotLiveRoomJoinManagement`（CMS 文案：进入秀场直播间） | `hot_live_room_join` | 观众 `joinRoom` 成功且 `live_room_cfgs.category=1`(秀场/Hot)，不含主播本人；每次 +1；`JoinRoom` Pub `ShowcaseLiveRoomJoinTrackingEvent`，`module/tracking` 队列异步消费 | `POST /trackingEvent/getHotLiveRoomJoinTrend` |
| 同上页 **Tab「直播首帧画面渲染完成」** | `live_first_frame_rendered` | 秀场(category=1)、直播中；App `POST /liveRoom/reportLiveFirstFrameRendered`（Sofie `POST /sofie/studio/notifyLiveFirstFrameReady`）；不含主播；每次成功 +1 | `POST /trackingEvent/getLiveFirstFrameRenderTrend` |
| 同上页 **Tab「退出直播间」** | `hot_live_room_leave` | 秀场(category=1)；`leaveRoom` 成功 +1，不含主播；仅主动退房 API（非踢出/心跳清理） | `POST /trackingEvent/getHotLiveRoomLeaveTrend` |

表：`daily_hot_live_room_join_stats` / `weekly_*` / `monthly_*`（仅 `count`）；首帧：`daily_live_first_frame_render_stats` / `weekly_*` / `monthly_*`；退房：`daily_hot_live_room_leave_stats` / `weekly_*` / `monthly_*`。

| 页面 | event | 口径 | CMS API |
|------|-------|------|---------|
| `TrackingCall1v1InitiateManagement`（1v1视频通话） | `call_1v1_initiate` | 秀场 category=1；`liveRoomCall` 成功且 `call_order.source=1`（直播间来源）+1 | `POST /trackingEvent/getCall1v1InitiateTrend` |

表：`daily_call_1v1_initiate_stats` / `weekly_*` / `monthly_*`。

## 新增一个事件的 checklist

1. Go：常量 `event_key` + 业务点 `Incr` +（如需）stat 表/DAO
2. Go：CMS controller + DTO trend 响应（daily/weekly/monthly points + total）
3. CMS：`tracking/<slug>.vue`（薄包装）+ i18n + api module 方法
4. CMS：`permission-menu-tree` + `permission-api-paths` + 路由 + 侧栏
5. 角色权限里勾选新页（或超管默认可见）

## 与已移除功能的关系

- 已去掉 **Firebase Analytics CMS 配置页** 及 Go `firebaseanalytics` 模块；后续埋点走 **本约定**，不恢复 GA4 配置 UI。
