import 'dart:async';
import 'dart:ui';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/cache/image_disk_cache.dart';
import 'package:sara_live_app/core/http/api_exception.dart';
import 'package:sara_live_app/core/widgets/cached_avatar.dart';
import 'package:sara_live_app/data/gift/gift_models.dart';
import 'package:sara_live_app/data/live_room/live_room_models.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/ui/widgets/live_room_cover_backdrop.dart';

/// 观众观看秀场：沉浸式拉流 + 公屏 + 送礼。
class LiveRoomWatchPage extends StatefulWidget {
  const LiveRoomWatchPage({
    super.key,
    required this.scope,
    required this.roomId,
    required this.title,
    this.coverUrl = '',
    this.anchorName = '',
    this.anchorAvatar = '',
  });

  final AppScope scope;
  final int roomId;
  final String title;
  final String coverUrl;
  final String anchorName;
  final String anchorAvatar;

  @override
  State<LiveRoomWatchPage> createState() => _LiveRoomWatchPageState();
}

class _LiveRoomWatchPageState extends State<LiveRoomWatchPage>
    with SingleTickerProviderStateMixin {
  final _chatController = TextEditingController();
  final _chatScroll = ScrollController();
  final _chatFocus = FocusNode();

  bool _entering = true;
  String? _enterError;
  List<AppGiftItem> _gifts = const [];
  bool _giftsLoading = false;
  bool _giftPanelOpen = false;

  int _onlineCount = 0;
  List<OnlineUserItem> _onlinePreview = const [];
  String _roomTitle = '';
  String _coverUrl = '';
  Timer? _onlineRefreshTimer;

  late final AnimationController _livePulse;

  @override
  void initState() {
    super.initState();
    _roomTitle = widget.title;
    _coverUrl = widget.coverUrl;
    _livePulse = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1200),
    )..repeat(reverse: true);
    widget.scope.liveRoomViewer.addListener(_onViewer);
    unawaited(_enterRoom());
  }

  @override
  void dispose() {
    _onlineRefreshTimer?.cancel();
    _livePulse.dispose();
    widget.scope.liveRoomViewer.removeListener(_onViewer);
    _chatController.dispose();
    _chatScroll.dispose();
    _chatFocus.dispose();
    unawaited(widget.scope.liveRoomViewer.leaveRoom());
    super.dispose();
  }

  void _onViewer() {
    if (!mounted) return;
    setState(() {});
    _scrollChatToEnd();
    if (widget.scope.liveRoomViewer.anchorStopped) {
      _showSnack(AppLocalizations.of(context).liveAnchorEnded);
    }
  }

  void _scrollChatToEnd() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!_chatScroll.hasClients) return;
      _chatScroll.animateTo(
        _chatScroll.position.maxScrollExtent,
        duration: const Duration(milliseconds: 200),
        curve: Curves.easeOut,
      );
    });
  }

  Future<void> _enterRoom() async {
    setState(() {
      _entering = true;
      _enterError = null;
    });
    unawaited(_refreshRoomMeta(silent: true));
    try {
      await widget.scope.liveRoomViewer.enterRoom(
        roomId: widget.roomId,
        title: widget.title,
      );
      if (!mounted) return;
      setState(() => _entering = false);
      unawaited(_refreshRoomMeta());
      _startOnlineRefresh();
    } on ApiException catch (e) {
      if (!mounted) return;
      setState(() {
        _entering = false;
        _enterError = e.message;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _entering = false;
        _enterError = e.toString();
      });
    }
  }

  void _startOnlineRefresh() {
    _onlineRefreshTimer?.cancel();
    _onlineRefreshTimer = Timer.periodic(const Duration(seconds: 25), (_) {
      unawaited(_refreshRoomMeta(silent: true));
    });
  }

  Future<void> _refreshRoomMeta({bool silent = false}) async {
    try {
      final room = await widget.scope.liveRoomRepository.getLiveRoom(
        anchorId: widget.roomId,
        forceRefresh: true,
      );
      final online = await widget.scope.liveRoomRepository.getOnlineList(
        anchorId: widget.roomId,
        page: 1,
        pageSize: 12,
        forceRefresh: true,
      );
      if (!mounted) return;
      setState(() {
        if (room != null) {
          _onlineCount = room.onlineCount;
          if (room.title.isNotEmpty) _roomTitle = room.title;
          if (room.cover.isNotEmpty) _coverUrl = room.cover;
        } else {
          _onlineCount = online.total;
        }
        _onlinePreview = online.list.take(3).toList();
        if (_onlineCount == 0 && online.total > 0) {
          _onlineCount = online.total;
        }
      });
    } catch (_) {
      if (!silent && mounted) {
        // 元数据失败不阻断观看
      }
    }
  }

  void _showSnack(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }

  Future<void> _sendChat() async {
    final text = _chatController.text;
    if (text.trim().isEmpty) return;
    try {
      await widget.scope.liveRoomViewer.sendChat(text);
      _chatController.clear();
    } on ApiException catch (e) {
      _showSnack(e.message);
    } catch (e) {
      _showSnack(e.toString());
    }
  }

  Future<void> _loadGifts() async {
    if (_giftsLoading) return;
    setState(() => _giftsLoading = true);
    try {
      final list = await widget.scope.giftApi.fetchAppGiftList();
      if (!mounted) return;
      setState(() {
        _gifts = list;
        _giftsLoading = false;
      });
    } on ApiException catch (e) {
      if (!mounted) return;
      setState(() => _giftsLoading = false);
      _showSnack(e.message);
    } catch (e) {
      if (!mounted) return;
      setState(() => _giftsLoading = false);
    }
  }

  Future<void> _toggleGiftPanel() async {
    if (_giftPanelOpen) {
      setState(() => _giftPanelOpen = false);
      return;
    }
    if (_gifts.isEmpty && !_giftsLoading) {
      await _loadGifts();
    }
    if (!mounted) return;
    setState(() => _giftPanelOpen = true);
  }

  Future<void> _sendGift(AppGiftItem gift) async {
    final giftId = int.tryParse(gift.id);
    if (giftId == null || giftId <= 0) {
      _showSnack('Invalid gift');
      return;
    }
    final l10n = AppLocalizations.of(context);
    try {
      final res = await widget.scope.liveRoomViewer.sendGift(
        giftId: giftId,
        count: 1,
      );
      if (!mounted) return;
      setState(() => _giftPanelOpen = false);
      _showSnack(l10n.liveGiftSent(gift.name, res.diamond));
    } on ApiException catch (e) {
      _showSnack(e.message);
    } catch (e) {
      _showSnack(e.toString());
    }
  }

  Future<void> _openAudienceSheet() async {
    if (!mounted) return;
    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      backgroundColor: Theme.of(context).colorScheme.surface,
      builder: (ctx) => _AudienceListSheet(
        scope: widget.scope,
        roomId: widget.roomId,
        initialTotal: _onlineCount,
      ),
    );
    unawaited(_refreshRoomMeta(silent: true));
  }

  Future<void> _leaveAndPop() async {
    await widget.scope.liveRoomViewer.leaveRoom();
    if (mounted) {
      Navigator.of(context).pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final viewer = widget.scope.liveRoomViewer;
    final theme = Theme.of(context);
    final anchorName = widget.anchorName.trim();
    final rawTitle = _roomTitle.isNotEmpty
        ? _roomTitle
        : (viewer.roomTitle.isNotEmpty ? viewer.roomTitle : widget.title);
    final displayTitle = _liveRoomHeadline(
      title: rawTitle,
      anchorName: anchorName,
    );
    final cover = _coverUrl.isNotEmpty ? _coverUrl : widget.coverUrl;
    final disabled =
        viewer.anchorStopped || _entering || !viewer.isWatching;

    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) async {
        if (didPop) return;
        await _leaveAndPop();
      },
      child: Scaffold(
        backgroundColor: Colors.black,
        extendBodyBehindAppBar: true,
        body: Stack(
          fit: StackFit.expand,
          children: [
            Positioned.fill(
              child: _enterError == null
                  ? ValueListenableBuilder<bool>(
                      valueListenable: viewer.rtcSession.remoteVideoReady,
                      builder: (context, ready, _) {
                        return Stack(
                          fit: StackFit.expand,
                          children: [
                            LiveRoomCoverBackdrop(
                              coverUrl: cover,
                              fallbackAvatarUrl: widget.anchorAvatar,
                              showSubtitle: false,
                              videoForegroundActive:
                                  ready && !viewer.anchorStopped,
                            ),
                            if (ready)
                              Positioned.fill(
                                child: ColoredBox(
                                  color: Colors.transparent,
                                  child: viewer.rtcSession.remoteView(),
                                ),
                              ),
                          ],
                        );
                      },
                    )
                  : LiveRoomCoverBackdrop(
                      coverUrl: cover,
                      fallbackAvatarUrl: widget.anchorAvatar,
                      showSubtitle: false,
                    ),
            ),
            const _VignetteOverlay(),
            if (_enterError == null)
              Positioned(
                top: 0,
                left: 0,
                right: 0,
                child: _WatchTopBar(
                  livePulse: _livePulse,
                  connecting: _entering,
                  anchorName: anchorName,
                  anchorAvatar: widget.anchorAvatar,
                  headline: displayTitle,
                  onlineCount: _onlineCount,
                  onlinePreview: _onlinePreview,
                  onClose: _leaveAndPop,
                  onAudienceTap: _openAudienceSheet,
                ),
              ),
            if (_enterError != null)
              Center(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        _enterError!,
                        textAlign: TextAlign.center,
                        style: theme.textTheme.bodyLarge?.copyWith(
                          color: Colors.white70,
                        ),
                      ),
                      const SizedBox(height: 16),
                      FilledButton(
                        onPressed: _enterRoom,
                        child: Text(l10n.commonRetry),
                      ),
                    ],
                  ),
                ),
              ),
            if (_enterError == null)
              Positioned(
                left: 0,
                right: 0,
                bottom: 0,
                child: _ChatPanel(
                  lines: viewer.chatLines,
                  scrollController: _chatScroll,
                  chatController: _chatController,
                  chatFocus: _chatFocus,
                  onSend: _sendChat,
                  anchorStopped: disabled,
                  giftPanelOpen: _giftPanelOpen,
                  giftsLoading: _giftsLoading,
                  gifts: _gifts,
                  onToggleGift: _toggleGiftPanel,
                  onSendGift: _sendGift,
                ),
              ),
            if (viewer.anchorStopped && !_entering && _enterError == null)
              Positioned.fill(
                child: ColoredBox(
                  color: Colors.black54,
                  child: Center(
                    child: Text(
                      l10n.liveAnchorEnded,
                      style: theme.textTheme.titleMedium?.copyWith(
                        color: Colors.white,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}

/// 直播间标题：无标题时用主播昵称。
String _liveRoomHeadline({
  required String title,
  required String anchorName,
}) {
  final t = title.trim();
  if (t.isNotEmpty) return t;
  return anchorName.trim();
}

class _VignetteOverlay extends StatelessWidget {
  const _VignetteOverlay();

  @override
  Widget build(BuildContext context) {
    return IgnorePointer(
      child: DecoratedBox(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              Colors.black.withValues(alpha: 0.55),
              Colors.transparent,
              Colors.transparent,
              Colors.black.withValues(alpha: 0.75),
            ],
            stops: const [0, 0.22, 0.55, 1],
          ),
        ),
      ),
    );
  }
}

class _WatchTopBar extends StatelessWidget {
  const _WatchTopBar({
    required this.livePulse,
    required this.connecting,
    required this.anchorName,
    required this.anchorAvatar,
    required this.headline,
    required this.onlineCount,
    required this.onlinePreview,
    required this.onClose,
    required this.onAudienceTap,
  });

  final AnimationController livePulse;
  final bool connecting;
  final String anchorName;
  final String anchorAvatar;
  final String headline;
  final int onlineCount;
  final List<OnlineUserItem> onlinePreview;
  final VoidCallback onClose;
  final VoidCallback onAudienceTap;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final line = headline.isNotEmpty
        ? headline
        : (anchorName.isNotEmpty ? anchorName : 'Live');

    return SafeArea(
      bottom: false,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(8, 6, 12, 0),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            IconButton(
              onPressed: onClose,
              icon: const Icon(Icons.close_rounded, color: Colors.white),
              style: IconButton.styleFrom(
                backgroundColor: Colors.black38,
              ),
            ),
            const SizedBox(width: 4),
            Expanded(
              child: ClipRRect(
                borderRadius: BorderRadius.circular(999),
                child: BackdropFilter(
                  filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                  child: Container(
                    padding: const EdgeInsets.fromLTRB(6, 6, 10, 6),
                    color: Colors.black.withValues(alpha: 0.35),
                    child: Row(
                      children: [
                        CachedAvatar(
                          url: anchorAvatar,
                          radius: 18,
                          fallbackText: line.isNotEmpty ? line[0] : '?',
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Row(
                            children: [
                              Flexible(
                                child: Text(
                                  line,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  style: const TextStyle(
                                    color: Colors.white,
                                    fontWeight: FontWeight.w800,
                                    fontSize: 14,
                                  ),
                                ),
                              ),
                              const SizedBox(width: 6),
                              if (connecting)
                                _ConnectingBadge(
                                  label: l10n.liveJoiningRoom,
                                )
                              else
                                _LiveBadge(
                                  pulse: livePulse,
                                  label: l10n.liveBadgeOnAir,
                                ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 8),
            _AudienceEntryChip(
              count: onlineCount,
              preview: onlinePreview,
              onTap: onAudienceTap,
            ),
          ],
        ),
      ),
    );
  }
}

class _ConnectingBadge extends StatelessWidget {
  const _ConnectingBadge({required this.label});

  final String label;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.18),
        borderRadius: BorderRadius.circular(4),
        border: Border.all(color: Colors.white30),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const SizedBox(
            width: 10,
            height: 10,
            child: CircularProgressIndicator(
              strokeWidth: 1.5,
              color: Colors.white,
            ),
          ),
          const SizedBox(width: 4),
          Text(
            label,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 9,
              fontWeight: FontWeight.w700,
            ),
          ),
        ],
      ),
    );
  }
}

class _LiveBadge extends StatelessWidget {
  const _LiveBadge({required this.pulse, required this.label});

  final AnimationController pulse;
  final String label;

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: pulse,
      builder: (context, child) {
        return Container(
          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
          decoration: BoxDecoration(
            color: Color.lerp(
              const Color(0xFFE11D48),
              const Color(0xFFFB7185),
              pulse.value,
            ),
            borderRadius: BorderRadius.circular(4),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 5,
                height: 5,
                decoration: const BoxDecoration(
                  color: Colors.white,
                  shape: BoxShape.circle,
                ),
              ),
              const SizedBox(width: 4),
              Text(
                label,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 9,
                  fontWeight: FontWeight.w800,
                  letterSpacing: 0.3,
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _AudienceEntryChip extends StatelessWidget {
  const _AudienceEntryChip({
    required this.count,
    required this.preview,
    required this.onTap,
  });

  final int count;
  final List<OnlineUserItem> preview;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(999),
        child: Ink(
          decoration: BoxDecoration(
            color: Colors.black.withValues(alpha: 0.4),
            borderRadius: BorderRadius.circular(999),
            border: Border.all(color: Colors.white24),
          ),
          padding: const EdgeInsets.fromLTRB(6, 6, 10, 6),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (preview.isNotEmpty)
                SizedBox(
                  width: preview.length * 14.0 + 10,
                  height: 28,
                  child: Stack(
                    clipBehavior: Clip.none,
                    children: [
                      for (var i = 0; i < preview.length; i++)
                        Positioned(
                          left: i * 14.0,
                          child: CircleAvatar(
                            radius: 14,
                            backgroundColor: Colors.white24,
                            child: CircleAvatar(
                              radius: 12,
                              backgroundImage: preview[i].avatar.isNotEmpty
                                  ? CachedNetworkImageProvider(
                                      preview[i].avatar,
                                      cacheManager:
                                          ImageDiskCache.instance.manager,
                                    )
                                  : null,
                              child: preview[i].avatar.isEmpty
                                  ? Text(
                                      preview[i].nickname.isNotEmpty
                                          ? preview[i].nickname[0]
                                          : '?',
                                      style: const TextStyle(fontSize: 10),
                                    )
                                  : null,
                            ),
                          ),
                        ),
                    ],
                  ),
                )
              else
                const Icon(
                  Icons.people_alt_outlined,
                  color: Colors.white,
                  size: 22,
                ),
              const SizedBox(width: 6),
              Text(
                '$count',
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w800,
                  fontSize: 13,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ChatPanel extends StatelessWidget {
  _ChatPanel({
    required this.lines,
    required this.scrollController,
    required this.chatController,
    required this.chatFocus,
    required this.onSend,
    required this.anchorStopped,
    required this.giftPanelOpen,
    required this.giftsLoading,
    required this.gifts,
    required this.onToggleGift,
    required this.onSendGift,
  });

  final List<LiveRoomChatLine> lines;
  final ScrollController scrollController;
  final TextEditingController chatController;
  final FocusNode chatFocus;
  final VoidCallback onSend;
  final bool anchorStopped;
  final bool giftPanelOpen;
  final bool giftsLoading;
  final List<AppGiftItem> gifts;
  final VoidCallback onToggleGift;
  final Future<void> Function(AppGiftItem gift) onSendGift;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);

    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.bottomCenter,
          end: Alignment.topCenter,
          colors: [
            Colors.black.withValues(alpha: 0.88),
            Colors.black.withValues(alpha: 0.5),
            Colors.transparent,
          ],
          stops: const [0, 0.55, 1],
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(12, 16, 12, 10),
        child: SafeArea(
          top: false,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              ConstrainedBox(
                constraints: BoxConstraints(
                  maxHeight: giftPanelOpen ? 120 : 200,
                ),
                child: ListView.builder(
                  controller: scrollController,
                  padding: EdgeInsets.zero,
                  itemCount: lines.length,
                  itemBuilder: (context, index) {
                    final line = lines[index];
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 6, right: 8),
                      child: _ChatLineBubble(line: line),
                    );
                  },
                ),
              ),
              if (giftPanelOpen)
                Padding(
                  padding: const EdgeInsets.only(top: 8, bottom: 10),
                  child: _GiftBottomPanel(
                    loading: giftsLoading,
                    gifts: gifts,
                    onSendGift: onSendGift,
                  ),
                ),
              Row(
                children: [
                  Material(
                    color: giftPanelOpen
                        ? Colors.white24
                        : const Color(0xFFE11D48),
                    borderRadius: BorderRadius.circular(999),
                    child: InkWell(
                      onTap: anchorStopped ? null : onToggleGift,
                      borderRadius: BorderRadius.circular(999),
                      child: Padding(
                        padding: const EdgeInsets.all(11),
                        child: Icon(
                          giftPanelOpen
                              ? Icons.keyboard_arrow_down_rounded
                              : Icons.card_giftcard_rounded,
                          color: Colors.white,
                          size: 22,
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: TextField(
                      controller: chatController,
                      focusNode: chatFocus,
                      enabled: !anchorStopped,
                      style: const TextStyle(color: Colors.white),
                      decoration: InputDecoration(
                        hintText: l10n.liveChatHint,
                        hintStyle: TextStyle(
                          color: Colors.white.withValues(alpha: 0.45),
                        ),
                        filled: true,
                        fillColor: Colors.black.withValues(alpha: 0.45),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(24),
                          borderSide: BorderSide(
                            color: Colors.white.withValues(alpha: 0.12),
                          ),
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(24),
                          borderSide: BorderSide(
                            color: Colors.white.withValues(alpha: 0.12),
                          ),
                        ),
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 11,
                        ),
                      ),
                      textInputAction: TextInputAction.send,
                      onSubmitted: (_) => onSend(),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Material(
                    color: const Color(0xFFE11D48),
                    borderRadius: BorderRadius.circular(999),
                    child: InkWell(
                      onTap: anchorStopped ? null : onSend,
                      borderRadius: BorderRadius.circular(999),
                      child: const Padding(
                        padding: EdgeInsets.all(11),
                        child: Icon(
                          Icons.send_rounded,
                          color: Colors.white,
                          size: 22,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _GiftBottomPanel extends StatelessWidget {
  const _GiftBottomPanel({
    required this.loading,
    required this.gifts,
    required this.onSendGift,
  });

  final bool loading;
  final List<AppGiftItem> gifts;
  final Future<void> Function(AppGiftItem gift) onSendGift;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return ClipRRect(
      borderRadius: BorderRadius.circular(14),
      child: ColoredBox(
        color: Colors.black.withValues(alpha: 0.45),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(12, 10, 12, 6),
              child: Text(
                l10n.liveSendGift,
                style: theme.textTheme.titleSmall?.copyWith(
                  color: Colors.white,
                  fontWeight: FontWeight.w800,
                ),
              ),
            ),
            SizedBox(
              height: 196,
              child: loading && gifts.isEmpty
                  ? const Center(
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : gifts.isEmpty
                      ? Center(
                          child: Text(
                            l10n.liveGiftListEmpty,
                            style: TextStyle(
                              color: Colors.white.withValues(alpha: 0.65),
                            ),
                          ),
                        )
                      : GridView.builder(
                          padding: const EdgeInsets.fromLTRB(10, 0, 10, 10),
                          gridDelegate:
                              const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 4,
                            mainAxisSpacing: 8,
                            crossAxisSpacing: 8,
                            childAspectRatio: 0.72,
                          ),
                          itemCount: gifts.length,
                          itemBuilder: (context, index) {
                            final gift = gifts[index];
                            return _GiftGridTile(
                              gift: gift,
                              onSend: () => onSendGift(gift),
                            );
                          },
                        ),
            ),
          ],
        ),
      ),
    );
  }
}

class _GiftGridTile extends StatelessWidget {
  const _GiftGridTile({required this.gift, required this.onSend});

  final AppGiftItem gift;
  final VoidCallback onSend;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.white.withValues(alpha: 0.1),
      borderRadius: BorderRadius.circular(10),
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onSend,
        child: Padding(
          padding: const EdgeInsets.all(4),
          child: Column(
            children: [
              Expanded(
                child: gift.icon.isNotEmpty
                    ? CachedNetworkImage(
                        imageUrl: gift.icon,
                        fit: BoxFit.contain,
                        cacheManager: ImageDiskCache.instance.manager,
                        errorWidget: (_, __, ___) => const Icon(
                          Icons.redeem_outlined,
                          color: Colors.white70,
                          size: 28,
                        ),
                      )
                    : const Icon(
                        Icons.redeem_outlined,
                        color: Colors.white70,
                        size: 28,
                      ),
              ),
              Text(
                gift.name,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 10,
                  fontWeight: FontWeight.w600,
                ),
              ),
              Text(
                '${gift.price}',
                style: const TextStyle(
                  color: Color(0xFFFFB4B4),
                  fontSize: 10,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ChatLineBubble extends StatelessWidget {
  const _ChatLineBubble({required this.line});

  final LiveRoomChatLine line;

  @override
  Widget build(BuildContext context) {
    final name = line.senderName.isNotEmpty ? line.senderName : '…';
    final text = line.isGift ? (line.giftSummary ?? '') : line.content;
    if (text.isEmpty) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: Colors.black.withValues(alpha: 0.38),
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
    );
  }
}

class _AudienceListSheet extends StatefulWidget {
  const _AudienceListSheet({
    required this.scope,
    required this.roomId,
    required this.initialTotal,
  });

  final AppScope scope;
  final int roomId;
  final int initialTotal;

  @override
  State<_AudienceListSheet> createState() => _AudienceListSheetState();
}

class _AudienceListSheetState extends State<_AudienceListSheet> {
  bool _loading = true;
  List<OnlineUserItem> _items = const [];
  int _total = 0;

  @override
  void initState() {
    super.initState();
    _total = widget.initialTotal;
    unawaited(_load());
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      final res = await widget.scope.liveRoomRepository.getOnlineList(
        anchorId: widget.roomId,
        page: 1,
        pageSize: 50,
        forceRefresh: true,
      );
      if (!mounted) return;
      setState(() {
        _items = res.list;
        _total = res.total;
        _loading = false;
      });
    } catch (_) {
      if (!mounted) return;
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return DraggableScrollableSheet(
      expand: false,
      initialChildSize: 0.55,
      minChildSize: 0.35,
      maxChildSize: 0.9,
      builder: (_, scrollController) {
        return SafeArea(
          top: false,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 4, 20, 8),
                child: Text(
                  l10n.liveWatchAudienceTitle,
                  style: theme.textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.w800,
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 0, 20, 12),
                child: Text(
                  l10n.liveWatchAudienceSubtitle(_total),
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurface.withValues(alpha: 0.55),
                  ),
                ),
              ),
              Expanded(
                child: _loading && _items.isEmpty
                    ? const Center(child: CircularProgressIndicator())
                    : _items.isEmpty
                        ? Center(child: Text(l10n.liveWatchAudienceEmpty))
                        : ListView.separated(
                            controller: scrollController,
                            padding: const EdgeInsets.fromLTRB(12, 0, 12, 24),
                            itemCount: _items.length,
                            separatorBuilder: (_, __) =>
                                const SizedBox(height: 8),
                            itemBuilder: (context, index) {
                              final item = _items[index];
                              final name = item.nickname.isNotEmpty
                                  ? item.nickname
                                  : 'ID ${item.userId}';
                              return ListTile(
                                tileColor: theme.colorScheme
                                    .surfaceContainerHighest
                                    .withValues(alpha: 0.5),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                leading: CachedAvatar(
                                  url: item.avatar,
                                  radius: 22,
                                  fallbackText:
                                      name.isNotEmpty ? name[0] : '?',
                                ),
                                title: Text(
                                  name,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                ),
                                subtitle: item.vipLevel > 0
                                    ? Text('VIP ${item.vipLevel}')
                                    : null,
                              );
                            },
                          ),
              ),
            ],
          ),
        );
      },
    );
  }
}
