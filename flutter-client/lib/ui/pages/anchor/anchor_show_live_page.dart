import 'dart:async';
import 'dart:typed_data';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/cache/image_disk_cache.dart';
import 'package:sara_live_app/data/live_room/live_room_models.dart';
import 'package:sara_live_app/ui/contracts/anchor_show_live_ui.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/ui/pages/anchor/anchor_show_live_broadcast_page.dart';
import 'package:sara_live_app/viewmodel/anchor_show_live_view_model.dart';

/// 主播 Tab：秀场直播间创建与开播 / 下播。
class AnchorShowLivePage extends StatefulWidget {
  const AnchorShowLivePage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<AnchorShowLivePage> createState() => _AnchorShowLivePageState();
}

class _AnchorShowLivePageState extends State<AnchorShowLivePage>
    with AutomaticKeepAliveClientMixin
    implements AnchorShowLiveUi {
  late final AnchorShowLiveViewModel _vm;
  final _titleController = TextEditingController();
  final _noticeController = TextEditingController();

  bool _pageLoading = true;
  bool _busy = false;
  bool _fieldsDirty = false;
  LiveRoomDetail? _room;
  Uint8List? _pickedCoverBytes;
  String? _pickedCoverFilename;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _vm = AnchorShowLiveViewModel(scope: widget.scope, ui: this);
    WidgetsBinding.instance.addPostFrameCallback((_) => _vm.load());
  }

  @override
  void dispose() {
    _titleController.dispose();
    _noticeController.dispose();
    _vm.dispose();
    super.dispose();
  }

  void _applyRoomFields(LiveRoomDetail? room) {
    if (room == null) return;
    _titleController.text = room.title;
    _noticeController.text = room.notice;
  }

  Future<void> _pickCover() async {
    if (_busy) return;
    final picker = ImagePicker();
    final file = await picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 1920,
      maxHeight: 1920,
      imageQuality: 88,
    );
    if (file == null) return;
    final name = file.name.toLowerCase();
    final okExt = name.endsWith('.jpg') ||
        name.endsWith('.jpeg') ||
        name.endsWith('.png') ||
        name.endsWith('.webp') ||
        name.endsWith('.gif') ||
        name.endsWith('.bmp');
    if (!okExt) {
      showShowLiveError('仅支持 jpg/png/webp/gif/bmp');
      return;
    }
    final bytes = await file.readAsBytes();
    if (bytes.length > 1024 * 1024) {
      showShowLiveError('封面请小于 1MB');
      return;
    }
    var filename = file.name;
    if (!filename.contains('.')) {
      filename = '$filename.jpg';
    }
    if (!mounted) return;
    setState(() {
      _pickedCoverBytes = bytes;
      _pickedCoverFilename = filename;
    });
  }

  void _clearPickedCover() {
    setState(() {
      _pickedCoverBytes = null;
      _pickedCoverFilename = null;
    });
  }

  @override
  void showShowLiveLoading(bool loading) {
    if (!mounted) return;
    setState(() => _pageLoading = loading);
  }

  @override
  void renderShowLive({
    required LiveRoomDetail? room,
    required bool busy,
    required bool publishing,
  }) {
    if (!mounted) return;
    setState(() {
      _room = room;
      _busy = busy;
      _pageLoading = false;
      if (!busy) {
        _pickedCoverBytes = null;
        _pickedCoverFilename = null;
      }
    });
    if (!busy && !_fieldsDirty) {
      _applyRoomFields(room);
    }
  }

  @override
  void showShowLiveError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  void showShowLiveMessage(String message) {
    showShowLiveError(message);
  }

  Future<void> _openBroadcast(LiveRoomDetail room) async {
    final stopped = await Navigator.of(context).push<bool>(
      MaterialPageRoute(
        builder: (_) => AnchorShowLiveBroadcastPage(
          scope: widget.scope,
          room: room,
        ),
        fullscreenDialog: true,
      ),
    );
    if (!mounted) return;
    if (stopped == true) {
      await _vm.load();
    } else {
      await _vm.load();
    }
  }

  @override
  void openShowLiveBroadcast(LiveRoomDetail room) {
    unawaited(_openBroadcast(room));
  }

  Widget _coverPreview(ThemeData theme, ColorScheme cs) {
    const aspect = 16 / 9;
    final maxW = 480.0;
    final height = maxW / aspect;

    Widget child;
    if (_pickedCoverBytes != null) {
      child = Image.memory(
        _pickedCoverBytes!,
        width: maxW,
        height: height,
        fit: BoxFit.cover,
      );
    } else if (_room != null && _room!.cover.isNotEmpty) {
      child = CachedNetworkImage(
        imageUrl: _room!.cover,
        cacheManager: ImageDiskCache.instance.manager,
        width: maxW,
        height: height,
        fit: BoxFit.cover,
        placeholder: (_, __) => Container(
          color: cs.surfaceContainerHigh,
          child: const Center(child: CircularProgressIndicator(strokeWidth: 2)),
        ),
        errorWidget: (_, __, ___) => _coverPlaceholder(cs),
      );
    } else {
      child = _coverPlaceholder(cs);
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(
          '直播封面',
          style: theme.textTheme.titleSmall?.copyWith(
            fontWeight: FontWeight.w700,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '创建直播间或开始直播前可设置，观众在列表中看到此封面',
          style: theme.textTheme.bodySmall?.copyWith(
            color: cs.onSurface.withValues(alpha: 0.55),
          ),
        ),
        const SizedBox(height: 10),
        ClipRRect(
          borderRadius: BorderRadius.circular(12),
          child: SizedBox(
            width: maxW,
            height: height,
            child: Stack(
              fit: StackFit.expand,
              children: [
                child,
                if (!_busy && !(_room?.isLive ?? false))
                  Material(
                    color: Colors.black26,
                    child: InkWell(
                      onTap: _pickCover,
                      child: Center(
                        child: Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 16,
                            vertical: 10,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.black54,
                            borderRadius: BorderRadius.circular(999),
                          ),
                          child: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              const Icon(
                                Icons.photo_camera_outlined,
                                color: Colors.white,
                                size: 20,
                              ),
                              const SizedBox(width: 8),
                              Text(
                                _pickedCoverBytes != null ||
                                        (_room?.cover.isNotEmpty ?? false)
                                    ? '更换封面'
                                    : '选择封面',
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ),
        if (_pickedCoverBytes != null && !_busy) ...[
          const SizedBox(height: 8),
          Align(
            alignment: Alignment.centerRight,
            child: TextButton(
              onPressed: _clearPickedCover,
              child: const Text('取消本次选择'),
            ),
          ),
        ],
      ],
    );
  }

  Widget _coverPlaceholder(ColorScheme cs) {
    return Container(
      color: cs.surfaceContainerHigh,
      alignment: Alignment.center,
      child: Icon(
        Icons.image_outlined,
        size: 48,
        color: cs.onSurface.withValues(alpha: 0.35),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final hasRoom = _room != null;
    final isLive = _room?.isLive ?? false;

    return Scaffold(
      body: AppPageFrame(
        child: _pageLoading
            ? const Center(child: CircularProgressIndicator())
            : SingleChildScrollView(
                padding: const EdgeInsets.fromLTRB(16, 16, 16, 32),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Text(
                      '秀场直播',
                      style: theme.textTheme.headlineSmall?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                    const SizedBox(height: 8),
                    if (hasRoom)
                      Row(
                        children: [
                          Icon(
                            isLive ? Icons.circle : Icons.circle_outlined,
                            size: 12,
                            color: isLive
                                ? const Color(0xFF22C55E)
                                : cs.outline,
                          ),
                          const SizedBox(width: 6),
                          Text(
                            isLive
                                ? '直播中 · 在线 ${_room!.onlineCount}'
                                : '未开播',
                            style: theme.textTheme.titleSmall,
                          ),
                        ],
                      )
                    else
                      Text(
                        '尚未创建秀场直播间，先填写标题并创建。',
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: cs.onSurface.withValues(alpha: 0.65),
                        ),
                      ),
                    const SizedBox(height: 20),
                    if (isLive)
                      FilledButton.icon(
                        onPressed: _busy || _room == null
                            ? null
                            : () => _openBroadcast(_room!),
                        icon: const Icon(Icons.fullscreen),
                        label: const Text('进入直播间'),
                      ),
                    if (isLive) const SizedBox(height: 20),
                    if (!isLive) _coverPreview(theme, cs),
                    if (!isLive) const SizedBox(height: 20),
                    TextField(
                      controller: _titleController,
                      enabled: !_busy,
                      onChanged: (_) => _fieldsDirty = true,
                      decoration: const InputDecoration(
                        labelText: '直播间标题',
                        border: OutlineInputBorder(),
                      ),
                    ),
                    const SizedBox(height: 12),
                    TextField(
                      controller: _noticeController,
                      enabled: !_busy,
                      onChanged: (_) => _fieldsDirty = true,
                      decoration: const InputDecoration(
                        labelText: '公告（可选）',
                        border: OutlineInputBorder(),
                      ),
                      maxLines: 2,
                    ),
                    if (hasRoom)
                      Padding(
                        padding: const EdgeInsets.only(top: 8),
                        child: Text(
                          '标题与公告可随时修改并保存，与老 App 一致。',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: cs.onSurface.withValues(alpha: 0.55),
                          ),
                        ),
                      ),
                    const SizedBox(height: 24),
                    if (!hasRoom)
                      FilledButton(
                        onPressed: _busy
                            ? null
                            : () async {
                                _fieldsDirty = false;
                                await _vm.createRoom(
                                  title: _titleController.text,
                                  notice: _noticeController.text,
                                  coverBytes: _pickedCoverBytes,
                                  coverFilename: _pickedCoverFilename,
                                );
                              },
                        child: _busy
                            ? const SizedBox(
                                width: 22,
                                height: 22,
                                child:
                                    CircularProgressIndicator(strokeWidth: 2),
                              )
                            : const Text('创建秀场直播间'),
                      )
                    else ...[
                      OutlinedButton(
                        onPressed: _busy
                            ? null
                            : () async {
                                _fieldsDirty = false;
                                await _vm.saveProfile(
                                  title: _titleController.text,
                                  notice: _noticeController.text,
                                  coverBytes: _pickedCoverBytes,
                                  coverFilename: _pickedCoverFilename,
                                );
                              },
                        child: const Text('保存资料'),
                      ),
                      const SizedBox(height: 12),
                      if (!isLive)
                        FilledButton(
                          onPressed: _busy
                              ? null
                              : () async {
                                  _fieldsDirty = false;
                                  await _vm.startLive(
                                    title: _titleController.text,
                                    notice: _noticeController.text,
                                    coverBytes: _pickedCoverBytes,
                                    coverFilename: _pickedCoverFilename,
                                  );
                                },
                          child: const Text('开始直播'),
                        ),
                      if (!isLive) const SizedBox(height: 12),
                    ],
                  ],
                ),
              ),
      ),
    );
  }
}
