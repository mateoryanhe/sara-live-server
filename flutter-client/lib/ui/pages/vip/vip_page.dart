import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/cache/image_disk_cache.dart';
import 'package:sara_live_app/data/vip/vip_models.dart';
import 'package:sara_live_app/ui/contracts/vip_ui.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/viewmodel/vip_view_model.dart';

class VipPage extends StatefulWidget {
  const VipPage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<VipPage> createState() => _VipPageState();
}

class _VipPageState extends State<VipPage> implements VipUi {
  late final VipViewModel _vm;
  bool _loading = false;
  VipDetail? _detail;
  List<VipCfgItem> _cfgList = const [];
  String? _error;

  PageController? _pageController;
  final ScrollController _chipScroll = ScrollController();
  int _pageIndex = 0;
  bool _pagerReady = false;

  @override
  void initState() {
    super.initState();
    _vm = VipViewModel(scope: widget.scope, ui: this);
    _vm.onPushRefresh = () {
      if (!mounted) return;
      _pagerReady = false;
    };
    WidgetsBinding.instance.addPostFrameCallback((_) => _vm.load());
  }

  @override
  void dispose() {
    _vm.dispose();
    _pageController?.dispose();
    _chipScroll.dispose();
    super.dispose();
  }

  @override
  void showVipLoading() {
    if (!mounted) return;
    setState(() {
      _loading = true;
      _error = null;
    });
  }

  @override
  void hideVipLoading() {
    if (!mounted) return;
    setState(() => _loading = false);
  }

  @override
  void renderVip({
    required VipDetail detail,
    required List<VipCfgItem> cfgList,
  }) {
    if (!mounted) return;
    final nextIndex = _defaultPageIndex(detail, cfgList);
    setState(() {
      _detail = detail;
      _cfgList = cfgList;
      _error = null;
      if (!_pagerReady && cfgList.isNotEmpty) {
        _pageIndex = nextIndex;
        _pageController?.dispose();
        _pageController = PageController(initialPage: nextIndex);
        _pagerReady = true;
      } else if (cfgList.isNotEmpty) {
        _pageIndex = _pageIndex.clamp(0, cfgList.length - 1);
      }
    });
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollChipIntoView(_pageIndex);
    });
  }

  @override
  void showVipError(String message) {
    if (!mounted) return;
    setState(() => _error = message);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  int _defaultPageIndex(VipDetail detail, List<VipCfgItem> list) {
    if (list.isEmpty) return 0;
    if (detail.vipLevel > 0) {
      final i = list.indexWhere((c) => c.level == detail.vipLevel);
      if (i >= 0) return i;
    }
    if (detail.nextLevel > 0) {
      final i = list.indexWhere((c) => c.level == detail.nextLevel);
      if (i >= 0) return i;
    }
    return 0;
  }

  String _fmt(double v) {
    if (v == v.roundToDouble()) return '${v.toInt()}';
    return v.toStringAsFixed(2);
  }

  Future<void> _onRefresh() async {
    _pagerReady = false;
    await _vm.load(forceRefresh: true);
  }

  void _goToPage(int index) {
    if (_pageController == null || index < 0 || index >= _cfgList.length) {
      return;
    }
    setState(() => _pageIndex = index);
    _pageController!.animateToPage(
      index,
      duration: const Duration(milliseconds: 280),
      curve: Curves.easeOutCubic,
    );
    _scrollChipIntoView(index);
  }

  void _scrollChipIntoView(int index) {
    if (!_chipScroll.hasClients || index < 0) return;
    // 近似：每芯片约 88px + 间距
    final target = (index * 92.0) - 24;
    _chipScroll.animateTo(
      target.clamp(0.0, _chipScroll.position.maxScrollExtent),
      duration: const Duration(milliseconds: 240),
      curve: Curves.easeOut,
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final d = _detail;

    return Scaffold(
      appBar: AppBar(title: const Text('VIP')),
      body: AppPageFrame(
        child: _loading && d == null
            ? const Center(child: CircularProgressIndicator())
            : (_error != null && d == null)
                ? Padding(
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(_error!, textAlign: TextAlign.center),
                        const SizedBox(height: 12),
                        FilledButton(
                          onPressed: () => _vm.load(forceRefresh: true),
                          child: const Text('重试'),
                        ),
                      ],
                    ),
                  )
                : d == null
                    ? const SizedBox.shrink()
                    : Column(
                        children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
                            child: _VipHeroCard(detail: d, fmt: _fmt),
                          ),
                          if (_cfgList.isNotEmpty) ...[
                            Padding(
                              padding: const EdgeInsets.fromLTRB(16, 4, 16, 0),
                              child: Row(
                                children: [
                                  Text(
                                    '左右滑动查看等级',
                                    style: theme.textTheme.labelMedium?.copyWith(
                                      color: cs.onSurface.withValues(alpha: 0.5),
                                      fontWeight: FontWeight.w600,
                                    ),
                                  ),
                                  const Spacer(),
                                  Text(
                                    '${_pageIndex + 1}/${_cfgList.length}',
                                    style: theme.textTheme.labelMedium?.copyWith(
                                      color: cs.onSurface.withValues(alpha: 0.45),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            const SizedBox(height: 8),
                            SizedBox(
                              height: 40,
                              child: ListView.separated(
                                controller: _chipScroll,
                                scrollDirection: Axis.horizontal,
                                padding:
                                    const EdgeInsets.symmetric(horizontal: 16),
                                itemCount: _cfgList.length,
                                separatorBuilder: (_, __) =>
                                    const SizedBox(width: 8),
                                itemBuilder: (context, index) {
                                  final cfg = _cfgList[index];
                                  final selected = index == _pageIndex;
                                  final unlocked = d.vipLevel >= cfg.level;
                                  final label = cfg.levelName.isNotEmpty
                                      ? cfg.levelName
                                      : 'VIP${cfg.level}';
                                  return ChoiceChip(
                                    selected: selected,
                                    avatar: unlocked
                                        ? null
                                        : Icon(
                                            Icons.lock_outline_rounded,
                                            size: 14,
                                            color: selected
                                                ? cs.onPrimary
                                                : cs.onSurface
                                                    .withValues(alpha: 0.45),
                                          ),
                                    label: Text(
                                      label,
                                      style: TextStyle(
                                        fontWeight: selected
                                            ? FontWeight.w700
                                            : FontWeight.w500,
                                        fontSize: 13,
                                      ),
                                    ),
                                    onSelected: (_) => _goToPage(index),
                                  );
                                },
                              ),
                            ),
                            const SizedBox(height: 8),
                            Expanded(
                              child: PageView.builder(
                                controller: _pageController,
                                itemCount: _cfgList.length,
                                onPageChanged: (i) {
                                  setState(() => _pageIndex = i);
                                  _scrollChipIntoView(i);
                                },
                                itemBuilder: (context, index) {
                                  final cfg = _cfgList[index];
                                  return RefreshIndicator(
                                    onRefresh: _onRefresh,
                                    child: _VipLevelPage(
                                      cfg: cfg,
                                      unlocked: d.vipLevel >= cfg.level,
                                      isCurrent: d.vipLevel == cfg.level,
                                      fmt: _fmt,
                                    ),
                                  );
                                },
                              ),
                            ),
                            Padding(
                              padding: const EdgeInsets.fromLTRB(16, 4, 16, 12),
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  for (var i = 0; i < _cfgList.length; i++)
                                    AnimatedContainer(
                                      duration:
                                          const Duration(milliseconds: 200),
                                      margin: const EdgeInsets.symmetric(
                                          horizontal: 3),
                                      width: i == _pageIndex ? 16 : 6,
                                      height: 6,
                                      decoration: BoxDecoration(
                                        color: i == _pageIndex
                                            ? cs.primary
                                            : cs.onSurface
                                                .withValues(alpha: 0.2),
                                        borderRadius: BorderRadius.circular(99),
                                      ),
                                    ),
                                ],
                              ),
                            ),
                          ] else
                            Expanded(
                              child: RefreshIndicator(
                                onRefresh: _onRefresh,
                                child: ListView(
                                  physics:
                                      const AlwaysScrollableScrollPhysics(),
                                  padding: const EdgeInsets.all(24),
                                  children: [
                                    Text(
                                      '暂无 VIP 等级配置',
                                      textAlign: TextAlign.center,
                                      style: theme.textTheme.bodyMedium
                                          ?.copyWith(
                                        color: cs.onSurface
                                            .withValues(alpha: 0.55),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                        ],
                      ),
      ),
    );
  }
}

class _VipHeroCard extends StatelessWidget {
  const _VipHeroCard({required this.detail, required this.fmt});

  final VipDetail detail;
  final String Function(double) fmt;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final icon = detail.currentCfg?.levelIcon ?? '';
    final title = detail.hasVip
        ? (detail.levelName.isNotEmpty
            ? detail.levelName
            : 'Lv.${detail.vipLevel}')
        : '尚未开通 VIP';
    final progress = !detail.isMaxLevel && detail.nextUpgradeLimit > 0
        ? (detail.totalRecharge / detail.nextUpgradeLimit).clamp(0.0, 1.0)
        : 1.0;

    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(20),
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            cs.primary.withValues(alpha: 0.28),
            cs.surface,
          ],
        ),
        border: Border.all(color: cs.onSurface.withValues(alpha: 0.06)),
      ),
      padding: const EdgeInsets.all(18),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              _VipIcon(url: icon, size: 56),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: theme.textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      detail.hasVip
                          ? '当前等级 Lv.${detail.vipLevel}'
                          : '累计充值可自动升级',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: cs.onSurface.withValues(alpha: 0.6),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Text(
            '累计充值 ${fmt(detail.totalRecharge)} 金币',
            style: theme.textTheme.bodyMedium?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          if (detail.isMaxLevel)
            Text(
              '已达最高等级',
              style: theme.textTheme.bodySmall?.copyWith(color: cs.primary),
            )
          else ...[
            ClipRRect(
              borderRadius: BorderRadius.circular(999),
              child: LinearProgressIndicator(
                value: progress,
                minHeight: 8,
                backgroundColor: cs.onSurface.withValues(alpha: 0.08),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '距 Lv.${detail.nextLevel} 还需 ${fmt(detail.rechargeToNextLevel)} 金币'
              '（门槛 ${fmt(detail.nextUpgradeLimit)}）',
              style: theme.textTheme.bodySmall?.copyWith(
                color: cs.onSurface.withValues(alpha: 0.65),
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _VipLevelPage extends StatelessWidget {
  const _VipLevelPage({
    required this.cfg,
    required this.unlocked,
    required this.isCurrent,
    required this.fmt,
  });

  final VipCfgItem cfg;
  final bool unlocked;
  final bool isCurrent;
  final String Function(double) fmt;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final title = cfg.levelName.isEmpty ? 'VIP${cfg.level}' : cfg.levelName;
    final privileges = cfg.privilegeList;

    return ListView(
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.fromLTRB(16, 4, 16, 24),
      children: [
        Card(
          clipBehavior: Clip.antiAlias,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
            side: isCurrent
                ? BorderSide(
                    color: cs.primary.withValues(alpha: 0.55),
                    width: 1.2,
                  )
                : BorderSide.none,
          ),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 12),
            child: Row(
              children: [
                _VipIcon(url: cfg.levelIcon, size: 52),
                const SizedBox(width: 14),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Flexible(
                            child: Text(
                              title,
                              style: theme.textTheme.titleMedium?.copyWith(
                                fontWeight: FontWeight.w800,
                              ),
                            ),
                          ),
                          if (isCurrent) ...[
                            const SizedBox(width: 8),
                            Container(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 8,
                                vertical: 2,
                              ),
                              decoration: BoxDecoration(
                                color: cs.primary.withValues(alpha: 0.16),
                                borderRadius: BorderRadius.circular(999),
                              ),
                              child: Text(
                                '当前',
                                style: theme.textTheme.labelSmall?.copyWith(
                                  color: cs.primary,
                                  fontWeight: FontWeight.w700,
                                ),
                              ),
                            ),
                          ],
                        ],
                      ),
                      const SizedBox(height: 4),
                      Text(
                        'Lv.${cfg.level} · 门槛 ${fmt(cfg.upgradeRechargeLimit)} 金币'
                        '${unlocked ? ' · 已达成' : ' · 未解锁'}',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: cs.onSurface.withValues(alpha: 0.6),
                        ),
                      ),
                    ],
                  ),
                ),
                Icon(
                  unlocked
                      ? Icons.check_circle
                      : Icons.lock_outline_rounded,
                  color: unlocked
                      ? cs.primary
                      : cs.onSurface.withValues(alpha: 0.35),
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 12),
        Text(
          '本等级特权',
          style: theme.textTheme.titleSmall?.copyWith(
            color: cs.onSurface.withValues(alpha: 0.55),
            fontWeight: FontWeight.w700,
          ),
        ),
        const SizedBox(height: 8),
        if (privileges.isEmpty)
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 24),
            child: Text(
              '暂无特权配置',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodySmall?.copyWith(
                color: cs.onSurface.withValues(alpha: 0.45),
              ),
            ),
          )
        else
          Card(
            clipBehavior: Clip.antiAlias,
            child: Opacity(
              opacity: unlocked ? 1 : 0.55,
              child: Column(
                children: [
                  for (var i = 0; i < privileges.length; i++) ...[
                    if (i > 0)
                      Divider(
                        height: 1,
                        indent: 56,
                        color: cs.onSurface.withValues(alpha: 0.06),
                      ),
                    _PrivilegeTile(item: privileges[i]),
                  ],
                ],
              ),
            ),
          ),
      ],
    );
  }
}

class _PrivilegeTile extends StatelessWidget {
  const _PrivilegeTile({required this.item});

  final VipPrivilegeItem item;

  @override
  Widget build(BuildContext context) {
    final title = item.title.resolve();
    final desc = item.desc.resolve();
    return ListTile(
      leading: item.icon.isEmpty
          ? Icon(
              Icons.stars_rounded,
              color: Theme.of(context).colorScheme.primary,
            )
          : _VipIcon(url: item.icon, size: 36),
      title: Text(
        title.isEmpty ? item.typeLabel : title,
        style: const TextStyle(fontWeight: FontWeight.w600),
      ),
      subtitle: desc.isEmpty ? Text(item.typeLabel) : Text(desc),
    );
  }
}

class _VipIcon extends StatelessWidget {
  const _VipIcon({required this.url, required this.size});

  final String url;
  final double size;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    if (url.isEmpty) {
      return Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          color: cs.primary.withValues(alpha: 0.15),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Icon(Icons.workspace_premium_rounded, color: cs.primary),
      );
    }
    return ClipRRect(
      borderRadius: BorderRadius.circular(12),
      child: CachedNetworkImage(
        imageUrl: url,
        cacheManager: ImageDiskCache.instance.manager,
        width: size,
        height: size,
        fit: BoxFit.cover,
        placeholder: (_, __) => SizedBox(
          width: size,
          height: size,
          child: const Center(
            child: SizedBox(
              width: 16,
              height: 16,
              child: CircularProgressIndicator(strokeWidth: 2),
            ),
          ),
        ),
        errorWidget: (_, __, ___) => Icon(
          Icons.workspace_premium_rounded,
          color: cs.primary,
          size: size * 0.55,
        ),
      ),
    );
  }
}
