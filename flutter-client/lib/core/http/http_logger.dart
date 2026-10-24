import 'dart:convert';

import 'http_log_config.dart';
import 'http_log_sink_io.dart'
    if (dart.library.html) 'http_log_sink_web.dart';

/// Records HTTP requests/responses for debugging.
class HttpLogger {
  HttpLogger._();

  static bool _initialized = false;

  static Future<void> init() async {
    if (!HttpLogSettings.enabled || _initialized) return;
    _initialized = true;
    await initHttpLogSink(
      platform: HttpLogConfig.platformTag,
      logDir: HttpLogSettings.logDir,
    );
  }

  static Future<void> log({
    required String method,
    required Uri uri,
    required Map<String, String> headers,
    required String? requestBody,
    required int? statusCode,
    required String? responseBody,
    required Duration elapsed,
    Object? error,
    String? note,
  }) async {
    if (!HttpLogSettings.enabled) return;
    if (!_initialized) {
      await init();
    }

    final buf = StringBuffer()
      ..writeln('======== HTTP ${HttpLogConfig.platformTag} ========')
      ..writeln('time    : ${DateTime.now().toIso8601String()}')
      ..writeln('elapsed : ${elapsed.inMilliseconds} ms')
      ..writeln('method  : $method')
      ..writeln('url     : $uri');
    if (note != null && note.isNotEmpty) {
      buf.writeln('note    : $note');
    }
    buf
      ..writeln('--- request headers ---')
      ..writeln(_formatHeaders(headers))
      ..writeln('--- request body ---')
      ..writeln(_truncate(requestBody ?? ''))
      ..writeln('--- response ---')
      ..writeln('status  : ${statusCode ?? '-'}')
      ..writeln('--- response body ---')
      ..writeln(_truncate(responseBody ?? ''));
    if (error != null) {
      buf
        ..writeln('--- error ---')
        ..writeln('$error');
    }
    buf.writeln('========================================');

    await writeHttpLogSink(buf.toString());
  }

  static String _formatHeaders(Map<String, String> headers) {
    final lines = <String>[];
    for (final e in headers.entries) {
      lines.add('${e.key}: ${_redactHeader(e.key, e.value)}');
    }
    return lines.join('\n');
  }

  static String _redactHeader(String name, String value) {
    final lower = name.toLowerCase();
    if (lower == 'authorization') {
      if (value.length <= 12) return '***';
      return '${value.substring(0, 8)}...(${value.length} chars)';
    }
    return value;
  }

  static String _truncate(String raw) {
    if (raw.isEmpty) return '(empty)';
    try {
      final decoded = jsonDecode(raw);
      final pretty = const JsonEncoder.withIndent('  ').convert(decoded);
      return _clip(pretty);
    } catch (_) {
      return _clip(raw);
    }
  }

  static String _clip(String s) {
    final max = HttpLogConfig.maxBodyChars;
    if (s.length <= max) return s;
    return '${s.substring(0, max)}\n... [truncated ${s.length - max} chars]';
  }
}
