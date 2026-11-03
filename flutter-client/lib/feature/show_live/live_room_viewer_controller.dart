import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/http/api_exception.dart';
import '../../core/http/auth_token.dart';
import '../../core/rtc/rtc_session.dart';
import '../../core/ws/push_cmd.dart';
import '../../core/ws/ws_client.dart';
import '../../data/agora/agora_api.dart';
import '../../data/call/call_models.dart';
import '../../data/live_room/live_room_api.dart';
import '../../data/live_room/live_room_models.dart';
import '../../data/user/user_repository.dart';
import 'live_room_gift_effect_player.dart';
import 'live_room_heart_reporter.dart';

/// 观众进房：join HTTP + 声网 Subscriber + 公屏 WS。
///
/// 自行订阅 [AppBus]：`cmd=8/9/22/23` 本房公屏；`cmd=10/11` 钱包（与 [UserRepository] 并列）。
class LiveRoomViewerController extends ChangeNotifier {
  LiveRoomViewerController({
    required LiveRoomApi liveRoomApi,
    required AgoraApi agoraApi,
    required UserRepository userRepository,
    required AuthTokenProvider tokenProvider,
    required AppBus appBus,
    required LiveRoomGiftEffectPlayer giftEffectPlayer,
  })  : _liveRoomApi = liveRoomApi,
        _agoraApi = agoraApi,
        _userRepository = userRepository,
        _tokenProvider = tokenProvider,
        _appBus = appBus,
        _giftEffectPlayer = giftEffectPlayer,
        _liveHeart = LiveRoomHeartReporter(liveRoomApi) {
    _appBus.add(_onBus);
  }

  final LiveRoomApi _liveRoomApi;
  final AgoraApi _agoraApi;
  final UserRepository _userRepository;
  final AuthTokenProvider _tokenProvider;
  final AppBus _appBus;
  final LiveRoomGiftEffectPlayer _giftEffectPlayer;
  final RtcSession _rtc = RtcSession();
  final LiveRoomHeartReporter _liveHeart;

  int? _roomId;
  String _roomTitle = '';
  bool _watching = false;
  bool _anchorStopped = false;
  final List<LiveRoomChatLine> _chatLines = [];
  double _walletGold = 0;
  double _walletDiamond = 0;

  int? get roomId => _roomId;

  String get roomTitle => _roomTitle;

  bool get isWatching => _watching;

  bool get anchorStopped => _anchorStopped;

  List<LiveRoomChatLine> get chatLines => List.unmodifiable(_chatLines);

  double get walletGold => _walletGold;

  double get walletDiamond => _walletDiamond;

  RtcSession get rtcSession => _rtc;

  int? _currentUserId() {
    final id = int.tryParse(_tokenProvider()?.userId ?? '');
    if (id == null || id <= 0) return null;
    return id;
  }

  Future<void> _syncWalletFromPush(WsBusEvent event) async {
    final userId = _currentUserId();
    if (userId == null) return;
    try {
      final profile = await _userRepository.applyWalletBalancePush(
        userId: userId,
        cmd: event.cmd,
        data: event.data,
      );
      if (profile == null) return;
      _walletGold = profile.gold;
      _walletDiamond = profile.diamond;
      notifyListeners();
    } catch (_) {}
  }

  Future<void> _syncWalletFromProfileCache() async {
    final userId = _currentUserId();
    if (userId == null) return;
    try {
      final profile = await _userRepository.getProfile(
        userId: userId,
        forceRefresh: false,
      );
      _walletGold = profile.gold;
      _walletDiamond = profile.diamond;
      notifyListeners();
    } catch (_) {}
  }

  Future<void> _syncWalletAfterWsOnline() async {
    final userId = _currentUserId();
    if (userId == null) return;
    try {
      final profile = await _userRepository.refreshWalletFromServer(userId);
      if (profile == null) return;
      _walletGold = profile.gold;
      _walletDiamond = profile.diamond;
      notifyListeners();
    } catch (_) {}
  }

  void _onBus(AppBusEvent event) {
    switch (event) {
      case AppLifecycleBusEvent(:final kind):
        if (kind == AppLifecycleKind.resumed) {
          unawaited(_syncWalletAfterWsOnline());
        }
      case WsConnectionBusEvent(:final state):
        if (state != WsConnectionState.connected) {
          return;
        }
        if (_watching && _roomId != null) {
          unawaited(_rejoinAfterWsReconnect());
        }
        unawaited(_syncWalletAfterWsOnline());
        return;
      case WsBusEvent(:final cmd, :final data):
        _onWsPush(cmd, data);
      default:
        break;
    }
  }

  void _onWsPush(int cmd, Object? data) {
    switch (cmd) {
      case PushCmd.diamond:
      case PushCmd.gold:
      case PushCmd.rechargeSuccess:
        unawaited(_syncWalletFromPush(WsBusEvent(cmd: cmd, data: data)));
        return;
      default:
        break;
    }

    final roomId = _roomId;
    if (roomId == null || !_watching) return;
    final map = pushDataMap(data);
    if (map == null) return;

    switch (cmd) {
      case PushCmd.liveRoomChat:
      case PushCmd.liveRoomPaidDanmaku:
        if (!liveRoomPushMatchesRoom(map, roomId)) return;
        _appendChat(
          senderId: '${map['senderId'] ?? ''}',
          senderName: '${map['senderName'] ?? ''}',
          senderAvatar: '${map['senderAvatar'] ?? ''}',
          content: '${map['content'] ?? ''}',
          vipLevel: _jsonInt(map['vipLevel']),
        );
      case PushCmd.liveRoomGift:
        if (!liveRoomPushMatchesRoom(map, roomId)) return;
        _giftEffectPlayer.enqueueFromPush(map, roomId);
        final name = '${map['senderName'] ?? ''}';
        final gift = '${map['giftName'] ?? '礼物'}';
        _appendChat(
          senderId: '${map['senderId'] ?? ''}',
          senderName: name,
          senderAvatar: '${map['senderAvatar'] ?? ''}',
          content: '',
          vipLevel: _jsonInt(map['vipLevel']),
          isGift: true,
          giftName: gift,
          giftIcon: '${map['giftIcon'] ?? ''}',
          giftCount: _jsonInt(map['count']),
        );
      case PushCmd.liveRoomStopLive:
        if (!liveRoomPushMatchesRoom(map, roomId) &&
            '${map['anchorId'] ?? ''}' != '$roomId') {
          return;
        }
        _anchorStopped = true;
        notifyListeners();
      default:
        break;
    }
  }

  Future<void> _rejoinAfterWsReconnect() async {
    final roomId = _roomId;
    if (roomId == null || !_watching) return;
    try {
      final join = await _liveRoomApi.joinRoom(roomId: roomId);
      if (join.kickBanned) {
        _anchorStopped = true;
        notifyListeners();
      }
    } catch (_) {}
  }

  void _appendChat({
    String senderId = '',
    required String senderName,
    String senderAvatar = '',
    required String content,
    int vipLevel = 0,
    bool isGift = false,
    String giftName = '',
    String giftIcon = '',
    int giftCount = 0,
  }) {
    if (isGift) {
      if (giftName.isEmpty && giftIcon.isEmpty) return;
    } else if (senderName.isEmpty && content.isEmpty) {
      return;
    }
    if (!isGift &&
        _chatLines.isNotEmpty &&
        senderId.isNotEmpty &&
        _chatLines.last.senderId == senderId &&
        _chatLines.last.content == content &&
        !_chatLines.last.isGift) {
      return;
    }
    _chatLines.add(
      LiveRoomChatLine(
        senderId: senderId,
        senderName: senderName,
        senderAvatar: senderAvatar,
        content: content,
        vipLevel: vipLevel,
        isGift: isGift,
        giftName: giftName,
        giftIcon: giftIcon,
        giftCount: giftCount > 0 ? giftCount : 1,
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
    _giftEffectPlayer.bindRoom(roomId);

    var joined = false;
    try {
      final join = await _liveRoomApi.joinRoom(roomId: roomId);
      if (join.kickBanned) {
        final remain = join.kickRemainSeconds;
        final hint = remain > 0
            ? '您已被踢出，${remain}s 内无法进入'
            : '您已被踢出，暂时无法进入';
        throw ApiException(code: 0, message: hint);
      }
      joined = true;
      _liveHeart.start(roomId);

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
      unawaited(_syncWalletFromProfileCache());
      notifyListeners();
    } catch (e) {
      if (joined) {
        await leaveRoom();
      } else {
        _roomId = null;
        _roomTitle = '';
        _watching = false;
      }
      rethrow;
    }
  }

  Future<void> sendChat(String content) async {
    final roomId = _roomId;
    if (roomId == null || !_watching) return;
    final text = content.trim();
    if (text.isEmpty) return;

    final userId = _currentUserId();
    if (userId != null) {
      try {
        final profile = await _userRepository.getProfile(
          userId: userId,
          forceRefresh: false,
        );
        _appendChat(
          senderId: '$userId',
          senderName: profile.nickname,
          senderAvatar: profile.avatar,
          content: text,
          vipLevel: profile.vipLevel,
        );
      } catch (_) {
        _appendChat(
          senderId: '$userId',
          senderName: '',
          content: text,
        );
      }
    }

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
    _giftEffectPlayer.bindRoom(null);
    _liveHeart.stop();
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
