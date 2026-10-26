import 'dart:async';

import 'package:shared_preferences/shared_preferences.dart';

import '../app/app_config.dart';
import '../app/server_endpoint_store.dart';
import '../core/bus/app_bus.dart';
import '../core/bus/app_bus_event.dart';
import '../core/cache/api_disk_cache.dart';
import '../core/cache/business_object_cache.dart';
import '../core/cache/device_tier.dart';
import '../core/cache/device_tier_detect.dart';
import '../core/cache/image_disk_cache.dart';
import '../core/cache/memory_cache_store.dart';
import '../core/http/api_client.dart';
import '../core/http/api_exception.dart';
import '../core/http/auth_token.dart';
import '../core/http/http_logger.dart';
import '../core/ws/push_cmd.dart';
import '../core/ws/ws_client.dart';
import '../core/ws/ws_config.dart';
import '../data/agora/agora_api.dart';
import '../data/auth/auth_api.dart';
import '../data/auth/device_info_builder.dart';
import '../data/auth/session_store.dart';
import '../data/call/call_api.dart';
import '../data/gift/gift_api.dart';
import '../feature/show_live/live_room_viewer_controller.dart';
import '../feature/show_live/live_room_anchor_controller.dart';
import '../feature/show_live/anchor_live_publish_prefs.dart';
import '../feature/show_live/show_live_publisher.dart';
import '../data/live_follow/live_follow_api.dart';
import '../data/live_room/live_room_api.dart';
import '../data/live_room/live_room_repository.dart';
import '../data/one_to_one/one_to_one_api.dart';
import '../data/one_to_one/one_to_one_repository.dart';
import '../feature/call/call_controller.dart';
import '../data/user/user_api.dart';
import '../data/user/user_repository.dart';
import '../data/vip/vip_api.dart';
import '../data/vip/vip_repository.dart';

/// 进程级依赖（数据层），UI/ViewModel 从这里取。
///
/// 统一事件：[appBus] — `AppLifecycleBusEvent` / `WsBusEvent`；模块同一套 `add`/`remove`。
class AppScope {
  AppScope._({
    required this.tier,
    required this.session,
    required this.apiClient,
    required this.authApi,
    required this.deviceInfoBuilder,
    required this.userRepository,
    required this.vipRepository,
    required this.oneToOneRepository,
    required this.liveRoomRepository,
    required this.callController,
    required this.showLivePublisher,
    required this.liveRoomViewer,
    required this.liveRoomAnchor,
    required this.liveFollowApi,
    required this.giftApi,
    required this.cache,
    required this.appBus,
    required this.wsClient,
    required this.serverEndpoint,
  });

  final DeviceTier tier;
  final SessionStore session;
  final ApiClient apiClient;
  final AuthApi authApi;
  final DeviceInfoBuilder deviceInfoBuilder;
  final UserRepository userRepository;
  final VipRepository vipRepository;
  final OneToOneRepository oneToOneRepository;
  final LiveRoomRepository liveRoomRepository;
  final CallController callController;
  final ShowLivePublisher showLivePublisher;
  final LiveRoomViewerController liveRoomViewer;
  final LiveRoomAnchorController liveRoomAnchor;
  final LiveFollowApi liveFollowApi;
  final GiftApi giftApi;
  final BusinessObjectCache cache;
  final AppBus appBus;
  final WsClient wsClient;
  final ServerEndpointStore serverEndpoint;

  int _accountSwitchDepth = 0;

  /// 设备码切换进行中：勿因旧 token 的 `code=3` 清会话并回 Splash（会与 deviceLogin 并发）。
  bool get isAccountSwitchInProgress => _accountSwitchDepth > 0;

  /// 包裹账号切换流程，抑制 [ApiClient.onUnauthorized]。
  Future<T> runAccountSwitch<T>(Future<T> Function() action) async {
    _accountSwitchDepth++;
    try {
      return await action();
    } finally {
      _accountSwitchDepth--;
    }
  }

  static Future<AppScope> create() async {
    final prefs = await SharedPreferences.getInstance();
    final session = SessionStore(prefs);
    await session.load();
    final serverEndpoint = ServerEndpointStore(prefs);
    await serverEndpoint.load();

    await HttpLogger.init();

    final tier = await detectDeviceTier();
    ImageDiskCache.init(tier);
    await ApiDiskCache.init(tier);
    final memory = MemoryCacheStore(tier: tier);
    final cache = BusinessObjectCache(
      memory: memory,
      disk: ApiDiskCache.instance,
    );
    final appBus = AppBus();
    // 冷启动：先清 L1+L2，再发 started（避免 splash 读到旧磁盘缓存）
    cache.bindAppLifecycle(appBus);
    await cache.clearAll();

    late AppScope scope;
    final client = ApiClient(
      baseUrlProvider: () => serverEndpoint.apiBaseUrl,
      appPackage: AppConfig.appPackage,
      tokenProvider: () => scope.session.token,
      onUnauthorizedWithToken: (ApiException _, String? failedToken) {
        if (scope.isAccountSwitchInProgress ||
            scope.session.isDeviceCodeLoginPending) {
          return;
        }
        // Keep the token that failed. The callback is dispatched on a later
        // microtask, so a device-code switch may complete before it runs;
        // never log out the newly authenticated account because of an old
        // request from the previous account.
        if (failedToken == null) {
          return;
        }
        Future(() async {
          if (scope.isAccountSwitchInProgress ||
              scope.session.isDeviceCodeLoginPending) {
            return;
          }
          if (scope.session.token?.headerValue != failedToken) {
            return;
          }
          await scope.handleLoggedOut();
        });
      },
    );

    final wsClient = WsClient(
      config: WsConfig(wsBaseUrlProvider: () => serverEndpoint.wsBaseUrl),
      tokenProvider: () => scope.session.token,
      onMessages: (msg) => appBus.publish(WsBusEvent.fromPush(msg)),
      onState: (state) {
        if (state == WsConnectionState.connected) {
          appBus.publish(const WsBusEvent(cmd: PushCmd.wsConnected));
        }
      },
    );

    AuthToken? tokenProvider() => scope.session.token;

    final liveRoomApi = LiveRoomApi(client);
    final liveFollowApi = LiveFollowApi(client);
    final agoraApi = AgoraApi(client);
    final liveRoomRepository = LiveRoomRepository(
      api: liveRoomApi,
      cache: cache,
      appBus: appBus,
      tokenProvider: tokenProvider,
    );

    scope = AppScope._(
      tier: tier,
      session: session,
      apiClient: client,
      authApi: AuthApi(client),
      deviceInfoBuilder: DeviceInfoBuilder(),
      userRepository: UserRepository(
        api: UserApi(client),
        cache: cache,
        appBus: appBus,
        tokenProvider: tokenProvider,
      ),
      vipRepository: VipRepository(
        api: VipApi(client),
        cache: cache,
        appBus: appBus,
        tokenProvider: tokenProvider,
      ),
      oneToOneRepository: OneToOneRepository(
        api: OneToOneApi(client),
        cache: cache,
        appBus: appBus,
        tokenProvider: tokenProvider,
      ),
      liveRoomRepository: liveRoomRepository,
      callController: CallController(
        api: CallApi(client),
        appBus: appBus,
      ),
      showLivePublisher: ShowLivePublisher(
        agoraApi: agoraApi,
        liveRoomApi: liveRoomApi,
        publishPrefs: AnchorLivePublishPrefs(prefs),
      ),
      liveRoomViewer: LiveRoomViewerController(
        liveRoomApi: liveRoomApi,
        agoraApi: agoraApi,
        appBus: appBus,
      ),
      liveFollowApi: liveFollowApi,
      liveRoomAnchor: LiveRoomAnchorController(
        liveRoomApi: liveRoomApi,
        liveRoomRepository: liveRoomRepository,
        liveFollowApi: liveFollowApi,
        appBus: appBus,
      ),
      giftApi: GiftApi(client),
      cache: cache,
      appBus: appBus,
      wsClient: wsClient,
      serverEndpoint: serverEndpoint,
    );

    appBus.publish(AppLifecycleBusEvent.started());
    return scope;
  }

  /// 切换 API 根地址（内网调试）；会清会话，需重新登录。
  Future<void> applyServerEndpointOverride(String? normalizedOverride) async {
    await serverEndpoint.setOverride(normalizedOverride);
    // A different server can reuse user IDs with different data, so do not
    // carry any API object cache across an endpoint change.
    await handleLoggedOut(clearAllCache: true);
  }

  /// 登录成功后连接 WS。
  Future<void> handleLoggedIn() async {
    await wsClient.connect();
  }

  /// Start an explicit device-code account switch.
  ///
  /// The old token is cleared before the editor is shown, and only the cache
  /// namespace belonging to that user is invalidated. The device code itself
  /// is retained as the editor's initial value until a replacement login
  /// succeeds.
  Future<String> beginDeviceCodeLogin() async {
    final initialDeviceId =
        session.storedDeviceId ?? await session.getOrCreateDeviceId();
    final currentUserId = _currentUserId;

    await runAccountSwitch(() async {
      await session.beginDeviceCodeLogin();
      await _tearDownClientSession(userId: currentUserId);
      await session.clearToken();
    });
    return initialDeviceId;
  }

  /// Authenticate with a device code entered in the account-switch page.
  Future<DeviceLoginResult> loginWithDeviceCode(String deviceId) async {
    final normalized = deviceId.trim();
    if (normalized.isEmpty) {
      throw ArgumentError.value(deviceId, 'deviceId', 'must not be empty');
    }

    return runAccountSwitch(() async {
      // Keep this operation safe if the page was opened directly or restored
      // from a deep link rather than through [beginDeviceCodeLogin].
      if (!session.isDeviceCodeLoginPending) {
        final currentUserId = _currentUserId;
        await session.beginDeviceCodeLogin();
        await _tearDownClientSession(userId: currentUserId);
        if (session.isLoggedIn) {
          await session.clearToken();
        }
      }

      final deviceInfo = await deviceInfoBuilder.build(deviceId: normalized);
      final result = await authApi.deviceLogin(deviceInfo);
      await session.completeDeviceCodeLogin(
        deviceId: normalized,
        token: result.token,
      );
      await handleLoggedIn();
      return result;
    });
  }

  /// Leave the editor without attempting another account login.
  Future<void> cancelDeviceCodeLogin() async {
    await session.cancelDeviceCodeLogin();
  }

  /// Tear down active client sessions while retaining the current token.
  ///
  /// Kept for callers that need to perform their own immediate login after
  /// cleanup; the device-code page uses [beginDeviceCodeLogin] instead.
  Future<void> prepareForAccountSwitch() async {
    await _tearDownClientSession(userId: _currentUserId);
  }

  /// 登出 / 鉴权失效：断 WS、清会话与缓存。
  Future<void> handleLoggedOut({bool clearAllCache = false}) async {
    final currentUserId = _currentUserId;
    await _tearDownClientSession(
      userId: currentUserId,
      clearAllCache: clearAllCache,
    );
    try {
      await session.clearToken();
    } catch (_) {}
  }

  int? get _currentUserId {
    final id = int.tryParse(session.token?.userId ?? '');
    if (id == null || id <= 0) return null;
    return id;
  }

  Future<void> _tearDownClientSession({
    int? userId,
    bool clearAllCache = false,
  }) async {
    try {
      await wsClient.disconnect();
    } catch (_) {}
    try {
      await callController.hangUp();
    } catch (_) {}
    try {
      await showLivePublisher.stopPublish();
    } catch (_) {}
    liveRoomAnchor.exitSession();
    try {
      await liveRoomViewer.leaveRoom();
    } catch (_) {}
    try {
      if (!clearAllCache && userId != null && userId > 0) {
        await cache.invalidateUser(userId);
      } else {
        await cache.clearAll();
      }
    } catch (_) {}
  }

  void dispose() {
    userRepository.dispose();
    vipRepository.dispose();
    oneToOneRepository.dispose();
    liveRoomRepository.dispose();
    callController.dispose();
    showLivePublisher.dispose();
    liveRoomViewer.dispose();
    liveRoomAnchor.dispose();
    cache.disposeLifecycle();
    if (ApiDiskCache.isInitialized) {
      ApiDiskCache.instance.dispose();
    }
    unawaited(wsClient.disconnect());
    apiClient.close();
  }
}
