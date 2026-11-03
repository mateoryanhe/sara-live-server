import 'dart:async';

import '../app/app_scope.dart';
import '../core/bus/app_bus_event.dart';
import '../core/http/api_exception.dart';
import '../core/ws/push_cmd.dart';
import '../core/ws/ws_client.dart';
import '../data/live_room/live_room_models.dart';
import '../ui/contracts/anchor_audience_ui.dart';

class AnchorAudienceViewModel {
  AnchorAudienceViewModel({
    required AppScope scope,
    required AnchorAudienceUi ui,
  })  : _scope = scope,
        _ui = ui {
    _scope.appBus.add(_onBus);
  }

  static const _pageSize = 30;

  final AppScope _scope;
  final AnchorAudienceUi _ui;

  final List<OnlineUserItem> _items = [];
  int _total = 0;
  int _page = 1;
  bool _hasMore = false;
  bool _loadingMore = false;
  bool _busy = false;

  int? get _anchorId {
    final raw = _scope.session.token?.userId;
    if (raw == null) return null;
    return int.tryParse(raw);
  }

  void _onBus(AppBusEvent event) {
    switch (event) {
      case WsConnectionBusEvent(:final state):
        if (state == WsConnectionState.connected) {
          unawaited(refresh(silent: true));
        }
      case WsBusEvent(:final cmd, :final data):
        if (cmd != PushCmd.liveRoomAudienceListRefresh) {
          return;
        }
        final anchorId = _anchorId;
        if (anchorId == null) return;
        final roomId = _roomIdFromPush(data);
        if (roomId != null && roomId != anchorId) return;
        unawaited(refresh(silent: true));
      default:
        break;
    }
  }

  int? _roomIdFromPush(Object? data) {
    if (data is! Map) return null;
    final raw = data['roomId'];
    if (raw is num) return raw.toInt();
    return int.tryParse('$raw');
  }

  Future<void> loadInitial() => refresh();

  Future<void> refresh({bool silent = false, bool forceRefresh = true}) async {
    if (_busy) return;
    _busy = true;
    if (!silent) {
      _ui.showAudienceLoading(initial: true);
    }
    try {
      final anchorId = _anchorId;
      if (anchorId == null || anchorId <= 0) {
        _ui.showAudienceError('未登录');
        return;
      }
      final page = await _scope.liveRoomRepository.getOnlineList(
        anchorId: anchorId,
        page: 1,
        pageSize: _pageSize,
        forceRefresh: forceRefresh,
      );
      _items
        ..clear()
        ..addAll(page.list);
      _total = page.total;
      _page = page.page;
      _hasMore = page.hasMore;
      _ui.renderAudienceList(
        items: List.unmodifiable(_items),
        total: _total,
        loadingMore: false,
      );
    } on ApiException catch (e) {
      _ui.showAudienceError(e.message);
    } catch (e) {
      _ui.showAudienceError(e.toString());
    } finally {
      _busy = false;
      if (!silent) {
        _ui.hideAudienceLoading();
      }
    }
  }

  Future<void> loadMore() async {
    if (_busy || _loadingMore || !_hasMore) return;
    final anchorId = _anchorId;
    if (anchorId == null) return;

    _loadingMore = true;
    _ui.renderAudienceList(
      items: List.unmodifiable(_items),
      total: _total,
      loadingMore: true,
    );
    _busy = true;
    try {
      final page = await _scope.liveRoomRepository.getOnlineList(
        anchorId: anchorId,
        page: _page + 1,
        pageSize: _pageSize,
        forceRefresh: true,
      );
      _items.addAll(page.list);
      _total = page.total;
      _page = page.page;
      _hasMore = page.hasMore;
      _ui.renderAudienceList(
        items: List.unmodifiable(_items),
        total: _total,
        loadingMore: false,
      );
    } on ApiException catch (e) {
      _ui.showAudienceError(e.message);
    } finally {
      _loadingMore = false;
      _busy = false;
    }
  }

  void dispose() {
    _scope.appBus.remove(_onBus);
  }
}
