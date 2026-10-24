import 'dart:async';

import '../app/app_scope.dart';
import '../core/bus/app_bus_event.dart';
import '../core/http/api_exception.dart';
import '../core/ws/push_cmd.dart';
import '../data/vip/vip_models.dart';
import '../ui/contracts/vip_ui.dart';

/// VIP 详情：自行订阅推送，失效后由 HTTP 补偿刷 UI。
class VipViewModel {
  VipViewModel({
    required AppScope scope,
    required VipUi ui,
  })  : _scope = scope,
        _ui = ui {
    _scope.appBus.add(_onBus);
  }

  final AppScope _scope;
  final VipUi _ui;

  /// 推送后回调，供页面重置 PageView 等。
  void Function()? onPushRefresh;

  int? get _userId {
    final token = _scope.session.token;
    if (token == null) return null;
    return int.tryParse(token.userId);
  }

  void _onBus(AppBusEvent event) {
    if (event is! WsBusEvent) return;
    switch (event.cmd) {
      case PushCmd.vipLevel:
      case PushCmd.wsConnected:
        onPushRefresh?.call();
        unawaited(load(forceRefresh: true, silent: true));
        break;
    }
  }

  Future<void> load({bool forceRefresh = false, bool silent = false}) async {
    final userId = _userId;
    if (userId == null || userId <= 0) {
      _ui.showVipError('未登录');
      return;
    }

    if (!silent) {
      _ui.showVipLoading();
    }
    try {
      final results = await Future.wait([
        _scope.vipRepository.getVipDetail(
          userId: userId,
          forceRefresh: forceRefresh,
        ),
        _scope.vipRepository.getVipCfgList(
          userId: userId,
          forceRefresh: forceRefresh,
        ),
      ]);
      _ui.renderVip(
        detail: results[0] as VipDetail,
        cfgList: results[1] as List<VipCfgItem>,
      );
    } on ApiException catch (e) {
      _ui.showVipError(e.message);
    } catch (e) {
      _ui.showVipError(e.toString());
    } finally {
      if (!silent) {
        _ui.hideVipLoading();
      }
    }
  }

  void dispose() {
    _scope.appBus.remove(_onBus);
  }
}
