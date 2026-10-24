import 'dart:typed_data';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:sara_live_app/core/cache/image_disk_cache.dart';

/// 直播间视频层下方的封面与暗角，进房/推流前先有画面。
class LiveRoomCoverBackdrop extends StatelessWidget {
  const LiveRoomCoverBackdrop({
    super.key,
    this.coverUrl = '',
    this.coverBytes,
    this.fallbackAvatarUrl = '',
    this.subtitle,
    this.showSubtitle = true,
    this.videoForegroundActive = false,
  });

  final String coverUrl;
  final Uint8List? coverBytes;
  /// 封面为空或加载失败时，用主播头像铺满背景。
  final String fallbackAvatarUrl;
  final String? subtitle;
  final bool showSubtitle;
  /// 上层已有全屏视频时，略加深底部渐变，封面仍保留在视频下方。
  final bool videoForegroundActive;

  @override
  Widget build(BuildContext context) {
    final topAlpha = videoForegroundActive ? 0.35 : 0.28;
    final bottomAlpha = videoForegroundActive ? 0.82 : 0.65;

    return SizedBox.expand(
      child: Stack(
        fit: StackFit.expand,
        children: [
          _CoverImage(
            coverUrl: coverUrl,
            coverBytes: coverBytes,
            fallbackAvatarUrl: fallbackAvatarUrl,
          ),
          DecoratedBox(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                colors: [
                  Colors.black.withValues(alpha: topAlpha),
                  Colors.black.withValues(alpha: 0.12),
                  Colors.black.withValues(alpha: 0.28),
                  Colors.black.withValues(alpha: bottomAlpha),
                ],
                stops: const [0, 0.35, 0.65, 1],
              ),
            ),
          ),
          if (showSubtitle &&
              subtitle != null &&
              subtitle!.trim().isNotEmpty)
            Positioned(
              left: 16,
              right: 16,
              bottom: 120,
              child: Text(
                subtitle!.trim(),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: Colors.white.withValues(alpha: 0.85),
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  shadows: const [
                    Shadow(color: Colors.black54, blurRadius: 8),
                  ],
                ),
              ),
            ),
        ],
      ),
    );
  }
}

class _CoverImage extends StatelessWidget {
  const _CoverImage({
    required this.coverUrl,
    this.coverBytes,
    this.fallbackAvatarUrl = '',
  });

  final String coverUrl;
  final Uint8List? coverBytes;
  final String fallbackAvatarUrl;

  @override
  Widget build(BuildContext context) {
    if (coverBytes != null && coverBytes!.isNotEmpty) {
      return Image.memory(
        coverBytes!,
        fit: BoxFit.cover,
        width: double.infinity,
        height: double.infinity,
      );
    }
    final cover = coverUrl.trim();
    if (cover.isNotEmpty) {
      return CachedNetworkImage(
        imageUrl: cover,
        fit: BoxFit.cover,
        width: double.infinity,
        height: double.infinity,
        cacheManager: ImageDiskCache.instance.manager,
        placeholder: (_, __) => _CoverImage(
          coverUrl: '',
          fallbackAvatarUrl: fallbackAvatarUrl,
        ),
        errorWidget: (_, __, ___) => _CoverImage(
          coverUrl: '',
          fallbackAvatarUrl: fallbackAvatarUrl,
        ),
      );
    }
    final avatar = fallbackAvatarUrl.trim();
    if (avatar.isNotEmpty) {
      return CachedNetworkImage(
        imageUrl: avatar,
        fit: BoxFit.cover,
        width: double.infinity,
        height: double.infinity,
        cacheManager: ImageDiskCache.instance.manager,
        placeholder: (_, __) => const _CoverFallback(),
        errorWidget: (_, __, ___) => const _CoverFallback(),
      );
    }
    return const _CoverFallback();
  }
}

class _CoverFallback extends StatelessWidget {
  const _CoverFallback();

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [Color(0xFF1E293B), Color(0xFF0F172A)],
        ),
      ),
      alignment: Alignment.center,
      child: Icon(
        Icons.live_tv_rounded,
        size: 72,
        color: Colors.white.withValues(alpha: 0.18),
      ),
    );
  }
}
