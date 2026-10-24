import 'package:cached_network_image/cached_network_image.dart';
import 'dart:async';

import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/cache/image_disk_cache.dart';
import 'package:sara_live_app/core/widgets/cached_avatar.dart';
import 'package:sara_live_app/data/one_to_one/one_to_one_models.dart';
import 'package:sara_live_app/ui/contracts/one_to_one_ui.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/viewmodel/one_to_one_view_model.dart';

class OneToOnePage extends StatefulWidget {
  const OneToOnePage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<OneToOnePage> createState() => _OneToOnePageState();
}

class _OneToOnePageState extends State<OneToOnePage>
    with AutomaticKeepAliveClientMixin
    implements OneToOneUi {
  late final OneToOneViewModel _vm;
  final _scrollController = ScrollController();

  bool _initialLoading = true;
  bool _loadingMore = false;
  List<OneToOneRoomItem> _items = const [];
  int _total = 0;
  int _statusFilter = OneToOneStatusFilter.all;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _vm = OneToOneViewModel(scope: widget.scope, ui: this);
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) => _vm.loadInitial());
  }

  @override
  void dispose() {
    _scrollController.dispose();
    _vm.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients) return;
    final pos = _scrollController.position;
    if (pos.pixels >= pos.maxScrollExtent - 280) {
      _vm.loadMore();
    }
  }

  @override
  void showOneToOneLoading({required bool initial}) {
    if (!mounted) return;
    setState(() => _initialLoading = initial);
  }

  @override
  void hideOneToOneLoading() {
    if (!mounted) return;
    setState(() => _initialLoading = false);
  }

  @override
  void renderOneToOneList({
    required List<OneToOneRoomItem> items,
    required int total,
    required int statusFilter,
    required bool hasMore,
    required bool loadingMore,
  }) {
    if (!mounted) return;
    setState(() {
      _items = items;
      _total = total;
      _statusFilter = statusFilter;
      _loadingMore = loadingMore;
      _initialLoading = false;
    });
  }

  @override
  void showOneToOneError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }

  String _billingLabel(double billing) {
    if (billing <= 0) return '免费通话';
    final s = billing == billing.roundToDouble()
        ? '${billing.toInt()}'
        : billing.toStringAsFixed(1);
    return '$s 钻/分钟';
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final theme = Theme.of(context);
    final cs = theme.colorScheme;

    return Scaffold(
      body: AppPageFrame(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
              child: Text(
                '1v1',
                style: theme.textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.w800,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12),
              child: Wrap(
                spacing: 8,
                children: [
                  _FilterChip(
                    label: '全部',
                    selected: _statusFilter == OneToOneStatusFilter.all,
                    onTap: () =>
                        _vm.setStatusFilter(OneToOneStatusFilter.all),
                  ),
                  _FilterChip(
                    label: '在线',
                    selected: _statusFilter == OneToOneStatusFilter.online,
                    onTap: () =>
                        _vm.setStatusFilter(OneToOneStatusFilter.online),
                  ),
                  _FilterChip(
                    label: '离线',
                    selected: _statusFilter == OneToOneStatusFilter.offline,
                    onTap: () =>
                        _vm.setStatusFilter(OneToOneStatusFilter.offline),
                  ),
                ],
              ),
            ),
            if (_total > 0)
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 8, 16, 4),
                child: Text(
                  '共 $_total 位',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: cs.onSurface.withValues(alpha: 0.55),
                  ),
                ),
              ),
            Expanded(
              child: RefreshIndicator(
                onRefresh: () => _vm.refresh(forceRefresh: true),
                child: _initialLoading && _items.isEmpty
                    ? ListView(
                        physics: const AlwaysScrollableScrollPhysics(),
                        children: const [
                          SizedBox(height: 120),
                          Center(child: CircularProgressIndicator()),
                        ],
                      )
                    : _items.isEmpty
                        ? ListView(
                            physics: const AlwaysScrollableScrollPhysics(),
                            children: [
                              const SizedBox(height: 80),
                              Icon(
                                Icons.videocam_off_outlined,
                                size: 48,
                                color: cs.onSurface.withValues(alpha: 0.35),
                              ),
                              const SizedBox(height: 12),
                              Text(
                                '暂无 1v1 房间',
                                textAlign: TextAlign.center,
                                style: theme.textTheme.titleMedium,
                              ),
                            ],
                          )
                        : ListView.separated(
                            controller: _scrollController,
                            physics: const AlwaysScrollableScrollPhysics(),
                            padding: const EdgeInsets.fromLTRB(12, 8, 12, 24),
                            itemCount:
                                _items.length + (_loadingMore ? 1 : 0),
                            separatorBuilder: (_, __) =>
                                const SizedBox(height: 10),
                            itemBuilder: (context, index) {
                              if (index >= _items.length) {
                                return const Padding(
                                  padding: EdgeInsets.symmetric(vertical: 16),
                                  child: Center(
                                    child: CircularProgressIndicator(),
                                  ),
                                );
                              }
                              return _OneToOneRoomTile(
                                scope: widget.scope,
                                item: _items[index],
                                billingLabel: _billingLabel,
                                onError: showOneToOneError,
                              );
                            },
                          ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _FilterChip extends StatelessWidget {
  const _FilterChip({
    required this.label,
    required this.selected,
    required this.onTap,
  });

  final String label;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return FilterChip(
      label: Text(label),
      selected: selected,
      onSelected: (_) => onTap(),
      showCheckmark: false,
    );
  }
}

class _OneToOneRoomTile extends StatelessWidget {
  const _OneToOneRoomTile({
    required this.scope,
    required this.item,
    required this.billingLabel,
    required this.onError,
  });

  final AppScope scope;
  final OneToOneRoomItem item;
  final String Function(double) billingLabel;
  final void Function(String message) onError;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final name = item.anchorNickname.isNotEmpty
        ? item.anchorNickname
        : (item.title.isNotEmpty ? item.title : '主播 ${item.roomId}');

    return Material(
      color: cs.surfaceContainerHighest.withValues(alpha: 0.45),
      borderRadius: BorderRadius.circular(14),
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: () {
          if (!item.allowCallIcon) {
            onError('当前不可发起通话');
            return;
          }
          final anchorId = int.tryParse(item.roomId);
          if (anchorId == null || anchorId <= 0) {
            onError('无效的主播 ID');
            return;
          }
          final name = item.anchorNickname.isNotEmpty
              ? item.anchorNickname
              : item.title;
          unawaited(
            scope.callController.startOneToOneCall(
              targetUserId: anchorId,
              peerNickname: name,
              peerAvatar: item.anchorAvatar,
            ),
          );
        },
        child: Padding(
          padding: const EdgeInsets.all(10),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Stack(
                clipBehavior: Clip.none,
                children: [
                  _CoverThumb(url: item.cover, fallback: name),
                  Positioned(
                    right: -2,
                    bottom: -2,
                    child: CachedAvatar(
                      url: item.anchorAvatar,
                      radius: 18,
                      fallbackText: name.isNotEmpty ? name[0] : '?',
                    ),
                  ),
                  if (item.isWsOnline)
                    Positioned(
                      left: 4,
                      top: 4,
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 6,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: const Color(0xFF22C55E),
                          borderRadius: BorderRadius.circular(6),
                        ),
                        child: Text(
                          '在线',
                          style: theme.textTheme.labelSmall?.copyWith(
                            color: Colors.white,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ),
                    ),
                ],
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      name,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                    if (item.title.isNotEmpty &&
                        item.title != item.anchorNickname) ...[
                      const SizedBox(height: 2),
                      Text(
                        item.title,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: cs.onSurface.withValues(alpha: 0.6),
                        ),
                      ),
                    ],
                    const SizedBox(height: 6),
                    Text(
                      billingLabel(item.billing),
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: cs.primary,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
              if (item.allowCallIcon)
                Icon(
                  Icons.videocam_rounded,
                  color: cs.primary.withValues(alpha: 0.85),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

class _CoverThumb extends StatelessWidget {
  const _CoverThumb({required this.url, required this.fallback});

  final String url;
  final String fallback;

  @override
  Widget build(BuildContext context) {
    const w = 72.0;
    const h = 96.0;
    final cs = Theme.of(context).colorScheme;
    if (url.isEmpty) {
      return Container(
        width: w,
        height: h,
        decoration: BoxDecoration(
          color: cs.surfaceContainerHigh,
          borderRadius: BorderRadius.circular(10),
        ),
        alignment: Alignment.center,
        child: Icon(
          Icons.image_outlined,
          color: cs.onSurface.withValues(alpha: 0.35),
        ),
      );
    }
    return ClipRRect(
      borderRadius: BorderRadius.circular(10),
      child: CachedNetworkImage(
        imageUrl: url,
        cacheManager: ImageDiskCache.instance.manager,
        width: w,
        height: h,
        fit: BoxFit.cover,
        memCacheWidth: (w * 2).round(),
        memCacheHeight: (h * 2).round(),
        placeholder: (_, __) => Container(
          width: w,
          height: h,
          color: cs.surfaceContainerHigh,
          child: const Center(
            child: SizedBox(
              width: 22,
              height: 22,
              child: CircularProgressIndicator(strokeWidth: 2),
            ),
          ),
        ),
        errorWidget: (_, __, ___) => Container(
          width: w,
          height: h,
          color: cs.surfaceContainerHigh,
          alignment: Alignment.center,
          child: Text(
            fallback.isNotEmpty ? fallback[0] : '?',
            style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w700),
          ),
        ),
      ),
    );
  }
}
