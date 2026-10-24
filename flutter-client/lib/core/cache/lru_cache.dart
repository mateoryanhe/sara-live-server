import 'dart:collection';

/// 泛型 LRU：按访问顺序淘汰；可选字节预算（由调用方提供 sizeOf）。
class LruCache<K, V> {
  LruCache({
    required this.maxEntries,
    this.maxBytes = 0,
    int Function(K key, V value)? sizeOf,
  })  : assert(maxEntries > 0),
        _sizeOf = sizeOf ?? ((_, __) => 1);

  final int maxEntries;
  final int maxBytes;
  final int Function(K key, V value) _sizeOf;

  final LinkedHashMap<K, _Entry<V>> _map = LinkedHashMap<K, _Entry<V>>();
  int _bytes = 0;

  int get length => _map.length;

  int get currentBytes => _bytes;

  bool containsKey(K key) => _map.containsKey(key);

  V? get(K key) {
    final entry = _map.remove(key);
    if (entry == null) {
      return null;
    }
    // remove + put 移到 LinkedHashMap 尾部 = 最近使用
    _map[key] = entry;
    return entry.value;
  }

  V? operator [](K key) => get(key);

  void put(K key, V value) {
    final nextSize = _sizeOf(key, value);
    final old = _map.remove(key);
    if (old != null) {
      _bytes -= old.size;
    }
    _map[key] = _Entry(value, nextSize);
    _bytes += nextSize;
    _evictIfNeeded();
  }

  void operator []=(K key, V value) => put(key, value);

  V? remove(K key) {
    final old = _map.remove(key);
    if (old == null) {
      return null;
    }
    _bytes -= old.size;
    return old.value;
  }

  void removeWhere(bool Function(K key, V value) test) {
    final keys = <K>[];
    _map.forEach((k, e) {
      if (test(k, e.value)) {
        keys.add(k);
      }
    });
    for (final k in keys) {
      remove(k);
    }
  }

  void clear() {
    _map.clear();
    _bytes = 0;
  }

  void _evictIfNeeded() {
    while (_map.isNotEmpty &&
        (_map.length > maxEntries || (maxBytes > 0 && _bytes > maxBytes))) {
      final firstKey = _map.keys.first;
      final removed = _map.remove(firstKey)!;
      _bytes -= removed.size;
    }
  }
}

class _Entry<V> {
  _Entry(this.value, this.size);

  final V value;
  final int size;
}
