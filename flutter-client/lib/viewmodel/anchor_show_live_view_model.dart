import 'dart:async';

import '../app/app_scope.dart';
import '../core/http/api_exception.dart';
import '../data/live_room/live_room_models.dart';
import '../ui/contracts/anchor_show_live_ui.dart';

class AnchorShowLiveViewModel {
  AnchorShowLiveViewModel({
    required AppScope scope,
    required AnchorShowLiveUi ui,
  })  : _scope = scope,
        _ui = ui {
    _scope.showLivePublisher.addListener(_onPublisherChanged);
  }

  final AppScope _scope;
  final AnchorShowLiveUi _ui;

  LiveRoomDetail? _room;
  bool _loading = false;
  bool _busy = false;

  int? get _anchorId {
    final raw = _scope.session.token?.userId;
    if (raw == null) return null;
    return int.tryParse(raw);
  }

  void _onPublisherChanged() {
    _ui.renderShowLive(
      room: _room,
      busy: _busy,
      publishing: _scope.showLivePublisher.isPublishing,
    );
  }

  Future<void> _resumePublishIfLive() async {
    final anchorId = _anchorId;
    if (anchorId == null || _room == null || !_room!.isLive) return;
    if (_scope.showLivePublisher.isPublishing) return;
    if (_scope.callController.isActive) return;
    try {
      await _scope.showLivePublisher.startPublish(anchorId);
    } catch (e) {
      _ui.showShowLiveError('推流恢复失败：$e');
    }
  }

  Future<void> load() async {
    if (_loading) return;
    _loading = true;
    _ui.showShowLiveLoading(true);
    try {
      final anchorId = _anchorId;
      if (anchorId == null) {
        _ui.showShowLiveError('未登录');
        return;
      }
      _room = await _scope.liveRoomRepository.getLiveRoom(
        anchorId: anchorId,
        forceRefresh: true,
      );
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
      await _resumePublishIfLive();
    } on ApiException catch (e) {
      _ui.showShowLiveError(e.message);
    } catch (e) {
      _ui.showShowLiveError(e.toString());
    } finally {
      _loading = false;
      _ui.showShowLiveLoading(false);
    }
  }

  Future<void> createRoom({
    required String title,
    required String notice,
    List<int>? coverBytes,
    String? coverFilename,
  }) async {
    if (_busy) return;
    final anchorId = _anchorId;
    if (anchorId == null) {
      _ui.showShowLiveError('未登录');
      return;
    }
    if (title.trim().isEmpty) {
      _ui.showShowLiveError('请填写直播间标题');
      return;
    }
    _busy = true;
    _ui.renderShowLive(
      room: _room,
      busy: true,
      publishing: _scope.showLivePublisher.isPublishing,
    );
    try {
      _room = await _scope.liveRoomRepository.ensureShowLiveRoom(
        anchorId: anchorId,
        title: title.trim(),
        notice: notice.trim(),
        coverBytes: coverBytes,
        coverFilename: coverFilename,
      );
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
      _ui.showShowLiveMessage('秀场直播间已创建');
    } on ApiException catch (e) {
      _ui.showShowLiveError(e.message);
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } catch (e) {
      _ui.showShowLiveError(e.toString());
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } finally {
      _busy = false;
    }
  }

  Future<void> saveProfile({
    required String title,
    required String notice,
    List<int>? coverBytes,
    String? coverFilename,
  }) async {
    if (_busy) return;
    final anchorId = _anchorId;
    if (anchorId == null) {
      _ui.showShowLiveError('未登录');
      return;
    }
    if (_room == null) {
      _ui.showShowLiveError('请先创建秀场直播间');
      return;
    }
    if (title.trim().isEmpty) {
      _ui.showShowLiveError('请填写直播间标题');
      return;
    }
    _busy = true;
    _ui.renderShowLive(
      room: _room,
      busy: true,
      publishing: _scope.showLivePublisher.isPublishing,
    );
    try {
      _room = await _scope.liveRoomRepository.saveShowLiveProfile(
        anchorId: anchorId,
        title: title.trim(),
        notice: notice.trim(),
        coverBytes: coverBytes,
        coverFilename: coverFilename,
      );
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
      _ui.showShowLiveMessage('已保存');
    } on ApiException catch (e) {
      _ui.showShowLiveError(e.message);
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } catch (e) {
      _ui.showShowLiveError(e.toString());
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } finally {
      _busy = false;
    }
  }

  Future<void> startLive({
    required String title,
    required String notice,
    List<int>? coverBytes,
    String? coverFilename,
  }) async {
    if (_busy) return;
    final anchorId = _anchorId;
    if (anchorId == null) {
      _ui.showShowLiveError('未登录');
      return;
    }
    if (_room == null) {
      _ui.showShowLiveError('请先创建秀场直播间');
      return;
    }
    if (title.trim().isEmpty) {
      _ui.showShowLiveError('请填写直播间标题');
      return;
    }
    if (_scope.callController.isActive) {
      _ui.showShowLiveError('通话中无法开播');
      return;
    }
    _busy = true;
    _ui.renderShowLive(
      room: _room,
      busy: true,
      publishing: _scope.showLivePublisher.isPublishing,
    );
    try {
      _room = await _scope.liveRoomRepository.startLive(
        anchorId: anchorId,
        title: title.trim(),
        notice: notice.trim(),
        coverBytes: coverBytes,
        coverFilename: coverFilename,
      );
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
      if (_room != null) {
        _ui.openShowLiveBroadcast(_room!);
      }
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
      _ui.showShowLiveMessage('已开播');
    } on ApiException catch (e) {
      _ui.showShowLiveError(e.message);
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } catch (e) {
      _ui.showShowLiveError(e.toString());
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } finally {
      _busy = false;
    }
  }

  Future<void> stopLive() async {
    if (_busy) return;
    final anchorId = _anchorId;
    if (anchorId == null) {
      _ui.showShowLiveError('未登录');
      return;
    }
    _busy = true;
    _ui.renderShowLive(
      room: _room,
      busy: true,
      publishing: _scope.showLivePublisher.isPublishing,
    );
    try {
      await _scope.showLivePublisher.stopPublish();
      _room = await _scope.liveRoomRepository.stopLive(anchorId: anchorId);
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: false,
      );
      _ui.showShowLiveMessage('已下播');
    } on ApiException catch (e) {
      _ui.showShowLiveError(e.message);
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } catch (e) {
      _ui.showShowLiveError(e.toString());
      _ui.renderShowLive(
        room: _room,
        busy: false,
        publishing: _scope.showLivePublisher.isPublishing,
      );
    } finally {
      _busy = false;
    }
  }

  void dispose() {
    _scope.showLivePublisher.removeListener(_onPublisherChanged);
  }
}
