import 'package:flutter/material.dart';

/// 公屏 VIP 气泡/边框配色（无 CMS 图标时的兜底，按等级分档）。
class VipChatAppearance {
  const VipChatAppearance({
    required this.borderGradient,
    required this.fillColor,
    required this.textColor,
    required this.badgeGradient,
    required this.badgeTextColor,
  });

  final Gradient borderGradient;
  final Color fillColor;
  final Color textColor;
  final Gradient badgeGradient;
  final Color badgeTextColor;

  static VipChatAppearance forLevel(int vipLevel) {
    if (vipLevel <= 0) return vip0;
    final index = (vipLevel - 1).clamp(0, _vipTiers.length - 1);
    return _vipTiers[index];
  }

  /// 公屏非会员 / VIP0 统一灰阶样式。
  static const vip0 = VipChatAppearance(
    borderGradient: LinearGradient(
      colors: [Color(0xFF94A3B8), Color(0xFF64748B)],
    ),
    fillColor: Color(0xCC1E293B),
    textColor: Colors.white,
    badgeGradient: LinearGradient(
      colors: [Color(0xFF475569), Color(0xFF334155)],
    ),
    badgeTextColor: Colors.white,
  );

  static const _vipTiers = <VipChatAppearance>[
    VipChatAppearance(
      borderGradient: LinearGradient(
        colors: [Color(0xFF7DD3FC), Color(0xFF38BDF8)],
      ),
      fillColor: Color(0xCC0C4A6E),
      textColor: Color(0xFFE0F2FE),
      badgeGradient: LinearGradient(
        colors: [Color(0xFF0EA5E9), Color(0xFF0369A1)],
      ),
      badgeTextColor: Colors.white,
    ),
    VipChatAppearance(
      borderGradient: LinearGradient(
        colors: [Color(0xFF86EFAC), Color(0xFF22C55E)],
      ),
      fillColor: Color(0xCC14532D),
      textColor: Color(0xFFDCFCE7),
      badgeGradient: LinearGradient(
        colors: [Color(0xFF22C55E), Color(0xFF15803D)],
      ),
      badgeTextColor: Colors.white,
    ),
    VipChatAppearance(
      borderGradient: LinearGradient(
        colors: [Color(0xFFC4B5FD), Color(0xFF8B5CF6)],
      ),
      fillColor: Color(0xCC3B0764),
      textColor: Color(0xFFF3E8FF),
      badgeGradient: LinearGradient(
        colors: [Color(0xFF8B5CF6), Color(0xFF6D28D9)],
      ),
      badgeTextColor: Colors.white,
    ),
    VipChatAppearance(
      borderGradient: LinearGradient(
        colors: [Color(0xFFFDA4AF), Color(0xFFFB7185)],
      ),
      fillColor: Color(0xCC881337),
      textColor: Color(0xFFFFE4E6),
      badgeGradient: LinearGradient(
        colors: [Color(0xFFFB7185), Color(0xFFE11D48)],
      ),
      badgeTextColor: Colors.white,
    ),
    VipChatAppearance(
      borderGradient: LinearGradient(
        colors: [Color(0xFFFDE68A), Color(0xFFFBBF24), Color(0xFFF59E0B)],
      ),
      fillColor: Color(0xCC422006),
      textColor: Color(0xFFFFF7ED),
      badgeGradient: LinearGradient(
        colors: [Color(0xFFFBBF24), Color(0xFFD97706), Color(0xFFB45309)],
      ),
      badgeTextColor: Color(0xFF422006),
    ),
  ];
}
