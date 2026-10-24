import 'cache_key.dart';
import 'device_tier.dart';
import 'lru_cache.dart';

/// 进程内缓存仓库：按设备档位建 LRU。
///
/// **用户数据必须走 [putUser]/[getUser]/[invalidateUser]**，key 维度为 userId。
/// 全局（与登录无关）才用 [putGlobal]/[getGlobal]。
class MemoryCacheStore {
  MemoryCacheStore({
    required DeviceTier tier,
    int Function(String key, Object? value)? sizeOf,
  }) : _cache = LruCache<String, _TimedValue>(
          maxEntries: CacheCapacity.forTier(tier).maxEntries,
          maxBytes: CacheCapacity.forTier(tier).maxBytes,
          sizeOf: (key, timed) {
            if (sizeOf != null) {
              return sizeOf(key, timed.value);
            }
            return 1;
          },
        );

  final LruCache<String, _TimedValue> _cache;

  // --- 用户维度（主路径）---

  void putUser(int userId, String scope, Object? value, {Duration? ttl}) {
    _put(CacheKey.user(userId, scope), value, ttl: ttl);
  }

  T? getUser<T>(int userId, String scope) =>
      _get<T>(CacheKey.user(userId, scope));

  void invalidateUserScope(int userId, String scope) =>
      _cache.remove(CacheKey.user(userId, scope));

  /// 登出 / 踢人：清掉该 userId 下全部缓存。
  void invalidateUser(int userId) {
    final prefix = CacheKey.userPrefix(userId);
    _cache.removeWhere((key, _) => key.startsWith(prefix));
  }

  // --- 全局（非用户）---

  void putGlobal(String scope, Object? value, {Duration? ttl}) {
    _put(CacheKey.global(scope), value, ttl: ttl);
  }

  T? getGlobal<T>(String scope) => _get<T>(CacheKey.global(scope));

  void invalidateGlobal(String scope) =>
      _cache.remove(CacheKey.global(scope));

  void clear() => _cache.clear();

  int get length => _cache.length;

  void _put(String key, Object? value, {Duration? ttl}) {
    final expireAt = ttl == null ? null : DateTime.now().add(ttl);
    _cache.put(key, _TimedValue(value, expireAt));
  }

  T? _get<T>(String key) {
    final timed = _cache.get(key);
    if (timed == null) {
      return null;
    }
    if (timed.isExpired) {
      _cache.remove(key);
      return null;
    }
    final v = timed.value;
    if (v is T) {
      return v;
    }
    return null;
  }
}

class _TimedValue {
  _TimedValue(this.value, this.expireAt);

  final Object? value;
  final DateTime? expireAt;

  bool get isExpired =>
      expireAt != null && !DateTime.now().isBefore(expireAt!);
}
