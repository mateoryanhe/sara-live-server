import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/ws/push_cmd.dart';
import '../../core/ws/ws_client.dart';
import '../../data/call/call_models.dart';
import '../../data/live_follow/live_follow_api.dart';
import '../../data/live_room/live_room_api.dart';
import '../../data/live_room/live_room_models.dart';
import '../../data/live_room/live_room_repository.dart';
import 'live_room_gift_effect_player.dart';
import 'live_room_heart_reporter.dart';

/// 主播开播会话：公屏 WS、在线人数、观众列表与房管操作。
///
/// 自行订阅 [AppBus]：`cmd=8/9/22/24` 等本房推送（与观众端并列，互不调用）。
class LiveRoomAnchorController extends ChangeNotifier {
  LiveRoomAnchorController({
    required LiveRoomApi liveRoomApi,
    required LiveRoomRepository liveRoomRepository,
    required LiveFollowApi liveFollowApi,
    required AppBus appBus,
    required LiveRoomGiftEffectPlayer giftEffectPlayer,
  })  : _liveRoomApi = liveRoomApi,
        _liveRoomRepository = liveRoomRepository,
        _liveFollowApi = liveFollowApi,
        _appBus = appBus,
        _giftEffectPlayer = giftEffectPlayer,
        _liveHeart = LiveRoomHeartReporter(liveRoomApi) {
    _appBus.add(_onBus);
  }

  final LiveRoomApi _liveRoomApi;
  final LiveRoomRepository _liveRoomRepository;
  final LiveFollowApi _liveFollowApi;
  final AppBus _appBus;
  final LiveRoomGiftEffectPlayer _giftEffectPlayer;
  final LiveRoomHeartReporter _liveHeart;

  int? _roomId;
  String _title = '';
  String _notice = '';
  bool _active = false;
  int _onlineCount = 0;
  final List<LiveRoomChatLine> _chatLines = [];
  final List<OnlineUserItem> _onlinePreview = [];
  bool _loadingOnline = false;

  int? get roomId => _roomId;

  String get title => _title;

  String get notice => _notice;

  bool get isActive => _active;

  int get onlineCount => _onlineCount;

  List<LiveRoomChatLine> get chatLines => List.unmodifiable(_chatLines);

  List<OnlineUserItem> get onlinePreview => List.unmodifiable(_onlinePreview);

  bool get loadingOnline => _loadingOnline;

  Future<void> enterSession({
    required int roomId,
    required String title,
    required String notice,
    required int initialOnlineCount,
  }) async {
    _roomId = roomId;
    _title = title;
    _notice = notice;
    _onlineCount = initialOnlineCount;
    _active = true;
    _chatLines.clear();
    _giftEffectPlayer.bindRoom(roomId);
    _liveHeart.start(roomId);
    notifyListeners();
    await refreshOnlineList();
  }

  void exitSession() {
    _liveHeart.stop();
    _giftEffectPlayer.bindRoom(null);
    _active = false;
    _roomId = null;
    _title = '';
    _notice = '';
    _onlineCount = 0;
    _chatLines.clear();
    _onlinePreview.clear();
    notifyListeners();
  }

  void _onBus(AppBusEvent event) {
    switch (event) {
      case WsConnectionBusEvent(:final state):
        if (state == WsConnectionState.connected && _active && _roomId != null) {
          unawaited(refreshOnlineList(silent: true));
        }
        return;
      case WsBusEvent(:final cmd, :final data):
        if (!_active || _roomId == null) return;
        _onRoomPush(cmd, data);
      default:
        break;
    }
  }

  void _onRoomPush(int cmd, Object? data) {
    final roomId = _roomId!;
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
        _appendChat(
          senderId: '${map['senderId'] ?? ''}',
          senderName: '${map['senderName'] ?? ''}',
          senderAvatar: '${map['senderAvatar'] ?? ''}',
          content: '',
          vipLevel: _jsonInt(map['vipLevel']),
          isGift: true,
          giftName: '${map['giftName'] ?? '礼物'}',
          giftIcon: '${map['giftIcon'] ?? ''}',
          giftCount: _jsonInt(map['count']),
        );
      case PushCmd.liveRoomAudienceListRefresh:
        if (!liveRoomPushMatchesRoom(map, roomId)) return;
        _onlineCount = _jsonInt(map['onlineCount']);
        _applyOnlineListFromPush(map['list']);
        notifyListeners();
      default:
        break;
    }
  }

  void _applyOnlineListFromPush(Object? raw) {
    if (raw is! List) return;
    final items = <OnlineUserItem>[];
    for (final e in raw) {
      if (e is Map) {
        items.add(OnlineUserItem.fromJson(Map<String, dynamic>.from(e)));
      }
    }
    if (items.isEmpty) return;
    _onlinePreview
      ..clear()
      ..addAll(items.take(8));
  }

  void _appendChat({
    required String senderId,
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
    if (_chatLines.length > 300) {
      _chatLines.removeRange(0, _chatLines.length - 300);
    }
    notifyListeners();
  }

  Future<void> sendChat(String content) async {
    final roomId = _roomId;
    if (roomId == null) return;
    final text = content.trim();
    if (text.isEmpty) return;
    await _liveRoomApi.sendChat(roomId: roomId, content: text);
  }

  Future<OnlineUserListPage> fetchOnlinePage({
    required int page,
    int pageSize = 30,
  }) async {
    final roomId = _roomId;
    if (roomId == null) {
      throw StateError('未在开播会话中');
    }
    return _liveRoomRepository.getOnlineList(
      anchorId: roomId,
      page: page,
      pageSize: pageSize,
      forceRefresh: true,
    );
  }

  Future<void> refreshOnlineList({bool silent = false}) async {
    final roomId = _roomId;
    if (roomId == null) return;
    if (!silent) {
      _loadingOnline = true;
      notifyListeners();
    }
    try {
      final page = await _liveRoomRepository.getOnlineList(
        anchorId: roomId,
        page: 1,
        pageSize: 30,
        forceRefresh: true,
      );
      _onlineCount = page.total;
      _onlinePreview
        ..clear()
        ..addAll(page.list.take(8));
    } finally {
      if (!silent) {
        _loadingOnline = false;
        notifyListeners();
      } else {
        notifyListeners();
      }
    }
  }

  Future<void> kickUser(int userId) async {
    final roomId = _roomId;
    if (roomId == null) return;
    await _liveRoomRepository.kickAudience(anchorId: roomId, userId: userId);
    await refreshOnlineList(silent: true);
  }

  Future<void> setMute(int userId, bool muted) async {
    final roomId = _roomId;
    if (roomId == null) return;
    await _liveRoomRepository.setAudienceMute(
      anchorId: roomId,
      userId: userId,
      muted: muted,
    );
    await refreshOnlineList(silent: true);
  }

  Future<void> cancelKickBan(int userId) async {
    final roomId = _roomId;
    if (roomId == null) return;
    await _liveRoomRepository.cancelKickBan(anchorId: roomId, userId: userId);
  }

  Future<void> blockUser(int userId) async {
    await _liveFollowApi.blockUser(targetId: userId);
    final roomId = _roomId;
    if (roomId != null) {
      try {
        await _liveRoomRepository.kickAudience(anchorId: roomId, userId: userId);
      } catch (_) {}
      await refreshOnlineList(silent: true);
    }
  }

  Future<AudienceRestrictStatus> restrictStatus(int userId) async {
    final roomId = _roomId;
    if (roomId == null) {
      throw StateError('未在开播会话中');
    }
    return _liveRoomRepository.getAudienceRestrictStatus(
      roomId: roomId,
      userId: userId,
    );
  }

  int _jsonInt(Object? v) {
    if (v is num) return v.toInt();
    return int.tryParse('$v') ?? 0;
  }

  @override
  void dispose() {
    _appBus.remove(_onBus);
    exitSession();
    super.dispose();
  }
}
