/// 登录页/启动页 UI 能力（Passive View）。
abstract class LoginUi {
  void showLoginLoading();
  void hideLoginLoading();
  void showLoginError(String message);
  void onLoginSuccess({required bool isNewUser});
}
