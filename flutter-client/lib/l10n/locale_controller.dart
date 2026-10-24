import 'dart:ui';

import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 应用语言；切换后 [notifyListeners]，并持久化。
class LocaleController extends ChangeNotifier {
  LocaleController._(this._prefs, this._locale);

  static const _kLocale = 'app_locale';

  static const supportedLocales = <Locale>[
    Locale('en'),
    Locale('zh'),
  ];

  final SharedPreferences _prefs;
  Locale _locale;

  Locale get locale => _locale;

  static Future<LocaleController> create() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getString(_kLocale);
    if (raw != null) {
      final parsed = _parseStored(raw);
      if (parsed != null) {
        return LocaleController._(prefs, parsed);
      }
    }
    final system = PlatformDispatcher.instance.locale;
    return LocaleController._(prefs, _resolveFromSystem(system));
  }

  static Locale _resolveFromSystem(Locale system) {
    final lang = system.languageCode.toLowerCase();
    if (lang.startsWith('zh')) {
      return const Locale('zh');
    }
    return const Locale('en');
  }

  static Locale? _parseStored(String raw) {
    final parts = raw.split('_');
    if (parts.isEmpty || parts.first.isEmpty) return null;
    final lang = parts.first.toLowerCase();
    for (final loc in supportedLocales) {
      if (loc.languageCode == lang) return loc;
    }
    return null;
  }

  String get _storageKey => _locale.languageCode;

  Future<void> setLocale(Locale locale) async {
    final match = supportedLocales.firstWhere(
      (e) => e.languageCode == locale.languageCode,
      orElse: () => const Locale('en'),
    );
    if (_locale.languageCode == match.languageCode) return;
    _locale = match;
    await _prefs.setString(_kLocale, _storageKey);
    notifyListeners();
  }

  Future<void> setEnglish() => setLocale(const Locale('en'));

  Future<void> setChinese() => setLocale(const Locale('zh'));
}
