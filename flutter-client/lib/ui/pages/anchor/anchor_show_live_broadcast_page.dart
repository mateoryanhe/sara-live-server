import 'dart:async';

import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/http/api_exception.dart';
import 'package:sara_live_app/core/widgets/cached_avatar.dart';
import 'package:sara_live_app/data/live_room/live_room_models.dart';
import 'package:sara_live_app/feature/show_live/live_room_anchor_controller.dart';
import 'package:sara_live_app/ui/widgets/anchor_live_control_panel.dart'
    show AnchorLiveToolbar;
import 'package:sara_live_app/ui/widgets/live_room_cover_backdrop.dart';

/// 主播开播中的全屏「直播间」：预览、公屏、观众与房管。
class AnchorShowLiveBroadcastPage extends StatefulWidget {
  const AnchorShowLiveBroadcastPage({
    super.key,
    required this.scope,
    required this.room,
  });

  final AppScope scope;
  final LiveRoomDetail room;

  @override
  State<AnchorShowLiveBroadcastPage> createState() =>
      _AnchorShowLiveBroadcastPageState();
}

class _AnchorShowLiveBroadcastPageState
    extends State<AnchorShowLiveBroadcastPage> {
  final _chatController = TextEditingController();
  final _chatScroll = ScrollController();
  bool _stopping = false;
  String _coverUrl = '';
  String _avatarUrl = '';
  String _anchorNickname = '';

  LiveRoomAnchorController get _anchor => widget.scope.liveRoomAnchor;

  int get _roomId => int.tryParse(widget.room.roomId) ?? 0;

  @override
  void initState() {
    super.initState();
    _coverUrl = widget.room.cover;
    _anchor.addListener(_onAnchor);
    widget.scope.showLivePublisher.addListener(_onPublisher);
    unawaited(_bootstrapSession());
  }

  Future<void> _bootstrapSession() async {
    try {
      final profile = await widget.scope.userRepository.getProfile(
        userId: _roomId,
        forceRefresh: false,
      );
      if (mounted) {
        setState(() {
          if (profile.avatar.isNotEmpty) _avatarUrl = profile.avatar;
          if (profile.nickname.isNotEmpty) _anchorNickname = profile.nickname;
        });
      }
    } catch (_) {}
    try {
      final fresh = await widget.scope.liveRoomRepository.getLiveRoom(
        anchorId: _roomId,
        forceRefresh: true,
      );
      if (mounted && fresh != null && fresh.cover.isNotEmpty) {
        setState(() => _coverUrl = fresh.cover);
      }
    } catch (_) {}

    await _anchor.enterSession(
      roomId: _roomId,
      title: widget.room.title,
      notice: widget.room.notice,
      initialOnlineCount: widget.room.onlineCount,
    );
    final publisher = widget.scope.showLivePublisher;
    if (!publisher.isPublishing && widget.room.isLive && _roomId > 0) {
      try {
        await publisher.startPublish(_roomId);
      } catch (e) {
        if (mounted) {
          _snack(e is StateError ? e.message : '推流失败：$e');
        }
      }
    }
  }

  void _onPublisher() {
    if (!mounted) return;
    setState(() {});
  }

  @override
  void dispose() {
    _anchor.removeListener(_onAnchor);
    widget.scope.showLivePublisher.removeListener(_onPublisher);
    _chatController.dispose();
    _chatScroll.dispose();
    super.dispose();
  }

  void _onAnchor() {
    if (!mounted) return;
    setState(() {});
    _scrollChatToEnd();
  }

  void _scrollChatToEnd() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!_chatScroll.hasClients) return;
      _chatScroll.jumpTo(_chatScroll.position.maxScrollExtent);
    });
  }

  void _snack(String msg) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(msg)));
  }

  Future<void> _sendChat() async {
    try {
      await _anchor.sendChat(_chatController.text);
      _chatController.clear();
    } on ApiException catch (e) {
      _snack(e.message);
    } catch (e) {
      _snack(e.toString());
    }
  }

  Future<void> _openAudienceSheet() async {
    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => _AudienceSheet(
        anchor: _anchor,
        onSnack: _snack,
        onUserSelected: (u) async {
          Navigator.pop(ctx);
          await _onUserAction(u);
        },
      ),
    );
  }

  Future<void> _stopLiveAndExit() async {
    if (_stopping) return;
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('结束直播'),
        content: const Text('确定下播并离开直播间？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('取消'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('结束'),
          ),
        ],
      ),
    );
    if (confirm != true || !mounted) return;

    setState(() => _stopping = true);
    try {
      await widget.scope.showLivePublisher.stopPublish();
      await widget.scope.liveRoomRepository.stopLive(anchorId: _roomId);
      _anchor.exitSession();
      if (mounted) {
        Navigator.of(context).pop(true);
      }
    } on ApiException catch (e) {
      _snack(e.message);
    } catch (e) {
      _snack(e.toString());
    } finally {
      if (mounted) setState(() => _stopping = false);
    }
  }

  Future<void> _onUserAction(OnlineUserItem user) async {
    final userId = int.tryParse(user.userId);
    if (userId == null || userId <= 0) return;
    if (userId.toString() == widget.room.roomId) return;

    AudienceRestrictStatus? status;
    try {
      status = await _anchor.restrictStatus(userId);
    } catch (_) {}

    if (!mounted) return;
    await showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      builder: (ctx) {
        final name =
            user.nickname.isNotEmpty ? user.nickname : '用户 ${user.userId}';
        return SafeArea(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Row(
                  children: [
                    CachedAvatar(
                      url: user.avatar,
                      radius: 24,
                      fallbackText: name.isNotEmpty ? name[0] : '?',
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            name,
                            style: Theme.of(ctx).textTheme.titleMedium,
                          ),
                          Text(
                            'VIP ${user.vipLevel} · 贡献 ${user.totalReward.toStringAsFixed(0)} 钻',
                            style: Theme.of(ctx).textTheme.bodySmall,
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                if (status != null) ...[
                  const SizedBox(height: 12),
                  Text(
                    status.muted
                        ? '当前：已禁言'
                        : status.kickBanned
                            ? '当前：踢出限制 ${status.kickRemainSeconds}s'
                            : '当前：正常',
                    style: Theme.of(ctx).textTheme.bodySmall,
                  ),
                ],
                const SizedBox(height: 16),
                _ActionTile(
                  icon: user.muted ? Icons.record_voice_over : Icons.mic_off,
                  label: user.muted ? '解除禁言' : '禁言',
                  onTap: () async {
                    Navigator.pop(ctx);
                    try {
                      await _anchor.setMute(userId, !user.muted);
                      _snack(user.muted ? '已解除禁言' : '已禁言');
                    } on ApiException catch (e) {
                      _snack(e.message);
                    }
                  },
                ),
                _ActionTile(
                  icon: Icons.logout,
                  label: '踢出直播间',
                  onTap: () async {
                    Navigator.pop(ctx);
                    try {
                      await _anchor.kickUser(userId);
                      _snack('已踢出');
                    } on ApiException catch (e) {
                      _snack(e.message);
                    }
                  },
                ),
                _ActionTile(
                  icon: Icons.block,
                  label: '拉黑并踢出',
                  destructive: true,
                  onTap: () async {
                    Navigator.pop(ctx);
                    try {
                      await _anchor.blockUser(userId);
                      _snack('已拉黑');
                    } on ApiException catch (e) {
                      _snack(e.message);
                    }
                  },
                ),
                if (status?.kickBanned == true)
                  _ActionTile(
                    icon: Icons.lock_open,
                    label: '取消踢出限制',
                    onTap: () async {
                      Navigator.pop(ctx);
                      try {
                        await _anchor.cancelKickBan(userId);
                        _snack('已取消限制');
                      } on ApiException catch (e) {
                        _snack(e.message);
                      }
                    },
                  ),
              ],
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final publisher = widget.scope.showLivePublisher;
    final rawTitle =
        _anchor.title.isNotEmpty ? _anchor.title : widget.room.title;
    final title = rawTitle.trim().isNotEmpty
        ? rawTitle.trim()
        : _anchorNickname.trim();

    return PopScope(
      canPop: !_stopping,
      child: Scaffold(
        backgroundColor: Colors.black,
        body: Stack(
          fit: StackFit.expand,
          children: [
            LiveRoomCoverBackdrop(
              coverUrl: _coverUrl,
              fallbackAvatarUrl: _avatarUrl,
              showSubtitle: false,
              videoForegroundActive:
                  publisher.isPublishing && publisher.cameraEnabled,
            ),
            if (publisher.isPublishing && publisher.cameraEnabled)
              Positioned.fill(
                child: ColoredBox(
                  color: Colors.transparent,
                  child: publisher.rtcSession.localView(),
                ),
              )
            else if (!publisher.isPublishing)
              Center(
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 20,
                    vertical: 14,
                  ),
                  decoration: BoxDecoration(
                    color: Colors.black54,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      SizedBox(
                        width: 22,
                        height: 22,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      ),
                      SizedBox(width: 12),
                      Text(
                        '正在连接摄像头…',
                        style: TextStyle(color: Colors.white),
                      ),
                    ],
                  ),
                ),
              ),
            SafeArea(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(8, 4, 8, 0),
                    child: Row(
                      children: [
                        IconButton(
                          onPressed: _stopping
                              ? null
                              : () => Navigator.of(context).pop(false),
                          icon: const Icon(Icons.keyboard_arrow_down_rounded,
                              color: Colors.white, size: 32),
                          tooltip: '收起',
                        ),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Container(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 8,
                                      vertical: 3,
                                    ),
                                    decoration: BoxDecoration(
                                      color: const Color(0xFFEF4444),
                                      borderRadius: BorderRadius.circular(999),
                                    ),
                                    child: const Text(
                                      'LIVE',
                                      style: TextStyle(
                                        color: Colors.white,
                                        fontSize: 11,
                                        fontWeight: FontWeight.w800,
                                        letterSpacing: 0.6,
                                      ),
                                    ),
                                  ),
                                  const SizedBox(width: 8),
                                  Expanded(
                                    child: Text(
                                      title,
                                      maxLines: 1,
                                      overflow: TextOverflow.ellipsis,
                                      style: const TextStyle(
                                        color: Colors.white,
                                        fontWeight: FontWeight.w700,
                                        fontSize: 16,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              if (_anchor.notice.isNotEmpty)
                                Padding(
                                  padding: const EdgeInsets.only(top: 4),
                                  child: Text(
                                    _anchor.notice,
                                    maxLines: 1,
                                    overflow: TextOverflow.ellipsis,
                                    style: TextStyle(
                                      color:
                                          Colors.white.withValues(alpha: 0.75),
                                      fontSize: 12,
                                    ),
                                  ),
                                ),
                            ],
                          ),
                        ),
                        Material(
                          color: Colors.white.withValues(alpha: 0.14),
                          borderRadius: BorderRadius.circular(999),
                          child: InkWell(
                            onTap: _openAudienceSheet,
                            borderRadius: BorderRadius.circular(999),
                            child: Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 12,
                                vertical: 8,
                              ),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  const Icon(Icons.people_alt_outlined,
                                      color: Colors.white, size: 18),
                                  const SizedBox(width: 6),
                                  Text(
                                    '${_anchor.onlineCount}',
                                    style: const TextStyle(
                                      color: Colors.white,
                                      fontWeight: FontWeight.w700,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  if (_anchor.onlinePreview.isNotEmpty)
                    SizedBox(
                      height: 44,
                      child: ListView.separated(
                        scrollDirection: Axis.horizontal,
                        padding: const EdgeInsets.symmetric(horizontal: 12),
                        itemCount: _anchor.onlinePreview.length,
                        separatorBuilder: (_, __) => const SizedBox(width: 8),
                        itemBuilder: (context, index) {
                          final u = _anchor.onlinePreview[index];
                          final name = u.nickname.isNotEmpty
                              ? u.nickname
                              : u.userId;
                          return ActionChip(
                            avatar: CachedAvatar(
                              url: u.avatar,
                              radius: 12,
                              fallbackText:
                                  name.isNotEmpty ? name[0] : '?',
                            ),
                            label: Text(
                              name,
                              style: const TextStyle(fontSize: 12),
                            ),
                            onPressed: () => _onUserAction(u),
                          );
                        },
                      ),
                    ),
                  const Spacer(),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                    child: ConstrainedBox(
                      constraints: const BoxConstraints(maxHeight: 200),
                      child: ListView.builder(
                        controller: _chatScroll,
                        padding: EdgeInsets.zero,
                        itemCount: _anchor.chatLines.length,
                        itemBuilder: (context, index) {
                          final line = _anchor.chatLines[index];
                          return GestureDetector(
                            onLongPress: line.senderId.isNotEmpty &&
                                    line.senderId != widget.room.roomId
                                ? () {
                                    final uid = int.tryParse(line.senderId);
                                    if (uid == null) return;
                                    _onUserAction(
                                      OnlineUserItem(
                                        userId: line.senderId,
                                        nickname: line.senderName,
                                        avatar: '',
                                        vipLevel: line.vipLevel,
                                        totalReward: 0,
                                        joinedAt: '',
                                      ),
                                    );
                                  }
                                : null,
                            child: _ChatBubble(line: line),
                          );
                        },
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(12, 8, 12, 8),
                    child: Row(
                      children: [
                        Expanded(
                          child: TextField(
                            controller: _chatController,
                            style: const TextStyle(color: Colors.white),
                            decoration: InputDecoration(
                              hintText: '说点什么…',
                              hintStyle: TextStyle(
                                color: Colors.white.withValues(alpha: 0.45),
                              ),
                              filled: true,
                              fillColor: Colors.white.withValues(alpha: 0.12),
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(24),
                                borderSide: BorderSide.none,
                              ),
                              contentPadding: const EdgeInsets.symmetric(
                                horizontal: 16,
                                vertical: 10,
                              ),
                            ),
                            textInputAction: TextInputAction.send,
                            onSubmitted: (_) => _sendChat(),
                          ),
                        ),
                        const SizedBox(width: 8),
                        IconButton.filled(
                          onPressed: _sendChat,
                          icon: const Icon(Icons.send_rounded),
                        ),
                        const SizedBox(width: 4),
                        IconButton.filledTonal(
                          style: IconButton.styleFrom(
                            backgroundColor:
                                const Color(0xFFEF4444).withValues(alpha: 0.9),
                            foregroundColor: Colors.white,
                          ),
                          onPressed: _stopping ? null : _stopLiveAndExit,
                          icon: _stopping
                              ? const SizedBox(
                                  width: 20,
                                  height: 20,
                                  child: CircularProgressIndicator(
                                    strokeWidth: 2,
                                    color: Colors.white,
                                  ),
                                )
                              : const Icon(Icons.stop_rounded),
                          tooltip: '结束直播',
                        ),
                        if (publisher.isPublishing) ...[
                          const SizedBox(width: 4),
                          AnchorLiveToolbar(
                            micMuted: publisher.micMuted,
                            cameraEnabled: publisher.cameraEnabled,
                            beautyEnabled: publisher.beautyEnabled,
                            beautyLevel: publisher.beautyLevel,
                            onFlipCamera: () => publisher.switchCamera(),
                            onToggleCamera: () =>
                                publisher.toggleCameraEnabled(),
                            onToggleMic: () =>
                                publisher.toggleMicrophoneMuted(),
                            onBeautyEnabledChanged: (v) =>
                                publisher.setBeautyEnabled(v),
                            onBeautyLevelChanged: (v) =>
                                publisher.setBeautyLevel(v),
                          ),
                        ],
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ChatBubble extends StatelessWidget {
  const _ChatBubble({required this.line});

  final LiveRoomChatLine line;

  @override
  Widget build(BuildContext context) {
    final name = line.senderName.isNotEmpty ? line.senderName : '…';
    final text = line.isGift ? (line.giftSummary ?? '') : line.content;
    if (text.isEmpty) return const SizedBox.shrink();

    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        decoration: BoxDecoration(
          color: Colors.black.withValues(alpha: 0.35),
          borderRadius: BorderRadius.circular(10),
        ),
        child: RichText(
          text: TextSpan(
            style: const TextStyle(
              color: Colors.white,
              fontSize: 13,
              height: 1.35,
            ),
            children: [
              if (line.vipLevel > 0)
                TextSpan(
                  text: 'VIP${line.vipLevel} ',
                  style: const TextStyle(
                    color: Color(0xFFFFD166),
                    fontWeight: FontWeight.w700,
                  ),
                ),
              TextSpan(
                text: '$name: ',
                style: TextStyle(
                  color: line.isGift
                      ? const Color(0xFFFFD166)
                      : const Color(0xFF93C5FD),
                  fontWeight: FontWeight.w700,
                ),
              ),
              TextSpan(text: text),
            ],
          ),
        ),
      ),
    );
  }
}

class _ActionTile extends StatelessWidget {
  const _ActionTile({
    required this.icon,
    required this.label,
    required this.onTap,
    this.destructive = false,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool destructive;

  @override
  Widget build(BuildContext context) {
    final color = destructive ? Colors.red : null;
    return ListTile(
      leading: Icon(icon, color: color),
      title: Text(label, style: TextStyle(color: color)),
      onTap: onTap,
    );
  }
}

class _AudienceSheet extends StatefulWidget {
  const _AudienceSheet({
    required this.anchor,
    required this.onSnack,
    required this.onUserSelected,
  });

  final LiveRoomAnchorController anchor;
  final void Function(String message) onSnack;
  final Future<void> Function(OnlineUserItem user) onUserSelected;

  @override
  State<_AudienceSheet> createState() => _AudienceSheetState();
}

class _AudienceSheetState extends State<_AudienceSheet> {
  final _items = <OnlineUserItem>[];
  int _total = 0;
  int _page = 1;
  bool _hasMore = true;
  bool _loading = true;
  bool _loadingMore = false;

  @override
  void initState() {
    super.initState();
    unawaited(_loadPage(1, reset: true));
  }

  Future<void> _loadPage(int page, {bool reset = false}) async {
    if (_loadingMore && !reset) return;
    if (reset) {
      setState(() => _loading = true);
    } else {
      setState(() => _loadingMore = true);
    }
    try {
      final res = await widget.anchor.fetchOnlinePage(page: page);
      if (!mounted) return;
      setState(() {
        if (reset) {
          _items
            ..clear()
            ..addAll(res.list);
        } else {
          _items.addAll(res.list);
        }
        _total = res.total;
        _page = res.page;
        _hasMore = res.hasMore;
        _loading = false;
        _loadingMore = false;
      });
    } on ApiException catch (e) {
      widget.onSnack(e.message);
      if (mounted) {
        setState(() {
          _loading = false;
          _loadingMore = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return DraggableScrollableSheet(
      initialChildSize: 0.72,
      minChildSize: 0.4,
      maxChildSize: 0.92,
      builder: (context, scrollController) {
        return Container(
          decoration: BoxDecoration(
            color: theme.colorScheme.surface,
            borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 16, 12, 8),
                child: Row(
                  children: [
                    Text(
                      '在线观众',
                      style: theme.textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '$_total',
                      style: theme.textTheme.titleMedium?.copyWith(
                        color: theme.colorScheme.primary,
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                    const Spacer(),
                    IconButton(
                      onPressed: () => _loadPage(1, reset: true),
                      icon: const Icon(Icons.refresh),
                    ),
                  ],
                ),
              ),
              Expanded(
                child: _loading && _items.isEmpty
                    ? const Center(child: CircularProgressIndicator())
                    : _items.isEmpty
                        ? Center(
                            child: Text(
                              '暂无观众',
                              style: theme.textTheme.bodyLarge,
                            ),
                          )
                        : NotificationListener<ScrollNotification>(
                            onNotification: (n) {
                              if (_loadingMore || !_hasMore) return false;
                              if (n.metrics.pixels >=
                                  n.metrics.maxScrollExtent - 80) {
                                unawaited(_loadPage(_page + 1));
                              }
                              return false;
                            },
                            child: ListView.separated(
                            controller: scrollController,
                            padding: const EdgeInsets.fromLTRB(12, 0, 12, 24),
                            itemCount:
                                _items.length + (_loadingMore ? 1 : 0),
                            separatorBuilder: (_, __) =>
                                const SizedBox(height: 6),
                            itemBuilder: (context, index) {
                              if (index >= _items.length) {
                                return const Center(
                                  child: Padding(
                                    padding: EdgeInsets.all(16),
                                    child: CircularProgressIndicator(),
                                  ),
                                );
                              }
                              final u = _items[index];
                              final name = u.nickname.isNotEmpty
                                  ? u.nickname
                                  : '用户 ${u.userId}';
                              return Material(
                                color: theme.colorScheme.surfaceContainerHighest
                                    .withValues(alpha: 0.5),
                                borderRadius: BorderRadius.circular(14),
                                child: ListTile(
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(14),
                                  ),
                                  leading: CachedAvatar(
                                    url: u.avatar,
                                    radius: 22,
                                    fallbackText:
                                        name.isNotEmpty ? name[0] : '?',
                                  ),
                                  title: Row(
                                    children: [
                                      Expanded(
                                        child: Text(
                                          name,
                                          maxLines: 1,
                                          overflow: TextOverflow.ellipsis,
                                        ),
                                      ),
                                      if (u.muted)
                                        Container(
                                          margin:
                                              const EdgeInsets.only(left: 6),
                                          padding: const EdgeInsets.symmetric(
                                            horizontal: 6,
                                            vertical: 2,
                                          ),
                                          decoration: BoxDecoration(
                                            color: Colors.orange
                                                .withValues(alpha: 0.2),
                                            borderRadius:
                                                BorderRadius.circular(4),
                                          ),
                                          child: const Text(
                                            '禁言',
                                            style: TextStyle(
                                              fontSize: 10,
                                              color: Colors.orange,
                                            ),
                                          ),
                                        ),
                                    ],
                                  ),
                                  subtitle: Text(
                                    'VIP ${u.vipLevel} · 贡献 ${u.totalReward.toStringAsFixed(0)} 钻',
                                  ),
                                  trailing: const Icon(Icons.more_horiz),
                                  onTap: () => widget.onUserSelected(u),
                                ),
                              );
                            },
                          ),
                          ),
              ),
            ],
          ),
        );
      },
    );
  }
}
