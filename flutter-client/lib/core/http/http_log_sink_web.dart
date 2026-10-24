import 'package:flutter/foundation.dart';

Future<void> initHttpLogSink({
  required String platform,
  required String logDir,
}) async {
  if (logDir.trim().isNotEmpty) {
    debugPrint(
      '[HttpLog] web: logs go to DevTools Console + flutter run terminal '
      '(HTTP_LOG_DIR=$logDir is label only; use Android build for file sink)',
    );
  } else {
    debugPrint('[HttpLog] web: logs -> DevTools Console + flutter run terminal');
  }
}

Future<void> writeHttpLogSink(String block) async {
  debugPrint(block);
}
