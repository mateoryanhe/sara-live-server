import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import '../../core/cache/image_disk_cache.dart';
import '../../core/widgets/cached_avatar.dart';

/// 直播间公屏发言头像；[frameUrl] 来自 VIP 公屏评论特权 icon。
class LiveRoomChatAvatar extends StatelessWidget {
  const LiveRoomChatAvatar({
    super.key,
    required this.avatarUrl,
    required this.displayName,
    this.frameUrl,
    this.size = 32,
  });

  final String avatarUrl;
  final String displayName;
  final String? frameUrl;
  final double size;

  @override
  Widget build(BuildContext context) {
    final initial = displayName.isNotEmpty ? displayName.characters.first : '?';
    final avatarRadius = size * 0.34;
    final frame = frameUrl?.trim() ?? '';

    return SizedBox(
      width: size,
      height: size,
      child: Stack(
        alignment: Alignment.center,
        clipBehavior: Clip.none,
        children: [
          CachedAvatar(
            url: avatarUrl,
            radius: avatarRadius,
            fallbackText: initial,
          ),
          if (frame.isNotEmpty)
            Positioned.fill(
              child: CachedNetworkImage(
                imageUrl: frame,
                cacheManager: ImageDiskCache.instance.manager,
                fit: BoxFit.contain,
                errorWidget: (_, __, ___) => const SizedBox.shrink(),
              ),
            ),
        ],
      ),
    );
  }
}
