import 'dart:async';

import '../app/app_scope.dart';
import '../app/server_endpoint_store.dart';
import '../core/bus/app_bus_event.dart';
import '../core/http/api_exception.dart';
import '../core/ws/push_cmd.dart';
import '../ui/contracts/mine_ui.dart';

/// 「我的」：资料、VIP 摘要、改昵称/头像、退出。
/// 自行订阅 [AppBus]：VIP 推送 / WS 重连后静默刷新 UI。
class MineViewModel {
  MineViewModel({
    required AppScope scope,
    required MineUi ui,
  })  : _scope = scope,
        _ui = ui {
    _scope.appBus.add(_onBus);
  }

  final AppScope _scope;
  final MineUi _ui;

  /// 防止连按保存导致并发切换（会互相顶掉 token / deviceLogin）。
  bool _deviceIdSwitchInFlight = false;

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
        unawaited(load(forceRefresh: true, silent: true));
        break;
    }
  }

  Future<void> load({bool forceRefresh = false, bool silent = false}) async {
    final userId = _userId;
    if (userId == null || userId <= 0) {
      _ui.showMineError('未登录');
      _ui.onLoggedOut();
      return;
    }

    if (!silent) {
      _ui.showMineLoading();
    }
    try {
      final profileFuture = _scope.userRepository.getProfile(
        userId: userId,
        forceRefresh: forceRefresh,
      );
      final vipFuture = _scope.vipRepository.getVipDetail(
        userId: userId,
        forceRefresh: forceRefresh,
      );

      final profile = await profileFuture;
      _ui.renderProfile(profile);

      try {
        final vip = await vipFuture;
        _ui.renderVipDetail(vip);
      } on ApiException catch (e) {
        _ui.showMineError(e.message);
        if (e.isUnauthorized) {
          await logout();
        }
      } catch (e) {
        _ui.showMineError(e.toString());
      }
    } on ApiException catch (e) {
      _ui.showMineError(e.message);
      if (e.isUnauthorized) {
        await logout();
      }
    } catch (e) {
      _ui.showMineError(e.toString());
    } finally {
      if (!silent) {
        _ui.hideMineLoading();
      }
    }
  }

  Future<void> updateNickname(String nickname) async {
    final userId = _userId;
    if (userId == null || userId <= 0) {
      _ui.showMineError('未登录');
      return;
    }
    final trimmed = nickname.trim();
    if (trimmed.isEmpty) {
      _ui.showMineError('昵称不能为空');
      return;
    }

    _ui.showMineLoading();
    try {
      final profile = await _scope.userRepository.updateNickname(
        userId: userId,
        nickname: trimmed,
      );
      _ui.renderProfile(profile);
      _ui.showMineMessage('昵称已更新');
    } on ApiException catch (e) {
      _ui.showMineError(e.message);
      if (e.isUnauthorized) {
        await logout();
      }
    } catch (e) {
      _ui.showMineError(e.toString());
    } finally {
      _ui.hideMineLoading();
    }
  }

  Future<void> updateAvatar({
    required List<int> bytes,
    required String filename,
  }) async {
    final userId = _userId;
    if (userId == null || userId <= 0) {
      _ui.showMineError('未登录');
      return;
    }
    if (bytes.isEmpty) {
      _ui.showMineError('请选择图片');
      return;
    }
    if (bytes.length > 1024 * 1024) {
      _ui.showMineError('图片过大，请选择 1MB 以内的图片');
      return;
    }

    _ui.showMineLoading();
    try {
      final profile = await _scope.userRepository.updateAvatar(
        userId: userId,
        bytes: bytes,
        filename: filename,
      );
      _ui.renderProfile(profile);
      _ui.showMineMessage('头像已更新');
    } on ApiException catch (e) {
      _ui.showMineError(e.message);
      if (e.isUnauthorized) {
        await logout();
      }
    } catch (e) {
      _ui.showMineError(e.toString());
    } finally {
      _ui.hideMineLoading();
    }
  }

  /// 写入新设备码并走 `/auth/deviceLogin`（等同快捷登录另一账号）。
  ///
  /// [validationError] 由 UI 传入本地化校验文案；为 null 表示校验通过。
  Future<void> applyDeviceIdAndLogin(
    String deviceId, {
    String? validationError,
    required String successMessage,
    required String unchangedMessage,
  }) async {
    if (validationError != null) {
      _ui.showMineError(validationError);
      return;
    }
    final next = deviceId.trim();

    final current = _scope.session.storedDeviceId ??
        await _scope.session.getOrCreateDeviceId();
    if (next == current) {
      _ui.showMineMessage(unchangedMessage);
      return;
    }

    if (_deviceIdSwitchInFlight) {
      return;
    }
    _deviceIdSwitchInFlight = true;
    final previousDeviceId = current;

    _ui.showMineLoading();
    try {
      await _scope.runAccountSwitch(() async {
        // 先落盘新码，避免并发 Splash deviceLogin 仍用旧 UUID。
        await _scope.session.setDeviceId(next);
        await _scope.prepareForAccountSwitch();
        final deviceInfo =
            await _scope.deviceInfoBuilder.build(deviceId: next);
        final result = await _scope.authApi.deviceLogin(deviceInfo);
        await _scope.session.saveToken(result.token);
        await _scope.handleLoggedIn();
        await load(forceRefresh: true, silent: true);
        _ui.onDeviceIdApplied(next);
        _ui.showMineMessage(successMessage);
      });
    } on ApiException catch (e) {
      _ui.showMineError(e.message);
      await _restoreDeviceSession(previousDeviceId);
    } catch (e) {
      _ui.showMineError(e.toString());
      await _restoreDeviceSession(previousDeviceId);
    } finally {
      _deviceIdSwitchInFlight = false;
      _ui.hideMineLoading();
    }
  }

  Future<void> _restoreDeviceSession(String deviceId) async {
    try {
      await _scope.runAccountSwitch(() async {
        await _scope.session.setDeviceId(deviceId);
        final deviceInfo =
            await _scope.deviceInfoBuilder.build(deviceId: deviceId);
        final result = await _scope.authApi.deviceLogin(deviceInfo);
        await _scope.session.saveToken(result.token);
        await _scope.handleLoggedIn();
        _ui.onDeviceIdApplied(deviceId);
        await load(forceRefresh: true, silent: true);
      });
    } catch (_) {
      await _scope.handleLoggedOut();
    }
  }

  /// 保存 API 根地址或恢复编译期默认；会清登录态并回到启动页重新 deviceLogin。
  Future<void> applyApiBaseUrl({
    required String input,
    required bool resetToDefault,
    String? validationError,
    required String successMessage,
    required String unchangedMessage,
  }) async {
    if (validationError != null) {
      _ui.showMineError(validationError);
      return;
    }

    String? nextOverride;
    if (resetToDefault) {
      nextOverride = null;
    } else {
      final v = ServerEndpointStore.validateInput(input);
      switch (v) {
        case ApiBaseValidationOk(:final url):
          nextOverride = url;
        case ApiBaseValidationEmpty():
        case ApiBaseValidationInvalid():
          return;
      }
    }

    final store = _scope.serverEndpoint;
    final unchanged = resetToDefault
        ? !store.hasOverride
        : store.hasOverride && store.storedOverride == nextOverride;
    if (unchanged) {
      _ui.showMineMessage(unchangedMessage);
      return;
    }

    _ui.showMineLoading();
    try {
      await _scope.applyServerEndpointOverride(nextOverride);
      _ui.onApiBaseApplied(store.apiBaseUrl);
      _ui.showMineMessage(successMessage);
      _ui.onLoggedOut();
    } catch (e) {
      _ui.showMineError(e.toString());
    } finally {
      _ui.hideMineLoading();
    }
  }

  Future<void> logout() async {
    await _scope.handleLoggedOut();
    _ui.onLoggedOut();
  }

  void dispose() {
    _scope.appBus.remove(_onBus);
  }
}
