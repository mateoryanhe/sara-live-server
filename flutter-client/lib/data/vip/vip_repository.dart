import 'dart:async';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/cache/business_object_cache.dart';
import '../../core/http/auth_token.dart';
import '../../core/ws/push_cmd.dart';
import '../../core/ws/ws_client.dart';
import 'vip_api.dart';
import 'vip_level_push.dart';
import 'vip_models.dart';

/// VIP：L1 → L2 → HTTP；**自行订阅**推送做缓存失效。
class VipRepository {
  VipRepository({
    required VipApi api,
    required BusinessObjectCache cache,
    required AppBus appBus,
    required AuthTokenProvider tokenProvider,
  })  : _api = api,
        _cache = cache,
        _appBus = appBus,
        _tokenProvider = tokenProvider {
    _appBus.add(_onBus);
  }

  static const scopeDetail = 'vipDetail';
  static const scopeCfgList = 'vipCfgList';

  final VipApi _api;
  final BusinessObjectCache _cache;
  final AppBus _appBus;
  final AuthTokenProvider _tokenProvider;

  Future<VipDetail> getVipDetail({
    required int userId,
    bool forceRefresh = false,
  }) async {
    if (!forceRefresh) {
      final cached = await _cache.getUser<VipDetail>(
        userId,
        scopeDetail,
        decode: VipDetail.fromJson,
      );
      if (cached != null) return cached;
    }
    final fresh = await _api.getVipDetail();
    await _cache.putUser(
      userId,
      scopeDetail,
      fresh,
      encode: (v) => v.toJson(),
    );
    return fresh;
  }

  Future<List<VipCfgItem>> getVipCfgList({
    required int userId,
    bool forceRefresh = false,
  }) async {
    if (!forceRefresh) {
      final cached = await _cache.getUser<List<VipCfgItem>>(
        userId,
        scopeCfgList,
        decode: _decodeVipCfgList,
      );
      if (cached != null) return cached;
    }
    final fresh = await _api.vipCfgListForApp();
    await _cache.putUser(
      userId,
      scopeCfgList,
      fresh,
      encode: _encodeVipCfgList,
    );
    return fresh;
  }

  Future<void> invalidate(int userId) async {
    await _cache.invalidateUserScope(userId, scopeDetail);
    await _cache.invalidateUserScope(userId, scopeCfgList);
  }

  void _onBus(AppBusEvent event) {
    final userId = _currentUserId();
    if (userId == null) return;

    switch (event) {
      case WsConnectionBusEvent(:final state):
        if (state == WsConnectionState.connected) {
          unawaited(invalidate(userId));
        }
      case WsBusEvent(:final cmd, :final data):
        if (cmd == PushCmd.vipLevel) {
          VipLevelPush.fromData(data);
          unawaited(invalidate(userId));
        }
      default:
        break;
    }
  }

  int? _currentUserId() {
    final id = int.tryParse(_tokenProvider()?.userId ?? '');
    if (id == null || id <= 0) return null;
    return id;
  }

  void dispose() {
    _appBus.remove(_onBus);
  }
}

List<VipCfgItem> _decodeVipCfgList(Map<String, dynamic> json) {
  final raw = json['items'];
  if (raw is! List) return [];
  return [
    for (final item in raw)
      if (item is Map)
        VipCfgItem.fromJson(Map<String, dynamic>.from(item)),
  ];
}

Map<String, dynamic> _encodeVipCfgList(List<VipCfgItem> list) => {
      'items': list.map((e) => e.toJson()).toList(),
    };
