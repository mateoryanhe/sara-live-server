import 'dart:async';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/cache/image_disk_cache.dart';
import 'package:sara_live_app/core/http/api_exception.dart';
import 'package:sara_live_app/data/live_room/live_room_models.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/ui/pages/live/live_room_watch_page.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';

/// 观众 Tab：Hot 直播中房间列表。
class LiveRoomListPage extends StatefulWidget {
  const LiveRoomListPage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<LiveRoomListPage> createState() => _LiveRoomListPageState();
}

class _LiveRoomListPageState extends State<LiveRoomListPage>
    with AutomaticKeepAliveClientMixin {
  final _scrollController = ScrollController();

  bool _initialLoading = true;
  bool _loadingMore = false;
  List<LiveRoomListItem> _items = const [];
  int _pageIndex = 1;
  bool _hasMore = true;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) => _loadInitial());
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients || _loadingMore || !_hasMore) return;
    final pos = _scrollController.position;
    if (pos.pixels >= pos.maxScrollExtent - 280) {
      unawaited(_loadMore());
    }
  }

  Future<void> _loadInitial() async {
    setState(() => _initialLoading = true);
    try {
      final page = await widget.scope.liveRoomRepository.getHotLiveRoomList(
        pageIndex: 1,
        pageSize: 20,
      );
      if (!mounted) return;
      setState(() {
        _items = page.list;
        _pageIndex = page.pageIndex;
        _hasMore = page.hasMore;
        _initialLoading = false;
      });
    } on ApiException catch (e) {
      if (!mounted) return;
      setState(() => _initialLoading = false);
      _showError(e.message);
    } catch (e) {
      if (!mounted) return;
      setState(() => _initialLoading = false);
      _showError(e.toString());
    }
  }

  Future<void> _loadMore() async {
    if (_loadingMore || !_hasMore) return;
    setState(() => _loadingMore = true);
    try {
      final next = _pageIndex + 1;
      final page = await widget.scope.liveRoomRepository.getHotLiveRoomList(
        pageIndex: next,
        pageSize: 20,
      );
      if (!mounted) return;
      setState(() {
        _items = [..._items, ...page.list];
        _pageIndex = page.pageIndex;
        _hasMore = page.hasMore;
        _loadingMore = false;
      });
    } on ApiException catch (e) {
      if (!mounted) return;
      setState(() => _loadingMore = false);
      _showError(e.message);
    } catch (_) {
      if (!mounted) return;
      setState(() => _loadingMore = false);
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }

  Future<void> _openRoom(LiveRoomListItem item) async {
    if (widget.scope.callController.isActive) {
      _showError(AppLocalizations.of(context).liveWatchBlockedInCall);
      return;
    }
    final roomId = int.tryParse(item.roomId);
    if (roomId == null || roomId <= 0) {
      _showError('Invalid room');
      return;
    }
    await Navigator.of(context).push<void>(
      MaterialPageRoute(
        builder: (_) => LiveRoomWatchPage(
          scope: widget.scope,
          roomId: roomId,
          title: item.title.isNotEmpty ? item.title : item.anchorNickname,
          coverUrl: item.cover,
          anchorName: item.anchorNickname,
          anchorAvatar: item.anchorAvatar,
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return AppPageFrame(
      child: RefreshIndicator(
        onRefresh: _loadInitial,
        child: CustomScrollView(
          controller: _scrollController,
          physics: const AlwaysScrollableScrollPhysics(),
          slivers: [
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
                child: Text(
                  l10n.liveHotListTitle,
                  style: theme.textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.w800,
                  ),
                ),
              ),
            ),
            if (_initialLoading)
              const SliverFillRemaining(
                child: Center(child: CircularProgressIndicator()),
              )
            else if (_items.isEmpty)
              SliverFillRemaining(
                hasScrollBody: false,
                child: Center(
                  child: Text(
                    l10n.liveHotListEmpty,
                    style: theme.textTheme.bodyLarge?.copyWith(
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
                ),
              )
            else ...[
              SliverPadding(
                padding: const EdgeInsets.fromLTRB(12, 0, 12, 8),
                sliver: SliverGrid.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    mainAxisSpacing: 10,
                    crossAxisSpacing: 10,
                    childAspectRatio: 0.72,
                  ),
                  itemCount: _items.length,
                  itemBuilder: (context, index) {
                    final item = _items[index];
                    return _LiveRoomCard(
                      item: item,
                      onTap: () => _openRoom(item),
                    );
                  },
                ),
              ),
              if (_loadingMore)
                const SliverToBoxAdapter(
                  child: Padding(
                    padding: EdgeInsets.symmetric(vertical: 16),
                    child: Center(child: CircularProgressIndicator()),
                  ),
                ),
              const SliverToBoxAdapter(child: SizedBox(height: 16)),
            ],
          ],
        ),
      ),
    );
  }
}

class _LiveRoomCard extends StatelessWidget {
  const _LiveRoomCard({required this.item, required this.onTap});

  final LiveRoomListItem item;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final l10n = AppLocalizations.of(context);

    return Material(
      color: cs.surfaceContainerLow,
      borderRadius: BorderRadius.circular(14),
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Expanded(
              child: Stack(
                fit: StackFit.expand,
                children: [
                  _LiveRoomCoverImage(item: item, colorScheme: cs),
                  if (item.isLive)
                    Positioned(
                      left: 8,
                      top: 8,
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 6,
                          vertical: 3,
                        ),
                        decoration: BoxDecoration(
                          color: const Color(0xFFE11D48),
                          borderRadius: BorderRadius.circular(5),
                        ),
                        child: Text(
                          l10n.liveBadgeOnAir,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 10,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ),
                    ),
                  Positioned(
                    right: 8,
                    bottom: 8,
                    child: Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 6,
                        vertical: 3,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.black54,
                        borderRadius: BorderRadius.circular(5),
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Icon(
                            Icons.people_alt_outlined,
                            size: 12,
                            color: Colors.white,
                          ),
                          const SizedBox(width: 3),
                          Text(
                            '${item.onlineCount}',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 10,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(8, 8, 8, 10),
              child: Text(
                item.title.trim().isNotEmpty
                    ? item.title.trim()
                    : item.anchorNickname.trim(),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: theme.textTheme.labelLarge?.copyWith(
                  fontWeight: FontWeight.w700,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// 封面优先；无封面或加载失败时回退主播头像。
class _LiveRoomCoverImage extends StatelessWidget {
  const _LiveRoomCoverImage({
    required this.item,
    required this.colorScheme,
  });

  final LiveRoomListItem item;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    final cover = item.cover.trim();
    final avatar = item.anchorAvatar.trim();
    if (cover.isEmpty && avatar.isEmpty) {
      return _placeholder(colorScheme);
    }
    final primary = cover.isNotEmpty ? cover : avatar;
    return CachedNetworkImage(
      imageUrl: primary,
      fit: BoxFit.cover,
      cacheManager: ImageDiskCache.instance.manager,
      errorWidget: (_, __, ___) {
        if (cover.isNotEmpty && avatar.isNotEmpty) {
          return CachedNetworkImage(
            imageUrl: avatar,
            fit: BoxFit.cover,
            cacheManager: ImageDiskCache.instance.manager,
            errorWidget: (_, __, ___) => _placeholder(colorScheme),
          );
        }
        return _placeholder(colorScheme);
      },
    );
  }

  Widget _placeholder(ColorScheme cs) {
    return ColoredBox(
      color: cs.surfaceContainerHigh,
      child: Center(
        child: Icon(
          Icons.person_outline,
          color: cs.onSurface.withValues(alpha: 0.35),
          size: 36,
        ),
      ),
    );
  }
}
