import 'dart:async';
import 'dart:io';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_vap_player/flutter_vap_player.dart';
import 'package:video_player/video_player.dart';

import '../../core/cache/image_disk_cache.dart';
import '../../data/gift/gift_media_url.dart';
import '../../data/gift/gift_vap_config.dart';
import '../../feature/show_live/live_room_gift_effect_player.dart';
import 'gift_split_alpha_video.dart';

/// 全屏礼物特效层（队列当前项播完需调 [onFinished]）。
///
/// 测服礼物 mp4 几乎全是腾讯 VAP（含 `vapc`）；用 [flutter_vap_player] 按
/// rgbFrame/aFrame 做透明合成。无 vapc 时回退 video_player。
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
  VapPlayerController? _vap;
  StreamSubscription<VapEvent>? _vapEvents;
  Timer? _imageTimer;
  GiftVapConfig? _vapCfg;
  int _playGen = 0;

  @override
  void initState() {
    super.initState();
    widget.player.addListener(_onPlayer);
    _syncFromPlayer();
  }

  @override
  void dispose() {
    widget.player.removeListener(_onPlayer);
    _disposePlayers();
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
    _disposePlayers();
    if (state == null) return;

    final gen = ++_playGen;
    switch (state.mediaKind) {
      case GiftEffectMediaKind.video:
        unawaited(_playVideoOrVap(state.localPath, gen));
      case GiftEffectMediaKind.image:
        _imageTimer = Timer(const Duration(milliseconds: 2600), _finish);
      case GiftEffectMediaKind.unknown:
        _imageTimer = Timer(const Duration(milliseconds: 2600), _finish);
    }
  }

  Future<void> _playVideoOrVap(String path, int gen) async {
    final vapCfg = await GiftVapConfig.tryParseFile(path);
    if (!mounted || gen != _playGen || widget.player.current == null) return;

    if (vapCfg != null) {
      await _playVap(path, vapCfg, gen);
      return;
    }
    await _playPlainVideo(path, gen);
  }

  Future<void> _playVap(String path, GiftVapConfig cfg, int gen) async {
    final controller = VapPlayerController.file(
      File(path),
      options: const VapPlayerOptions(
        viewType: VapViewType.textureView,
        repeatCount: 0,
        mute: true,
        scaleType: VapScaleType.fitCenter,
      ),
    );
    _vap = controller;
    _vapCfg = cfg;
    _vapEvents = controller.events.listen((event) {
      switch (event) {
        case VapCompletedEvent():
          _finish();
        case VapErrorEvent():
          _finish();
        default:
          break;
      }
    });
    try {
      await controller.initialize();
      if (!mounted || gen != _playGen || widget.player.current == null) {
        await controller.dispose();
        return;
      }
      setState(() {});
      await controller.play();
    } catch (_) {
      _finish();
    }
  }

  Future<void> _playPlainVideo(String path, int gen) async {
    final controller = VideoPlayerController.file(File(path));
    _video = controller;
    try {
      await controller.initialize();
      if (!mounted || gen != _playGen || widget.player.current == null) {
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

  void _disposePlayers() {
    _vapEvents?.cancel();
    _vapEvents = null;
    unawaited(_vap?.dispose());
    _vap = null;
    _vapCfg = null;
    _video?.removeListener(_onVideoTick);
    unawaited(_video?.dispose());
    _video = null;
  }

  void _finish() {
    _imageTimer?.cancel();
    _imageTimer = null;
    _disposePlayers();
    widget.player.completeCurrent();
  }

  @override
  Widget build(BuildContext context) {
    final state = widget.player.current;
    if (state == null) return const SizedBox.shrink();

    final job = state.job;
    final label = job.giftName.isNotEmpty ? job.giftName : 'Gift';
    final sender = job.senderName.isNotEmpty ? job.senderName : '';
    final vap = _vap;
    final video = _video;

    return IgnorePointer(
      child: Stack(
        alignment: Alignment.center,
        children: [
          if (vap != null && vap.value.isInitialized)
            Positioned.fill(
              child: Center(
                child: AspectRatio(
                  aspectRatio: (_vapCfg?.outAspectRatio ??
                      (vap.value.size.width > 0 && vap.value.size.height > 0
                          ? vap.value.size.width / vap.value.size.height
                          : 9 / 16)),
                  child: VapPlayer(vap),
                ),
              ),
            )
          else if (video != null && video.value.isInitialized)
            Positioned.fill(
              child: GiftSplitAlphaVideo(controller: video),
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
                  sender.isNotEmpty
                      ? '$sender · $label ×${job.count}'
                      : '$label ×${job.count}',
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
