import 'package:flutter/foundation.dart';

import '../platform/platform_check.dart';

/// HTTP debug logging (`--dart-define`).
abstract final class HttpLogConfig {
  static const bool defineEnabled = bool.fromEnvironment(
    'HTTP_LOG_ENABLED',
    defaultValue: false,
  );

  static const String defineLogDir = String.fromEnvironment(
    'HTTP_LOG_DIR',
    defaultValue: '',
  );

  static const int maxBodyChars = int.fromEnvironment(
    'HTTP_LOG_MAX_BODY',
    defaultValue: 8192,
  );

  static String get platformTag {
    if (kIsWeb) return 'web';
    if (isAndroidOs) return 'android';
    if (isIOSOs) return 'ios';
    return 'io';
  }
}

/// Runtime toggle / directory (defaults from [HttpLogConfig]).
class HttpLogSettings {
  HttpLogSettings._();

  static bool enabled = HttpLogConfig.defineEnabled;
  static String logDir = HttpLogConfig.defineLogDir;

  static bool get writesToFile =>
      !kIsWeb && logDir.trim().isNotEmpty;
}
