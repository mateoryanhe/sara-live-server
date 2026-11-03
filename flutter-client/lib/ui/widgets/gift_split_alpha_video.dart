import 'dart:async';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

/// 礼物 MP4 常见为「左 RGB + 右 Alpha 灰度/白模」同帧宽 2:1，需只取左半并做透明合成。
class GiftSplitAlphaVideo extends StatelessWidget {
  const GiftSplitAlphaVideo({
    super.key,
    required this.controller,
    this.splitAlphaThreshold = 1.55,
  });

  final VideoPlayerController controller;
  final double splitAlphaThreshold;

  bool get _splitAlpha {
    if (!controller.value.isInitialized) return false;
    return controller.value.aspectRatio >= splitAlphaThreshold;
  }

  @override
  Widget build(BuildContext context) {
    if (!controller.value.isInitialized) {
      return const SizedBox.shrink();
    }

    if (!_splitAlpha) {
      return Center(
        child: AspectRatio(
          aspectRatio: controller.value.aspectRatio,
          child: VideoPlayer(controller),
        ),
      );
    }

    final fullAr = controller.value.aspectRatio;
    final colorAr = fullAr / 2;

    return Center(
      child: AspectRatio(
        aspectRatio: colorAr,
        child: Stack(
          fit: StackFit.expand,
          children: [
            _RgbHalf(controller: controller, fullAspectRatio: fullAr),
            _AlphaMaskHalf(controller: controller, fullAspectRatio: fullAr),
          ],
        ),
      ),
    );
  }
}

/// 左半：彩色画面（放大到与合成区域同尺寸）。
class _RgbHalf extends StatelessWidget {
  const _RgbHalf({
    required this.controller,
    required this.fullAspectRatio,
  });

  final VideoPlayerController controller;
  final double fullAspectRatio;

  @override
  Widget build(BuildContext context) {
    return ClipRect(
      child: Align(
        alignment: Alignment.centerLeft,
        widthFactor: 0.5,
        child: OverflowBox(
          alignment: Alignment.centerLeft,
          minWidth: 0,
          maxWidth: double.infinity,
          maxHeight: double.infinity,
          child: AspectRatio(
            aspectRatio: fullAspectRatio,
            child: VideoPlayer(controller),
          ),
        ),
      ),
    );
  }
}

/// 右半：作为 Alpha，与左半彩色做 [BlendMode.dstIn] 合成。
class _AlphaMaskHalf extends StatelessWidget {
  const _AlphaMaskHalf({
    required this.controller,
    required this.fullAspectRatio,
  });

  final VideoPlayerController controller;
  final double fullAspectRatio;

  @override
  Widget build(BuildContext context) {
    return ShaderMask(
      blendMode: BlendMode.dstIn,
      shaderCallback: (bounds) {
        return const LinearGradient(
          colors: [Colors.white, Colors.white],
        ).createShader(bounds);
      },
      child: ClipRect(
        child: Align(
          alignment: Alignment.centerRight,
          widthFactor: 0.5,
          child: OverflowBox(
            alignment: Alignment.centerRight,
            minWidth: 0,
            maxWidth: double.infinity,
            maxHeight: double.infinity,
            child: ColorFiltered(
              colorFilter: const ColorFilter.matrix(<double>[
                0.2126, 0.7152, 0.0722, 0, 0,
                0.2126, 0.7152, 0.0722, 0, 0,
                0.2126, 0.7152, 0.0722, 0, 0,
                0, 0, 0, 1, 0,
              ]),
              child: AspectRatio(
                aspectRatio: fullAspectRatio,
                child: VideoPlayer(controller),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

bool giftVideoLikelySplitAlpha(double aspectRatio) =>
    aspectRatio >= 1.55;

/// 静态礼物图（左 RGB + 右 Alpha）与 MP4 同布局。
class GiftSplitAlphaImageFile extends StatelessWidget {
  const GiftSplitAlphaImageFile({
    super.key,
    required this.path,
    this.splitAlphaThreshold = 1.55,
  });

  final String path;
  final double splitAlphaThreshold;

  @override
  Widget build(BuildContext context) {
    final image = Image.file(
      File(path),
      fit: BoxFit.contain,
      filterQuality: FilterQuality.medium,
    );

    return LayoutBuilder(
      builder: (context, constraints) {
        return FutureBuilder<ImageInfo>(
          future: _loadImageInfo(path),
          builder: (context, snap) {
            final w = snap.data?.image.width.toDouble() ?? 0;
            final h = snap.data?.image.height.toDouble() ?? 1;
            final ar = w / h;
            if (w <= 0 || ar < splitAlphaThreshold) {
              return Center(child: image);
            }
            final colorAr = ar / 2;
            return Center(
              child: AspectRatio(
                aspectRatio: colorAr,
                child: Stack(
                  fit: StackFit.expand,
                  children: [
                    _ImageRgbHalf(fullChild: image, fullAspectRatio: ar),
                    _ImageAlphaHalf(fullChild: image, fullAspectRatio: ar),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }

  static Future<ImageInfo> _loadImageInfo(String path) {
    final provider = FileImage(File(path));
    final stream = provider.resolve(const ImageConfiguration());
    final completer = Completer<ImageInfo>();
    late ImageStreamListener listener;
    listener = ImageStreamListener(
      (info, _) {
        stream.removeListener(listener);
        completer.complete(info);
      },
      onError: (e, st) {
        stream.removeListener(listener);
        completer.completeError(e, st);
      },
    );
    stream.addListener(listener);
    return completer.future;
  }
}

class _ImageRgbHalf extends StatelessWidget {
  const _ImageRgbHalf({
    required this.fullChild,
    required this.fullAspectRatio,
  });

  final Widget fullChild;
  final double fullAspectRatio;

  @override
  Widget build(BuildContext context) {
    return ClipRect(
      child: Align(
        alignment: Alignment.centerLeft,
        widthFactor: 0.5,
        child: OverflowBox(
          alignment: Alignment.centerLeft,
          maxWidth: double.infinity,
          maxHeight: double.infinity,
          child: AspectRatio(
            aspectRatio: fullAspectRatio,
            child: fullChild,
          ),
        ),
      ),
    );
  }
}

class _ImageAlphaHalf extends StatelessWidget {
  const _ImageAlphaHalf({
    required this.fullChild,
    required this.fullAspectRatio,
  });

  final Widget fullChild;
  final double fullAspectRatio;

  @override
  Widget build(BuildContext context) {
    return ShaderMask(
      blendMode: BlendMode.dstIn,
      shaderCallback: (bounds) {
        return const LinearGradient(
          colors: [Colors.white, Colors.white],
        ).createShader(bounds);
      },
      child: ClipRect(
        child: Align(
          alignment: Alignment.centerRight,
          widthFactor: 0.5,
          child: OverflowBox(
            alignment: Alignment.centerRight,
            maxWidth: double.infinity,
            maxHeight: double.infinity,
            child: ColorFiltered(
              colorFilter: const ColorFilter.matrix(<double>[
                0.2126, 0.7152, 0.0722, 0, 0,
                0.2126, 0.7152, 0.0722, 0, 0,
                0.2126, 0.7152, 0.0722, 0, 0,
                0, 0, 0, 1, 0,
              ]),
              child: AspectRatio(
                aspectRatio: fullAspectRatio,
                child: fullChild,
              ),
            ),
          ),
        ),
      ),
    );
  }
}
