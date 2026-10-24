import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'app_config.dart';

/// 本地覆盖 API 根地址（内网调试）；未设置时用 [AppConfig.defaultApiBaseUrl]。
class ServerEndpointStore extends ChangeNotifier {
  ServerEndpointStore(this._prefs);

  static const _kOverride = 'api_base_override';

  final SharedPreferences _prefs;
  String? _override;

  String? get storedOverride => _override;

  bool get hasOverride => _override != null && _override!.isNotEmpty;

  String get apiBaseUrl =>
      hasOverride ? _override! : AppConfig.defaultApiBaseUrl;

  String get wsBaseUrl => AppConfig.wsBaseUrlFor(apiBaseUrl);

  Future<void> load() async {
    final raw = _prefs.getString(_kOverride)?.trim();
    _override = (raw == null || raw.isEmpty) ? null : raw;
  }

  Future<void> setOverride(String? normalized) async {
    if (normalized == null || normalized.isEmpty) {
      await _prefs.remove(_kOverride);
      _override = null;
    } else {
      await _prefs.setString(_kOverride, normalized);
      _override = normalized;
    }
    notifyListeners();
  }

  /// 规范化用户输入（见 [ApiBaseValidation]）。
  static ApiBaseValidation validateInput(String raw) {
    var s = raw.trim();
    s = s.replaceAll(RegExp(r',+$'), '').trim();
    if (s.isEmpty) {
      return const ApiBaseValidation.empty();
    }
    if (!s.contains('://')) {
      s = 'http://$s';
    }
    final uri = Uri.tryParse(s);
    if (uri == null || uri.host.isEmpty) {
      return const ApiBaseValidation.invalid();
    }
    if (uri.scheme != 'http' && uri.scheme != 'https') {
      return const ApiBaseValidation.invalid();
    }
    final normalized = uri
        .replace(
          path: '',
          query: '',
          fragment: '',
        )
        .toString()
        .replaceAll(RegExp(r'/+$'), '');
    return ApiBaseValidation.ok(normalized);
  }
}

sealed class ApiBaseValidation {
  const ApiBaseValidation();

  const factory ApiBaseValidation.ok(String url) = ApiBaseValidationOk;
  const factory ApiBaseValidation.empty() = ApiBaseValidationEmpty;
  const factory ApiBaseValidation.invalid() = ApiBaseValidationInvalid;
}

final class ApiBaseValidationOk extends ApiBaseValidation {
  const ApiBaseValidationOk(this.url);
  final String url;
}

final class ApiBaseValidationEmpty extends ApiBaseValidation {
  const ApiBaseValidationEmpty();
}

final class ApiBaseValidationInvalid extends ApiBaseValidation {
  const ApiBaseValidationInvalid();
}
