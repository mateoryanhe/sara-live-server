import 'dart:async';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/cache/business_object_cache.dart';
import '../../core/http/auth_token.dart';
import '../../core/ws/push_cmd.dart';
import 'user_api.dart';
import 'user_profile.dart';

/// 用户资料：L1 内存 → L2 磁盘 → HTTP；订阅 VIP 相关推送失效 profile。
class UserRepository {
  UserRepository({
    required UserApi api,
    required BusinessObjectCache cache,
    required AppBus appBus,
    required AuthTokenProvider tokenProvider,
  })  : _api = api,
        _cache = cache,
        _appBus = appBus,
        _tokenProvider = tokenProvider {
    _appBus.add(_onBus);
  }

  static const scopeProfile = 'profile';

  final UserApi _api;
  final BusinessObjectCache _cache;
  final AppBus _appBus;
  final AuthTokenProvider _tokenProvider;

  Future<UserProfile> getProfile({
    required int userId,
    bool forceRefresh = false,
  }) async {
    if (!forceRefresh) {
      final cached = await _cache.getUser<UserProfile>(
        userId,
        scopeProfile,
        decode: UserProfile.fromJson,
      );
      if (cached != null) {
        return cached;
      }
    }
    final fresh = await _api.getUserInfo(userId: '$userId');
    final id = int.tryParse(fresh.userId) ?? userId;
    await _cache.putUser(
      id,
      scopeProfile,
      fresh,
      encode: (v) => v.toJson(),
    );
    return fresh;
  }

  Future<UserProfile> updateNickname({
    required int userId,
    required String nickname,
  }) async {
    final next = await _api.updateNickname(nickname);
    final current = await _cache.getUser<UserProfile>(
          userId,
          scopeProfile,
          decode: UserProfile.fromJson,
        ) ??
        await _api.getUserInfo(userId: '$userId');
    final updated = current.copyWith(nickname: next);
    await _cache.putUser(
      userId,
      scopeProfile,
      updated,
      encode: (v) => v.toJson(),
    );
    return updated;
  }

  Future<UserProfile> updateAvatar({
    required int userId,
    required List<int> bytes,
    required String filename,
  }) async {
    final avatarUrl = await _api.uploadAvatar(bytes: bytes, filename: filename);
    final current = await _cache.getUser<UserProfile>(
          userId,
          scopeProfile,
          decode: UserProfile.fromJson,
        ) ??
        await _api.getUserInfo(userId: '$userId');
    final updated = current.copyWith(avatar: avatarUrl);
    await _cache.putUser(
      userId,
      scopeProfile,
      updated,
      encode: (v) => v.toJson(),
    );
    return updated;
  }

  Future<void> invalidateProfile(int userId) async {
    await _cache.invalidateUserScope(userId, scopeProfile);
  }

  void _onBus(AppBusEvent event) {
    if (event is! WsBusEvent) return;
    final userId = _currentUserId();
    if (userId == null) return;
    switch (event.cmd) {
      case PushCmd.vipLevel:
      case PushCmd.wsConnected:
        unawaited(invalidateProfile(userId));
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
