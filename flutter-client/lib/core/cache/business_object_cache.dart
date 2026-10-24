import 'dart:async';

import '../bus/app_bus.dart';
import '../bus/app_bus_event.dart';
import 'api_disk_cache.dart';
import 'memory_cache_store.dart';

/// 业务对象 **L1 内存 + L2 磁盘** 缓存门面；Repository 统一经此读写。
class BusinessObjectCache {
  BusinessObjectCache({
    required this.memory,
    required ApiDiskCache disk,
  }) : _disk = disk;

  final MemoryCacheStore memory;
  final ApiDiskCache _disk;
  AppBus? _lifecycleBus;
  AppBusHandler? _lifecycleHandler;

  /// 订阅 [AppLifecycleKind.started]：冷启动时 [clearAll]，避免沿用上次的 L2 与线上一致性冲突。
  /// 须在 publish `started` 之前注册；[AppScope.create] 还会在 publish 前 await [clearAll]。
  void bindAppLifecycle(AppBus bus) {
    if (_lifecycleHandler != null) {
      return;
    }
    _lifecycleBus = bus;
    _lifecycleHandler = (AppBusEvent event) {
      if (event is! AppLifecycleBusEvent ||
          event.kind != AppLifecycleKind.started) {
        return;
      }
      unawaited(clearAll());
    };
    bus.add(_lifecycleHandler!);
  }

  void disposeLifecycle() {
    if (_lifecycleHandler != null && _lifecycleBus != null) {
      _lifecycleBus!.remove(_lifecycleHandler!);
    }
    _lifecycleHandler = null;
    _lifecycleBus = null;
  }

  Future<T?> getUser<T>(
    int userId,
    String scope, {
    required T Function(Map<String, dynamic> json) decode,
  }) async {
    final l1 = memory.getUser<T>(userId, scope);
    if (l1 != null) {
      return l1;
    }
    final raw = await _disk.readPayload(userId: userId, scope: scope);
    if (raw == null) {
      return null;
    }
    final decoded = decode(raw);
    memory.putUser(userId, scope, decoded);
    return decoded;
  }

  Future<T?> getGlobal<T>(
    String scope, {
    required T Function(Map<String, dynamic> json) decode,
  }) async {
    final l1 = memory.getGlobal<T>(scope);
    if (l1 != null) {
      return l1;
    }
    final raw = await _disk.readPayloadGlobal(scope: scope);
    if (raw == null) {
      return null;
    }
    final decoded = decode(raw);
    memory.putGlobal(scope, decoded);
    return decoded;
  }

  Future<void> putUser<T>(
    int userId,
    String scope,
    T value, {
    required Map<String, dynamic> Function(T value) encode,
    Duration? ttl,
  }) async {
    memory.putUser(userId, scope, value, ttl: ttl);
    await _disk.writePayload(
      userId: userId,
      scope: scope,
      payload: encode(value),
      ttl: ttl,
    );
  }

  Future<void> putGlobal<T>(
    String scope,
    T value, {
    required Map<String, dynamic> Function(T value) encode,
    Duration? ttl,
  }) async {
    memory.putGlobal(scope, value, ttl: ttl);
    await _disk.writePayloadGlobal(
      scope: scope,
      payload: encode(value),
      ttl: ttl,
    );
  }

  Future<void> invalidateUserScope(int userId, String scope) async {
    memory.invalidateUserScope(userId, scope);
    await _disk.removeUserScope(userId, scope);
  }

  Future<void> invalidateUser(int userId) async {
    memory.invalidateUser(userId);
    await _disk.removeUser(userId);
  }

  void clearMemory() => memory.clear();

  Future<void> clearAll() async {
    memory.clear();
    await _disk.clearAll();
  }
}
