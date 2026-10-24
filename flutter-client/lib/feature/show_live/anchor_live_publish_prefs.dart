import 'package:shared_preferences/shared_preferences.dart';

/// 主播直播工具栏偏好（本机记住麦克风/摄像头/美颜）。
class AnchorLivePublishPrefs {
  AnchorLivePublishPrefs(this._prefs);

  final SharedPreferences _prefs;

  static const _micMuted = 'anchor_live_mic_muted';
  static const _cameraEnabled = 'anchor_live_camera_enabled';
  static const _beautyEnabled = 'anchor_live_beauty_enabled';
  static const _beautyLevel = 'anchor_live_beauty_level';

  bool get micMuted => _prefs.getBool(_micMuted) ?? false;

  bool get cameraEnabled => _prefs.getBool(_cameraEnabled) ?? true;

  bool get beautyEnabled => _prefs.getBool(_beautyEnabled) ?? true;

  double get beautyLevel => _prefs.getDouble(_beautyLevel) ?? 0.65;

  Future<void> save({
    required bool micMuted,
    required bool cameraEnabled,
    required bool beautyEnabled,
    required double beautyLevel,
  }) async {
    await _prefs.setBool(_micMuted, micMuted);
    await _prefs.setBool(_cameraEnabled, cameraEnabled);
    await _prefs.setBool(_beautyEnabled, beautyEnabled);
    await _prefs.setDouble(_beautyLevel, beautyLevel.clamp(0.0, 1.0));
  }
}
