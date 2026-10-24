import 'dart:async';

import '../app/app_scope.dart';
import '../core/bus/app_bus_event.dart';
import '../core/http/api_exception.dart';
import '../core/ws/push_cmd.dart';
import '../data/one_to_one/one_to_one_models.dart';
import '../ui/contracts/one_to_one_ui.dart';

class OneToOneViewModel {
  OneToOneViewModel({
    required AppScope scope,
    required OneToOneUi ui,
  })  : _scope = scope,
        _ui = ui {
    _scope.appBus.add(_onBus);
  }

  static const _pageSize = 20;

  final AppScope _scope;
  final OneToOneUi _ui;

  final List<OneToOneRoomItem> _items = [];
  int _total = 0;
  int _page = 1;
  int _statusFilter = OneToOneStatusFilter.all;
  bool _hasMore = false;
  bool _loadingMore = false;
  bool _busy = false;

  int? get _userId {
    final token = _scope.session.token;
    if (token == null) return null;
    return int.tryParse(token.userId);
  }

  void _onBus(AppBusEvent event) {
    if (event is! WsBusEvent) return;
    if (event.cmd == PushCmd.wsConnected) {
      unawaited(refresh(silent: true));
    }
  }

  Future<void> loadInitial() => refresh(forceRefresh: false);

  Future<void> refresh({bool silent = false, bool forceRefresh = true}) async {
    if (_busy) return;
    _busy = true;
    if (!silent) {
      _ui.showOneToOneLoading(initial: true);
    }
    try {
      final userId = _userId;
      if (userId == null || userId <= 0) {
        _ui.showOneToOneError('未登录');
        return;
      }
      final page = await _scope.oneToOneRepository.getRoomList(
        userId: userId,
        page: 1,
        pageSize: _pageSize,
        statusFilter: _statusFilter,
        forceRefresh: forceRefresh,
      );
      _items
        ..clear()
        ..addAll(page.list);
      _total = page.total;
      _page = page.page;
      _hasMore = page.hasMore;
      _ui.renderOneToOneList(
        items: List.unmodifiable(_items),
        total: _total,
        statusFilter: _statusFilter,
        hasMore: _hasMore,
        loadingMore: false,
      );
    } on ApiException catch (e) {
      _ui.showOneToOneError(e.message);
    } catch (e) {
      _ui.showOneToOneError(e.toString());
    } finally {
      _busy = false;
      if (!silent) {
        _ui.hideOneToOneLoading();
      }
    }
  }

  Future<void> setStatusFilter(int filter) async {
    if (_statusFilter == filter || _busy) return;
    _statusFilter = filter;
    await refresh();
  }

  Future<void> loadMore() async {
    if (_busy || _loadingMore || !_hasMore) return;
    final userId = _userId;
    if (userId == null || userId <= 0) return;

    _loadingMore = true;
    _ui.renderOneToOneList(
      items: List.unmodifiable(_items),
      total: _total,
      statusFilter: _statusFilter,
      hasMore: _hasMore,
      loadingMore: true,
    );

    _busy = true;
    try {
      final nextPage = _page + 1;
      final page = await _scope.oneToOneRepository.getRoomList(
        userId: userId,
        page: nextPage,
        pageSize: _pageSize,
        statusFilter: _statusFilter,
        forceRefresh: true,
      );
      _items.addAll(page.list);
      _total = page.total;
      _page = page.page;
      _hasMore = page.hasMore;
      _ui.renderOneToOneList(
        items: List.unmodifiable(_items),
        total: _total,
        statusFilter: _statusFilter,
        hasMore: _hasMore,
        loadingMore: false,
      );
    } on ApiException catch (e) {
      _ui.showOneToOneError(e.message);
    } catch (e) {
      _ui.showOneToOneError(e.toString());
    } finally {
      _loadingMore = false;
      _busy = false;
    }
  }

  void dispose() {
    _scope.appBus.remove(_onBus);
  }
}
