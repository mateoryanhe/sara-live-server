/// 礼物 icon / animation 地址归一化（推送里可能是相对路径或文件名）。
String resolveGiftMediaUrl(String raw, String apiBaseUrl) {
  final t = raw.trim();
  if (t.isEmpty) return '';
  if (t.startsWith('http://') || t.startsWith('https://')) return t;
  final base = apiBaseUrl.replaceAll(RegExp(r'/$'), '');
  if (t.startsWith('/')) return '$base$t';
  return '$base/$t';
}

enum GiftEffectMediaKind {
  video,
  image,
  unknown,
}

GiftEffectMediaKind giftEffectMediaKindForUrl(String url) {
  final path = Uri.tryParse(url)?.path.toLowerCase() ?? url.toLowerCase();
  if (path.endsWith('.mp4') ||
      path.endsWith('.webm') ||
      path.endsWith('.mov')) {
    return GiftEffectMediaKind.video;
  }
  if (path.endsWith('.gif') ||
      path.endsWith('.png') ||
      path.endsWith('.jpg') ||
      path.endsWith('.jpeg') ||
      path.endsWith('.webp') ||
      path.endsWith('.apng')) {
    return GiftEffectMediaKind.image;
  }
  return GiftEffectMediaKind.unknown;
}
