import 'package:flutter/foundation.dart';
import 'package:sara_live_app/core/http/auth_token.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart';

/// 本地会话：deviceId 持久化 + 登录 token。
class SessionStore extends ChangeNotifier {
  SessionStore(this._prefs);

  static const _kDeviceId = 'device_id';
  static const _kAuthToken = 'auth_token';
  static const _kDeviceCodeLoginPending = 'device_code_login_pending';

  final SharedPreferences _prefs;
  AuthToken? _token;
  bool _deviceCodeLoginPending = false;

  AuthToken? get token => _token;

  bool get isLoggedIn => _token != null;

  /// Whether the user explicitly asked to switch accounts by device code.
  ///
  /// This flag is persisted so a process restart cannot silently reuse the
  /// previous device code while the user is in the middle of a switch.
  bool get isDeviceCodeLoginPending => _deviceCodeLoginPending;

  /// 已保存的设备码（未保存过则为 null）。
  String? get storedDeviceId {
    final raw = _prefs.getString(_kDeviceId);
    if (raw == null) return null;
    final trimmed = raw.trim();
    return trimmed.isEmpty ? null : trimmed;
  }

  Future<void> load() async {
    _deviceCodeLoginPending =
        _prefs.getBool(_kDeviceCodeLoginPending) ?? false;
    _token = AuthToken.tryParse(_prefs.getString(_kAuthToken));
    if (_deviceCodeLoginPending) {
      // A crash can happen after the pending marker is persisted but before
      // AppScope finishes clearing the old token. Never revive that session
      // on the next launch; the user must explicitly submit a device code.
      _token = null;
      await _prefs.remove(_kAuthToken);
    }
  }

  Future<void> saveToken(AuthToken token) async {
    _token = token;
    await _prefs.setString(_kAuthToken, token.headerValue);
    notifyListeners();
  }

  Future<void> clearToken() async {
    _token = null;
    await _prefs.remove(_kAuthToken);
    notifyListeners();
  }

  /// Mark the account-switch flow as requiring an explicit device-code login.
  ///
  /// The current device code is intentionally retained as the initial value
  /// in the editor. It is only replaced after the new device-code login has
  /// succeeded.
  Future<void> beginDeviceCodeLogin() async {
    if (_deviceCodeLoginPending) return;
    _deviceCodeLoginPending = true;
    await _prefs.setBool(_kDeviceCodeLoginPending, true);
    notifyListeners();
  }

  /// Persist a successful device-code login as one session transition.
  ///
  /// Updating the token, device code and pending marker together prevents the
  /// router from observing a half-switched session after a successful login.
  Future<void> completeDeviceCodeLogin({
    required String deviceId,
    required AuthToken token,
  }) async {
    final id = deviceId.trim();
    if (id.isEmpty) {
      throw ArgumentError.value(deviceId, 'deviceId', 'must not be empty');
    }
    _token = token;
    _deviceCodeLoginPending = false;
    await _prefs.setString(_kDeviceId, id);
    await _prefs.setString(_kAuthToken, token.headerValue);
    await _prefs.remove(_kDeviceCodeLoginPending);
    notifyListeners();
  }

  /// Leave the explicit device-code editor and return to normal startup.
  Future<void> cancelDeviceCodeLogin() async {
    if (!_deviceCodeLoginPending) return;
    _deviceCodeLoginPending = false;
    await _prefs.remove(_kDeviceCodeLoginPending);
    notifyListeners();
  }

  /// 稳定设备码：首次生成 UUID 并落盘；作为 Go `deviceInfo.deviceId` / open_id。
  Future<String> getOrCreateDeviceId() async {
    final existing = storedDeviceId;
    if (existing != null) {
      return existing;
    }
    final id = const Uuid().v4();
    await setDeviceId(id);
    return id;
  }

  /// 覆盖本地设备码（切换账号前调用）。
  Future<void> setDeviceId(String deviceId) async {
    final id = deviceId.trim();
    await _prefs.setString(_kDeviceId, id);
    notifyListeners();
  }
}
