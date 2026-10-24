import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';

import 'api_disk_capacity.dart';
import 'cache_key.dart';
import 'device_tier.dart';

/// API 业务对象 **L2 磁盘** JSON 缓存（L1 为 [MemoryCacheStore]）。
class ApiDiskCache {
  ApiDiskCache._(this.tier, this._root, this._capacity);

  final DeviceTier tier;
  final Directory _root;
  final ApiDiskCapacity _capacity;

  static ApiDiskCache? _instance;

  static ApiDiskCache get instance {
    final cur = _instance;
    if (cur == null) {
      throw StateError('ApiDiskCache.init(tier) must be called before use');
    }
    return cur;
  }

  static bool get isInitialized => _instance != null;

  /// 启动时初始化目录；重复调用返回已有实例。
  /// [rootOverride] 仅测试注入（[Directory]）。
  static Future<ApiDiskCache> init(
    DeviceTier tier, {
    Object? rootOverride,
  }) async {
    final existing = _instance;
    if (existing != null) {
      return existing;
    }
    final capacity = ApiDiskCapacity.forTier(tier);
    final overrideDir = rootOverride as Directory?;
    final root = overrideDir ??
        Directory(
          p.join(
            (await getApplicationSupportDirectory()).path,
            'api_object_cache',
          ),
        );
    await root.create(recursive: true);
    return _instance = ApiDiskCache._(tier, root, capacity);
  }

  Future<Map<String, dynamic>?> readPayload({
    required int userId,
    required String scope,
  }) =>
      _read(CacheKey.user(userId, scope));

  Future<Map<String, dynamic>?> readPayloadGlobal({required String scope}) =>
      _read(CacheKey.global(scope));

  Future<void> writePayload({
    required int userId,
    required String scope,
    required Map<String, dynamic> payload,
    Duration? ttl,
  }) =>
      _write(CacheKey.user(userId, scope), payload, ttl: ttl);

  Future<void> writePayloadGlobal({
    required String scope,
    required Map<String, dynamic> payload,
    Duration? ttl,
  }) =>
      _write(CacheKey.global(scope), payload, ttl: ttl);

  Future<void> removeUserScope(int userId, String scope) =>
      _removeFile(_fileForKey(CacheKey.user(userId, scope)));

  Future<void> removeUser(int userId) async {
    final prefix = _safeFilePrefix(CacheKey.userPrefix(userId));
    if (!_root.existsSync()) return;
    await for (final entity in _root.list()) {
      if (entity is! File) continue;
      final name = p.basename(entity.path);
      if (name.startsWith(prefix)) {
        await _removeFile(entity);
      }
    }
  }

  Future<void> clearAll() async {
    if (!_root.existsSync()) return;
    await for (final entity in _root.list()) {
      if (entity is File && entity.path.endsWith('.json')) {
        await _removeFile(entity);
      }
    }
  }

  Future<void> evictExpired() async {
    if (!_root.existsSync()) return;
    await for (final entity in _root.list()) {
      if (entity is! File || !entity.path.endsWith('.json')) continue;
      final envelope = await _readEnvelopeFile(entity);
      if (envelope == null || _isExpired(envelope)) {
        await _removeFile(entity);
      }
    }
  }

  void dispose() {
    if (identical(_instance, this)) {
      _instance = null;
    }
  }

  Future<Map<String, dynamic>?> _read(String key) async {
    final file = _fileForKey(key);
    final envelope = await _readEnvelopeFile(file);
    if (envelope == null) {
      return null;
    }
    if (_isExpired(envelope)) {
      await _removeFile(file);
      return null;
    }
    final payload = envelope['payload'];
    if (payload is! Map) {
      await _removeFile(file);
      return null;
    }
    return Map<String, dynamic>.from(payload);
  }

  Future<void> _write(
    String key,
    Map<String, dynamic> payload, {
    Duration? ttl,
  }) async {
    final effectiveTtl = ttl ?? _capacity.defaultTtl;
    final expireAt = DateTime.now().add(effectiveTtl);
    final envelope = <String, dynamic>{
      'v': 1,
      'expireAt': expireAt.toUtc().toIso8601String(),
      'payload': payload,
    };
    final file = _fileForKey(key);
    await file.parent.create(recursive: true);
    await file.writeAsString(jsonEncode(envelope), flush: true);
    await _enforceLimits();
  }

  Future<Map<String, dynamic>?> _readEnvelopeFile(File file) async {
    if (!file.existsSync()) return null;
    try {
      final text = await file.readAsString();
      final decoded = jsonDecode(text);
      if (decoded is! Map) return null;
      return Map<String, dynamic>.from(decoded);
    } catch (_) {
      await _removeFile(file);
      return null;
    }
  }

  bool _isExpired(Map<String, dynamic> envelope) {
    final raw = envelope['expireAt'];
    if (raw is! String || raw.isEmpty) return false;
    final expireAt = DateTime.tryParse(raw);
    if (expireAt == null) return false;
    return !DateTime.now().toUtc().isBefore(expireAt.toUtc());
  }

  File _fileForKey(String key) =>
      File(p.join(_root.path, '${_safeFilePrefix(key)}.json'));

  String _safeFilePrefix(String key) =>
      key.replaceAll(':', '__').replaceAll('/', '_');

  Future<void> _removeFile(File file) async {
    if (file.existsSync()) {
      await file.delete();
    }
  }

  Future<void> _enforceLimits() async {
    final files = <File>[];
    var totalBytes = 0;
    if (!_root.existsSync()) return;
    await for (final entity in _root.list()) {
      if (entity is! File || !entity.path.endsWith('.json')) continue;
      final len = await entity.length();
      files.add(entity);
      totalBytes += len;
    }
    files.sort(
      (a, b) => a.statSync().modified.compareTo(b.statSync().modified),
    );
    while (files.length > _capacity.maxFiles ||
        totalBytes > _capacity.maxBytes) {
      if (files.isEmpty) break;
      final oldest = files.removeAt(0);
      totalBytes -= await oldest.length();
      await _removeFile(oldest);
    }
  }
}
