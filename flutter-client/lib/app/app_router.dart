import 'package:go_router/go_router.dart';

import 'app_routes.dart';
import 'app_scope.dart';
import '../ui/pages/auth/device_code_login_page.dart';
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
        path: AppRoutes.deviceCodeLogin,
        name: 'device-code-login',
        builder: (context, state) => DeviceCodeLoginPage(scope: scope),
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
      final pendingDeviceCodeLogin =
          scope.session.isDeviceCodeLoginPending;

      // While switching accounts, never let the automatic splash login reuse
      // the previous device code. This also restores the editor after a
      // process restart in the middle of the flow.
      if (!loggedIn && pendingDeviceCodeLogin &&
          loc != AppRoutes.deviceCodeLogin) {
        return AppRoutes.deviceCodeLogin;
      }

      // A completed login always leaves the explicit editor for the main
      // shell, even if the session listener refreshes the router first.
      if (loggedIn &&
          (loc == AppRoutes.splash || loc == AppRoutes.deviceCodeLogin)) {
        return AppRoutes.main;
      }

      if (!loggedIn &&
          loc != AppRoutes.splash &&
          loc != AppRoutes.deviceCodeLogin) {
        return AppRoutes.splash;
      }
      return null;
    },
  );
}
