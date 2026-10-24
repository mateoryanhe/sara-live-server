import 'package:flutter/material.dart';

/// 主播直播底栏：单个「工具」按钮，点击展开（参考常见直播 Kit：主工具栏 + 美颜子页）。
class AnchorLiveToolbar extends StatelessWidget {
  const AnchorLiveToolbar({
    super.key,
    required this.micMuted,
    required this.cameraEnabled,
    required this.beautyEnabled,
    required this.beautyLevel,
    required this.onFlipCamera,
    required this.onToggleCamera,
    required this.onToggleMic,
    required this.onBeautyEnabledChanged,
    required this.onBeautyLevelChanged,
  });

  final bool micMuted;
  final bool cameraEnabled;
  final bool beautyEnabled;
  final double beautyLevel;
  final VoidCallback onFlipCamera;
  final VoidCallback onToggleCamera;
  final VoidCallback onToggleMic;
  final ValueChanged<bool> onBeautyEnabledChanged;
  final ValueChanged<double> onBeautyLevelChanged;

  bool get _hasActiveTweak =>
      micMuted || !cameraEnabled || !beautyEnabled;

  static ButtonStyle _fabStyle({required bool highlight}) {
    return IconButton.styleFrom(
      backgroundColor: highlight
          ? const Color(0xFF6366F1).withValues(alpha: 0.92)
          : Colors.white.withValues(alpha: 0.16),
      foregroundColor: Colors.white,
      minimumSize: const Size(40, 40),
      padding: const EdgeInsets.all(8),
    );
  }

  void _openToolsSheet(BuildContext context) {
    showModalBottomSheet<void>(
      context: context,
      backgroundColor: const Color(0xFF161616),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      showDragHandle: true,
      isScrollControlled: true,
      builder: (ctx) => _LiveToolsSheet(
        micMuted: micMuted,
        cameraEnabled: cameraEnabled,
        beautyEnabled: beautyEnabled,
        beautyLevel: beautyLevel,
        onFlipCamera: onFlipCamera,
        onToggleCamera: onToggleCamera,
        onToggleMic: onToggleMic,
        onBeautyEnabledChanged: onBeautyEnabledChanged,
        onBeautyLevelChanged: onBeautyLevelChanged,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return IconButton(
      style: _fabStyle(highlight: _hasActiveTweak),
      tooltip: '直播工具',
      onPressed: () => _openToolsSheet(context),
      icon: const Icon(Icons.grid_view_rounded),
    );
  }
}

enum _LiveToolsPage { main, beauty }

class _LiveToolsSheet extends StatefulWidget {
  const _LiveToolsSheet({
    required this.micMuted,
    required this.cameraEnabled,
    required this.beautyEnabled,
    required this.beautyLevel,
    required this.onFlipCamera,
    required this.onToggleCamera,
    required this.onToggleMic,
    required this.onBeautyEnabledChanged,
    required this.onBeautyLevelChanged,
  });

  final bool micMuted;
  final bool cameraEnabled;
  final bool beautyEnabled;
  final double beautyLevel;
  final VoidCallback onFlipCamera;
  final VoidCallback onToggleCamera;
  final VoidCallback onToggleMic;
  final ValueChanged<bool> onBeautyEnabledChanged;
  final ValueChanged<double> onBeautyLevelChanged;

  @override
  State<_LiveToolsSheet> createState() => _LiveToolsSheetState();
}

class _LiveToolsSheetState extends State<_LiveToolsSheet> {
  _LiveToolsPage _page = _LiveToolsPage.main;
  late bool _micMuted;
  late bool _cameraEnabled;
  late bool _beautyEnabled;
  late double _beautyLevel;

  @override
  void initState() {
    super.initState();
    _syncFromWidget();
  }

  void _syncFromWidget() {
    _micMuted = widget.micMuted;
    _cameraEnabled = widget.cameraEnabled;
    _beautyEnabled = widget.beautyEnabled;
    _beautyLevel = widget.beautyLevel;
  }

  void _applyBeautyPreset({required bool enabled, required double level}) {
    widget.onBeautyEnabledChanged(enabled);
    if (enabled) {
      widget.onBeautyLevelChanged(level);
    }
    setState(() {
      _beautyEnabled = enabled;
      if (enabled) {
        _beautyLevel = level;
      }
    });
  }

  Widget _buildMainPage() {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Text(
          '直播工具',
          style: TextStyle(
            color: Colors.white,
            fontSize: 16,
            fontWeight: FontWeight.w700,
          ),
        ),
        const SizedBox(height: 6),
        Text(
          '设置保存在本机，下次开播自动生效',
          style: TextStyle(
            color: Colors.white.withValues(alpha: 0.45),
            fontSize: 12,
          ),
        ),
        const SizedBox(height: 22),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _RoundLiveTool(
              icon: _cameraEnabled
                  ? Icons.videocam_rounded
                  : Icons.videocam_off_rounded,
              label: _cameraEnabled ? '摄像头' : '已关闭',
              selected: !_cameraEnabled,
              onTap: () {
                widget.onToggleCamera();
                setState(() => _cameraEnabled = !_cameraEnabled);
              },
            ),
            _RoundLiveTool(
              icon: Icons.cameraswitch_rounded,
              label: '翻转',
              enabled: _cameraEnabled,
              onTap: widget.onFlipCamera,
            ),
            _RoundLiveTool(
              icon: _micMuted ? Icons.mic_off_rounded : Icons.mic_rounded,
              label: _micMuted ? '已静音' : '麦克风',
              selected: _micMuted,
              onTap: () {
                widget.onToggleMic();
                setState(() => _micMuted = !_micMuted);
              },
            ),
            _RoundLiveTool(
              icon: Icons.auto_fix_high_rounded,
              label: '美颜',
              selected: _beautyEnabled,
              onTap: () => setState(() => _page = _LiveToolsPage.beauty),
            ),
          ],
        ),
        const SizedBox(height: 8),
      ],
    );
  }

  Widget _buildBeautyPage() {
    const presets = <({String name, bool on, double level})>[
      (name: '原图', on: false, level: 0),
      (name: '自然', on: true, level: 0.4),
      (name: '标准', on: true, level: 0.65),
      (name: '清晰', on: true, level: 0.85),
    ];

    int? selectedPresetIndex() {
      if (!_beautyEnabled) {
        return 0;
      }
      for (var i = 1; i < presets.length; i++) {
        final p = presets[i];
        if ((p.level - _beautyLevel).abs() < 0.06) {
          return i;
        }
      }
      return null;
    }

    final presetIndex = selectedPresetIndex();

    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Row(
          children: [
            IconButton(
              onPressed: () => setState(() => _page = _LiveToolsPage.main),
              icon: const Icon(Icons.arrow_back_ios_new_rounded,
                  color: Colors.white, size: 20),
              tooltip: '返回',
            ),
            const Expanded(
              child: Text(
                '美颜',
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ),
            const SizedBox(width: 48),
          ],
        ),
        const SizedBox(height: 4),
        SizedBox(
          height: 36,
          child: ListView.separated(
            scrollDirection: Axis.horizontal,
            itemCount: presets.length,
            separatorBuilder: (_, __) => const SizedBox(width: 8),
            itemBuilder: (context, index) {
              final p = presets[index];
              final selected = presetIndex == index;
              return ChoiceChip(
                label: Text(p.name),
                selected: selected,
                onSelected: (_) =>
                    _applyBeautyPreset(enabled: p.on, level: p.level),
                labelStyle: TextStyle(
                  color: selected ? Colors.white : Colors.white70,
                  fontWeight: selected ? FontWeight.w700 : FontWeight.w500,
                  fontSize: 13,
                ),
                selectedColor: const Color(0xFF6366F1),
                backgroundColor: Colors.white.withValues(alpha: 0.08),
                side: BorderSide(
                  color: selected
                      ? Colors.transparent
                      : Colors.white.withValues(alpha: 0.12),
                ),
                padding: const EdgeInsets.symmetric(horizontal: 4),
              );
            },
          ),
        ),
        const SizedBox(height: 20),
        Text(
          '精细调节',
          style: TextStyle(
            color: Colors.white.withValues(alpha: 0.55),
            fontSize: 12,
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 12),
        _BeautySliderRow(
          icon: Icons.tune_rounded,
          label: '强度',
          value: _beautyLevel,
          enabled: _beautyEnabled,
          onChanged: (v) {
            if (!_beautyEnabled) {
              _applyBeautyPreset(enabled: true, level: v);
              return;
            }
            widget.onBeautyLevelChanged(v);
            setState(() => _beautyLevel = v);
          },
        ),
        if (!_beautyEnabled)
          Padding(
            padding: const EdgeInsets.only(top: 12),
            child: Text(
              '选择上方预设或拖动滑条即可开启美颜',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: Colors.white.withValues(alpha: 0.4),
                fontSize: 12,
              ),
            ),
          ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(16, 0, 16, 20),
        child: AnimatedSwitcher(
          duration: const Duration(milliseconds: 220),
          switchInCurve: Curves.easeOut,
          switchOutCurve: Curves.easeIn,
          child: _page == _LiveToolsPage.main
              ? KeyedSubtree(
                  key: const ValueKey('main'),
                  child: _buildMainPage(),
                )
              : KeyedSubtree(
                  key: const ValueKey('beauty'),
                  child: _buildBeautyPage(),
                ),
        ),
      ),
    );
  }
}

class _RoundLiveTool extends StatelessWidget {
  const _RoundLiveTool({
    required this.icon,
    required this.label,
    required this.onTap,
    this.selected = false,
    this.enabled = true,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool selected;
  final bool enabled;

  @override
  Widget build(BuildContext context) {
    final circleBg = selected
        ? const Color(0xFF6366F1).withValues(alpha: 0.35)
        : Colors.white.withValues(alpha: 0.1);
    final borderSide = selected
        ? const BorderSide(color: Color(0xFF818CF8), width: 1.5)
        : BorderSide(color: Colors.white.withValues(alpha: 0.08));

    return Opacity(
      opacity: enabled ? 1 : 0.35,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Material(
            color: circleBg,
            shape: CircleBorder(side: borderSide),
            clipBehavior: Clip.antiAlias,
            child: InkWell(
              onTap: enabled ? onTap : null,
              customBorder: const CircleBorder(),
              child: SizedBox(
                width: 56,
                height: 56,
                child: Icon(icon, color: Colors.white, size: 26),
              ),
            ),
          ),
          const SizedBox(height: 8),
          SizedBox(
            width: 64,
            child: Text(
              label,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              textAlign: TextAlign.center,
              style: TextStyle(
                color: Colors.white.withValues(alpha: 0.85),
                fontSize: 11,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _BeautySliderRow extends StatelessWidget {
  const _BeautySliderRow({
    required this.icon,
    required this.label,
    required this.value,
    required this.enabled,
    required this.onChanged,
  });

  final IconData icon;
  final String label;
  final double value;
  final bool enabled;
  final ValueChanged<double> onChanged;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, size: 20, color: Colors.white.withValues(alpha: 0.7)),
        const SizedBox(width: 8),
        SizedBox(
          width: 36,
          child: Text(
            label,
            style: TextStyle(
              color: Colors.white.withValues(alpha: 0.75),
              fontSize: 13,
            ),
          ),
        ),
        Expanded(
          child: SliderTheme(
            data: SliderTheme.of(context).copyWith(
              trackHeight: 3,
              thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 7),
              overlayShape: SliderComponentShape.noOverlay,
              activeTrackColor: const Color(0xFF6366F1),
              inactiveTrackColor: Colors.white.withValues(alpha: 0.12),
              thumbColor: Colors.white,
            ),
            child: Slider(
              value: enabled ? value : 0,
              min: 0,
              max: 1,
              onChanged: enabled ? onChanged : (_) => onChanged(0.45),
            ),
          ),
        ),
        SizedBox(
          width: 36,
          child: Text(
            '${(value * 100).round()}',
            textAlign: TextAlign.right,
            style: TextStyle(
              color: Colors.white.withValues(alpha: 0.5),
              fontSize: 12,
              fontFeatures: const [FontFeature.tabularFigures()],
            ),
          ),
        ),
      ],
    );
  }
}
