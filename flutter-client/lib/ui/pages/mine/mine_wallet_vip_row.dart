import 'package:flutter/material.dart';

import '../../../data/vip/vip_models.dart';
import '../../../l10n/app_localizations.dart';

/// 「我的」顶栏：左侧 [MineWalletStrip] + 右侧 [MineVipCard] 同一行、两个独立 Card。
class MineWalletVipRow extends StatelessWidget {
  const MineWalletVipRow({
    super.key,
    required this.l10n,
    required this.gold,
    required this.diamond,
    required this.vipLevelFallback,
    this.vip,
    this.hideGoldDiamond = false,
    this.onVipTap,
  });

  final AppLocalizations l10n;
  final String gold;
  final String diamond;
  final int vipLevelFallback;
  final VipDetail? vip;
  final bool hideGoldDiamond;
  final VoidCallback? onVipTap;

  static const _gap = 10.0;

  @override
  Widget build(BuildContext context) {
    final vipCard = MineVipCard(
      l10n: l10n,
      value: vip?.mineTitleValue ?? 'Lv.$vipLevelFallback',
      subtitle: vip?.mineSummary,
      onTap: onVipTap,
    );

    if (hideGoldDiamond) {
      return IntrinsicHeight(
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [Expanded(child: vipCard)],
        ),
      );
    }

    return IntrinsicHeight(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Expanded(
            flex: 3,
            child: MineWalletStrip(
              l10n: l10n,
              gold: gold,
              diamond: diamond,
            ),
          ),
          const SizedBox(width: _gap),
          Expanded(flex: 2, child: vipCard),
        ],
      ),
    );
  }
}

/// 金币 + 钻石（独立 Card）。
class MineWalletStrip extends StatelessWidget {
  const MineWalletStrip({
    super.key,
    required this.l10n,
    required this.gold,
    required this.diamond,
  });

  final AppLocalizations l10n;
  final String gold;
  final String diamond;

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      margin: EdgeInsets.zero,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 8),
        child: Center(
          child: Row(
            children: [
              _BalanceMetric(
                label: l10n.walletGold,
                value: gold,
                icon: Icons.monetization_on_outlined,
              ),
              _BalanceMetric(
                label: l10n.walletDiamond,
                value: diamond,
                icon: Icons.diamond_outlined,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// VIP（独立 Card）。
class MineVipCard extends StatelessWidget {
  const MineVipCard({
    super.key,
    required this.l10n,
    required this.value,
    this.subtitle,
    this.onTap,
  });

  final AppLocalizations l10n;
  final String value;
  final String? subtitle;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      margin: EdgeInsets.zero,
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          child: Padding(
            padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 8),
            child: Center(
              child: _MetricColumn(
                icon: Icons.workspace_premium_outlined,
                value: value,
                label: l10n.vipTitle,
                subtitle: subtitle,
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _BalanceMetric extends StatelessWidget {
  const _BalanceMetric({
    required this.label,
    required this.value,
    required this.icon,
  });

  final String label;
  final String value;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: _MetricColumn(
        icon: icon,
        value: value,
        label: label,
      ),
    );
  }
}

class _MetricColumn extends StatelessWidget {
  const _MetricColumn({
    required this.icon,
    required this.value,
    required this.label,
    this.subtitle,
  });

  final IconData icon;
  final String value;
  final String label;
  final String? subtitle;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    return Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Icon(icon, size: 18, color: cs.primary.withValues(alpha: 0.9)),
        const SizedBox(height: 6),
        Text(
          value,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          textAlign: TextAlign.center,
          style: theme.textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w800,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          label,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          textAlign: TextAlign.center,
          style: theme.textTheme.labelSmall?.copyWith(
            color: cs.onSurface.withValues(alpha: 0.55),
          ),
        ),
        if (subtitle != null && subtitle!.isNotEmpty) ...[
          const SizedBox(height: 2),
          Text(
            subtitle!,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            textAlign: TextAlign.center,
            style: theme.textTheme.labelSmall?.copyWith(
              fontSize: 10,
              color: cs.onSurface.withValues(alpha: 0.45),
            ),
          ),
        ],
      ],
    );
  }
}
