import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/l10n/locale_controller.dart';
import 'package:sara_live_app/ui/pages/anchor/anchor_audience_page.dart';
import 'package:sara_live_app/ui/pages/anchor/anchor_one_to_one_setup_page.dart';
import 'package:sara_live_app/ui/pages/anchor/anchor_show_live_page.dart';
import 'package:sara_live_app/ui/pages/main/main_shell_tab.dart';
import 'package:sara_live_app/ui/pages/mine/mine_page.dart';
import 'package:sara_live_app/ui/theme/theme_controller.dart';

/// **主播** 底部 4 Tab：在线观众 / 1v1 房间 / 秀场开播 / 我的。
List<MainShellTab> buildAnchorMainShellTabs({
  required AppScope scope,
  required ThemeController themeController,
  required LocaleController localeController,
  required AppLocalizations l10n,
}) {
  return [
    MainShellTab(
      id: 'audience',
      label: l10n.tabOnlineAudience,
      icon: Icons.people_outline,
      selectedIcon: Icons.people_rounded,
      page: AnchorAudiencePage(scope: scope),
    ),
    MainShellTab(
      id: 'one_to_one_setup',
      label: l10n.tabOneToOneRoom,
      icon: Icons.videocam_outlined,
      selectedIcon: Icons.videocam_rounded,
      page: AnchorOneToOneSetupPage(scope: scope),
    ),
    MainShellTab(
      id: 'show_live',
      label: l10n.tabShowLive,
      icon: Icons.live_tv_outlined,
      selectedIcon: Icons.live_tv_rounded,
      page: AnchorShowLivePage(scope: scope),
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
