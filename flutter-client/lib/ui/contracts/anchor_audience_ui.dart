import '../../data/live_room/live_room_models.dart';

abstract class AnchorAudienceUi {
  void showAudienceLoading({required bool initial});

  void hideAudienceLoading();

  void renderAudienceList({
    required List<OnlineUserItem> items,
    required int total,
    required bool loadingMore,
  });

  void showAudienceError(String message);

  void showAudienceMessage(String message);
}
