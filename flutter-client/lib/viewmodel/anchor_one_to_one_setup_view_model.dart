import 'dart:async';

import '../app/app_scope.dart';
import '../core/http/api_exception.dart';
import '../data/live_room/live_room_models.dart';
import '../ui/contracts/anchor_one_to_one_setup_ui.dart';

class AnchorOneToOneSetupViewModel {
  AnchorOneToOneSetupViewModel({
    required AppScope scope,
    required AnchorOneToOneSetupUi ui,
  })  : _scope = scope,
        _ui = ui;

  final AppScope _scope;
  final AnchorOneToOneSetupUi _ui;

  OneToOneRoomConfig? _config;
  List<BillingTierItem> _billingTiers = const [];
  bool _loading = false;
  bool _saving = false;

  int? get _anchorId {
    final raw = _scope.session.token?.userId;
    if (raw == null) return null;
    return int.tryParse(raw);
  }

  Future<void> load() async {
    if (_loading) return;
    _loading = true;
    _ui.showOneToOneSetupLoading(true);
    try {
      final anchorId = _anchorId;
      if (anchorId == null) {
        _ui.showOneToOneSetupError('未登录');
        return;
      }
      final results = await Future.wait([
        _scope.liveRoomRepository.getOneToOneConfig(
          anchorId: anchorId,
          forceRefresh: true,
        ),
        _scope.liveRoomRepository.getBillingTierList(
          anchorId: anchorId,
          forceRefresh: true,
        ),
      ]);
      _config = results[0] as OneToOneRoomConfig?;
      _billingTiers = results[1] as List<BillingTierItem>;
      _mergeLegacyBillingTier();
      _ui.renderOneToOneSetup(
        config: _config,
        billingTiers: _billingTiers,
        saving: false,
      );
    } on ApiException catch (e) {
      _ui.showOneToOneSetupError(e.message);
    } catch (e) {
      _ui.showOneToOneSetupError(e.toString());
    } finally {
      _loading = false;
      _ui.showOneToOneSetupLoading(false);
    }
  }

  Future<void> save({
    required String title,
    required double billing,
  }) async {
    if (_saving) return;
    final anchorId = _anchorId;
    if (anchorId == null) {
      _ui.showOneToOneSetupError('未登录');
      return;
    }
    final trimmed = title.trim();
    if (trimmed.isEmpty) {
      _ui.showOneToOneSetupError('请填写 1v1 房间标题');
      return;
    }
    if (_billingTiers.isNotEmpty &&
        !_billingTiers.any((t) => t.pricePerMinute == billing)) {
      _ui.showOneToOneSetupError('请从服务器档位中选择通话单价');
      return;
    }

    _saving = true;
    _ui.renderOneToOneSetup(
      config: _config,
      billingTiers: _billingTiers,
      saving: true,
    );
    try {
      _config = await _scope.liveRoomRepository.saveOneToOneConfig(
        anchorId: anchorId,
        billing: billing,
        title: trimmed,
      );
      _ui.renderOneToOneSetup(
        config: _config,
        billingTiers: _billingTiers,
        saving: false,
      );
      _ui.showOneToOneSetupMessage('1v1 房间已保存');
    } on ApiException catch (e) {
      _ui.showOneToOneSetupError(e.message);
      _ui.renderOneToOneSetup(
        config: _config,
        billingTiers: _billingTiers,
        saving: false,
      );
    } catch (e) {
      _ui.showOneToOneSetupError(e.toString());
      _ui.renderOneToOneSetup(
        config: _config,
        billingTiers: _billingTiers,
        saving: false,
      );
    } finally {
      _saving = false;
    }
  }

  void _mergeLegacyBillingTier() {
    final cfg = _config;
    if (cfg == null) return;
    if (_billingTiers.any((t) => t.pricePerMinute == cfg.billing)) {
      return;
    }
    _billingTiers = [
      ..._billingTiers,
      BillingTierItem(
        id: 'legacy',
        pricePerMinute: cfg.billing,
        sort: 9999,
      ),
    ];
  }

  void dispose() {}
}
