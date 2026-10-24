import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/rtc/rtc_session.dart';
import '../../core/ws/push_cmd.dart';
import '../../data/agora/agora_api.dart';
import '../../data/call/call_models.dart';
import '../../data/live_room/live_room_api.dart';
import '../../data/live_room/live_room_models.dart';

/// 观众进房：join HTTP + 声网 Subscriber + 公屏 WS。
class LiveRoomViewerController extends ChangeNotifier {
  LiveRoomViewerController({
    required LiveRoomApi liveRoomApi,
    required AgoraApi agoraApi,
    required AppBus appBus,
  })  : _liveRoomApi = liveRoomApi,
        _agoraApi = agoraApi,
        _appBus = appBus {
    _appBus.add(_onBus);
  }

  final LiveRoomApi _liveRoomApi;
  final AgoraApi _agoraApi;
  final AppBus _appBus;
  final RtcSession _rtc = RtcSession();

  static const _liveHeartPeriod = Duration(seconds: 30);

  Timer? _liveHeartTimer;

  int? _roomId;
  String _roomTitle = '';
  bool _watching = false;
  bool _anchorStopped = false;
  final List<LiveRoomChatLine> _chatLines = [];

  int? get roomId => _roomId;

  String get roomTitle => _roomTitle;

  bool get isWatching => _watching;

  bool get anchorStopped => _anchorStopped;

  List<LiveRoomChatLine> get chatLines => List.unmodifiable(_chatLines);

  RtcSession get rtcSession => _rtc;

  void _onBus(AppBusEvent event) {
    if (event is! WsBusEvent) return;
    if (_roomId == null) return;
    final roomKey = '${_roomId!}';
    final map = pushDataMap(event.data);
    if (map == null) return;

    switch (event.cmd) {
      case PushCmd.liveRoomChat:
        if ('${map['roomId'] ?? ''}' != roomKey) return;
        _appendChat(
          senderName: '${map['senderName'] ?? ''}',
          content: '${map['content'] ?? ''}',
          vipLevel: _jsonInt(map['vipLevel']),
        );
      case PushCmd.liveRoomGift:
        if ('${map['roomId'] ?? ''}' != roomKey) return;
        final name = '${map['senderName'] ?? ''}';
        final gift = '${map['giftName'] ?? '礼物'}';
        final count = _jsonInt(map['count']);
        _appendChat(
          senderName: name,
          content: '',
          vipLevel: _jsonInt(map['vipLevel']),
          isGift: true,
          giftSummary: '送出 $gift x$count',
        );
      case PushCmd.liveRoomStopLive:
        if ('${map['roomId'] ?? ''}' != roomKey &&
            '${map['anchorId'] ?? ''}' != roomKey) {
          return;
        }
        _anchorStopped = true;
        notifyListeners();
    }
  }

  void _appendChat({
    required String senderName,
    required String content,
    int vipLevel = 0,
    bool isGift = false,
    String? giftSummary,
  }) {
    if (senderName.isEmpty && content.isEmpty && giftSummary == null) return;
    _chatLines.add(
      LiveRoomChatLine(
        senderName: senderName,
        content: content,
        vipLevel: vipLevel,
        isGift: isGift,
        giftSummary: giftSummary,
      ),
    );
    if (_chatLines.length > 200) {
      _chatLines.removeRange(0, _chatLines.length - 200);
    }
    notifyListeners();
  }

  Future<void> enterRoom({
    required int roomId,
    required String title,
  }) async {
    if (_watching && _roomId == roomId) return;
    await leaveRoom();
    _roomId = roomId;
    _roomTitle = title;
    _anchorStopped = false;
    _chatLines.clear();

    await _liveRoomApi.joinRoom(roomId: roomId);

    final channelName = '$roomId';
    final creds = await _agoraApi.liveRoomToken(
      channelName: channelName,
      role: AgoraRtcRole.subscriber,
    );
    await _rtc.joinSubscribe(
      appId: creds.appId,
      channel: creds.channelName.isNotEmpty ? creds.channelName : channelName,
      token: creds.token,
      uid: creds.rtcUid,
      userAccount: creds.userAccount,
    );
    _watching = true;
    _startLiveHeart();
    notifyListeners();
  }

  void _startLiveHeart() {
    _stopLiveHeart();
    unawaited(_tickLiveHeart());
    _liveHeartTimer = Timer.periodic(_liveHeartPeriod, (_) {
      unawaited(_tickLiveHeart());
    });
  }

  void _stopLiveHeart() {
    _liveHeartTimer?.cancel();
    _liveHeartTimer = null;
  }

  Future<void> _tickLiveHeart() async {
    final roomId = _roomId;
    if (roomId == null || !_watching) return;
    try {
      await _liveRoomApi.reportLiveStartStatus(roomId: roomId);
    } catch (_) {}
  }

  Future<void> sendChat(String content) async {
    final roomId = _roomId;
    if (roomId == null) return;
    final text = content.trim();
    if (text.isEmpty) return;
    await _liveRoomApi.sendChat(roomId: roomId, content: text);
  }

  Future<SendGiftResult> sendGift({
    required int giftId,
    required int count,
  }) async {
    final roomId = _roomId;
    if (roomId == null) {
      throw StateError('未在直播间内');
    }
    return _liveRoomApi.sendGift(
      roomId: roomId,
      giftId: giftId,
      count: count,
    );
  }

  Future<void> leaveRoom() async {
    final roomId = _roomId;
    _watching = false;
    _stopLiveHeart();
    await _rtc.leave();
    if (roomId != null) {
      try {
        await _liveRoomApi.leaveRoom(roomId: roomId);
      } catch (_) {}
    }
    _roomId = null;
    _roomTitle = '';
    _anchorStopped = false;
    notifyListeners();
  }

  int _jsonInt(Object? v) {
    if (v is num) return v.toInt();
    return int.tryParse('$v') ?? 0;
  }

  @override
  void dispose() {
    _appBus.remove(_onBus);
    unawaited(leaveRoom());
    unawaited(_rtc.disposeViews());
    super.dispose();
  }
}
