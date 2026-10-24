import 'package:flutter_cache_manager/flutter_cache_manager.dart';

import 'device_tier.dart';

/// 网络图片磁盘缓存容量（按设备档位）。
class ImageDiskCapacity {
  const ImageDiskCapacity({
    required this.maxObjects,
    required this.stalePeriod,
  });

  final int maxObjects;
  final Duration stalePeriod;

  static ImageDiskCapacity forTier(DeviceTier tier) {
    switch (tier) {
      case DeviceTier.low:
        return const ImageDiskCapacity(
          maxObjects: 80,
          stalePeriod: Duration(days: 7),
        );
      case DeviceTier.mid:
        return const ImageDiskCapacity(
          maxObjects: 250,
          stalePeriod: Duration(days: 14),
        );
      case DeviceTier.high:
        return const ImageDiskCapacity(
          maxObjects: 600,
          stalePeriod: Duration(days: 30),
        );
    }
  }
}

/// 进程内单例图片磁盘缓存（flutter_cache_manager），供 CachedNetworkImage 使用。
class ImageDiskCache {
  ImageDiskCache._(this.tier, this.manager);

  final DeviceTier tier;
  final CacheManager manager;

  static ImageDiskCache? _instance;

  static ImageDiskCache get instance {
    final cur = _instance;
    if (cur == null) {
      throw StateError('ImageDiskCache.init(tier) must be called at startup');
    }
    return cur;
  }

  static bool get isInitialized => _instance != null;

  /// 启动时按档位初始化；重复调用忽略。
  static ImageDiskCache init(DeviceTier tier) {
    final existing = _instance;
    if (existing != null) {
      return existing;
    }
    final cap = ImageDiskCapacity.forTier(tier);
    final manager = CacheManager(
      Config(
        'sara_live_images_${tier.name}',
        stalePeriod: cap.stalePeriod,
        maxNrOfCacheObjects: cap.maxObjects,
      ),
    );
    return _instance = ImageDiskCache._(tier, manager);
  }

  Future<void> emptyCache() => manager.emptyCache();
}
