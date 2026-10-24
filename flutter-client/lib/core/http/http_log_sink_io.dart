import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';

String _platform = 'io';
Directory? _dir;

Future<void> initHttpLogSink({
  required String platform,
  required String logDir,
}) async {
  _platform = platform;
  if (logDir.trim().isNotEmpty) {
    _dir = Directory(logDir.trim());
  } else {
    final root = await getApplicationSupportDirectory();
    _dir = Directory(p.join(root.path, 'http_logs', platform));
  }
  await _dir!.create(recursive: true);
}

Future<void> writeHttpLogSink(String block) async {
  debugPrint(block);
  final dir = _dir;
  if (dir == null) return;
  final day = DateTime.now().toIso8601String().substring(0, 10);
  final file = File(p.join(dir.path, 'http_${_platform}_$day.log'));
  await file.writeAsString('$block\n', mode: FileMode.append, flush: true);
}
