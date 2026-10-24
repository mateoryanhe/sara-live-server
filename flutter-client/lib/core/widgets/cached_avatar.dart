import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import '../cache/image_disk_cache.dart';

/// 网络头像：走 core [ImageDiskCache]（按设备档位限磁盘对象数）。
class CachedAvatar extends StatelessWidget {
  const CachedAvatar({
    super.key,
    required this.url,
    required this.radius,
    this.fallbackText = '?',
  });

  final String url;
  final double radius;
  final String fallbackText;

  @override
  Widget build(BuildContext context) {
    final size = radius * 2;
    if (url.isEmpty) {
      return CircleAvatar(
        radius: radius,
        child: Text(fallbackText),
      );
    }
    return ClipOval(
      child: CachedNetworkImage(
        imageUrl: url,
        cacheManager: ImageDiskCache.instance.manager,
        width: size,
        height: size,
        fit: BoxFit.cover,
        memCacheWidth: size.round() * 2,
        memCacheHeight: size.round() * 2,
        placeholder: (_, __) => SizedBox(
          width: size,
          height: size,
          child: const Center(
            child: SizedBox(
              width: 18,
              height: 18,
              child: CircularProgressIndicator(strokeWidth: 2),
            ),
          ),
        ),
        errorWidget: (_, __, ___) => CircleAvatar(
          radius: radius,
          child: Text(fallbackText),
        ),
      ),
    );
  }
}
