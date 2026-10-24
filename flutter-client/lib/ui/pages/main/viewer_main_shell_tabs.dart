import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/l10n/locale_controller.dart';
import 'package:sara_live_app/ui/pages/main/main_shell_tab.dart';
import 'package:sara_live_app/ui/pages/live/live_room_list_page.dart';
import 'package:sara_live_app/ui/pages/main/shell_tab_placeholder_page.dart';
import 'package:sara_live_app/ui/pages/mine/mine_page.dart';
import 'package:sara_live_app/ui/pages/one_to_one/one_to_one_page.dart';
import 'package:sara_live_app/ui/theme/theme_controller.dart';

/// **普通观众** 底部 4 Tab：1v1 / 直播 / 游戏 / 我的。
List<MainShellTab> buildViewerMainShellTabs({
  required AppScope scope,
  required ThemeController themeController,
  required LocaleController localeController,
  required AppLocalizations l10n,
}) {
  return [
    MainShellTab(
      id: 'one_to_one',
      label: l10n.tabOneToOne,
      icon: Icons.videocam_outlined,
      selectedIcon: Icons.videocam_rounded,
      page: OneToOnePage(scope: scope),
    ),
    MainShellTab(
      id: 'live',
      label: l10n.tabLive,
      icon: Icons.live_tv_outlined,
      selectedIcon: Icons.live_tv_rounded,
      page: LiveRoomListPage(scope: scope),
    ),
    MainShellTab(
      id: 'game',
      label: l10n.tabGame,
      icon: Icons.sports_esports_outlined,
      selectedIcon: Icons.sports_esports_rounded,
      page: const ShellTabPlaceholderPage(kind: ShellPlaceholderKind.game),
    ),
    MainShellTab(
      id: 'mine',
      label: l10n.tabMine,
      icon: Icons.person_outline_rounded,
      selectedIcon: Icons.person_rounded,
      page: MinePage(
        scope: scope,
        themeController: themeController,
        localeController: localeController,
      ),
    ),
  ];
}
