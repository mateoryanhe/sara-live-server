import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'app_themes.dart';

/// 全局 UI 皮肤；切换后 [notifyListeners]，并持久化。
class ThemeController extends ChangeNotifier {
  ThemeController._(this._prefs, this._style);

  static const _kStyle = 'ui_style';

  final SharedPreferences _prefs;
  UiStyle _style;

  UiStyle get style => _style;

  static Future<ThemeController> create() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getString(_kStyle);
    final style = UiStyle.values.firstWhere(
      (e) => e.name == raw,
      orElse: () => UiStyle.stageDark,
    );
    return ThemeController._(prefs, style);
  }

  Future<void> setStyle(UiStyle style) async {
    if (_style == style) return;
    _style = style;
    await _prefs.setString(_kStyle, style.name);
    notifyListeners();
  }
}
