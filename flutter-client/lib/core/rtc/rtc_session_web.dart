import 'dart:js_interop';
import 'dart:js_interop_unsafe';
import 'dart:ui_web' as ui_web;

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:web/web.dart' as web;

import 'rtc_video_platform.dart';

@JS('saraAgoraRtc.join')
external JSPromise _jsJoin(JSAny options);

@JS('saraAgoraRtc.leave')
external JSPromise _jsLeave();

/// Web：通过 `web/agora_rtc_bridge.js` + Agora Web SDK。
class RtcSession implements RtcVideoPlatform {
  static int _viewSeq = 0;
  late final String _localViewType;
  late final String _remoteViewType;
  late final String _localElementId;
  late final String _remoteElementId;
  bool _joined = false;
  final ValueNotifier<bool> _remoteVideoReady = ValueNotifier<bool>(false);

  @override
  ValueListenable<bool> get remoteVideoReady => _remoteVideoReady;

  RtcSession() {
    final seq = ++_viewSeq;
    _localViewType = 'sara-agora-local-$seq';
    _remoteViewType = 'sara-agora-remote-$seq';
    _localElementId = 'sara-agora-local-el-$seq';
    _remoteElementId = 'sara-agora-remote-el-$seq';
    _registerView(_localViewType, _localElementId);
    _registerView(_remoteViewType, _remoteElementId);
  }

  void _registerView(String viewType, String elementId) {
    ui_web.platformViewRegistry.registerViewFactory(viewType, (int _) {
      final div = web.HTMLDivElement()
        ..id = elementId
        ..style.width = '100%'
        ..style.height = '100%'
        ..style.backgroundColor = 'transparent';
      return div;
    });
  }

  void _bindRemoteVideoHook() {
    try {
      void onRemote() {
        notifyRemoteVideoReady();
      }

      web.window.setProperty('saraAgoraRtcOnRemoteVideo'.toJS, onRemote.toJS);
    } catch (e, st) {
      debugPrint('saraAgoraRtcOnRemoteVideo hook failed: $e\n$st');
    }
  }

  Future<void> join({
    required String appId,
    required String channel,
    required String token,
    required int uid,
    String userAccount = '',
    bool liveBroadcast = false,
  }) async {
    await leave();
    if (appId.isEmpty) {
      throw StateError('Agora appId empty');
    }
    final options = {
      'appId': appId,
      'channel': channel,
      'token': token,
      'uid': uid,
      'userAccount': userAccount,
      'localViewId': _localElementId,
      'remoteViewId': _remoteElementId,
    }.jsify();
    _bindRemoteVideoHook();
    await _jsJoin(options!).toDart;
    _joined = true;
    _remoteVideoReady.value = false;
  }

  Future<void> joinSubscribe({
    required String appId,
    required String channel,
    required String token,
    required int uid,
    String userAccount = '',
  }) async {
    await leave();
    if (appId.isEmpty) {
      throw StateError('Agora appId empty');
    }
    final options = {
      'appId': appId,
      'channel': channel,
      'token': token,
      'uid': uid,
      'userAccount': userAccount,
      'localViewId': _localElementId,
      'remoteViewId': _remoteElementId,
      'role': 'audience',
    }.jsify();
    _bindRemoteVideoHook();
    await _jsJoin(options!).toDart;
    _joined = true;
    _remoteVideoReady.value = false;
  }

  void notifyRemoteVideoReady() {
    if (_joined) {
      _remoteVideoReady.value = true;
    }
  }

  Future<void> leave() async {
    if (!_joined) return;
    _joined = false;
    _remoteVideoReady.value = false;
    try {
      await _jsLeave().toDart;
    } catch (_) {}
  }

  bool get micMuted => false;

  bool get cameraEnabled => true;

  bool get beautyEnabled => false;

  double get beautyLevel => 0.65;

  Future<void> switchCamera() async {}

  Future<void> setCameraEnabled(bool enabled) async {}

  Future<void> setMicrophoneMuted(bool muted) async {}

  Future<void> setBeautyEffect({
    required bool enabled,
    double? level,
  }) async {}

  @override
  Widget localView() => SizedBox.expand(
        child: HtmlElementView(viewType: _localViewType),
      );

  @override
  Widget remoteView() => SizedBox.expand(
        child: HtmlElementView(viewType: _remoteViewType),
      );

  @override
  Future<void> disposeViews() async => leave();
}
