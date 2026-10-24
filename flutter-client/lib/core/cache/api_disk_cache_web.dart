import 'device_tier.dart';

/// Web 端：**无 L2 磁盘**（`path_provider` / `dart:io` 不可用），读恒 miss、写 no-op；L1 内存仍生效。
class ApiDiskCache {
  ApiDiskCache._(this.tier);

  final DeviceTier tier;

  static ApiDiskCache? _instance;

  static ApiDiskCache get instance {
    final cur = _instance;
    if (cur == null) {
      throw StateError('ApiDiskCache.init(tier) must be called before use');
    }
    return cur;
  }

  static bool get isInitialized => _instance != null;

  static Future<ApiDiskCache> init(
    DeviceTier tier, {
    Object? rootOverride,
  }) async {
    final existing = _instance;
    if (existing != null) {
      return existing;
    }
    return _instance = ApiDiskCache._(tier);
  }

  Future<Map<String, dynamic>?> readPayload({
    required int userId,
    required String scope,
  }) async =>
      null;

  Future<Map<String, dynamic>?> readPayloadGlobal({required String scope}) async =>
      null;

  Future<void> writePayload({
    required int userId,
    required String scope,
    required Map<String, dynamic> payload,
    Duration? ttl,
  }) async {}

  Future<void> writePayloadGlobal({
    required String scope,
    required Map<String, dynamic> payload,
    Duration? ttl,
  }) async {}

  Future<void> removeUserScope(int userId, String scope) async {}

  Future<void> removeUser(int userId) async {}

  Future<void> clearAll() async {}

  Future<void> evictExpired() async {}

  void dispose() {
    if (identical(_instance, this)) {
      _instance = null;
    }
  }
}
