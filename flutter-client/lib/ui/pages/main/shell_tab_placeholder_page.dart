import 'package:flutter/material.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';

enum ShellPlaceholderKind { live, game }

/// Tab 占位页：业务列表/房间稍后接入。
class ShellTabPlaceholderPage extends StatelessWidget {
  const ShellTabPlaceholderPage({
    super.key,
    required this.kind,
  });

  final ShellPlaceholderKind kind;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context);
    final (title, subtitle, icon) = switch (kind) {
      ShellPlaceholderKind.live => (
          l10n.placeholderLiveTitle,
          l10n.placeholderLiveSubtitle,
          Icons.live_tv_rounded,
        ),
      ShellPlaceholderKind.game => (
          l10n.placeholderGameTitle,
          l10n.placeholderGameSubtitle,
          Icons.sports_esports_rounded,
        ),
    };

    return Scaffold(
      body: AppPageFrame(
        child: Center(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 28),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  icon,
                  size: 48,
                  color: theme.colorScheme.primary.withValues(alpha: 0.7),
                ),
                const SizedBox(height: 12),
                Text(
                  title,
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w700,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 6),
                Text(
                  subtitle,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurface.withValues(alpha: 0.55),
                  ),
                  textAlign: TextAlign.center,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
