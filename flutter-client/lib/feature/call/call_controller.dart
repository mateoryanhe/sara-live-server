import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/scheduler.dart';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/http/api_exception.dart';
import '../../core/media/call_media_permission.dart';
import '../../core/rtc/rtc_session.dart';
import '../../core/ws/push_cmd.dart';
import '../../data/call/call_api.dart';
import '../../data/call/call_models.dart';

enum CallPhase {
  idle,
  outgoingRinging,
  incomingRinging,
  connecting,
  inCall,
}

enum CallSide {
  caller,
  callee,
}

class CallPeerDisplay {
  const CallPeerDisplay({
    required this.userId,
    required this.nickname,
    required this.avatar,
  });

  final String userId;
  final String nickname;
  final String avatar;
}

/// 1v1 / 通话会话编排：HTTP + WS + 跨平台 RTC。
class CallController extends ChangeNotifier {
  CallController({
    required CallApi api,
    required AppBus appBus,
  })  : _api = api,
        _appBus = appBus {
    _appBus.add(_onBus);
  }

  final CallApi _api;
  final AppBus _appBus;
  final RtcSession _rtc = RtcSession();

  CallPhase phase = CallPhase.idle;
  CallSide? side;
  int source = CallSource.oneToOneRoom;
  String? orderId;
  CallRtcCredentials? rtc;
  CallPeerDisplay? peer;
  String statusLine = '';
  String? lastError;

  Timer? _dialTimeout;
  Timer? _heartTimer;
  Timer? _confirmTimer;
  bool _confirmLoopRunning = false;
  bool _teardownInProgress = false;

  bool get isActive => phase != CallPhase.idle;

  RtcSession get rtcSession => _rtc;

  void _onBus(AppBusEvent event) {
    if (event is! WsBusEvent) return;
    unawaited(_handlePush(event.cmd, event.data));
  }

  Future<void> startOneToOneCall({
    required int targetUserId,
    required String peerNickname,
    required String peerAvatar,
  }) async {
    if (isActive) {
      lastError = '当前已在通话中';
      notifyListeners();
      return;
    }
    final ok = await ensureCallMediaPermission();
    if (!ok) {
      lastError = '需要摄像头与麦克风权限';
      notifyListeners();
      return;
    }

    try {
      statusLine = '正在呼叫…';
      side = CallSide.caller;
      source = CallSource.oneToOneRoom;
      peer = CallPeerDisplay(
        userId: '$targetUserId',
        nickname: peerNickname,
        avatar: peerAvatar,
      );
      phase = CallPhase.outgoingRinging;
      notifyListeners();

      final creds = await _api.oneToOneRoomCall(targetId: targetUserId);
      orderId = creds.orderId;
      rtc = creds;
      await _joinRtc(creds);
      _armDialTimeout(creds.orderId);
      statusLine = '等待对方接听…';
      notifyListeners();
    } on ApiException catch (e) {
      await _resetSession(message: e.message);
    } catch (e) {
      await _resetSession(message: e.toString());
    }
  }

  Future<void> acceptIncoming() async {
    if (phase != CallPhase.incomingRinging || orderId == null) return;
    final ok = await ensureCallMediaPermission();
    if (!ok) {
      lastError = '需要摄像头与麦克风权限';
      notifyListeners();
      return;
    }
    try {
      statusLine = '接听中…';
      phase = CallPhase.connecting;
      notifyListeners();
      final creds = await _api.acceptCall(orderId: orderId!);
      rtc = creds;
      await _joinRtc(creds);
      await _beginConfirmHandshake();
    } on ApiException catch (e) {
      await _resetSession(message: e.message);
    } catch (e) {
      await _resetSession(message: e.toString());
    }
  }

  Future<void> rejectIncoming() async {
    final id = orderId;
    if (id == null) {
      await _resetSession();
      return;
    }
    try {
      await _api.rejectCall(orderId: id);
    } catch (_) {}
    await _resetSession(message: '已拒接');
  }

  Future<void> hangUp() async {
    if (_teardownInProgress) return;
    _teardownInProgress = true;
    final id = orderId;
    final currentPhase = phase;
    try {
      if (id != null) {
        if (currentPhase == CallPhase.outgoingRinging ||
            currentPhase == CallPhase.incomingRinging) {
          try {
            await _api.callTimeout(orderId: id);
          } catch (_) {
            try {
              await _api.endCall(orderId: id);
            } catch (_) {}
          }
        } else {
          try {
            await _api.endCall(orderId: id);
          } catch (_) {}
        }
      }
    } finally {
      await _resetSession(message: '通话已结束');
      _teardownInProgress = false;
    }
  }

  Future<void> _handlePush(int cmd, Object? data) async {
    final map = pushDataMap(data);
    if (map == null) return;

    switch (cmd) {
      case PushCmd.callRequest:
        await _onCallRequest(CallRequestPush.fromJson(map));
      case PushCmd.callAccepted:
        if (side == CallSide.caller && phase == CallPhase.outgoingRinging) {
          _dialTimeout?.cancel();
          statusLine = '对方已接听，建立连接…';
          phase = CallPhase.connecting;
          notifyListeners();
          await _beginConfirmHandshake();
        }
      case PushCmd.callStarted:
        _enterInCall('通话中');
      case PushCmd.callRejected:
        if (isActive) {
          await _resetSession(message: map['message']?.toString() ?? '对方拒接');
        }
      case PushCmd.callEnded:
        if (isActive) {
          await _resetSession(message: map['message']?.toString() ?? '通话结束');
        }
      case PushCmd.callTimeout:
        if (isActive) {
          await _resetSession(message: map['message']?.toString() ?? '无人接听');
        }
    }
  }

  Future<void> _onCallRequest(CallRequestPush push) async {
    if (isActive) return;
    if (push.source != CallSource.oneToOneRoom) return;

    side = CallSide.callee;
    source = push.source;
    orderId = push.orderId;
    peer = CallPeerDisplay(
      userId: push.callerId,
      nickname: push.callerNickname,
      avatar: push.callerAvatar,
    );
    rtc = CallRtcCredentials(
      orderId: push.orderId,
      appId: '',
      channelName: push.channelName,
      token: push.receiverToken,
      userAccount: '',
      expireAt: 0,
    );
    statusLine = push.message.isNotEmpty ? push.message : '视频通话邀请';
    phase = CallPhase.incomingRinging;
    notifyListeners();
  }

  Future<void> _waitForOverlayFrame() async {
    final completer = Completer<void>();
    SchedulerBinding.instance.addPostFrameCallback((_) {
      if (!completer.isCompleted) completer.complete();
    });
    return completer.future;
  }

  Future<void> _joinRtc(CallRtcCredentials creds) async {
    await _waitForOverlayFrame();
    await _rtc.join(
      appId: creds.appId,
      channel: creds.channelName,
      token: creds.token,
      uid: creds.rtcUid,
      userAccount: creds.userAccount,
    );
  }

  void _armDialTimeout(String id) {
    _dialTimeout?.cancel();
    _dialTimeout = Timer(const Duration(seconds: 45), () {
      unawaited(_onDialTimeout(id));
    });
  }

  Future<void> _onDialTimeout(String id) async {
    if (phase != CallPhase.outgoingRinging) return;
    try {
      await _api.callTimeout(orderId: id);
    } catch (_) {}
    await _resetSession(message: '呼叫超时');
  }

  Future<void> _beginConfirmHandshake() async {
    if (_confirmLoopRunning) return;
    _confirmLoopRunning = true;
    _confirmTimer?.cancel();
    _confirmTimer = Timer.periodic(const Duration(seconds: 1), (_) {
      unawaited(_confirmTick());
    });
    await _confirmTick();
  }

  Future<void> _confirmTick() async {
    final id = orderId;
    if (id == null || phase == CallPhase.idle || phase == CallPhase.inCall) {
      return;
    }
    try {
      await _api.confirmCall(orderId: id);
      final st = await _api.getCallConfirmStatus(orderId: id);
      if (st.callStarted) {
        _enterInCall('通话中');
      } else {
        statusLine = st.peerConfirmed ? '等待双方确认…' : '等待对方确认…';
        notifyListeners();
      }
    } catch (_) {}
  }

  void _enterInCall(String line) {
    if (phase == CallPhase.inCall) return;
    _dialTimeout?.cancel();
    _stopConfirmLoop();
    phase = CallPhase.inCall;
    statusLine = line;
    _startHeart();
    notifyListeners();
  }

  void _startHeart() {
    _heartTimer?.cancel();
    _heartTimer = Timer.periodic(const Duration(seconds: 10), (_) {
      final id = orderId;
      if (id == null || phase != CallPhase.inCall) return;
      unawaited(_api.callHeart(orderId: id).catchError((_) {}));
    });
  }

  void _stopConfirmLoop() {
    _confirmLoopRunning = false;
    _confirmTimer?.cancel();
    _confirmTimer = null;
  }

  Future<void> _resetSession({String? message}) async {
    _dialTimeout?.cancel();
    _heartTimer?.cancel();
    _stopConfirmLoop();
    await _rtc.disposeViews();
    phase = CallPhase.idle;
    side = null;
    orderId = null;
    rtc = null;
    peer = null;
    if (message != null) {
      statusLine = message;
      lastError = message;
    } else {
      statusLine = '';
    }
    notifyListeners();
  }

  @override
  void dispose() {
    _appBus.remove(_onBus);
    _dialTimeout?.cancel();
    _heartTimer?.cancel();
    _stopConfirmLoop();
    unawaited(_rtc.disposeViews());
    super.dispose();
  }
}
