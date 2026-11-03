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
import 'package:sara_live_app/data/vip/vip_privilege_frames.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/ui/widgets/live_room_chat_line_tile.dart';
import 'package:sara_live_app/ui/widgets/live_room_cover_backdrop.dart';
import 'package:sara_live_app/ui/widgets/live_room_gift_effect_overlay.dart';

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

class _LiveRoomWatchPageState extends State<LiveRoomWatchPage> {
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
  String _coverUrl = '';
  Timer? _onlineRefreshTimer;
  VipPrivilegeFrames _vipFrames = VipPrivilegeFrames.empty;

  bool? _followingAnchor;
  bool _followBusy = false;

  @override
  void initState() {
    super.initState();
    _coverUrl = widget.coverUrl;
    widget.scope.liveRoomViewer.addListener(_onViewer);
    _chatFocus.addListener(_onChatFocusChange);
    unawaited(_enterRoom());
  }

  void _onChatFocusChange() {
    if (!mounted) return;
    if (_chatFocus.hasFocus) {
      if (_giftPanelOpen) {
        setState(() => _giftPanelOpen = false);
      } else {
        setState(() {});
      }
      _scrollChatToEnd();
    } else {
      setState(() {});
    }
  }

  void _dismissLiveInput() {
    _chatFocus.unfocus();
    if (_giftPanelOpen && mounted) {
      setState(() => _giftPanelOpen = false);
    }
  }

  @override
  void dispose() {
    _onlineRefreshTimer?.cancel();
    widget.scope.liveRoomViewer.removeListener(_onViewer);
    _chatFocus.removeListener(_onChatFocusChange);
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

  Future<void> _loadVipFrames() async {
    final userId =
        int.tryParse(widget.scope.session.token?.userId ?? '') ?? 0;
    if (userId <= 0) return;
    try {
      final cfgs = await widget.scope.vipRepository.getVipCfgList(
        userId: userId,
        forceRefresh: false,
      );
      if (!mounted) return;
      setState(
        () => _vipFrames = VipPrivilegeFrames.fromCfgList(cfgs),
      );
    } catch (_) {}
  }

  int? get _selfUserId =>
      int.tryParse(widget.scope.session.token?.userId ?? '');

  bool get _isSelfAnchor {
    final self = _selfUserId;
    return self != null && self > 0 && self == widget.roomId;
  }

  Future<void> _refreshFollowState() async {
    if (_isSelfAnchor) return;
    try {
      final following = await widget.scope.liveFollowApi.isFollowing(
        anchorId: widget.roomId,
      );
      if (!mounted) return;
      setState(() => _followingAnchor = following);
    } catch (_) {
      if (!mounted) return;
      setState(() => _followingAnchor ??= false);
    }
  }

  Future<void> _onFollowAnchorTap() async {
    if (_followBusy || _isSelfAnchor) return;
    setState(() => _followBusy = true);
    try {
      final following = _followingAnchor == true;
      final res = following
          ? await widget.scope.liveFollowApi.unfollow(
              anchorId: widget.roomId,
            )
          : await widget.scope.liveFollowApi.follow(
              anchorId: widget.roomId,
            );
      if (!mounted) return;
      setState(() => _followingAnchor = res.following);
    } on ApiException catch (e) {
      if (mounted) _showSnack(e.message);
    } catch (e) {
      if (mounted) _showSnack('$e');
    } finally {
      if (mounted) setState(() => _followBusy = false);
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
      unawaited(_loadVipFrames());
      unawaited(_refreshFollowState());
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
          if (room.cover.isNotEmpty) _coverUrl = room.cover;
        } else {
          _onlineCount = online.total;
        }
        _onlinePreview = online.list.take(8).toList();
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
      _chatFocus.unfocus();
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
      unawaited(
        widget.scope.liveRoomGiftEffectPlayer.prefetchGiftAnimations(
          list.map((g) => g.animation).where((u) => u.isNotEmpty),
        ),
      );
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
    final anchorDisplayName = anchorName.isNotEmpty
        ? anchorName
        : (widget.title.trim().isNotEmpty ? widget.title.trim() : 'Live');
    final cover = _coverUrl.isNotEmpty ? _coverUrl : widget.coverUrl;
    final disabled =
        viewer.anchorStopped || _entering || !viewer.isWatching;
    final media = MediaQuery.of(context);
    final keyboardBottom = media.viewInsets.bottom;

    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) async {
        if (didPop) return;
        await _leaveAndPop();
      },
      child: Scaffold(
        backgroundColor: Colors.black,
        // 抖音式：画面全屏不随键盘缩小，底部公屏/输入用 viewInsets 上移。
        resizeToAvoidBottomInset: false,
        extendBodyBehindAppBar: true,
        body: Stack(
          fit: StackFit.expand,
          children: [
            Positioned.fill(
              child: GestureDetector(
                behavior: HitTestBehavior.translucent,
                onTap: _dismissLiveInput,
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
            ),
            const _VignetteOverlay(),
            if (_enterError == null)
              Positioned(
                top: 0,
                left: 0,
                right: 0,
                child: _WatchTopBar(
                  connecting: _entering,
                  anchorName: anchorDisplayName,
                  anchorAvatar: widget.anchorAvatar,
                  onlineCount: _onlineCount,
                  onlinePreview: _onlinePreview,
                  showFollow: !_isSelfAnchor,
                  following: _followingAnchor == true,
                  followBusy: _followBusy,
                  onFollowTap: _onFollowAnchorTap,
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
                child: AnimatedPadding(
                  duration: const Duration(milliseconds: 120),
                  curve: Curves.easeOut,
                  padding: EdgeInsets.only(bottom: keyboardBottom),
                  child: _ChatPanel(
                    lines: viewer.chatLines,
                    vipFrames: _vipFrames,
                    walletDiamond: viewer.walletDiamond,
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
            LiveRoomGiftEffectOverlay(
              player: widget.scope.liveRoomGiftEffectPlayer,
            ),
          ],
        ),
      ),
    );
  }
}

String _formatBalance(double v) {
  if (v == v.roundToDouble()) {
    return v.toInt().toString();
  }
  return v.toStringAsFixed(2);
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

/// 顶部主播条 / 在线观众条统一高度。
const _liveTopBarPillHeight = 44.0;

class _WatchTopBar extends StatelessWidget {
  const _WatchTopBar({
    required this.connecting,
    required this.anchorName,
    required this.anchorAvatar,
    required this.onlineCount,
    required this.onlinePreview,
    required this.showFollow,
    required this.following,
    required this.followBusy,
    required this.onFollowTap,
    required this.onClose,
    required this.onAudienceTap,
  });

  final bool connecting;
  final String anchorName;
  final String anchorAvatar;
  final int onlineCount;
  final List<OnlineUserItem> onlinePreview;
  final bool showFollow;
  final bool following;
  final bool followBusy;
  final VoidCallback onFollowTap;
  final VoidCallback onClose;
  final VoidCallback onAudienceTap;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);

    return SafeArea(
      bottom: false,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(6, 6, 10, 0),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            IconButton(
              visualDensity: VisualDensity.compact,
              padding: EdgeInsets.zero,
              constraints: const BoxConstraints(minWidth: 40, minHeight: 40),
              onPressed: onClose,
              icon: const Icon(Icons.close_rounded, color: Colors.white),
              style: IconButton.styleFrom(
                backgroundColor: Colors.black38,
              ),
            ),
            const SizedBox(width: 4),
            Expanded(
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  Flexible(
                    flex: 5,
                    fit: FlexFit.loose,
                    child: SizedBox(
                      height: _liveTopBarPillHeight,
                      child: _AnchorHostCard(
                        anchorName: anchorName,
                        anchorAvatar: anchorAvatar,
                        connecting: connecting,
                        connectingLabel: l10n.liveJoiningRoom,
                        showFollow: showFollow,
                        following: following,
                        followBusy: followBusy,
                        followLabel: l10n.liveFollowAnchor,
                        followingLabel: l10n.liveFollowingAnchor,
                        onFollowTap: onFollowTap,
                      ),
                    ),
                  ),
                  const SizedBox(width: 6),
                  Expanded(
                    flex: 6,
                    child: SizedBox(
                      height: _liveTopBarPillHeight,
                      child: _AudienceEntryChip(
                        count: onlineCount,
                        preview: onlinePreview,
                        onTap: onAudienceTap,
                      ),
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

class _AnchorHostCard extends StatelessWidget {
  const _AnchorHostCard({
    required this.anchorName,
    required this.anchorAvatar,
    required this.connecting,
    required this.connectingLabel,
    required this.showFollow,
    required this.following,
    required this.followBusy,
    required this.followLabel,
    required this.followingLabel,
    required this.onFollowTap,
  });

  final String anchorName;
  final String anchorAvatar;
  final bool connecting;
  final String connectingLabel;
  final bool showFollow;
  final bool following;
  final bool followBusy;
  final String followLabel;
  final String followingLabel;
  final VoidCallback onFollowTap;

  @override
  Widget build(BuildContext context) {
    final initial =
        anchorName.isNotEmpty ? anchorName.characters.first : '?';

    return ClipRRect(
      borderRadius: BorderRadius.circular(14),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
        child: Container(
          width: double.infinity,
          height: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 8),
          alignment: Alignment.center,
          color: Colors.black.withValues(alpha: 0.38),
          child: Row(
            children: [
              CachedAvatar(
                url: anchorAvatar,
                radius: 16,
                fallbackText: initial,
              ),
              const SizedBox(width: 6),
              Expanded(
                child: Text(
                  anchorName,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w800,
                    fontSize: 13,
                  ),
                ),
              ),
              const SizedBox(width: 4),
              if (connecting)
                _ConnectingBadge(label: connectingLabel)
              else if (showFollow)
                _FollowAnchorChip(
                  following: following,
                  loading: followBusy,
                  followLabel: followLabel,
                  followingLabel: followingLabel,
                  onTap: onFollowTap,
                ),
            ],
          ),
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

class _FollowAnchorChip extends StatelessWidget {
  const _FollowAnchorChip({
    required this.following,
    required this.loading,
    required this.followLabel,
    required this.followingLabel,
    required this.onTap,
  });

  final bool following;
  final bool loading;
  final String followLabel;
  final String followingLabel;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final label = following ? followingLabel : followLabel;
    return Material(
      color: following
          ? Colors.white.withValues(alpha: 0.16)
          : const Color(0xFFE11D48),
      borderRadius: BorderRadius.circular(999),
      child: InkWell(
        onTap: loading ? null : onTap,
        borderRadius: BorderRadius.circular(999),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
          child: loading
              ? const SizedBox(
                  width: 14,
                  height: 14,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: Colors.white,
                  ),
                )
              : Text(
                  label,
                  style: TextStyle(
                    color: following
                        ? Colors.white.withValues(alpha: 0.85)
                        : Colors.white,
                    fontSize: 11,
                    fontWeight: FontWeight.w800,
                  ),
                ),
        ),
      ),
    );
  }
}

class _AudienceEntryChip extends StatelessWidget {
  const _AudienceEntryChip({
    required this.count,
    required this.preview,
    required this.onTap,
  });

  static const _maxAvatars = 8;
  static const _avatarDiameter = 26.0;
  static const _overlapStep = 10.0;

  final int count;
  final List<OnlineUserItem> preview;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final shown = preview.take(_maxAvatars).toList();
    final stackWidth = shown.isEmpty
        ? 0.0
        : _avatarDiameter + (shown.length - 1) * _overlapStep;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(14),
        child: Ink(
          width: double.infinity,
          height: double.infinity,
          decoration: BoxDecoration(
            color: Colors.black.withValues(alpha: 0.4),
            borderRadius: BorderRadius.circular(14),
            border: Border.all(color: Colors.white24),
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
              if (shown.isNotEmpty)
                SizedBox(
                  width: stackWidth,
                  height: _avatarDiameter,
                  child: Stack(
                    clipBehavior: Clip.none,
                    children: [
                      for (var i = 0; i < shown.length; i++)
                        Positioned(
                          left: i * _overlapStep,
                          child: _PreviewAvatar(item: shown[i]),
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
              const SizedBox(width: 8),
              Expanded(
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    Text(
                      '$count',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.w800,
                        fontSize: 14,
                      ),
                    ),
                    const SizedBox(width: 2),
                    Icon(
                      Icons.chevron_right_rounded,
                      color: Colors.white.withValues(alpha: 0.7),
                      size: 20,
                    ),
                  ],
                ),
              ),
            ],
            ),
          ),
        ),
      ),
    );
  }
}

class _PreviewAvatar extends StatelessWidget {
  const _PreviewAvatar({required this.item});

  final OnlineUserItem item;

  @override
  Widget build(BuildContext context) {
    const r = 13.0;
    return Container(
      width: r * 2,
      height: r * 2,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        border: Border.all(color: Colors.white54, width: 1.2),
      ),
      child: CircleAvatar(
        radius: r - 1.2,
        backgroundColor: Colors.white24,
        backgroundImage: item.avatar.isNotEmpty
            ? CachedNetworkImageProvider(
                item.avatar,
                cacheManager: ImageDiskCache.instance.manager,
              )
            : null,
        child: item.avatar.isEmpty
            ? Text(
                item.nickname.isNotEmpty ? item.nickname.characters.first : '?',
                style: const TextStyle(fontSize: 10, color: Colors.white),
              )
            : null,
      ),
    );
  }
}

class _ChatPanel extends StatelessWidget {
  _ChatPanel({
    required this.lines,
    required this.vipFrames,
    required this.walletDiamond,
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
  final VipPrivilegeFrames vipFrames;
  final double walletDiamond;
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
    final media = MediaQuery.of(context);
    final safe = media.padding;
    final keyboard = media.viewInsets.bottom;
    final typing = keyboard > 0 || chatFocus.hasFocus;

    final chatMaxHeight = typing ? 140.0 : 200.0;

    return Padding(
      padding: EdgeInsets.fromLTRB(
        12,
        0,
        12,
        10 + (keyboard > 0 ? 0 : safe.bottom),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (lines.isNotEmpty)
            ConstrainedBox(
              constraints: BoxConstraints(maxHeight: chatMaxHeight),
              child: ListView.builder(
                controller: scrollController,
                padding: EdgeInsets.zero,
                shrinkWrap: true,
                itemCount: lines.length,
                itemBuilder: (context, index) {
                  final line = lines[index];
                  return Padding(
                    padding: const EdgeInsets.only(bottom: 6),
                    child: LiveRoomChatLineTile(
                      line: line,
                      vipFrames: vipFrames,
                    ),
                  );
                },
              ),
            ),
          if (lines.isNotEmpty) const SizedBox(height: 8),
          if (giftPanelOpen)
            Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: _GiftBottomPanel(
                loading: giftsLoading,
                gifts: gifts,
                diamondBalance: walletDiamond,
                onSendGift: onSendGift,
              ),
            ),
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: chatController,
                  focusNode: chatFocus,
                  enabled: !anchorStopped,
                  minLines: 1,
                  maxLines: typing ? 3 : 1,
                  style: const TextStyle(color: Colors.white, fontSize: 15),
                  decoration: InputDecoration(
                    hintText: l10n.liveChatHint,
                    hintStyle: TextStyle(
                      color: Colors.white.withValues(alpha: 0.45),
                      fontSize: 15,
                    ),
                    filled: true,
                    fillColor: Colors.black.withValues(alpha: 0.35),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(22),
                      borderSide: BorderSide.none,
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(22),
                      borderSide: BorderSide.none,
                    ),
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(22),
                      borderSide: BorderSide(
                        color: Colors.white.withValues(alpha: 0.25),
                      ),
                    ),
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 14,
                      vertical: 10,
                    ),
                    isDense: true,
                  ),
                  textInputAction: TextInputAction.send,
                  onSubmitted: (_) => onSend(),
                ),
              ),
              const SizedBox(width: 8),
              Material(
                color: typing
                    ? const Color(0xFFE11D48)
                    : Colors.black.withValues(alpha: 0.35),
                borderRadius: BorderRadius.circular(999),
                child: InkWell(
                  onTap: anchorStopped ? null : onSend,
                  borderRadius: BorderRadius.circular(999),
                  child: const Padding(
                    padding: EdgeInsets.all(10),
                    child: Icon(
                      Icons.send_rounded,
                      color: Colors.white,
                      size: 22,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Material(
                color: giftPanelOpen
                    ? Colors.white24
                    : Colors.black.withValues(alpha: 0.35),
                borderRadius: BorderRadius.circular(999),
                child: InkWell(
                  onTap: anchorStopped ? null : onToggleGift,
                  borderRadius: BorderRadius.circular(999),
                  child: Padding(
                    padding: const EdgeInsets.all(10),
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
            ],
          ),
        ],
      ),
    );
  }
}

class _GiftBottomPanel extends StatelessWidget {
  const _GiftBottomPanel({
    required this.loading,
    required this.gifts,
    required this.diamondBalance,
    required this.onSendGift,
  });

  final bool loading;
  final List<AppGiftItem> gifts;
  final double diamondBalance;
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
              child: Row(
                children: [
                  Expanded(
                    child: Text(
                      l10n.liveSendGift,
                      style: theme.textTheme.titleSmall?.copyWith(
                        color: Colors.white,
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                  ),
                  Icon(
                    Icons.diamond_outlined,
                    size: 16,
                    color: Colors.white.withValues(alpha: 0.85),
                  ),
                  const SizedBox(width: 4),
                  Text(
                    _formatBalance(diamondBalance),
                    style: theme.textTheme.labelLarge?.copyWith(
                      color: const Color(0xFF93C5FD),
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                ],
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
