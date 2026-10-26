import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/l10n/locale_controller.dart';
import 'package:sara_live_app/main.dart';
import 'package:sara_live_app/ui/theme/theme_controller.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  testWidgets('app boots to splash', (tester) async {
    SharedPreferences.setMockInitialValues({});
    final scope = await AppScope.create();
    final theme = await ThemeController.create();
    final locale = await LocaleController.create();
    await tester.pumpWidget(
      SaraLiveApp(
        scope: scope,
        themeController: theme,
        localeController: locale,
      ),
    );
    expect(find.text('Sara Live'), findsOneWidget);
    scope.dispose();
  });
}
