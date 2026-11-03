import 'dart:async';

import '../../data/live_room/live_room_api.dart';

/// 在房心跳：`POST /liveRoom/reportLiveStartStatus`（Go 建议约 30s）。
///
/// 观众刷新在线 [HeartTime]；主播刷新房间 [HeartTime]，避免被判定离线/下播。
class LiveRoomHeartReporter {
  LiveRoomHeartReporter(this._api);

  final LiveRoomApi _api;

  static const period = Duration(seconds: 30);

  Timer? _timer;
  int? _roomId;

  void start(int roomId) {
    if (roomId <= 0) return;
    stop();
    _roomId = roomId;
    unawaited(_tick());
    _timer = Timer.periodic(period, (_) {
      unawaited(_tick());
    });
  }

  void stop() {
    _timer?.cancel();
    _timer = null;
    _roomId = null;
  }

  Future<void> _tick() async {
    final roomId = _roomId;
    if (roomId == null) return;
    try {
      await _api.reportLiveStartStatus(roomId: roomId);
    } catch (_) {}
  }
}
