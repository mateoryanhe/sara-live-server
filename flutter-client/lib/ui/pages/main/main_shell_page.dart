import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/l10n/locale_controller.dart';
import 'package:sara_live_app/ui/pages/main/anchor_main_shell_tabs.dart';
import 'package:sara_live_app/ui/pages/main/main_shell_role.dart';
import 'package:sara_live_app/ui/pages/main/main_shell_tab.dart';
import 'package:sara_live_app/ui/pages/main/viewer_main_shell_tabs.dart';
import 'package:sara_live_app/ui/theme/theme_controller.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';

/// 主壳：按 [MainShellRole] 切换观众 / 主播底部 Tab。
class MainShellPage extends StatefulWidget {
  const MainShellPage({
    super.key,
    required this.scope,
    required this.themeController,
    required this.localeController,
  });

  final AppScope scope;
  final ThemeController themeController;
  final LocaleController localeController;

  @override
  State<MainShellPage> createState() => _MainShellPageState();
}

class _MainShellPageState extends State<MainShellPage> {
  MainShellRole _role = MainShellRole.viewer;
  int _index = 0;
  List<MainShellTab>? _cachedTabs;
  MainShellRole? _cachedRole;

  @override
  void initState() {
    super.initState();
    widget.scope.session.addListener(_onSessionChanged);
    _resolveRoleFromProfile();
  }

  @override
  void dispose() {
    widget.scope.session.removeListener(_onSessionChanged);
    super.dispose();
  }

  void _onSessionChanged() {
    _cachedRole = null;
    _resolveRoleFromProfile();
  }

  List<MainShellTab> _freshTabs(AppLocalizations l10n, MainShellRole role) {
    switch (role) {
      case MainShellRole.viewer:
        return buildViewerMainShellTabs(
          scope: widget.scope,
          themeController: widget.themeController,
          localeController: widget.localeController,
          l10n: l10n,
        );
      case MainShellRole.anchor:
        return buildAnchorMainShellTabs(
          scope: widget.scope,
          themeController: widget.themeController,
          localeController: widget.localeController,
          l10n: l10n,
        );
    }
  }

  /// 角色变化时重建页面；仅语言变化时保留 [MainShellTab.page] 实例。
  List<MainShellTab> _tabs(AppLocalizations l10n) {
    final fresh = _freshTabs(l10n, _role);
    if (_cachedTabs == null || _cachedRole != _role) {
      _cachedTabs = fresh;
      _cachedRole = _role;
      return fresh;
    }
    final merged = <MainShellTab>[];
    for (var i = 0; i < fresh.length; i++) {
      merged.add(
        MainShellTab(
          id: fresh[i].id,
          label: fresh[i].label,
          icon: fresh[i].icon,
          selectedIcon: fresh[i].selectedIcon,
          page: _cachedTabs![i].page,
        ),
      );
    }
    _cachedTabs = merged;
    return merged;
  }

  Future<void> _resolveRoleFromProfile() async {
    final userId =
        int.tryParse(widget.scope.session.token?.userId ?? '') ?? 0;
    if (userId <= 0) {
      return;
    }
    try {
      final profile =
          await widget.scope.userRepository.getProfile(userId: userId);
      if (!mounted) {
        return;
      }
      final role =
          profile.isAnchor ? MainShellRole.anchor : MainShellRole.viewer;
      if (role == _role) {
        return;
      }
      setState(() {
        _role = role;
        _cachedRole = null;
        _index = 0;
      });
    } catch (_) {
      // 保持默认观众 Tab，避免阻塞进主壳
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final tabs = _tabs(l10n);
    final index = _index.clamp(0, tabs.length - 1);

    return Scaffold(
      body: IndexedStack(
        index: index,
        children: [for (final tab in tabs) tab.page],
      ),
      bottomNavigationBar: AppPageFrame(
        expand: false,
        child: NavigationBar(
          selectedIndex: index,
          onDestinationSelected: (i) => setState(() => _index = i),
          destinations: [for (final tab in tabs) tab.toDestination()],
        ),
      ),
    );
  }
}
