import 'dart:async';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/cache/business_object_cache.dart';
import '../../core/http/auth_token.dart';
import '../../core/ws/push_cmd.dart';
import 'one_to_one_api.dart';
import 'one_to_one_models.dart';

/// 1v1 房间列表：L1 → L2 → HTTP；WS 重连后失效首页缓存。
class OneToOneRepository {
  OneToOneRepository({
    required OneToOneApi api,
    required BusinessObjectCache cache,
    required AppBus appBus,
    required AuthTokenProvider tokenProvider,
  })  : _api = api,
        _cache = cache,
        _appBus = appBus,
        _tokenProvider = tokenProvider {
    _appBus.add(_onBus);
  }

  final OneToOneApi _api;
  final BusinessObjectCache _cache;
  final AppBus _appBus;
  final AuthTokenProvider _tokenProvider;

  static String listScope(int page, int statusFilter) =>
      'oneToOneList_p${page}_f$statusFilter';

  void _onBus(AppBusEvent event) {
    if (event is! WsBusEvent) return;
    if (event.cmd != PushCmd.wsConnected) return;
    final userId = _currentUserId();
    if (userId == null) return;
    unawaited(invalidateListCache(userId));
  }

  int? _currentUserId() {
    final raw = _tokenProvider()?.userId;
    if (raw == null) return null;
    return int.tryParse(raw);
  }

  Future<void> invalidateListCache(int userId) async {
    for (final filter in [
      OneToOneStatusFilter.all,
      OneToOneStatusFilter.online,
      OneToOneStatusFilter.offline,
    ]) {
      await _cache.invalidateUserScope(
        userId,
        listScope(1, filter),
      );
    }
  }

  Future<OneToOneRoomListPage> getRoomList({
    required int userId,
    int page = 1,
    int pageSize = 20,
    int statusFilter = OneToOneStatusFilter.all,
    bool forceRefresh = false,
  }) async {
    final scope = listScope(page, statusFilter);
    if (!forceRefresh) {
      final cached = await _cache.getUser<OneToOneRoomListPage>(
        userId,
        scope,
        decode: OneToOneRoomListPage.fromJson,
      );
      if (cached != null) {
        return cached;
      }
    }
    final fresh = await _api.fetchRoomList(
      page: page,
      pageSize: pageSize,
      statusFilter: statusFilter,
    );
    await _cache.putUser(
      userId,
      scope,
      fresh,
      encode: (v) => v.toJson(),
    );
    return fresh;
  }

  void dispose() {
    _appBus.remove(_onBus);
  }
}
