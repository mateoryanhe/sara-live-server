import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/data/live_room/live_room_models.dart';
import 'package:sara_live_app/ui/contracts/anchor_one_to_one_setup_ui.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/viewmodel/anchor_one_to_one_setup_view_model.dart';

/// 主播 Tab：开通 / 更新 1v1 房间（标题 + 服务器通话单价档位）。
class AnchorOneToOneSetupPage extends StatefulWidget {
  const AnchorOneToOneSetupPage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<AnchorOneToOneSetupPage> createState() =>
      _AnchorOneToOneSetupPageState();
}

class _AnchorOneToOneSetupPageState extends State<AnchorOneToOneSetupPage>
    with AutomaticKeepAliveClientMixin
    implements AnchorOneToOneSetupUi {
  late final AnchorOneToOneSetupViewModel _vm;
  final _titleController = TextEditingController();

  bool _pageLoading = true;
  bool _saving = false;
  OneToOneRoomConfig? _config;
  List<BillingTierItem> _billingTiers = const [];
  double? _selectedBilling;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _vm = AnchorOneToOneSetupViewModel(scope: widget.scope, ui: this);
    WidgetsBinding.instance.addPostFrameCallback((_) => _vm.load());
  }

  @override
  void dispose() {
    _titleController.dispose();
    _vm.dispose();
    super.dispose();
  }

  void _syncFieldsFromConfig(OneToOneRoomConfig? cfg) {
    if (cfg == null) return;
    if (_titleController.text.isEmpty && cfg.title.isNotEmpty) {
      _titleController.text = cfg.title;
    }
    if (_selectedBilling == null && cfg.billing >= 0) {
      _selectedBilling = cfg.billing;
    }
  }

  List<BillingTierItem> get _sortedBillingTiers {
    final list = List<BillingTierItem>.from(_billingTiers);
    list.sort((a, b) {
      final bySort = a.sort.compareTo(b.sort);
      if (bySort != 0) return bySort;
      return b.pricePerMinute.compareTo(a.pricePerMinute);
    });
    return list;
  }

  @override
  void showOneToOneSetupLoading(bool loading) {
    if (!mounted) return;
    setState(() => _pageLoading = loading);
  }

  @override
  void renderOneToOneSetup({
    required OneToOneRoomConfig? config,
    required List<BillingTierItem> billingTiers,
    required bool saving,
  }) {
    if (!mounted) return;
    setState(() {
      _config = config;
      _billingTiers = billingTiers;
      _saving = saving;
      _pageLoading = false;
      if (_selectedBilling == null && billingTiers.isNotEmpty) {
        _selectedBilling = billingTiers.first.pricePerMinute;
      }
    });
    _syncFieldsFromConfig(config);
  }

  @override
  void showOneToOneSetupError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  void showOneToOneSetupMessage(String message) {
    showOneToOneSetupError(message);
  }

  String _billingLabel(double v) {
    if (v <= 0) return '免费通话';
    final s = v == v.roundToDouble() ? '${v.toInt()}' : v.toStringAsFixed(1);
    return '$s 钻 / 分钟';
  }

  Widget _sectionCard({
    required String title,
    required IconData icon,
    required Widget child,
  }) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    return Card(
      elevation: 0,
      color: cs.surfaceContainerHighest.withValues(alpha: 0.55),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(16, 16, 16, 18),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                Icon(icon, size: 20, color: cs.primary),
                const SizedBox(width: 8),
                Text(
                  title,
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w700,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 14),
            child,
          ],
        ),
      ),
    );
  }

  Widget _billingTierTile(BillingTierItem tier) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final selected = _selectedBilling == tier.pricePerMinute;
    final price = tier.pricePerMinute;

    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: Material(
        color: selected
            ? cs.primaryContainer.withValues(alpha: 0.65)
            : cs.surface.withValues(alpha: 0.85),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(14),
          side: BorderSide(
            color: selected
                ? cs.primary.withValues(alpha: 0.85)
                : cs.outline.withValues(alpha: 0.25),
            width: selected ? 2 : 1,
          ),
        ),
        clipBehavior: Clip.antiAlias,
        child: InkWell(
          onTap: _saving
              ? null
              : () => setState(() => _selectedBilling = price),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 18),
            child: Row(
              children: [
                Icon(
                  Icons.diamond_outlined,
                  color: selected
                      ? cs.primary
                      : cs.onSurface.withValues(alpha: 0.5),
                ),
                const SizedBox(width: 14),
                Expanded(
                  child: Text(
                    _billingLabel(price),
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w700,
                      color: selected ? cs.onPrimaryContainer : cs.onSurface,
                    ),
                  ),
                ),
                Icon(
                  selected
                      ? Icons.radio_button_checked
                      : Icons.radio_button_off,
                  color: selected
                      ? cs.primary
                      : cs.onSurface.withValues(alpha: 0.35),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final hasRoom = _config != null;

    return Scaffold(
      body: AppPageFrame(
        child: _pageLoading
            ? const Center(child: CircularProgressIndicator())
            : SingleChildScrollView(
                padding: const EdgeInsets.fromLTRB(16, 12, 16, 28),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    const SizedBox(height: 4),
                    Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Container(
                                width: 52,
                                height: 52,
                                decoration: BoxDecoration(
                                  borderRadius: BorderRadius.circular(14),
                                  gradient: LinearGradient(
                                    colors: [
                                      cs.primary,
                                      cs.primary.withValues(alpha: 0.55),
                                    ],
                                  ),
                                ),
                                child: Icon(
                                  Icons.videocam_rounded,
                                  color: cs.onPrimary,
                                  size: 28,
                                ),
                              ),
                              const SizedBox(width: 14),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      '1v1 房间',
                                      style:
                                          theme.textTheme.headlineSmall?.copyWith(
                                        fontWeight: FontWeight.w800,
                                      ),
                                    ),
                                    const SizedBox(height: 6),
                                    Text(
                                      hasRoom
                                          ? '已开通 · ${_config!.isOnShelf ? '已上架' : '未上架'}'
                                          : '填写标题并选择通话单价即可创建',
                                      style: theme.textTheme.bodyMedium?.copyWith(
                                        color: cs.onSurface.withValues(alpha: 0.62),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                    const SizedBox(height: 22),
                    _sectionCard(
                            icon: Icons.edit_outlined,
                            title: '房间标题',
                            child: TextField(
                              controller: _titleController,
                              enabled: !_saving,
                              textInputAction: TextInputAction.done,
                              decoration: InputDecoration(
                                hintText: '给观众的 1v1 房间起个名字',
                                filled: true,
                                fillColor: cs.surface.withValues(alpha: 0.9),
                                border: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                contentPadding: const EdgeInsets.symmetric(
                                  horizontal: 16,
                                  vertical: 16,
                                ),
                              ),
                            ),
                          ),
                    const SizedBox(height: 16),
                    _sectionCard(
                            icon: Icons.payments_outlined,
                            title: '通话单价',
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.stretch,
                              children: [
                                Text(
                                  '以下为服务器已上架档位，单列选择',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: cs.onSurface.withValues(alpha: 0.55),
                                  ),
                                ),
                                const SizedBox(height: 12),
                                if (_billingTiers.isEmpty)
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                      vertical: 24,
                                    ),
                                    child: Column(
                                      children: [
                                        Icon(
                                          Icons.info_outline,
                                          size: 36,
                                          color: cs.error.withValues(alpha: 0.8),
                                        ),
                                        const SizedBox(height: 10),
                                        Text(
                                          '暂无已上架档位',
                                          textAlign: TextAlign.center,
                                          style: theme.textTheme.titleSmall,
                                        ),
                                        const SizedBox(height: 4),
                                        Text(
                                          '请在 CMS 配置并上架 1v1 视频通话计费',
                                          textAlign: TextAlign.center,
                                          style: theme.textTheme.bodySmall?.copyWith(
                                            color: cs.onSurface.withValues(alpha: 0.55),
                                          ),
                                        ),
                                      ],
                                    ),
                                  )
                                else
                                  for (final tier in _sortedBillingTiers)
                                    _billingTierTile(tier),
                              ],
                            ),
                          ),
                    const SizedBox(height: 32),
                    FilledButton(
                            style: FilledButton.styleFrom(
                              minimumSize: const Size.fromHeight(52),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(14),
                              ),
                            ),
                            onPressed: _saving ||
                                    _billingTiers.isEmpty ||
                                    _selectedBilling == null
                                ? null
                                : () {
                                    _vm.save(
                                      title: _titleController.text,
                                      billing: _selectedBilling!,
                                    );
                                  },
                            child: _saving
                                ? const SizedBox(
                                    width: 22,
                                    height: 22,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                    ),
                                  )
                                : Text(
                                    hasRoom ? '保存 1v1 房间设置' : '创建 1v1 房间',
                                    style: const TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.w700,
                                    ),
                                  ),
                    ),
                  ],
                ),
              ),
      ),
    );
  }
}
