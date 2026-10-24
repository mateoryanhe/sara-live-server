import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:go_router/go_router.dart';

import 'app/app_router.dart';
import 'app/app_scope.dart';
import 'core/bus/app_bus_event.dart';
import 'l10n/app_localizations.dart';
import 'l10n/locale_controller.dart';
import 'ui/pages/call/call_session_overlay.dart';
import 'ui/theme/app_themes.dart';
import 'ui/theme/theme_controller.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final scope = await AppScope.create();
  final themeController = await ThemeController.create();
  final localeController = await LocaleController.create();
  ErrorWidget.builder = (details) {
    final isZh = localeController.locale.languageCode == 'zh';
    final title = isZh ? 'UI 渲染错误' : 'UI render error';
    return Material(
      color: const Color(0xFFFFF5F5),
      child: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(20),
          child: SelectableText(
            '$title\n\n${details.exceptionAsString()}',
            style: const TextStyle(color: Color(0xFFB91C1C), fontSize: 14),
          ),
        ),
      ),
    );
  };
  runApp(SaraLiveApp(
    scope: scope,
    themeController: themeController,
    localeController: localeController,
  ));
}

class SaraLiveApp extends StatefulWidget {
  const SaraLiveApp({
    super.key,
    required this.scope,
    required this.themeController,
    required this.localeController,
  });

  final AppScope scope;
  final ThemeController themeController;
  final LocaleController localeController;

  @override
  State<SaraLiveApp> createState() => _SaraLiveAppState();
}

class _SaraLiveAppState extends State<SaraLiveApp> with WidgetsBindingObserver {
  late final GoRouter _router = createAppRouter(
    widget.scope,
    widget.themeController,
    widget.localeController,
  );

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    final bus = widget.scope.appBus;
    switch (state) {
      case AppLifecycleState.paused:
        bus.publish(AppLifecycleBusEvent.paused());
      case AppLifecycleState.resumed:
        bus.publish(AppLifecycleBusEvent.resumed());
      case AppLifecycleState.detached:
      case AppLifecycleState.inactive:
      case AppLifecycleState.hidden:
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    // 勿在外层 ListenableBuilder 包 MaterialApp.router：hot reload / 换肤时会
    // 整棵重建 MaterialApp，Theme 仍 notify 旧子树 → framework 断言失败。
    final style = widget.themeController.style;
    final locale = widget.localeController.locale;
    return MaterialApp.router(
      title: 'Sara Live',
      theme: AppThemes.of(style),
      themeMode: ThemeMode.light,
      locale: locale,
      supportedLocales: LocaleController.supportedLocales,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      routerConfig: _router,
      builder: (context, child) {
        return _AppShellOverlay(
          scope: widget.scope,
          themeController: widget.themeController,
          localeController: widget.localeController,
          navChild: child,
        );
      },
    );
  }
}

/// 换肤 / 语言 / 通话浮层：只重建 shell，不重建 [MaterialApp.router]。
class _AppShellOverlay extends StatelessWidget {
  const _AppShellOverlay({
    required this.scope,
    required this.themeController,
    required this.localeController,
    required this.navChild,
  });

  final AppScope scope;
  final ThemeController themeController;
  final LocaleController localeController;
  final Widget? navChild;

  @override
  Widget build(BuildContext context) {
    return ListenableBuilder(
      listenable: Listenable.merge([
        themeController,
        localeController,
        scope.callController,
      ]),
      builder: (context, _) {
        final page = navChild ??
            const ColoredBox(
              color: Color(0xFFF8FAFC),
              child: Center(child: CircularProgressIndicator()),
            );
        final themed = Theme(
          data: AppThemes.of(themeController.style),
          child: Localizations.override(
            context: context,
            locale: localeController.locale,
            child: page,
          ),
        );
        if (!scope.callController.isActive) {
          return themed;
        }
        return Stack(
          fit: StackFit.expand,
          children: [
            Positioned.fill(child: themed),
            Positioned.fill(
              child: CallSessionOverlay(scope: scope),
            ),
          ],
        );
      },
    );
  }
}
