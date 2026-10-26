import 'package:flutter_test/flutter_test.dart';
import 'package:sara_live_app/core/http/auth_token.dart';
import 'package:sara_live_app/data/auth/session_store.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  test('device-code switch survives restart and completes atomically', () async {
    SharedPreferences.setMockInitialValues({
      'device_id': 'old-device',
      'auth_token': '10001.old-secret',
    });
    final firstPrefs = await SharedPreferences.getInstance();
    final first = SessionStore(firstPrefs);
    await first.load();

    await first.beginDeviceCodeLogin();
    expect(first.isLoggedIn, isTrue);
    expect(first.isDeviceCodeLoginPending, isTrue);

    final second = SessionStore(firstPrefs);
    await second.load();
    expect(second.storedDeviceId, 'old-device');
    expect(second.isDeviceCodeLoginPending, isTrue);
    expect(second.isLoggedIn, isFalse);

    await second.completeDeviceCodeLogin(
      deviceId: 'new-device',
      token: const AuthToken(userId: '20002', secret: 'new-secret'),
    );
    expect(second.isLoggedIn, isTrue);
    expect(second.token?.userId, '20002');
    expect(second.storedDeviceId, 'new-device');
    expect(second.isDeviceCodeLoginPending, isFalse);
  });
}
