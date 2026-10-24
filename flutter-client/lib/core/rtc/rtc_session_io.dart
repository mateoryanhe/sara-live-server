import 'dart:async';

import 'package:agora_rtc_engine/agora_rtc_engine.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

import 'rtc_video_platform.dart';

void _logAgoraRtc(String message) {
  debugPrint('[AgoraRtc] $message');
}

/// Android / 桌面：Agora RTC Engine。
class RtcSession implements RtcVideoPlatform {
  RtcEngineEx? _engine;
  final ValueNotifier<int?> remoteUid = ValueNotifier<int?>(null);
  final ValueNotifier<bool> _remoteVideoReady = ValueNotifier<bool>(false);
  bool _micMuted = false;
  bool _cameraEnabled = true;
  bool _beautyEnabled = false;
  double _beautyLevel = 0.65;

  @override
  ValueListenable<bool> get remoteVideoReady => _remoteVideoReady;

  void _trackRemoteUid(int? uid) {
    if (uid != null && uid > 0) {
      remoteUid.value = uid;
      _remoteVideoReady.value = true;
    }
  }

  RtcEngineEventHandler _eventHandler(Completer<void> joinCompleter) {
    return RtcEngineEventHandler(
      onJoinChannelSuccess: (connection, elapsed) {
        _logAgoraRtc(
          'onJoinChannelSuccess channel=${connection.channelId} uid=${connection.localUid} elapsed=${elapsed}ms',
        );
        if (!joinCompleter.isCompleted) {
          joinCompleter.complete();
        }
      },
      onError: (err, msg) {
        _logAgoraRtc('onError err=$err msg=$msg');
        if (!joinCompleter.isCompleted) {
          joinCompleter.completeError(
            StateError('Agora error $err${msg.isNotEmpty ? ': $msg' : ''}'),
          );
        }
      },
      onConnectionStateChanged: (connection, state, reason) {
        if (state != ConnectionStateType.connectionStateFailed) {
          return;
        }
        if (!joinCompleter.isCompleted) {
          joinCompleter.completeError(
            StateError('Agora connection failed ($reason)'),
          );
        }
      },
      onUserJoined: (connection, remoteUid, elapsed) {
        _trackRemoteUid(remoteUid);
      },
      onUserOffline: (connection, remoteUid, reason) {
        if (this.remoteUid.value == remoteUid) {
          this.remoteUid.value = null;
        }
      },
      onRemoteVideoStateChanged: (
        connection,
        remoteUid,
        state,
        reason,
        elapsed,
      ) {
        if (state == RemoteVideoState.remoteVideoStateStarting ||
            state == RemoteVideoState.remoteVideoStateDecoding) {
          _trackRemoteUid(remoteUid);
        }
      },
    );
  }

  Future<void> _awaitJoin(RtcEngineEx engine, Completer<void> joinCompleter) async {
    try {
      await joinCompleter.future.timeout(const Duration(seconds: 15));
    } catch (e) {
      try {
        await engine.leaveChannel();
        await engine.release();
      } catch (_) {}
      if (e is TimeoutException) {
        throw StateError('Agora join timeout');
      }
      rethrow;
    }
  }

  Future<void> _joinChannelInternal({
    required RtcEngineEx engine,
    required String token,
    required String channel,
    required String userAccount,
    required int uid,
    required ChannelMediaOptions options,
    required Completer<void> joinCompleter,
  }) async {
    final account = userAccount.trim();
    if (account.isNotEmpty) {
      await engine.joinChannelWithUserAccount(
        token: token,
        channelId: channel,
        userAccount: account,
        options: options,
      );
    } else {
      if (uid <= 0) {
        throw StateError('Agora uid/userAccount empty');
      }
      await engine.joinChannel(
        token: token,
        channelId: channel,
        uid: uid,
        options: options,
      );
    }
    await _awaitJoin(engine, joinCompleter);
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
    final profile = liveBroadcast
        ? ChannelProfileType.channelProfileLiveBroadcasting
        : ChannelProfileType.channelProfileCommunication;
    final engine = createAgoraRtcEngineEx();
    final joinCompleter = Completer<void>();
    await engine.initialize(
      RtcEngineContext(
        appId: appId,
        channelProfile: profile,
      ),
    );
    engine.registerEventHandler(_eventHandler(joinCompleter));
    await engine.enableVideo();
    await engine.enableAudio();
    await engine.setClientRole(role: ClientRoleType.clientRoleBroadcaster);
    await engine.startPreview();
    final options = ChannelMediaOptions(
      channelProfile: profile,
      clientRoleType: ClientRoleType.clientRoleBroadcaster,
      publishCameraTrack: true,
      publishMicrophoneTrack: true,
      autoSubscribeAudio: !liveBroadcast,
      autoSubscribeVideo: !liveBroadcast,
    );
    await _joinChannelInternal(
      engine: engine,
      token: token,
      channel: channel,
      userAccount: userAccount,
      uid: uid,
      options: options,
      joinCompleter: joinCompleter,
    );
    _engine = engine;
    await _applyPublishAVState(engine);
  }

  /// 秀场观众：订阅主播画面，不推流。
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
    const profile = ChannelProfileType.channelProfileLiveBroadcasting;
    final engine = createAgoraRtcEngineEx();
    final joinCompleter = Completer<void>();
    await engine.initialize(
      RtcEngineContext(
        appId: appId,
        channelProfile: profile,
      ),
    );
    engine.registerEventHandler(_eventHandler(joinCompleter));
    await engine.enableVideo();
    await engine.enableAudio();
    await engine.setClientRole(role: ClientRoleType.clientRoleAudience);
    const options = ChannelMediaOptions(
      channelProfile: profile,
      clientRoleType: ClientRoleType.clientRoleAudience,
      publishCameraTrack: false,
      publishMicrophoneTrack: false,
      autoSubscribeAudio: true,
      autoSubscribeVideo: true,
    );
    await _joinChannelInternal(
      engine: engine,
      token: token,
      channel: channel,
      userAccount: userAccount,
      uid: uid,
      options: options,
      joinCompleter: joinCompleter,
    );
    _engine = engine;
  }

  bool get micMuted => _micMuted;

  bool get cameraEnabled => _cameraEnabled;

  bool get beautyEnabled => _beautyEnabled;

  double get beautyLevel => _beautyLevel;

  Future<void> switchCamera() async {
    if (!_cameraEnabled) return;
    final engine = _engine;
    if (engine == null) return;
    await engine.switchCamera();
  }

  Future<void> setCameraEnabled(bool enabled) async {
    _cameraEnabled = enabled;
    final engine = _engine;
    if (engine == null) return;
    await _applyCameraPublishState(engine);
    try {
      await engine.setBeautyEffectOptions(
        enabled: _beautyEnabled && _cameraEnabled,
        options: BeautyOptions(
          lighteningContrastLevel:
              LighteningContrastLevel.lighteningContrastNormal,
          lighteningLevel: _beautyLevel,
          smoothnessLevel: _beautyLevel,
          rednessLevel: _beautyLevel * 0.35,
        ),
      );
    } catch (e) {
      _logAgoraRtc('beauty after camera toggle failed: $e');
    }
  }

  Future<void> setMicrophoneMuted(bool muted) async {
    final engine = _engine;
    if (engine == null) {
      _micMuted = muted;
      return;
    }
    await engine.muteLocalAudioStream(muted);
    _micMuted = muted;
  }

  Future<void> setBeautyEffect({
    required bool enabled,
    double? level,
  }) async {
    if (level != null) {
      _beautyLevel = level.clamp(0.0, 1.0);
    }
    _beautyEnabled = enabled;
    final engine = _engine;
    if (engine == null) return;
    await engine.setBeautyEffectOptions(
      enabled: enabled,
      options: BeautyOptions(
        lighteningContrastLevel:
            LighteningContrastLevel.lighteningContrastNormal,
        lighteningLevel: _beautyLevel,
        smoothnessLevel: _beautyLevel,
        rednessLevel: _beautyLevel * 0.35,
      ),
    );
  }

  Future<void> _applyCameraPublishState(RtcEngineEx engine) async {
    if (_cameraEnabled) {
      await engine.enableLocalVideo(true);
      await engine.startPreview();
      await engine.muteLocalVideoStream(false);
    } else {
      await engine.muteLocalVideoStream(true);
      await engine.enableLocalVideo(false);
      await engine.stopPreview();
    }
  }

  Future<void> _applyPublishAVState(RtcEngineEx engine) async {
    try {
      await engine.muteLocalAudioStream(_micMuted);
      await _applyCameraPublishState(engine);
      await engine.setBeautyEffectOptions(
        enabled: _beautyEnabled && _cameraEnabled,
        options: BeautyOptions(
          lighteningContrastLevel:
              LighteningContrastLevel.lighteningContrastNormal,
          lighteningLevel: _beautyLevel,
          smoothnessLevel: _beautyLevel,
          rednessLevel: _beautyLevel * 0.35,
        ),
      );
    } catch (e) {
      _logAgoraRtc('apply publish AV state failed: $e');
    }
  }

  Future<void> leave() async {
    remoteUid.value = null;
    _remoteVideoReady.value = false;
    final engine = _engine;
    _engine = null;
    if (engine == null) return;
    try {
      await engine.leaveChannel();
      await engine.stopPreview();
      await engine.release();
    } catch (_) {}
  }

  @override
  Widget localView() {
    final engine = _engine;
    if (engine == null) {
      return const SizedBox.expand(
        child: ColoredBox(color: Color(0xFF1A1A1A)),
      );
    }
    if (!_cameraEnabled) {
      return const SizedBox.expand(
        child: ColoredBox(
          color: Color(0xFF1A1A1A),
          child: Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.videocam_off_outlined, color: Colors.white54, size: 52),
                SizedBox(height: 10),
                Text(
                  '摄像头已关闭',
                  style: TextStyle(color: Colors.white54, fontSize: 14),
                ),
              ],
            ),
          ),
        ),
      );
    }
    return AgoraVideoView(
      controller: VideoViewController(
        rtcEngine: engine,
        canvas: const VideoCanvas(uid: 0),
      ),
    );
  }

  @override
  Widget remoteView() {
    return ValueListenableBuilder<int?>(
      valueListenable: remoteUid,
      builder: (context, uid, _) {
        final engine = _engine;
        if (engine == null || uid == null) {
          return const SizedBox.expand(
            child: ColoredBox(color: Color(0xFF0D0D0D)),
          );
        }
        return AgoraVideoView(
          controller: VideoViewController(
            rtcEngine: engine,
            canvas: VideoCanvas(uid: uid),
          ),
        );
      },
    );
  }

  @override
  Future<void> disposeViews() async => leave();
}
