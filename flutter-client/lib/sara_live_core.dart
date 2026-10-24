/// Sara Live 客户端 core 层导出。
///
/// # 原则（先定死，再扩业务）
///
/// ## A. 分层：UI / ViewModel / 数据
///
/// 1. **UI 层：纯粹 UI**
///    - 只负责展示与手势转发；不写业务分支、不直接打 HTTP/WS、不碰 cache 策略
///    - **向 ViewModel 提供「操作 UI」的方法**（接口/回调），如 `showLoading` / `renderXxx` / `showError`
///    - ViewModel 只依赖该接口，不依赖具体 Widget 页面类
/// 2. **ViewModel 层：编排逻辑**
///    - 调用 **UI 接口** + **数据层**（Repository / cache / http / ws），完成业务
///    - 读：触发 load；写：调 API 后更新/失效缓存并刷新 UI
///    - **可订阅 AppBus**：收到推送 → 数据层 invalidate/patch → 再调 UI 刷新（或走一遍 cache→HTTP）
///    - 进页/登录订阅，离页/登出取消；断线重连后对关键数据做 HTTP 补偿
/// 3. **数据层（本库 `core/`）**：HTTP / WS / Cache / Bus 分目录、单向依赖
///    - `http`：请求与信封解析，不碰业务缓存策略
///    - `ws`：连接、心跳、帧解析；经 AppBus 发布，不直接知道业务 key 语义
///    - `cache`：业务 **L1 内存 LRU + L2 磁盘 JSON**（[BusinessObjectCache]；Web 仅 L1）+ 网络图片磁盘，均按档位限容量
///    - `bus`：统一事件总线 [AppBus]（启动 / WS 用类型区分）
///
/// ## B. 数据与协议
///
/// 4. **读请求：L1 内存 → L2 磁盘 → HTTP**（Repository 经 [BusinessObjectCache]）
///    - hit 直接返回；miss 拉网后 **同时写 L1 + L2**
///    - 强制刷新用 `forceRefresh` / `invalidate` 后再拉；写成功后更新或失效对应 key
/// 5. **缓存不是真相源**：写操作、登录态、余额/下单等不以「仅缓存」为准；可先展示缓存，以服务器结果最终覆盖
/// 6. **统一事件 [AppBus]**：`add` / `remove` / `publish`；**Repository 与 ViewModel 各自订阅**
///    - [AppLifecycleBusEvent]：`started`（[AppScope.create]，sticky；[BusinessObjectCache.bindAppLifecycle] 触发 `clearAll`）、`paused` / `resumed`（根 Widget）
///    - [WsBusEvent]：`WsClient` 下行 + 本机 `PushCmd.wsConnected`；订 `cmd=12` 等失效/刷 UI
///    - 登录 `handleLoggedIn` 连 `/ws`；登出断开
/// 7. **用户维度（key 以 userId 为命名空间）**
///    - 用户数据经 [BusinessObjectCache]（内部 `CacheKey.user` + L1/L2）
///    - 格式 `user:{userId}:{scope}`；登出/踢人调用 `invalidateUser(userId)`
/// 8. **对齐现有 Go App 协议**（见各文件注释），Flutter 不传 `X-H5-Client`
/// 9. **当前 App 登录**：仅设备码 `POST /auth/deviceLogin`；启动页自动登录后进主壳，「我的」拉 `/userInfo/get` + `/vip/getVipDetail`
/// 10. **网络图片走 core 磁盘缓存**：`ImageDiskCache`（按档位限对象数）+ `CachedAvatar`，勿直接 `NetworkImage`
/// 11. **路由用 go_router**：路径见 `AppRoutes`；登录成功 `context.go(main)`，退出 `context.go(splash)`；VIP `AppRoutes.vip`；主壳 [MainShellPage] 按 `UserProfile.isAnchor` 切换观众/主播底部 4 Tab（观众：1v1/直播/游戏/我的，默认 1v1）
/// 12. **1v1 视频通话**：[CallController] 编排 `POST /call/oneToOneRoomCall` → WS(25–32) → `acceptCall` / `confirmCall` / `callHeart` / `endCall`；RTC：**Web**=`web/agora_rtc_bridge.js`+Agora Web SDK，**Android**=`agora_rtc_engine`+`permission_handler`
/// 13. **秀场开播**：[ShowLivePublisher] `startLive` 后 `POST /agora/liveRoomToken`(Publisher) → 同套 RTC 推流；下播先 leave 再 `stopLive`
/// 14. **UI 皮肤可切换**：`UiStyle`，「我的 → 设置 → 界面风格」；无顶栏换肤/刷新，下拉刷新
/// 15. **VIP**：横滑等级页；推送与缓存/VM 解耦订阅，无中心 `vipRevision`
library;

export 'core/bus/app_bus.dart';
export 'core/bus/app_bus_event.dart';
export 'core/cache/api_disk_cache.dart';
export 'core/cache/api_disk_capacity.dart';
export 'core/cache/business_object_cache.dart';
export 'core/cache/cache_key.dart';
export 'core/cache/device_tier.dart';
export 'core/cache/device_tier_detect.dart';
export 'core/cache/image_disk_cache.dart';
export 'core/cache/lru_cache.dart';
export 'core/cache/memory_cache_store.dart';
export 'core/widgets/cached_avatar.dart';
export 'core/http/api_client.dart';
export 'core/http/api_exception.dart';
export 'core/http/api_response.dart';
export 'core/http/auth_token.dart';
export 'core/ws/push_cmd.dart';
export 'core/ws/push_message.dart';
export 'core/ws/ws_client.dart';
export 'core/ws/ws_config.dart';
