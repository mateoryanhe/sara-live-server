import 'package:flutter/material.dart';

/// 底部导航一项：标签 + 页面（由 [IndexedStack] 保活）。
class MainShellTab {
  const MainShellTab({
    required this.id,
    required this.label,
    required this.icon,
    required this.selectedIcon,
    required this.page,
  });

  final String id;
  final String label;
  final IconData icon;
  final IconData selectedIcon;
  final Widget page;

  NavigationDestination toDestination() => NavigationDestination(
        icon: Icon(icon),
        selectedIcon: Icon(selectedIcon),
        label: label,
      );
}
