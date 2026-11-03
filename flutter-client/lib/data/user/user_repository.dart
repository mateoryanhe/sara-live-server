import 'dart:async';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/cache/business_object_cache.dart';
import '../../core/http/auth_token.dart';
import '../../core/ws/push_cmd.dart';
import '../../core/ws/ws_client.dart';
import '../call/call_models.dart';
import 'user_api.dart';
import 'user_profile.dart';
import 'wallet_balance_push.dart';

/// 用户资料：L1 内存 → L2 磁盘 → HTTP；订阅 VIP / 钱包余额推送更新 profile 缓存。
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

  Future<UserProfile?>? _profileWalletResyncFuture;

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

  /// 回前台 / WS 重连：失效 profile 缓存后强制 `/userInfo/get`，对齐服务端金币/钻石。
  ///
  /// 并发调用合并为同一 in-flight；Repository 与 UI 层可共用，避免重复请求。
  Future<UserProfile?> refreshWalletFromServer(int userId) {
    return _profileWalletResyncFuture ??=
        _refreshWalletFromServerJob(userId).whenComplete(() {
      _profileWalletResyncFuture = null;
    });
  }

  /// 同 [refreshWalletFromServer]（WS 上线补偿命名保留）。
  Future<UserProfile?> syncProfileAfterWsOnline(int userId) =>
      refreshWalletFromServer(userId);

  Future<UserProfile?> _refreshWalletFromServerJob(int userId) async {
    await invalidateProfile(userId);
    try {
      return await getProfile(userId: userId, forceRefresh: true);
    } catch (_) {
      return null;
    }
  }

  /// WS cmd=10/11：patch L1 内存中的 profile（同步，供同一次 publish 内其它订阅者读缓存）。
  UserProfile? tryPatchWalletBalanceInMemory({
    required int userId,
    required int cmd,
    required Object? data,
  }) {
    final push = WalletBalancePush.fromCmd(cmd, data);
    final current = _cache.memory.getUser<UserProfile>(userId, scopeProfile);
    if (current == null) return null;
    final UserProfile updated;
    switch (push) {
      case WalletDiamondPush(:final balance):
        updated = current.copyWith(diamond: balance);
      case WalletGoldPush(:final balance):
        updated = current.copyWith(gold: balance);
      case WalletBalancePushUnknown():
        return null;
    }
    _cache.memory.putUser(userId, scopeProfile, updated);
    unawaited(
      _cache.putUser(
        userId,
        scopeProfile,
        updated,
        encode: (v) => v.toJson(),
      ),
    );
    return updated;
  }

  /// WS cmd=10/11/40：用推送余额更新 profile 缓存（无缓存时先拉 `/userInfo/get`）。
  Future<UserProfile?> applyWalletBalancePush({
    required int userId,
    required int cmd,
    required Object? data,
  }) async {
    if (cmd == PushCmd.rechargeSuccess) {
      final map = pushDataMap(data);
      if (map == null) return null;
      return _patchProfileWallet(
        userId,
        gold: _parseBalance(map['goldBalance']),
      );
    }

    final push = WalletBalancePush.fromCmd(cmd, data);
    switch (push) {
      case WalletDiamondPush(:final balance):
        return _patchProfileWallet(userId, diamond: balance);
      case WalletGoldPush(:final balance):
        return _patchProfileWallet(userId, gold: balance);
      case WalletBalancePushUnknown():
        return null;
    }
  }

  Future<UserProfile?> _ensureProfileCached(int userId) async {
    final cached = await _cache.getUser<UserProfile>(
      userId,
      scopeProfile,
      decode: UserProfile.fromJson,
    );
    if (cached != null) {
      return cached;
    }
    try {
      final fresh = await _api.getUserInfo(userId: '$userId');
      final id = int.tryParse(fresh.userId) ?? userId;
      await _cache.putUser(
        id,
        scopeProfile,
        fresh,
        encode: (v) => v.toJson(),
      );
      return fresh;
    } catch (_) {
      return null;
    }
  }

  Future<UserProfile?> _patchProfileWallet(
    int userId, {
    double? gold,
    double? diamond,
  }) async {
    final current = await _ensureProfileCached(userId);
    if (current == null) return null;
    final updated = current.copyWith(
      gold: gold,
      diamond: diamond,
    );
    await _cache.putUser(
      userId,
      scopeProfile,
      updated,
      encode: (v) => v.toJson(),
    );
    return updated;
  }

  void _onBus(AppBusEvent event) {
    switch (event) {
      case AppLifecycleBusEvent(:final kind):
        if (kind == AppLifecycleKind.resumed) {
          final userId = _currentUserId();
          if (userId != null) {
            unawaited(refreshWalletFromServer(userId));
          }
        }
      case WsConnectionBusEvent(:final state):
        if (state != WsConnectionState.connected) {
          return;
        }
        final userId = _currentUserId();
        if (userId == null) {
          return;
        }
        unawaited(refreshWalletFromServer(userId));
      case WsBusEvent(:final cmd, :final data):
        _onWsPush(cmd, data);
      default:
        break;
    }
  }

  void _onWsPush(int cmd, Object? data) {
    final userId = _currentUserId();
    if (userId == null) return;
    switch (cmd) {
      case PushCmd.diamond:
      case PushCmd.gold:
        tryPatchWalletBalanceInMemory(
          userId: userId,
          cmd: cmd,
          data: data,
        );
        unawaited(
          applyWalletBalancePush(
            userId: userId,
            cmd: cmd,
            data: data,
          ),
        );
      case PushCmd.rechargeSuccess:
        unawaited(
          applyWalletBalancePush(
            userId: userId,
            cmd: cmd,
            data: data,
          ),
        );
      case PushCmd.vipLevel:
        unawaited(invalidateProfile(userId));
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

double _parseBalance(Object? v) {
  if (v is num) return v.toDouble();
  return double.tryParse('$v') ?? 0;
}
