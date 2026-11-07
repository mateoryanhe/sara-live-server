import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../core/cache/gift_effect_cache.dart';
import '../../data/gift/gift_media_url.dart';
import 'live_room_gift_effect_job.dart';

/// 直播间礼物特效：cmd=8 入队，定时消费；优先播放已缓存资源。
class LiveRoomGiftEffectPlayer extends ChangeNotifier {
  LiveRoomGiftEffectPlayer({required this.apiBaseUrl});

  final String apiBaseUrl;

  final List<LiveRoomGiftEffectJob> _queue = [];
  LiveRoomGiftEffectPlayState? _current;
  Timer? _pumpTimer;
  bool _downloadBusy = false;
  int? _boundRoomId;

  LiveRoomGiftEffectPlayState? get current => _current;

  int get pendingCount => _queue.length;

  void bindRoom(int? roomId) {
    if (_boundRoomId == roomId) return;
    _boundRoomId = roomId;
    _queue.clear();
    _current = null;
    _downloadBusy = false;
    _stopPumpTimer();
    notifyListeners();
  }

  void enqueueFromPush(Map<String, dynamic> map, int roomId) {
    if (_boundRoomId != null && _boundRoomId != roomId) return;
    final job = LiveRoomGiftEffectJob.fromPushMap(map, roomId, apiBaseUrl);
    // 无 giftAnim 不播全屏（避免把 icon 当特效；测服如柠檬茶/魔法棒）。
    if (job.effectUrl.isEmpty) return;
    _queue.add(job);
    unawaited(GiftEffectCache.instance.prefetch(job.playUrl));
    if (job.iconUrl.isNotEmpty) {
      unawaited(GiftEffectCache.instance.prefetch(job.iconUrl));
    }
    _ensurePumpTimer();
    unawaited(_pump());
    notifyListeners();
  }

  /// 预缓存礼物列表里的 animation（打开礼物面板时调用）。
  Future<void> prefetchGiftAnimations(Iterable<String> urls) async {
    for (final raw in urls) {
      final url = resolveGiftMediaUrl(raw, apiBaseUrl);
      if (url.isEmpty) continue;
      unawaited(GiftEffectCache.instance.prefetch(url));
    }
  }

  void completeCurrent() {
    _current = null;
    notifyListeners();
    unawaited(_pump());
  }

  void _ensurePumpTimer() {
    if (_pumpTimer != null) return;
    _pumpTimer = Timer.periodic(
      const Duration(milliseconds: 450),
      (_) => unawaited(_pump()),
    );
  }

  void _stopPumpTimer() {
    _pumpTimer?.cancel();
    _pumpTimer = null;
  }

  Future<void> _pump() async {
    if (_current != null) return;
    if (_queue.isEmpty) {
      _stopPumpTimer();
      return;
    }
    _ensurePumpTimer();

    final cachedIndex = await _indexFirstCached();
    if (cachedIndex != null) {
      final job = _queue.removeAt(cachedIndex);
      await _startPlay(job);
      return;
    }

    if (_downloadBusy) return;
    _downloadBusy = true;
    try {
      final job = _queue.first;
      final file = await GiftEffectCache.instance.getFile(job.playUrl);
      _queue.removeAt(0);
      await _startPlay(job, filePath: file.path);
    } catch (_) {
      if (_queue.isNotEmpty) {
        _queue.removeAt(0);
      }
    } finally {
      _downloadBusy = false;
    }
    if (_queue.isNotEmpty) {
      unawaited(_pump());
    }
  }

  Future<int?> _indexFirstCached() async {
    for (var i = 0; i < _queue.length; i++) {
      if (await GiftEffectCache.instance.isCached(_queue[i].playUrl)) {
        return i;
      }
    }
    return null;
  }

  Future<void> _startPlay(
    LiveRoomGiftEffectJob job, {
    String? filePath,
  }) async {
    final url = job.playUrl;
    final path = filePath ?? (await GiftEffectCache.instance.getFile(url)).path;
    _current = LiveRoomGiftEffectPlayState(
      job: job,
      localPath: path,
      mediaKind: giftEffectMediaKindForUrl(url),
    );
    notifyListeners();
  }

  @override
  void dispose() {
    _stopPumpTimer();
    super.dispose();
  }
}
