import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/media/call_media_permission.dart';
import '../../core/rtc/rtc_session.dart';
import '../../data/agora/agora_api.dart';
import '../../data/live_room/live_room_api.dart';
import 'anchor_live_publish_prefs.dart';

/// 秀场主播推流（独立 [RtcSession]，不与 1v1 [CallController] 共用实例）。
class ShowLivePublisher extends ChangeNotifier {
  ShowLivePublisher({
    required AgoraApi agoraApi,
    required LiveRoomApi liveRoomApi,
    required AnchorLivePublishPrefs publishPrefs,
  })  : _agoraApi = agoraApi,
        _liveRoomApi = liveRoomApi,
        _publishPrefs = publishPrefs {
    _loadPrefsFromDisk();
  }

  final AgoraApi _agoraApi;
  final LiveRoomApi _liveRoomApi;
  final AnchorLivePublishPrefs _publishPrefs;
  final RtcSession _rtc = RtcSession();

  static const _liveHeartPeriod = Duration(seconds: 30);

  bool _publishing = false;
  String? _lastError;
  Timer? _liveHeartTimer;
  int? _liveHeartRoomId;
  Future<void>? _startPublishFuture;
  bool _micMuted = false;
  bool _cameraEnabled = true;
  bool _beautyEnabled = true;
  double _beautyLevel = 0.65;

  void _loadPrefsFromDisk() {
    _micMuted = _publishPrefs.micMuted;
    _cameraEnabled = _publishPrefs.cameraEnabled;
    _beautyEnabled = _publishPrefs.beautyEnabled;
    _beautyLevel = _publishPrefs.beautyLevel;
  }

  void _persistPrefs() {
    unawaited(
      _publishPrefs.save(
        micMuted: _micMuted,
        cameraEnabled: _cameraEnabled,
        beautyEnabled: _beautyEnabled,
        beautyLevel: _beautyLevel,
      ),
    );
  }

  bool get isPublishing => _publishing;

  String? get lastError => _lastError;

  bool get micMuted => _micMuted;

  bool get cameraEnabled => _cameraEnabled;

  bool get beautyEnabled => _beautyEnabled;

  double get beautyLevel => _beautyLevel;

  RtcSession get rtcSession => _rtc;

  /// 频道名与线上一致：`roomId` = 主播 userId 字符串。
  Future<void> startPublish(int anchorId) async {
    if (_publishing) return;
    final inFlight = _startPublishFuture;
    if (inFlight != null) {
      await inFlight;
      return;
    }
    final task = _startPublishInner(anchorId);
    _startPublishFuture = task;
    try {
      await task;
    } finally {
      if (identical(_startPublishFuture, task)) {
        _startPublishFuture = null;
      }
    }
  }

  Future<void> _startPublishInner(int anchorId) async {
    if (_publishing) return;
    _lastError = null;
    final ok = await ensureCallMediaPermission();
    if (!ok) {
      _lastError = '需要摄像头与麦克风权限';
      notifyListeners();
      throw StateError(_lastError!);
    }
    final channelName = '$anchorId';
    final creds = await _agoraApi.liveRoomToken(
      channelName: channelName,
      role: AgoraRtcRole.publisher,
    );
    if (creds.appId.isEmpty || creds.token.isEmpty) {
      _lastError = '声网 Token 无效，请检查服务端声网配置';
      notifyListeners();
      throw StateError(_lastError!);
    }
    await _rtc.setMicrophoneMuted(_micMuted);
    await _rtc.setCameraEnabled(_cameraEnabled);
    await _rtc.setBeautyEffect(
      enabled: _beautyEnabled,
      level: _beautyLevel,
    );
    await _rtc.join(
      appId: creds.appId,
      channel: creds.channelName.isNotEmpty ? creds.channelName : channelName,
      token: creds.token,
      uid: creds.rtcUid,
      userAccount: creds.userAccount,
      liveBroadcast: true,
    );
    _publishing = true;
    _micMuted = _rtc.micMuted;
    _cameraEnabled = _rtc.cameraEnabled;
    _beautyEnabled = _rtc.beautyEnabled;
    _beautyLevel = _rtc.beautyLevel;
    _startLiveHeart(anchorId);
    notifyListeners();
  }

  Future<void> switchCamera() async {
    if (!_publishing) return;
    await _rtc.switchCamera();
  }

  Future<void> toggleMicrophoneMuted() async {
    final next = !_micMuted;
    if (_publishing) {
      await _rtc.setMicrophoneMuted(next);
      _micMuted = _rtc.micMuted;
    } else {
      _micMuted = next;
    }
    _persistPrefs();
    notifyListeners();
  }

  Future<void> toggleCameraEnabled() async {
    final next = !_cameraEnabled;
    if (_publishing) {
      await _rtc.setCameraEnabled(next);
      _cameraEnabled = _rtc.cameraEnabled;
    } else {
      _cameraEnabled = next;
    }
    _persistPrefs();
    notifyListeners();
  }

  Future<void> setBeautyEnabled(bool enabled) async {
    if (_publishing) {
      await _rtc.setBeautyEffect(enabled: enabled, level: _beautyLevel);
      _beautyEnabled = _rtc.beautyEnabled;
    } else {
      _beautyEnabled = enabled;
    }
    _persistPrefs();
    notifyListeners();
  }

  Future<void> setBeautyLevel(double level) async {
    _beautyLevel = level.clamp(0.0, 1.0);
    if (_publishing && _beautyEnabled) {
      await _rtc.setBeautyEffect(enabled: true, level: _beautyLevel);
      _beautyLevel = _rtc.beautyLevel;
    }
    _persistPrefs();
    notifyListeners();
  }

  Future<void> stopPublish() async {
    _stopLiveHeart();
    if (!_publishing) {
      await _rtc.leave();
      return;
    }
    await _rtc.leave();
    _publishing = false;
    notifyListeners();
  }

  void _startLiveHeart(int roomId) {
    _stopLiveHeart();
    _liveHeartRoomId = roomId;
    unawaited(_tickLiveHeart());
    _liveHeartTimer = Timer.periodic(_liveHeartPeriod, (_) {
      unawaited(_tickLiveHeart());
    });
  }

  void _stopLiveHeart() {
    _liveHeartTimer?.cancel();
    _liveHeartTimer = null;
    _liveHeartRoomId = null;
  }

  Future<void> _tickLiveHeart() async {
    final roomId = _liveHeartRoomId;
    if (roomId == null || !_publishing) return;
    try {
      await _liveRoomApi.reportLiveStartStatus(roomId: roomId);
    } catch (_) {}
  }

  @override
  void dispose() {
    _stopLiveHeart();
    unawaited(_rtc.disposeViews());
    super.dispose();
  }
}
