import '../../data/live_room/live_room_models.dart';

abstract class AnchorOneToOneSetupUi {
  void showOneToOneSetupLoading(bool loading);

  void renderOneToOneSetup({
    required OneToOneRoomConfig? config,
    required List<BillingTierItem> billingTiers,
    required bool saving,
  });

  void showOneToOneSetupError(String message);

  void showOneToOneSetupMessage(String message);
}
