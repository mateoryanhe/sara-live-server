import 'dart:convert';
import 'dart:io';
import 'dart:typed_data';

/// 腾讯 VAP（Video Animation Player）配置：从 mp4 的 `vapc` box 解析。
class GiftVapConfig {
  const GiftVapConfig({
    required this.outW,
    required this.outH,
    required this.videoW,
    required this.videoH,
    required this.rgbFrame,
    required this.aFrame,
  });

  final int outW;
  final int outH;
  final int videoW;
  final int videoH;

  /// [x, y, w, h] 彩色区域（相对 videoW/videoH）。
  final List<int> rgbFrame;

  /// [x, y, w, h] Alpha 灰模区域（相对 videoW/videoH）。
  final List<int> aFrame;

  double get outAspectRatio => outW <= 0 || outH <= 0 ? 1 : outW / outH;

  static Future<GiftVapConfig?> tryParseFile(String path) async {
    final file = File(path);
    if (!await file.exists()) return null;
    final len = await file.length();
    if (len < 16) return null;
    // vapc 通常紧挨 moov，读前 256KB 足够；过大则再扫全文件头段。
    final headLen = len < 256 * 1024 ? len : 256 * 1024;
    final bytes = await file.openRead(0, headLen).fold<BytesBuilder>(
      BytesBuilder(copy: false),
      (b, chunk) {
        b.add(chunk);
        return b;
      },
    );
    return tryParseBytes(bytes.takeBytes());
  }

  static GiftVapConfig? tryParseBytes(Uint8List data) {
    var offset = 0;
    while (offset + 8 <= data.length) {
      final size = _u32(data, offset);
      if (size < 8) break;
      final type = String.fromCharCodes(data.sublist(offset + 4, offset + 8));
      if (type == 'vapc') {
        final end = offset + size;
        if (end > data.length) return null;
        final payload = data.sublist(offset + 8, end);
        try {
          final text = utf8.decode(payload);
          final json = jsonDecode(text);
          if (json is! Map) return null;
          final info = json['info'];
          if (info is! Map) return null;
          final rgb = _intList(info['rgbFrame']);
          final alpha = _intList(info['aFrame']);
          if (rgb.length != 4 || alpha.length != 4) return null;
          final videoW = _int(info['videoW']);
          final videoH = _int(info['videoH']);
          final outW = _int(info['w']);
          final outH = _int(info['h']);
          if (videoW <= 0 || videoH <= 0 || outW <= 0 || outH <= 0) {
            return null;
          }
          return GiftVapConfig(
            outW: outW,
            outH: outH,
            videoW: videoW,
            videoH: videoH,
            rgbFrame: rgb,
            aFrame: alpha,
          );
        } catch (_) {
          return null;
        }
      }
      // size==0 表示到文件末尾；size==1 是 64-bit largesize，礼物文件里基本不会出现。
      if (size == 0) break;
      if (size == 1) {
        if (offset + 16 > data.length) break;
        final large = _u64(data, offset + 8);
        if (large < 16) break;
        offset += large;
        continue;
      }
      offset += size;
    }
    return null;
  }

  static int _u32(Uint8List d, int i) =>
      (d[i] << 24) | (d[i + 1] << 16) | (d[i + 2] << 8) | d[i + 3];

  static int _u64(Uint8List d, int i) {
    final hi = _u32(d, i);
    final lo = _u32(d, i + 4);
    return (hi << 32) + lo;
  }

  static int _int(Object? v) {
    if (v is num) return v.toInt();
    return int.tryParse('$v') ?? 0;
  }

  static List<int> _intList(Object? v) {
    if (v is! List) return const [];
    return v.map(_int).toList(growable: false);
  }
}
