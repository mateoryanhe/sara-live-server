import 'vip_chat_bubble_skin.dart';
import 'vip_models.dart';

/// 对齐 Go `vipcfgdto.AppVipPrivilegeType*`。
abstract final class VipPrivilegeType {
  static const int entryEffect = 2;
  static const int commentEffect = 3;
  static const int customerService = 4;
}

/// 由 App VIP 配置解析公屏头像框等资源（按等级查 [privilegeList]）。
class VipPrivilegeFrames {
  VipPrivilegeFrames._(this._cfgByLevel);

  factory VipPrivilegeFrames.fromCfgList(List<VipCfgItem> list) {
    final map = <int, VipCfgItem>{};
    for (final item in list) {
      if (item.level > 0) {
        map[item.level] = item;
      }
    }
    return VipPrivilegeFrames._(map);
  }

  static final empty = VipPrivilegeFrames._({});

  final Map<int, VipCfgItem> _cfgByLevel;

  VipCfgItem? cfgForLevel(int vipLevel) => _cfgByLevel[vipLevel];

  /// CMS 配置的等级图标 URL（`levelIcon`）。
  String? levelIconUrl(int vipLevel) {
    if (vipLevel <= 0) return null;
    final icon = _cfgByLevel[vipLevel]?.levelIcon.trim() ?? '';
    return icon.isNotEmpty ? icon : null;
  }

  String levelLabel(int vipLevel) {
    if (vipLevel <= 0) return 'VIP0';
    final cfg = _cfgByLevel[vipLevel];
    if (cfg != null && cfg.levelName.trim().isNotEmpty) {
      return cfg.levelName.trim();
    }
    return 'VIP$vipLevel';
  }

  /// 公屏评论特权（type=3）的 [VipPrivilegeItem.icon] 作为头像框。
  String? chatAvatarFrameUrl(int vipLevel) {
    if (vipLevel <= 0) return null;
    final cfg = _cfgByLevel[vipLevel];
    if (cfg == null) return null;
    for (final p in cfg.privilegeList) {
      if (p.privilegeType == VipPrivilegeType.commentEffect &&
          p.icon.trim().isNotEmpty) {
        return p.icon.trim();
      }
    }
    return null;
  }

  /// 公屏第三列文字气泡框：全员 APK 内置点九图（VIP0 为 frame1）。
  String? chatMessageBubbleFrameUrl(int vipLevel) {
    return VipChatBubbleSkin.frameAssetForLevel(vipLevel);
  }
}
