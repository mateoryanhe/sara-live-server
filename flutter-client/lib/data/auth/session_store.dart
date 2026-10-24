import 'package:flutter/foundation.dart';
import 'package:sara_live_app/core/http/auth_token.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart';

/// 本地会话：deviceId 持久化 + 登录 token。
class SessionStore extends ChangeNotifier {
  SessionStore(this._prefs);

  static const _kDeviceId = 'device_id';
  static const _kAuthToken = 'auth_token';

  final SharedPreferences _prefs;
  AuthToken? _token;

  AuthToken? get token => _token;

  bool get isLoggedIn => _token != null;

  /// 已保存的设备码（未保存过则为 null）。
  String? get storedDeviceId {
    final raw = _prefs.getString(_kDeviceId);
    if (raw == null) return null;
    final trimmed = raw.trim();
    return trimmed.isEmpty ? null : trimmed;
  }

  Future<void> load() async {
    _token = AuthToken.tryParse(_prefs.getString(_kAuthToken));
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
