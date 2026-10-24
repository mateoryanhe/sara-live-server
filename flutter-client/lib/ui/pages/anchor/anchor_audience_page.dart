import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/widgets/cached_avatar.dart';
import 'package:sara_live_app/data/live_room/live_room_models.dart';
import 'package:sara_live_app/ui/contracts/anchor_audience_ui.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/viewmodel/anchor_audience_view_model.dart';

/// 主播 Tab：本房间在线观众列表。
class AnchorAudiencePage extends StatefulWidget {
  const AnchorAudiencePage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<AnchorAudiencePage> createState() => _AnchorAudiencePageState();
}

class _AnchorAudiencePageState extends State<AnchorAudiencePage>
    with AutomaticKeepAliveClientMixin
    implements AnchorAudienceUi {
  late final AnchorAudienceViewModel _vm;
  final _scrollController = ScrollController();

  bool _initialLoading = true;
  List<OnlineUserItem> _items = const [];
  int _total = 0;
  bool _loadingMore = false;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _vm = AnchorAudienceViewModel(scope: widget.scope, ui: this);
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
    if (pos.pixels >= pos.maxScrollExtent - 200) {
      _vm.loadMore();
    }
  }

  @override
  void showAudienceLoading({required bool initial}) {
    if (!mounted) return;
    setState(() => _initialLoading = initial);
  }

  @override
  void hideAudienceLoading() {
    if (!mounted) return;
    setState(() => _initialLoading = false);
  }

  @override
  void renderAudienceList({
    required List<OnlineUserItem> items,
    required int total,
    required bool loadingMore,
  }) {
    if (!mounted) return;
    setState(() {
      _items = items;
      _total = total;
      _loadingMore = loadingMore;
      _initialLoading = false;
    });
  }

  @override
  void showAudienceError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  void showAudienceMessage(String message) {
    showAudienceError(message);
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
                '在线观众',
                style: theme.textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.w800,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
              child: Text(
                '本直播间当前在线 $_total 人（不含主播）',
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
                                Icons.people_outline,
                                size: 48,
                                color: cs.onSurface.withValues(alpha: 0.35),
                              ),
                              const SizedBox(height: 12),
                              Text(
                                '暂无在线观众',
                                textAlign: TextAlign.center,
                                style: theme.textTheme.titleMedium,
                              ),
                              const SizedBox(height: 8),
                              Text(
                                '观众进入你的直播间后会出现在这里',
                                textAlign: TextAlign.center,
                                style: theme.textTheme.bodySmall?.copyWith(
                                  color: cs.onSurface.withValues(alpha: 0.5),
                                ),
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
                                const SizedBox(height: 8),
                            itemBuilder: (context, index) {
                              if (index >= _items.length) {
                                return const Padding(
                                  padding: EdgeInsets.symmetric(vertical: 16),
                                  child: Center(
                                    child: CircularProgressIndicator(),
                                  ),
                                );
                              }
                              final item = _items[index];
                              final name = item.nickname.isNotEmpty
                                  ? item.nickname
                                  : '用户 ${item.userId}';
                              return ListTile(
                                tileColor: cs.surfaceContainerHighest
                                    .withValues(alpha: 0.45),
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
                                subtitle: Text(
                                  item.totalReward > 0
                                      ? '贡献 ${item.totalReward.toStringAsFixed(0)} 钻'
                                      : 'VIP ${item.vipLevel}',
                                  style: theme.textTheme.bodySmall,
                                ),
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
