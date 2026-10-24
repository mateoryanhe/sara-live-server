import 'dart:async';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/cache/business_object_cache.dart';
import '../../core/http/api_exception.dart';
import '../../core/http/auth_token.dart';
import '../../core/ws/push_cmd.dart';
import 'live_room_api.dart';
import 'live_room_models.dart';

class LiveRoomRepository {
  LiveRoomRepository({
    required LiveRoomApi api,
    required BusinessObjectCache cache,
    required AppBus appBus,
    required AuthTokenProvider tokenProvider,
  })  : _api = api,
        _cache = cache,
        _appBus = appBus,
        _tokenProvider = tokenProvider {
    _appBus.add(_onBus);
  }

  final LiveRoomApi _api;
  final BusinessObjectCache _cache;
  final AppBus _appBus;
  final AuthTokenProvider _tokenProvider;

  static String onlineListScope(int roomId, int page) =>
      'liveOnline_p${page}_r$roomId';

  static String oneToOneScope(int anchorId) => 'oneToOneCfg_$anchorId';
  static String liveRoomScope(int roomId) => 'liveRoom_$roomId';

  void _onBus(AppBusEvent event) {
    if (event is! WsBusEvent) return;
    if (event.cmd != PushCmd.wsConnected) return;
    final userId = _currentUserId();
    if (userId == null) return;
    unawaited(invalidateAnchorCaches(userId));
  }

  int? _currentUserId() {
    final raw = _tokenProvider()?.userId;
    if (raw == null) return null;
    return int.tryParse(raw);
  }

  Future<void> invalidateAnchorCaches(int anchorId) async {
    await _cache.invalidateUserScope(anchorId, oneToOneScope(anchorId));
    await _cache.invalidateUserScope(anchorId, liveRoomScope(anchorId));
    for (var p = 1; p <= 3; p++) {
      await _cache.invalidateUserScope(
        anchorId,
        onlineListScope(anchorId, p),
      );
    }
  }

  Future<OnlineUserListPage> getOnlineList({
    required int anchorId,
    int page = 1,
    int pageSize = 20,
    bool forceRefresh = false,
  }) async {
    final scope = onlineListScope(anchorId, page);
    if (!forceRefresh) {
      final cached = await _cache.getUser<OnlineUserListPage>(
        anchorId,
        scope,
        decode: OnlineUserListPage.fromJson,
      );
      if (cached != null) return cached;
    }
    final fresh = await _api.fetchOnlineList(
      roomId: anchorId,
      page: page,
      pageSize: pageSize,
    );
    await _cache.putUser(
      anchorId,
      scope,
      fresh,
      encode: (v) => v.toJson(),
    );
    return fresh;
  }

  Future<OneToOneRoomConfig?> getOneToOneConfig({
    required int anchorId,
    bool forceRefresh = false,
  }) async {
    final scope = oneToOneScope(anchorId);
    if (!forceRefresh) {
      final cached = await _cache.getUser<OneToOneRoomConfig>(
        anchorId,
        scope,
        decode: OneToOneRoomConfig.fromJson,
      );
      if (cached != null) return cached;
    }
    try {
      final fresh = await _api.fetchOneToOneRoom(anchorId: anchorId);
      await _cache.putUser(
        anchorId,
        scope,
        fresh,
        encode: (v) => v.toJson(),
      );
      return fresh;
    } on ApiException catch (e) {
      if (e.code == 168) return null;
      rethrow;
    }
  }

  Future<List<BillingTierItem>> getBillingTierList({
    required int anchorId,
    bool forceRefresh = false,
  }) async {
    const scope = 'billingTierList';
    if (!forceRefresh) {
      final cached = await _cache.getUser<List<BillingTierItem>>(
        anchorId,
        scope,
        decode: _decodeBillingTierList,
      );
      if (cached != null) return cached;
    }
    final fresh = await _api.fetchBillingTierList();
    await _cache.putUser(
      anchorId,
      scope,
      fresh,
      encode: _encodeBillingTierList,
    );
    return fresh;
  }

  Future<OneToOneRoomConfig> saveOneToOneConfig({
    required int anchorId,
    required double billing,
    required String title,
  }) async {
    await _api.createOneToOneRoom(
      billing: billing,
      title: title,
    );
    await _cache.invalidateUserScope(anchorId, oneToOneScope(anchorId));
    final cfg = await _api.fetchOneToOneRoom(anchorId: anchorId);
    await _cache.putUser(
      anchorId,
      oneToOneScope(anchorId),
      cfg,
      encode: (v) => v.toJson(),
    );
    return cfg;
  }

  Future<LiveRoomDetail?> getLiveRoom({
    required int anchorId,
    bool forceRefresh = false,
  }) async {
    final scope = liveRoomScope(anchorId);
    if (!forceRefresh) {
      final cached = await _cache.getUser<LiveRoomDetail>(
        anchorId,
        scope,
        decode: LiveRoomDetail.fromJson,
      );
      if (cached != null) return cached;
    }
    try {
      final fresh = await _api.fetchRoom(roomId: anchorId);
      await _cache.putUser(
        anchorId,
        scope,
        fresh,
        encode: (v) => v.toJson(),
      );
      return fresh;
    } on ApiException catch (e) {
      if (e.code == 52) return null;
      rethrow;
    }
  }

  Future<LiveRoomDetail> ensureShowLiveRoom({
    required int anchorId,
    required String title,
    required String notice,
    List<int>? coverBytes,
    String? coverFilename,
  }) async {
    await _api.createShowLiveRoom(
      title: title,
      notice: notice,
      coverBytes: coverBytes,
      coverFilename: coverFilename,
    );
    await _cache.invalidateUserScope(anchorId, liveRoomScope(anchorId));
    final room = await _api.fetchRoom(roomId: anchorId);
    await _cache.putUser(
      anchorId,
      liveRoomScope(anchorId),
      room,
      encode: (v) => v.toJson(),
    );
    return room;
  }

  /// 保存秀场资料（可多次调用）：标题走 create 合并，公告走 updateNotice（含清空）。
  Future<LiveRoomDetail> saveShowLiveProfile({
    required int anchorId,
    required String title,
    required String notice,
    List<int>? coverBytes,
    String? coverFilename,
  }) async {
    await _api.updateNotice(notice: notice);
    await _api.createShowLiveRoom(
      title: title,
      notice: notice,
      coverBytes: coverBytes,
      coverFilename: coverFilename,
    );
    await _cache.invalidateUserScope(anchorId, liveRoomScope(anchorId));
    final room = await _api.fetchRoom(roomId: anchorId);
    await _cache.putUser(
      anchorId,
      liveRoomScope(anchorId),
      room,
      encode: (v) => v.toJson(),
    );
    return room;
  }

  /// 开播前同步标题/公告/封面，再 startLive。
  Future<LiveRoomDetail> startLive({
    required int anchorId,
    required String title,
    required String notice,
    List<int>? coverBytes,
    String? coverFilename,
  }) async {
    await saveShowLiveProfile(
      anchorId: anchorId,
      title: title,
      notice: notice,
      coverBytes: coverBytes,
      coverFilename: coverFilename,
    );
    await _api.startLive();
    await _cache.invalidateUserScope(anchorId, liveRoomScope(anchorId));
    final room = await _api.fetchRoom(roomId: anchorId);
    await _cache.putUser(
      anchorId,
      liveRoomScope(anchorId),
      room,
      encode: (v) => v.toJson(),
    );
    return room;
  }

  Future<HotLiveRoomListPage> getHotLiveRoomList({
    int pageIndex = 1,
    int pageSize = 20,
  }) {
    return _api.fetchHotLiveRoomList(
      pageIndex: pageIndex,
      pageSize: pageSize,
    );
  }

  Future<void> kickAudience({
    required int anchorId,
    required int userId,
  }) async {
    await _api.kickAudience(userId: userId);
    await _invalidateOnlineListCache(anchorId);
  }

  Future<void> setAudienceMute({
    required int anchorId,
    required int userId,
    required bool muted,
  }) async {
    if (muted) {
      await _api.setAudienceMute(userId: userId, muted: true);
    } else {
      await _api.cancelAudienceMute(userId: userId);
    }
    await _invalidateOnlineListCache(anchorId);
  }

  Future<void> cancelKickBan({
    required int anchorId,
    required int userId,
  }) async {
    await _api.cancelKickBan(userId: userId);
    await _invalidateOnlineListCache(anchorId);
  }

  Future<AudienceRestrictStatus> getAudienceRestrictStatus({
    required int roomId,
    required int userId,
  }) {
    return _api.fetchAudienceRestrictStatus(roomId: roomId, userId: userId);
  }

  Future<void> _invalidateOnlineListCache(int anchorId) async {
    for (var p = 1; p <= 5; p++) {
      await _cache.invalidateUserScope(anchorId, onlineListScope(anchorId, p));
    }
  }

  Future<LiveRoomDetail> stopLive({required int anchorId}) async {
    await _api.stopLive();
    await _cache.invalidateUserScope(anchorId, liveRoomScope(anchorId));
    final room = await _api.fetchRoom(roomId: anchorId);
    await _cache.putUser(
      anchorId,
      liveRoomScope(anchorId),
      room,
      encode: (v) => v.toJson(),
    );
    return room;
  }

  static List<BillingTierItem> _decodeBillingTierList(Map<String, dynamic> json) {
    final list = json['list'];
    if (list is! List) return const [];
    return [
      for (final e in list)
        if (e is Map)
          BillingTierItem.fromJson(Map<String, dynamic>.from(e)),
    ];
  }

  static Map<String, dynamic> _encodeBillingTierList(List<BillingTierItem> list) {
    return {
      'list': [
        for (final i in list)
          {
            'id': i.id,
            'pricePerMinute': i.pricePerMinute,
            'sort': i.sort,
          },
      ],
    };
  }

  void dispose() {
    _appBus.remove(_onBus);
  }
}
