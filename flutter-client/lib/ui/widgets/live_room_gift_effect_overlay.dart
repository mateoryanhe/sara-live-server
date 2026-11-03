import 'dart:async';
import 'dart:io';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

import '../../core/cache/image_disk_cache.dart';
import '../../data/gift/gift_media_url.dart';
import '../../feature/show_live/live_room_gift_effect_player.dart';
import 'gift_split_alpha_video.dart';

/// 全屏礼物特效层（队列当前项播完需调 [onFinished]）。
class LiveRoomGiftEffectOverlay extends StatefulWidget {
  const LiveRoomGiftEffectOverlay({
    super.key,
    required this.player,
  });

  final LiveRoomGiftEffectPlayer player;

  @override
  State<LiveRoomGiftEffectOverlay> createState() =>
      _LiveRoomGiftEffectOverlayState();
}

class _LiveRoomGiftEffectOverlayState extends State<LiveRoomGiftEffectOverlay> {
  VideoPlayerController? _video;
  Timer? _imageTimer;

  @override
  void initState() {
    super.initState();
    widget.player.addListener(_onPlayer);
    _syncFromPlayer();
  }

  @override
  void dispose() {
    widget.player.removeListener(_onPlayer);
    _disposeVideo();
    _imageTimer?.cancel();
    super.dispose();
  }

  void _onPlayer() {
    if (!mounted) return;
    setState(() {});
    _syncFromPlayer();
  }

  void _syncFromPlayer() {
    final state = widget.player.current;
    _imageTimer?.cancel();
    _imageTimer = null;
    _disposeVideo();
    if (state == null) return;

    switch (state.mediaKind) {
      case GiftEffectMediaKind.video:
        unawaited(_playVideo(state.localPath));
      case GiftEffectMediaKind.image:
        _imageTimer = Timer(const Duration(milliseconds: 2600), _finish);
      case GiftEffectMediaKind.unknown:
        _imageTimer = Timer(const Duration(milliseconds: 2600), _finish);
    }
  }

  Future<void> _playVideo(String path) async {
    final controller = VideoPlayerController.file(File(path));
    _video = controller;
    try {
      await controller.initialize();
      if (!mounted || widget.player.current == null) {
        await controller.dispose();
        return;
      }
      controller.addListener(_onVideoTick);
      await controller.setLooping(false);
      await controller.setVolume(0);
      await controller.play();
      setState(() {});
    } catch (_) {
      _finish();
    }
  }

  void _onVideoTick() {
    final c = _video;
    if (c == null || !c.value.isInitialized) return;
    if (c.value.position >= c.value.duration) {
      _finish();
    }
  }

  void _disposeVideo() {
    _video?.removeListener(_onVideoTick);
    unawaited(_video?.dispose());
    _video = null;
  }

  void _finish() {
    _imageTimer?.cancel();
    _imageTimer = null;
    _disposeVideo();
    widget.player.completeCurrent();
  }

  @override
  Widget build(BuildContext context) {
    final state = widget.player.current;
    if (state == null) return const SizedBox.shrink();

    final job = state.job;
    final label = job.giftName.isNotEmpty ? job.giftName : 'Gift';
    final sender = job.senderName.isNotEmpty ? job.senderName : '';

    return IgnorePointer(
      child: Stack(
        alignment: Alignment.center,
        children: [
          if (state.mediaKind == GiftEffectMediaKind.video &&
              _video != null &&
              _video!.value.isInitialized)
            Positioned.fill(
              child: GiftSplitAlphaVideo(controller: _video!),
            )
          else if (state.mediaKind == GiftEffectMediaKind.image &&
              state.localPath.isNotEmpty)
            Positioned.fill(
              child: GiftSplitAlphaImageFile(path: state.localPath),
            )
          else
            _FallbackGiftBurst(
              iconUrl: job.iconUrl,
              giftName: label,
              count: job.count,
            ),
          Positioned(
            bottom: 120,
            child: DecoratedBox(
              decoration: BoxDecoration(
                color: Colors.black.withValues(alpha: 0.45),
                borderRadius: BorderRadius.circular(999),
              ),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
                child: Text(
                  sender.isNotEmpty ? '$sender · $label ×${job.count}' : '$label ×${job.count}',
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w700,
                    fontSize: 14,
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _FallbackGiftBurst extends StatelessWidget {
  const _FallbackGiftBurst({
    required this.iconUrl,
    required this.giftName,
    required this.count,
  });

  final String iconUrl;
  final String giftName;
  final int count;

  @override
  Widget build(BuildContext context) {
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0.6, end: 1),
      duration: const Duration(milliseconds: 500),
      curve: Curves.easeOutBack,
      builder: (context, scale, child) {
        return Transform.scale(scale: scale, child: child);
      },
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 120,
            height: 120,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              boxShadow: [
                BoxShadow(
                  color: const Color(0xFFF59E0B).withValues(alpha: 0.55),
                  blurRadius: 32,
                  spreadRadius: 8,
                ),
              ],
            ),
            child: iconUrl.isNotEmpty
                ? ClipOval(
                    child: CachedNetworkImage(
                      imageUrl: iconUrl,
                      cacheManager: ImageDiskCache.instance.manager,
                      fit: BoxFit.cover,
                      errorWidget: (_, __, ___) => const Icon(
                        Icons.card_giftcard_rounded,
                        size: 64,
                        color: Color(0xFFFFD166),
                      ),
                    ),
                  )
                : const Icon(
                    Icons.card_giftcard_rounded,
                    size: 64,
                    color: Color(0xFFFFD166),
                  ),
          ),
          const SizedBox(height: 12),
          Text(
            giftName,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 18,
              fontWeight: FontWeight.w800,
              shadows: [Shadow(color: Colors.black54, blurRadius: 8)],
            ),
          ),
          Text(
            '×$count',
            style: const TextStyle(
              color: Color(0xFFFFE566),
              fontSize: 22,
              fontWeight: FontWeight.w900,
            ),
          ),
        ],
      ),
    );
  }
}
