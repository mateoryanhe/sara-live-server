import 'package:flutter_test/flutter_test.dart';
import 'package:sara_live_app/core/cache/cache_key.dart';
import 'package:sara_live_app/core/cache/device_tier.dart';
import 'package:sara_live_app/core/cache/lru_cache.dart';
import 'package:sara_live_app/core/cache/memory_cache_store.dart';

void main() {
  test('lru evicts least recently used', () {
    final cache = LruCache<String, int>(maxEntries: 2);
    cache.put('a', 1);
    cache.put('b', 2);
    expect(cache.get('a'), 1); // a becomes most recent
    cache.put('c', 3); // evict b
    expect(cache.containsKey('b'), isFalse);
    expect(cache.get('a'), 1);
    expect(cache.get('c'), 3);
  });

  test('user cache is scoped by userId', () {
    final store = MemoryCacheStore(tier: DeviceTier.mid);
    store.putUser(1, 'profile', {'name': 'a'});
    store.putUser(2, 'profile', {'name': 'b'});

    expect(CacheKey.user(1, 'profile'), 'user:1:profile');
    expect(store.getUser<Map>(1, 'profile')?['name'], 'a');
    expect(store.getUser<Map>(2, 'profile')?['name'], 'b');

    store.invalidateUser(1);
    expect(store.getUser<Map>(1, 'profile'), isNull);
    expect(store.getUser<Map>(2, 'profile')?['name'], 'b');
  });
}
