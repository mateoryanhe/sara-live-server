import '../../data/gift/gift_media_url.dart';

/// 待播放的直播间礼物特效任务（来自 cmd=8 推送）。
class LiveRoomGiftEffectJob {
  LiveRoomGiftEffectJob({
    required this.roomId,
    required this.giftId,
    required this.giftName,
    required this.effectUrl,
    required this.iconUrl,
    required this.count,
    required this.senderName,
    required this.enqueuedAt,
  });

  final int roomId;
  final String giftId;
  final String giftName;
  final String effectUrl;
  final String iconUrl;
  final int count;
  final String senderName;
  final DateTime enqueuedAt;

  /// 仅播礼物动画；无 animation 的礼物不进全屏特效（公屏仍展示 icon）。
  String get playUrl => effectUrl;

  factory LiveRoomGiftEffectJob.fromPushMap(
    Map<String, dynamic> map,
    int roomId,
    String apiBaseUrl,
  ) {
    final anim = resolveGiftMediaUrl('${map['giftAnim'] ?? ''}', apiBaseUrl);
    final icon = resolveGiftMediaUrl('${map['giftIcon'] ?? ''}', apiBaseUrl);
    return LiveRoomGiftEffectJob(
      roomId: roomId,
      giftId: '${map['giftId'] ?? ''}',
      giftName: '${map['giftName'] ?? ''}',
      effectUrl: anim,
      iconUrl: icon,
      count: _int(map['count']),
      senderName: '${map['senderName'] ?? ''}',
      enqueuedAt: DateTime.now(),
    );
  }

  static int _int(Object? v) {
    if (v is num) return v.toInt();
    return int.tryParse('$v') ?? 1;
  }
}

class LiveRoomGiftEffectPlayState {
  LiveRoomGiftEffectPlayState({
    required this.job,
    required this.localPath,
    required this.mediaKind,
  });

  final LiveRoomGiftEffectJob job;
  final String localPath;
  final GiftEffectMediaKind mediaKind;
}
