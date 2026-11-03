import 'dart:convert';

import '../http/http_log_config.dart';
import '../http/http_logger.dart';
import '../http/http_log_sink_io.dart'
    if (dart.library.html) '../http/http_log_sink_web.dart';

/// WebSocket push / connection debug logs (same switch & sink as [HttpLogger]).
class WsLogger {
  WsLogger._();

  static Future<void> logConnection({
    required String state,
    Uri? uri,
    int reconnectAttempt = 0,
    Duration? reconnectIn,
    String? note,
  }) async {
    if (!HttpLogSettings.enabled) return;
    await HttpLogger.ensureInitialized();

    final buf = StringBuffer()
      ..writeln('======== WS ${HttpLogConfig.platformTag} ========')
      ..writeln('time    : ${DateTime.now().toIso8601String()}')
      ..writeln('state   : $state');
    if (reconnectAttempt > 0) {
      buf.writeln('reconnectAttempt : $reconnectAttempt');
    }
    if (reconnectIn != null) {
      buf.writeln('reconnectIn      : ${reconnectIn.inMilliseconds} ms');
    }
    if (uri != null) {
      buf.writeln('url     : ${_redactUri(uri)}');
    }
    if (note != null && note.isNotEmpty) {
      buf.writeln('note    : $note');
    }
    buf.writeln('========================================');

    await writeHttpLogSink(buf.toString());
  }

  static Future<void> logPush({
    required int cmd,
    Object? data,
  }) async {
    if (!HttpLogSettings.enabled) return;
    await HttpLogger.ensureInitialized();

    final buf = StringBuffer()
      ..writeln('======== WS ${HttpLogConfig.platformTag} ========')
      ..writeln('time    : ${DateTime.now().toIso8601String()}')
      ..writeln('direction : recv')
      ..writeln('cmd     : $cmd')
      ..writeln('--- data ---')
      ..writeln(_formatData(data))
      ..writeln('========================================');

    await writeHttpLogSink(buf.toString());
  }

  static Future<void> logError({
    required String phase,
    required Object error,
    StackTrace? stackTrace,
  }) async {
    if (!HttpLogSettings.enabled) return;
    await HttpLogger.ensureInitialized();

    final buf = StringBuffer()
      ..writeln('======== WS ${HttpLogConfig.platformTag} ========')
      ..writeln('time    : ${DateTime.now().toIso8601String()}')
      ..writeln('phase   : $phase')
      ..writeln('--- error ---')
      ..writeln('$error');
    if (stackTrace != null) {
      buf.writeln('--- stack ---');
      buf.writeln('$stackTrace');
    }
    buf.writeln('========================================');

    await writeHttpLogSink(buf.toString());
  }

  static String _redactUri(Uri uri) {
    final auth = uri.queryParameters['Authorization'];
    if (auth == null || auth.isEmpty) {
      return uri.replace(queryParameters: {}).toString();
    }
    final redacted = auth.length <= 12 ? '***' : '${auth.substring(0, 8)}...';
    return uri
        .replace(
          queryParameters: {'Authorization': redacted},
        )
        .toString();
  }

  static String _formatData(Object? data) {
    if (data == null) return '(null)';
    try {
      final pretty = const JsonEncoder.withIndent('  ').convert(data);
      return _clip(pretty);
    } catch (_) {
      return _clip('$data');
    }
  }

  static String _clip(String s) {
    final max = HttpLogConfig.maxBodyChars;
    if (s.length <= max) return s;
    return '${s.substring(0, max)}\n... [truncated ${s.length - max} chars]';
  }
}
