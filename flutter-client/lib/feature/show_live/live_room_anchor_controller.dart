import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/bus/app_bus.dart';
import '../../core/bus/app_bus_event.dart';
import '../../core/ws/push_cmd.dart';
import '../../data/call/call_models.dart';
import '../../data/live_follow/live_follow_api.dart';
import '../../data/live_room/live_room_api.dart';
import '../../data/live_room/live_room_models.dart';
import '../../data/live_room/live_room_repository.dart';

/// 主播开播会话：公屏 WS、在线人数、观众列表与房管操作。
class LiveRoomAnchorController extends ChangeNotifier {
  LiveRoomAnchorController({
    required LiveRoomApi liveRoomApi,
    required LiveRoomRepository liveRoomRepository,
    required LiveFollowApi liveFollowApi,
    required AppBus appBus,
  })  : _liveRoomApi = liveRoomApi,
        _liveRoomRepository = liveRoomRepository,
        _liveFollowApi = liveFollowApi,
        _appBus = appBus {
    _appBus.add(_onBus);
  }

  final LiveRoomApi _liveRoomApi;
  final LiveRoomRepository _liveRoomRepository;
  final LiveFollowApi _liveFollowApi;
  final AppBus _appBus;

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
    notifyListeners();
    await refreshOnlineList();
  }

  void exitSession() {
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
    if (event is! WsBusEvent || !_active || _roomId == null) return;
    final roomKey = '${_roomId!}';
    final map = pushDataMap(event.data);
    if (map == null) return;

    switch (event.cmd) {
      case PushCmd.liveRoomChat:
        if ('${map['roomId'] ?? ''}' != roomKey) return;
        _appendChat(
          senderId: '${map['senderId'] ?? ''}',
          senderName: '${map['senderName'] ?? ''}',
          content: '${map['content'] ?? ''}',
          vipLevel: _jsonInt(map['vipLevel']),
        );
      case PushCmd.liveRoomGift:
        if ('${map['roomId'] ?? ''}' != roomKey) return;
        _appendChat(
          senderId: '${map['senderId'] ?? ''}',
          senderName: '${map['senderName'] ?? ''}',
          content: '',
          vipLevel: _jsonInt(map['vipLevel']),
          isGift: true,
          giftSummary:
              '送出 ${map['giftName'] ?? '礼物'} x${_jsonInt(map['count'])}',
        );
      case PushCmd.liveRoomAudienceListRefresh:
        if ('${map['roomId'] ?? ''}' != roomKey) return;
        _onlineCount = _jsonInt(map['onlineCount']);
        _applyOnlineListFromPush(map['list']);
        notifyListeners();
      case PushCmd.wsConnected:
        unawaited(refreshOnlineList(silent: true));
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
    required String content,
    int vipLevel = 0,
    bool isGift = false,
    String? giftSummary,
  }) {
    if (senderName.isEmpty && content.isEmpty && giftSummary == null) return;
    _chatLines.add(
      LiveRoomChatLine(
        senderId: senderId,
        senderName: senderName,
        content: content,
        vipLevel: vipLevel,
        isGift: isGift,
        giftSummary: giftSummary,
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
