import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';
import 'package:sara_live_app/app/app_routes.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/app/server_endpoint_store.dart';
import 'package:sara_live_app/core/widgets/cached_avatar.dart';
import 'package:sara_live_app/data/user/user_profile.dart';
import 'package:sara_live_app/data/vip/vip_models.dart';
import 'package:sara_live_app/ui/contracts/mine_ui.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/l10n/locale_controller.dart';
import 'package:sara_live_app/l10n/locale_picker.dart';
import 'package:sara_live_app/ui/theme/app_themes.dart';
import 'package:sara_live_app/ui/theme/theme_controller.dart';
import 'package:sara_live_app/ui/theme/ui_style_picker.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/viewmodel/mine_view_model.dart';

class MinePage extends StatefulWidget {
  const MinePage({
    super.key,
    required this.scope,
    required this.themeController,
    required this.localeController,
  });

  final AppScope scope;
  final ThemeController themeController;
  final LocaleController localeController;

  @override
  State<MinePage> createState() => _MinePageState();
}

class _MinePageState extends State<MinePage> implements MineUi {
  late final MineViewModel _vm;
  bool _loading = false;
  UserProfile? _profile;
  VipDetail? _vip;
  String? _error;
  String? _deviceId;
  String _apiBaseDisplay = '';

  @override
  void initState() {
    super.initState();
    _vm = MineViewModel(scope: widget.scope, ui: this);
    _apiBaseDisplay = widget.scope.serverEndpoint.apiBaseUrl;
    widget.scope.serverEndpoint.addListener(_onServerEndpointChanged);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _vm.load();
      _refreshDeviceId();
    });
  }

  void _onServerEndpointChanged() {
    if (!mounted) return;
    setState(
      () => _apiBaseDisplay = widget.scope.serverEndpoint.apiBaseUrl,
    );
  }

  @override
  void dispose() {
    widget.scope.serverEndpoint.removeListener(_onServerEndpointChanged);
    _vm.dispose();
    super.dispose();
  }

  Future<void> _refreshDeviceId() async {
    final id = widget.scope.session.storedDeviceId ??
        await widget.scope.session.getOrCreateDeviceId();
    if (!mounted) return;
    setState(() => _deviceId = id);
  }

  String _deviceIdSubtitle(AppLocalizations l10n) {
    final id = _deviceId;
    if (id == null || id.isEmpty) return l10n.deviceCodeNotSet;
    if (id.length <= 24) return id;
    return '${id.substring(0, 10)}…${id.substring(id.length - 8)}';
  }

  @override
  void showMineLoading() {
    if (!mounted) return;
    setState(() {
      _loading = true;
      _error = null;
    });
  }

  @override
  void hideMineLoading() {
    if (!mounted) return;
    setState(() => _loading = false);
  }

  @override
  void renderProfile(UserProfile profile) {
    if (!mounted) return;
    setState(() {
      _profile = profile;
      _error = null;
    });
  }

  @override
  void renderVipDetail(VipDetail detail) {
    if (!mounted) return;
    setState(() => _vip = detail);
  }

  @override
  void showMineError(String message) {
    if (!mounted) return;
    setState(() => _error = message);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  void showMineMessage(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  void onLoggedOut() {
    if (!mounted) return;
    context.go(AppRoutes.splash);
  }

  @override
  void onDeviceIdApplied(String deviceId) {
    if (!mounted) return;
    setState(() => _deviceId = deviceId);
  }

  @override
  void onApiBaseApplied(String apiBaseUrl) {
    if (!mounted) return;
    setState(() => _apiBaseDisplay = apiBaseUrl);
  }

  String _apiBaseSubtitle(AppLocalizations l10n) {
    final url = _apiBaseDisplay;
    if (url.isEmpty) return l10n.apiBaseNotSet;
    if (widget.scope.serverEndpoint.hasOverride) {
      return url;
    }
    return l10n.apiBaseUsingDefault(url);
  }

  Future<void> _editApiBase() async {
    final l10n = AppLocalizations.of(context);
    final store = widget.scope.serverEndpoint;
    final initial = store.hasOverride
        ? store.storedOverride!
        : store.apiBaseUrl;
    final result = await showDialog<({String? text, bool reset})>(
      context: context,
      builder: (ctx) => _ApiBaseEditDialog(initial: initial),
    );
    if (result == null) return;

    String? validationError;
    if (!result.reset) {
      final v = result.text ?? '';
      if (v.isEmpty) {
        validationError = l10n.apiBaseEmpty;
      } else {
        final parsed = ServerEndpointStore.validateInput(v);
        if (parsed is ApiBaseValidationInvalid) {
          validationError = l10n.apiBaseInvalid;
        } else if (parsed is ApiBaseValidationEmpty) {
          validationError = l10n.apiBaseEmpty;
        }
      }
    }

    await _vm.applyApiBaseUrl(
      input: result.text ?? '',
      resetToDefault: result.reset,
      validationError: validationError,
      successMessage: l10n.apiBaseSwitchSuccess,
      unchangedMessage: l10n.apiBaseUnchanged,
    );
  }

  Future<void> _editDeviceCode() async {
    final started = await _vm.beginDeviceCodeLogin();
    if (!started || !mounted) return;
    // The VM clears the old token and user-scoped cache before navigation.
    // The router also enforces this location while the pending flag is set,
    // including after a process restart.
    context.go(AppRoutes.deviceCodeLogin);
  }

  Future<void> _editNickname() async {
    final current = _profile?.nickname ?? '';
    final next = await showDialog<String>(
      context: context,
      builder: (ctx) => _NicknameEditDialog(initial: current),
    );
    if (next == null || next.isEmpty || next == current) return;
    await _vm.updateNickname(next);
  }

  Future<void> _pickAndUploadAvatar() async {
    final l10n = AppLocalizations.of(context);
    final picker = ImagePicker();
    final file = await picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 1024,
      maxHeight: 1024,
      imageQuality: 85,
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
      showMineError(l10n.avatarFormatError);
      return;
    }
    final bytes = await file.readAsBytes();
    var filename = file.name;
    if (!filename.contains('.')) {
      filename = '$filename.jpg';
    }
    await _vm.updateAvatar(bytes: bytes, filename: filename);
  }

  String _fmt(double v) {
    if (v == v.roundToDouble()) return '${v.toInt()}';
    return v.toStringAsFixed(2);
  }

  String _localeSubtitle(AppLocalizations l10n) {
    return widget.localeController.locale.languageCode == 'zh'
        ? l10n.languageChinese
        : l10n.languageEnglish;
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final p = _profile;

    return Scaffold(
      body: AppPageFrame(
        child: RefreshIndicator(
          onRefresh: () => _vm.load(forceRefresh: true),
          child: CustomScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            slivers: [
              // 无顶栏：底栏已有「我的」；刷新靠下拉 RefreshIndicator
              if (_loading && p == null)
                SliverFillRemaining(
                  hasScrollBody: false,
                  child: Padding(
                    padding: EdgeInsets.only(
                      top: MediaQuery.paddingOf(context).top,
                    ),
                    child: const Center(child: CircularProgressIndicator()),
                  ),
                ),
              if (_error != null && p == null)
                SliverFillRemaining(
                  hasScrollBody: false,
                  child: Padding(
                    padding: EdgeInsets.fromLTRB(
                      24,
                      24 + MediaQuery.paddingOf(context).top,
                      24,
                      24,
                    ),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(_error!, textAlign: TextAlign.center),
                        const SizedBox(height: 12),
                        FilledButton(
                          onPressed: () => _vm.load(forceRefresh: true),
                          child: Text(l10n.commonRetry),
                        ),
                      ],
                    ),
                  ),
                ),
              if (p != null) ...[
                // 顶栏去掉后，给状态栏留出安全区
                SliverToBoxAdapter(
                  child: SizedBox(height: MediaQuery.paddingOf(context).top),
                ),
                SliverToBoxAdapter(
                  child: _HeroHeader(
                    l10n: l10n,
                    profile: p,
                    busy: _loading,
                    onEditAvatar: _pickAndUploadAvatar,
                    onEditNickname: _editNickname,
                  ),
                ),
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
                    child: _WalletStrip(
                      l10n: l10n,
                      gold: _fmt(p.gold),
                      diamond: _fmt(p.diamond),
                      follow: '${p.followCount}',
                      fans: '${p.followerCount}',
                      hideGoldDiamond: p.isAnchor,
                    ),
                  ),
                ),
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 20, 16, 8),
                    child: Text(
                      l10n.accountSection,
                      style: theme.textTheme.titleSmall?.copyWith(
                        color: cs.onSurface.withValues(alpha: 0.55),
                        fontWeight: FontWeight.w700,
                        letterSpacing: 0.4,
                      ),
                    ),
                  ),
                ),
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: Card(
                      clipBehavior: Clip.antiAlias,
                      child: Column(
                        children: [
                          _InfoTile(
                            icon: Icons.workspace_premium_outlined,
                            title: l10n.vipTitle,
                            value: _vip?.mineTitleValue ?? 'Lv.${p.vipLevel}',
                            subtitle: _vip?.mineSummary,
                            trailingIcon: Icons.chevron_right_rounded,
                            onTap: () => context.push(AppRoutes.vip),
                          ),
                          Divider(
                            height: 1,
                            indent: 56,
                            color: cs.onSurface.withValues(alpha: 0.06),
                          ),
                          _InfoTile(
                            icon: Icons.qr_code_2_rounded,
                            title: l10n.inviteCodeTitle,
                            value: p.shareCode.isEmpty ? '-' : p.shareCode,
                            trailingIcon: p.shareCode.isEmpty
                                ? null
                                : Icons.copy_rounded,
                            onTap: p.shareCode.isEmpty
                                ? null
                                : () async {
                                    await Clipboard.setData(
                                      ClipboardData(text: p.shareCode),
                                    );
                                    if (!context.mounted) return;
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                        content: Text(l10n.inviteCodeCopied),
                                      ),
                                    );
                                  },
                          ),
                          Divider(
                            height: 1,
                            indent: 56,
                            color: cs.onSurface.withValues(alpha: 0.06),
                          ),
                          _InfoTile(
                            icon: Icons.videocam_outlined,
                            title: l10n.anchorRoleTitle,
                            value: p.isAnchor
                                ? l10n.anchorRoleEnabled
                                : l10n.anchorRoleDisabled,
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 20, 16, 8),
                    child: Text(
                      l10n.settingsSection,
                      style: theme.textTheme.titleSmall?.copyWith(
                        color: cs.onSurface.withValues(alpha: 0.55),
                        fontWeight: FontWeight.w700,
                        letterSpacing: 0.4,
                      ),
                    ),
                  ),
                ),
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: Card(
                      child: Column(
                        children: [
                          ListTile(
                            leading: Icon(
                              Icons.phonelink_setup_rounded,
                              color: cs.primary,
                            ),
                            title: Text(l10n.deviceCodeTitle),
                            subtitle: Text(_deviceIdSubtitle(l10n)),
                            trailing:
                                const Icon(Icons.chevron_right_rounded),
                            onTap: _loading ? null : _editDeviceCode,
                          ),
                          Divider(
                            height: 1,
                            indent: 56,
                            color: cs.onSurface.withValues(alpha: 0.06),
                          ),
                          ListTile(
                            leading: Icon(
                              Icons.dns_outlined,
                              color: cs.primary,
                            ),
                            title: Text(l10n.apiBaseTitle),
                            subtitle: Text(
                              _apiBaseSubtitle(l10n),
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                            ),
                            trailing:
                                const Icon(Icons.chevron_right_rounded),
                            onTap: _loading ? null : _editApiBase,
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                ListenableBuilder(
                  listenable: widget.localeController,
                  builder: (context, _) {
                    return SliverToBoxAdapter(
                      child: Padding(
                        padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
                        child: Card(
                          child: ListTile(
                            leading: Icon(
                              Icons.language_rounded,
                              color: cs.primary,
                            ),
                            title: Text(l10n.languageTitle),
                            subtitle: Text(_localeSubtitle(l10n)),
                            trailing:
                                const Icon(Icons.chevron_right_rounded),
                            onTap: () => showLocalePicker(
                              context,
                              controller: widget.localeController,
                            ),
                          ),
                        ),
                      ),
                    );
                  },
                ),
                ListenableBuilder(
                  listenable: widget.themeController,
                  builder: (context, _) {
                    final style = widget.themeController.style;
                    return SliverToBoxAdapter(
                      child: Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        child: Card(
                          margin: const EdgeInsets.only(top: 12),
                          child: ListTile(
                            leading: Icon(
                              Icons.palette_outlined,
                              color: cs.primary,
                            ),
                            title: Text(l10n.uiStyleTitle),
                            subtitle:
                                Text('${style.label} · ${style.subtitle}'),
                            trailing: const Icon(Icons.chevron_right_rounded),
                            onTap: () => showUiStylePicker(
                              context,
                              controller: widget.themeController,
                            ),
                          ),
                        ),
                      ),
                    );
                  },
                ),
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 28, 16, 40),
                    child: OutlinedButton(
                      onPressed: _loading ? null : _vm.logout,
                      style: OutlinedButton.styleFrom(
                        foregroundColor: cs.error,
                        side: BorderSide(color: cs.error.withValues(alpha: 0.45)),
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                      ),
                      child: Text(l10n.logout),
                    ),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

class _HeroHeader extends StatelessWidget {
  const _HeroHeader({
    required this.l10n,
    required this.profile,
    required this.busy,
    required this.onEditAvatar,
    required this.onEditNickname,
  });

  final AppLocalizations l10n;
  final UserProfile profile;
  final bool busy;
  final VoidCallback onEditAvatar;
  final VoidCallback onEditNickname;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final uid = int.tryParse(profile.userId) ?? 0;
    final name = profile.nickname.isEmpty
        ? l10n.userDisplayName(uid)
        : profile.nickname;

    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(20),
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              cs.primary.withValues(alpha: 0.22),
              cs.surface,
              cs.surface,
            ],
          ),
          border: Border.all(color: cs.onSurface.withValues(alpha: 0.06)),
        ),
        padding: const EdgeInsets.all(18),
        child: Row(
          children: [
            GestureDetector(
              onTap: busy ? null : onEditAvatar,
              child: Stack(
                clipBehavior: Clip.none,
                children: [
                  Container(
                    padding: const EdgeInsets.all(3),
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      border: Border.all(
                        color: cs.primary.withValues(alpha: 0.55),
                        width: 2,
                      ),
                    ),
                    child: CachedAvatar(
                      url: profile.avatar,
                      radius: 36,
                      fallbackText: name.substring(0, 1),
                    ),
                  ),
                  Positioned(
                    right: 0,
                    bottom: 0,
                    child: Container(
                      padding: const EdgeInsets.all(5),
                      decoration: BoxDecoration(
                        color: cs.primary,
                        shape: BoxShape.circle,
                        border: Border.all(color: cs.surface, width: 2),
                      ),
                      child: Icon(
                        Icons.camera_alt_rounded,
                        size: 12,
                        color: cs.onPrimary,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  InkWell(
                    onTap: busy ? null : onEditNickname,
                    borderRadius: BorderRadius.circular(8),
                    child: Row(
                      children: [
                        Flexible(
                          child: Text(
                            name,
                            style: theme.textTheme.headlineSmall?.copyWith(
                              fontWeight: FontWeight.w800,
                              height: 1.1,
                            ),
                          ),
                        ),
                        const SizedBox(width: 6),
                        Icon(
                          Icons.edit_rounded,
                          size: 16,
                          color: cs.onSurface.withValues(alpha: 0.45),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 8),
                  Builder(
                    builder: (context) {
                      final pretty = profile.prettyId.isNotEmpty &&
                              profile.prettyId != '0'
                          ? profile.prettyId
                          : '${profile.userId}';
                      return InkWell(
                        onTap: () async {
                          await Clipboard.setData(ClipboardData(text: pretty));
                          if (!context.mounted) return;
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(content: Text(l10n.commonCopied)),
                          );
                        },
                        borderRadius: BorderRadius.circular(999),
                        child: _Chip(
                          text: pretty,
                          emphasize: true,
                          trailing: Icons.copy_rounded,
                        ),
                      );
                    },
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

class _Chip extends StatelessWidget {
  const _Chip({
    required this.text,
    this.emphasize = false,
    this.trailing,
  });

  final String text;
  final bool emphasize;
  final IconData? trailing;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final color =
        emphasize ? cs.primary : cs.onSurface.withValues(alpha: 0.75);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: emphasize
            ? cs.primary.withValues(alpha: 0.16)
            : cs.onSurface.withValues(alpha: 0.06),
        borderRadius: BorderRadius.circular(999),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            text,
            style: Theme.of(context).textTheme.labelMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: color,
                ),
          ),
          if (trailing != null) ...[
            const SizedBox(width: 6),
            Icon(trailing, size: 14, color: color.withValues(alpha: 0.85)),
          ],
        ],
      ),
    );
  }
}

class _WalletStrip extends StatelessWidget {
  const _WalletStrip({
    required this.l10n,
    required this.gold,
    required this.diamond,
    required this.follow,
    required this.fans,
    this.hideGoldDiamond = false,
  });

  final AppLocalizations l10n;
  final String gold;
  final String diamond;
  final String follow;
  final String fans;
  final bool hideGoldDiamond;

  @override
  Widget build(BuildContext context) {
    final metrics = <Widget>[
      if (!hideGoldDiamond) ...[
        _Metric(
          label: l10n.walletGold,
          value: gold,
          icon: Icons.monetization_on_outlined,
        ),
        _Metric(
          label: l10n.walletDiamond,
          value: diamond,
          icon: Icons.diamond_outlined,
        ),
      ],
      _Metric(
        label: l10n.walletFollow,
        value: follow,
        icon: Icons.person_add_alt_1_outlined,
      ),
      _Metric(
        label: l10n.walletFans,
        value: fans,
        icon: Icons.groups_outlined,
      ),
    ];
    return Card(
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 8),
        child: Row(children: metrics),
      ),
    );
  }
}

class _Metric extends StatelessWidget {
  const _Metric({
    required this.label,
    required this.value,
    required this.icon,
  });

  final String label;
  final String value;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    return Expanded(
      child: Column(
        children: [
          Icon(icon, size: 18, color: cs.primary.withValues(alpha: 0.9)),
          const SizedBox(height: 6),
          Text(
            value,
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.w800,
            ),
          ),
          const SizedBox(height: 2),
          Text(
            label,
            style: theme.textTheme.labelSmall?.copyWith(
              color: cs.onSurface.withValues(alpha: 0.55),
            ),
          ),
        ],
      ),
    );
  }
}

class _InfoTile extends StatelessWidget {
  const _InfoTile({
    required this.icon,
    required this.title,
    required this.value,
    this.subtitle,
    this.trailingIcon,
    this.onTap,
  });

  final IconData icon;
  final String title;
  final String value;
  final String? subtitle;
  final IconData? trailingIcon;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final showSubtitle =
        subtitle != null && subtitle!.isNotEmpty && subtitle != value;
    return ListTile(
      leading: Icon(icon, color: cs.primary),
      title: Text(title),
      subtitle: showSubtitle
          ? Text(
              subtitle!,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: cs.onSurface.withValues(alpha: 0.55),
                  ),
            )
          : null,
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 140),
            child: Text(
              value,
              overflow: TextOverflow.ellipsis,
              textAlign: TextAlign.right,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: cs.onSurface.withValues(alpha: 0.72),
                    fontWeight: FontWeight.w600,
                  ),
            ),
          ),
          if (trailingIcon != null) ...[
            const SizedBox(width: 4),
            Icon(
              trailingIcon,
              size: 18,
              color: cs.onSurface.withValues(alpha: 0.35),
            ),
          ],
        ],
      ),
      onTap: onTap,
    );
  }
}

class _ApiBaseEditDialog extends StatefulWidget {
  const _ApiBaseEditDialog({required this.initial});

  final String initial;

  @override
  State<_ApiBaseEditDialog> createState() => _ApiBaseEditDialogState();
}

class _ApiBaseEditDialogState extends State<_ApiBaseEditDialog> {
  late final TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: widget.initial);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _save() {
    FocusManager.instance.primaryFocus?.unfocus();
    Navigator.pop(context, (text: _controller.text.trim(), reset: false));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final muted = Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.65);
    return AlertDialog(
      scrollable: true,
      title: Text(l10n.apiBaseEditTitle),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            l10n.apiBaseWarning,
            style: Theme.of(context).textTheme.bodySmall?.copyWith(color: muted),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _controller,
            autofocus: true,
            maxLength: 256,
            keyboardType: TextInputType.url,
            decoration: InputDecoration(
              hintText: l10n.apiBaseHint,
              counterText: '',
            ),
            textInputAction: TextInputAction.done,
            onSubmitted: (_) => _save(),
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () =>
              Navigator.pop(context, (text: null, reset: true)),
          child: Text(l10n.apiBaseResetDefault),
        ),
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: Text(l10n.commonCancel),
        ),
        FilledButton(
          onPressed: _save,
          child: Text(l10n.commonSave),
        ),
      ],
    );
  }
}

class _NicknameEditDialog extends StatefulWidget {
  const _NicknameEditDialog({required this.initial});

  final String initial;

  @override
  State<_NicknameEditDialog> createState() => _NicknameEditDialogState();
}

class _NicknameEditDialogState extends State<_NicknameEditDialog> {
  late final TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: widget.initial);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _save() {
    FocusManager.instance.primaryFocus?.unfocus();
    Navigator.pop(context, _controller.text.trim());
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    return AlertDialog(
      scrollable: true,
      title: Text(l10n.editNicknameTitle),
      content: TextField(
        controller: _controller,
        autofocus: true,
        maxLength: 20,
        decoration: InputDecoration(
          hintText: l10n.editNicknameHint,
          counterText: '',
        ),
        textInputAction: TextInputAction.done,
        onSubmitted: (_) => _save(),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: Text(l10n.commonCancel),
        ),
        FilledButton(
          onPressed: _save,
          child: Text(l10n.commonSave),
        ),
      ],
    );
  }
}
