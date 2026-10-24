import 'package:go_router/go_router.dart';

import 'app_routes.dart';
import 'app_scope.dart';
import '../ui/pages/main/main_shell_page.dart';
import '../ui/pages/splash/splash_page.dart';
import '../ui/pages/vip/vip_page.dart';
import '../l10n/locale_controller.dart';
import '../ui/theme/theme_controller.dart';

/// 创建应用路由；[scope] / [themeController] / [localeController] 注入各页面。
GoRouter createAppRouter(
  AppScope scope,
  ThemeController themeController,
  LocaleController localeController,
) {
  return GoRouter(
    initialLocation: AppRoutes.splash,
    refreshListenable: scope.session,
    routes: [
      GoRoute(
        path: AppRoutes.splash,
        name: 'splash',
        builder: (context, state) => SplashPage(scope: scope),
      ),
      GoRoute(
        path: AppRoutes.main,
        name: 'main',
        builder: (context, state) => MainShellPage(
          scope: scope,
          themeController: themeController,
          localeController: localeController,
        ),
      ),
      GoRoute(
        path: AppRoutes.vip,
        name: 'vip',
        builder: (context, state) => VipPage(scope: scope),
      ),
    ],
    redirect: (context, state) {
      final loggedIn = scope.session.isLoggedIn;
      final loc = state.matchedLocation;
      if (!loggedIn && loc != AppRoutes.splash) {
        return AppRoutes.splash;
      }
      return null;
    },
  );
}
