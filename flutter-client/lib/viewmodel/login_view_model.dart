import '../app/app_scope.dart';
import '../core/http/api_exception.dart';
import '../ui/contracts/login_ui.dart';

/// 设备码登录编排。
class LoginViewModel {
  LoginViewModel({
    required AppScope scope,
    required LoginUi ui,
  })  : _scope = scope,
        _ui = ui;

  final AppScope _scope;
  final LoginUi _ui;

  /// 有本地 token 则直接成功；否则走 `/auth/deviceLogin`。
  Future<void> start({bool forceRelogin = false}) async {
    if (_scope.isAccountSwitchInProgress ||
        _scope.session.isDeviceCodeLoginPending) {
      _ui.hideLoginLoading();
      return;
    }
    _ui.showLoginLoading();
    try {
      if (!forceRelogin && _scope.session.isLoggedIn) {
        await _scope.handleLoggedIn();
        _ui.onLoginSuccess(isNewUser: false);
        return;
      }

      final deviceId = await _scope.session.getOrCreateDeviceId();
      final deviceInfo =
          await _scope.deviceInfoBuilder.build(deviceId: deviceId);
      final result = await _scope.authApi.deviceLogin(deviceInfo);
      await _scope.session.saveToken(result.token);
      await _scope.handleLoggedIn();
      _ui.onLoginSuccess(isNewUser: result.isNewUser);
    } on ApiException catch (e) {
      _ui.showLoginError(e.message);
    } catch (e) {
      _ui.showLoginError(e.toString());
    } finally {
      _ui.hideLoginLoading();
    }
  }

  /// Log in with the device code entered on the explicit account-switch page.
  Future<void> loginWithDeviceId(String deviceId) async {
    if (_scope.isAccountSwitchInProgress) {
      return;
    }
    final normalized = deviceId.trim();
    if (normalized.isEmpty) {
      _ui.showLoginError('设备码不能为空');
      return;
    }

    _ui.showLoginLoading();
    try {
      final result = await _scope.loginWithDeviceCode(normalized);
      _ui.onLoginSuccess(isNewUser: result.isNewUser);
    } on ApiException catch (e) {
      _ui.showLoginError(e.message);
    } catch (e) {
      _ui.showLoginError(e.toString());
    } finally {
      _ui.hideLoginLoading();
    }
  }
}
