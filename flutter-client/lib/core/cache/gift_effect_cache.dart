import 'dart:io';

import 'package:flutter_cache_manager/flutter_cache_manager.dart';

import 'device_tier.dart';

/// 礼物特效资源磁盘缓存（mp4/gif/webp 等，独立于 [ImageDiskCache]）。
class GiftEffectCache {
  GiftEffectCache._(this.tier, this.manager);

  final DeviceTier tier;
  final CacheManager manager;

  static GiftEffectCache? _instance;

  static GiftEffectCache get instance {
    final cur = _instance;
    if (cur == null) {
      throw StateError('GiftEffectCache.init(tier) must be called at startup');
    }
    return cur;
  }

  static GiftEffectCache init(DeviceTier tier) {
    final existing = _instance;
    if (existing != null) return existing;
    final maxObjects = switch (tier) {
      DeviceTier.low => 24,
      DeviceTier.mid => 48,
      DeviceTier.high => 80,
    };
    final manager = CacheManager(
      Config(
        'sara_live_gift_effects_${tier.name}',
        stalePeriod: const Duration(days: 30),
        maxNrOfCacheObjects: maxObjects,
      ),
    );
    return _instance = GiftEffectCache._(tier, manager);
  }

  Future<bool> isCached(String url) async {
    if (url.trim().isEmpty) return false;
    final info = await manager.getFileFromCache(url);
    return info != null;
  }

  Future<File> getFile(String url) => manager.getSingleFile(url);

  Future<void> prefetch(String url) async {
    if (url.trim().isEmpty) return;
    await manager.downloadFile(url);
  }
}
