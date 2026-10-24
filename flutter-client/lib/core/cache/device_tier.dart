/// 设备性能档位，用于约束内存 LRU 与图片磁盘缓存容量。
enum DeviceTier {
  low,
  mid,
  high,
}

/// 按物理内存（MB）粗分档；启动时算一次即可。
DeviceTier deviceTierFromTotalMemoryMb(int totalMemoryMb) {
  if (totalMemoryMb <= 0) {
    return DeviceTier.mid;
  }
  if (totalMemoryMb < 3072) {
    return DeviceTier.low;
  }
  if (totalMemoryMb < 6144) {
    return DeviceTier.mid;
  }
  return DeviceTier.high;
}

/// 各档 LRU 默认上限（条目数 + 估算字节预算）。
class CacheCapacity {
  const CacheCapacity({
    required this.maxEntries,
    required this.maxBytes,
  });

  final int maxEntries;
  final int maxBytes;

  static CacheCapacity forTier(DeviceTier tier) {
    switch (tier) {
      case DeviceTier.low:
        return const CacheCapacity(maxEntries: 80, maxBytes: 8 << 20);
      case DeviceTier.mid:
        return const CacheCapacity(maxEntries: 200, maxBytes: 24 << 20);
      case DeviceTier.high:
        return const CacheCapacity(maxEntries: 500, maxBytes: 48 << 20);
    }
  }
}
