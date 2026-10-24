import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// 可切换的专业 UI 皮肤（灵感：直播 App 暗色舞台 / Stream Cart 字体体系）。
enum UiStyle {
  /// 暗色舞台：炭黑 + 琥珀点缀（偏直播间）
  stageDark,

  /// 夜场霓虹：深蓝黑 + 青绿点缀
  neonNight,

  /// 日间工作室：浅灰底 + 深蓝强调（偏运营后台质感）
  dayStudio,

  /// 墨纸极简：近黑白、细分割线、大字号
  inkPaper,
}

extension UiStyleX on UiStyle {
  String get label => switch (this) {
        UiStyle.stageDark => '暗色舞台',
        UiStyle.neonNight => '夜场霓虹',
        UiStyle.dayStudio => '日间工作室',
        UiStyle.inkPaper => '墨纸极简',
      };

  String get subtitle => switch (this) {
        UiStyle.stageDark => '直播间黑场 · 琥珀强调',
        UiStyle.neonNight => '深蓝夜色 · 青绿霓虹',
        UiStyle.dayStudio => '浅色专业 · 清晰层级',
        UiStyle.inkPaper => '黑白编辑风 · 大字号',
      };

  bool get isDark => this != UiStyle.dayStudio;
}

class AppThemes {
  AppThemes._();

  static ThemeData of(UiStyle style) {
    switch (style) {
      case UiStyle.stageDark:
        return _build(
          brightness: Brightness.dark,
          seed: const Color(0xFFE8A838),
          surface: const Color(0xFF121212),
          scaffold: const Color(0xFF0A0A0A),
          card: const Color(0xFF1C1C1E),
          font: GoogleFonts.plusJakartaSansTextTheme,
        );
      case UiStyle.neonNight:
        return _build(
          brightness: Brightness.dark,
          seed: const Color(0xFF2DD4BF),
          surface: const Color(0xFF0F172A),
          scaffold: const Color(0xFF020617),
          card: const Color(0xFF1E293B),
          font: GoogleFonts.manropeTextTheme,
        );
      case UiStyle.dayStudio:
        return _build(
          brightness: Brightness.light,
          seed: const Color(0xFF0F4C81),
          surface: const Color(0xFFF4F6F8),
          scaffold: const Color(0xFFF8FAFC),
          card: Colors.white,
          font: GoogleFonts.beVietnamProTextTheme,
        );
      case UiStyle.inkPaper:
        return _build(
          brightness: Brightness.dark,
          seed: const Color(0xFFF5F5F4),
          surface: const Color(0xFF111111),
          scaffold: const Color(0xFF000000),
          card: const Color(0xFF161616),
          font: GoogleFonts.dmSansTextTheme,
          sharp: true,
        );
    }
  }

  static ThemeData _build({
    required Brightness brightness,
    required Color seed,
    required Color surface,
    required Color scaffold,
    required Color card,
    required TextTheme Function(TextTheme?) font,
    bool sharp = false,
  }) {
    final baseScheme = ColorScheme.fromSeed(
      seedColor: seed,
      brightness: brightness,
    );
    final scheme = baseScheme.copyWith(
      surface: surface,
      primary: seed,
    );

    final radius = sharp ? 8.0 : 16.0;
    final textTheme = font(ThemeData(brightness: brightness).textTheme);

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      colorScheme: scheme,
      scaffoldBackgroundColor: scaffold,
      textTheme: textTheme,
      appBarTheme: AppBarTheme(
        backgroundColor: scaffold,
        foregroundColor: scheme.onSurface,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        titleTextStyle: textTheme.titleLarge?.copyWith(
          fontWeight: FontWeight.w700,
          color: scheme.onSurface,
        ),
      ),
      cardTheme: CardThemeData(
        color: card,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radius),
          side: BorderSide(
            color: scheme.onSurface.withValues(alpha: 0.06),
          ),
        ),
        margin: EdgeInsets.zero,
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: card,
        indicatorColor: seed.withValues(alpha: 0.22),
        labelTextStyle: WidgetStateProperty.resolveWith((states) {
          final selected = states.contains(WidgetState.selected);
          return textTheme.labelMedium?.copyWith(
            fontWeight: selected ? FontWeight.w700 : FontWeight.w500,
          );
        }),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(radius),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
          textStyle: textTheme.labelLarge?.copyWith(fontWeight: FontWeight.w700),
        ),
      ),
      listTileTheme: ListTileThemeData(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radius),
        ),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 2),
      ),
      dividerTheme: DividerThemeData(
        color: scheme.onSurface.withValues(alpha: 0.08),
        space: 1,
      ),
    );
  }
}
