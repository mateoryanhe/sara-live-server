import '../../data/live_room/live_room_models.dart';

abstract class AnchorShowLiveUi {
  void showShowLiveLoading(bool loading);

  void renderShowLive({
    required LiveRoomDetail? room,
    required bool busy,
    required bool publishing,
  });

  void showShowLiveError(String message);

  void showShowLiveMessage(String message);

  /// 开播成功后进入全屏直播间。
  void openShowLiveBroadcast(LiveRoomDetail room);
}
