import '../../data/user/user_profile.dart';
import '../../data/vip/vip_models.dart';

/// 「我的」页 UI 能力。
abstract class MineUi {
  void showMineLoading();
  void hideMineLoading();
  void renderProfile(UserProfile profile);
  void renderVipDetail(VipDetail detail);
  void showMineError(String message);
  void showMineMessage(String message);
  void onLoggedOut();
  void onDeviceIdApplied(String deviceId);

  void onApiBaseApplied(String apiBaseUrl);
}
