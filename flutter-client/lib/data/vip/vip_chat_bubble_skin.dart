import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import '../../core/cache/image_disk_cache.dart';

/// 公屏 VIP 聊天气泡框（点九拉伸参数对齐 `bubbleframe` 128×112，资源打进 APK）。
abstract final class VipChatBubbleSkin {
  static const _assetDir = 'assets/vip/chat_bubbles';

  /// 公屏点九框：VIP0 固定 `bubble_frame1`；VIP1+ 按等级在 12 套间轮换。
  static String frameAssetForLevel(int vipLevel) {
    final idx = vipLevel <= 0 ? 1 : ((vipLevel - 1) % 12) + 1;
    return '$_assetDir/bubble_frame$idx.png';
  }

  static const centerSlice = Rect.fromLTRB(48, 40, 80, 72);

  static const contentPadding = EdgeInsets.fromLTRB(16, 11, 16, 11);

  static bool isBundledAsset(String source) =>
      source.startsWith('assets/');

  static ImageProvider frameImageProvider(String source) {
    if (isBundledAsset(source)) {
      return AssetImage(source);
    }
    return CachedNetworkImageProvider(
      source,
      cacheManager: ImageDiskCache.instance.manager,
    );
  }
}

bool vipChatBubbleUrlLooksLikeImage(String url) {
  final lower = url.trim().toLowerCase();
  if (lower.isEmpty) return false;
  return lower.endsWith('.png') ||
      lower.endsWith('.webp') ||
      lower.endsWith('.jpg') ||
      lower.endsWith('.jpeg') ||
      lower.endsWith('.gif');
}
