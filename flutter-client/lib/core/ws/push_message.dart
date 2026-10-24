import 'dart:convert';

/// 对齐 Go `httpserver.PushResp`：`{cmd,data}`。
/// 服务端常一次下发 **数组** `[{cmd,data}, ...]`（Binary/Text JSON）。
class PushMessage {
  const PushMessage({required this.cmd, this.data});

  final int cmd;
  final Object? data;

  factory PushMessage.fromJson(Map<String, dynamic> json) {
    return PushMessage(
      cmd: (json['cmd'] as num?)?.toInt() ?? 0,
      data: json['data'],
    );
  }

  static List<PushMessage> parseFrame(Object? decoded) {
    if (decoded is List) {
      return decoded
          .whereType<Map>()
          .map((e) => PushMessage.fromJson(Map<String, dynamic>.from(e)))
          .toList();
    }
    if (decoded is Map) {
      return [PushMessage.fromJson(Map<String, dynamic>.from(decoded))];
    }
    return const [];
  }

  static List<PushMessage> parseBytes(List<int> bytes) {
    final decoded = jsonDecode(utf8.decode(bytes));
    return parseFrame(decoded);
  }

  static List<PushMessage> parseText(String text) {
    final decoded = jsonDecode(text);
    return parseFrame(decoded);
  }
}

/// 推送命令回调；业务层在此 invalidate / patch 缓存。
typedef PushHandler = void Function(PushMessage message);
