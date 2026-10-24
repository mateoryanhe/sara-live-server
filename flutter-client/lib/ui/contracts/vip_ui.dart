import '../../data/vip/vip_models.dart';

/// VIP 详情页 UI 能力。
abstract class VipUi {
  void showVipLoading();
  void hideVipLoading();
  void renderVip({
    required VipDetail detail,
    required List<VipCfgItem> cfgList,
  });
  void showVipError(String message);
}
