/// WebSocket 连接配置，对齐 Go：`GET /ws?Authorization={userId}.{token}`。
class WsConfig {
  WsConfig({
    required String Function() wsBaseUrlProvider,
    this.reconnectDelay = const Duration(seconds: 2),
    this.maxReconnectDelay = const Duration(seconds: 30),
  }) : _wsBaseUrlProvider = wsBaseUrlProvider;

  final String Function() _wsBaseUrlProvider;

  /// 如 `wss://api.example.com`（不要带 path；客户端拼 `/ws`）。
  String get wsBaseUrl => _wsBaseUrlProvider();
  final Duration reconnectDelay;
  final Duration maxReconnectDelay;

  Uri buildUri(String authorizationQuery) {
    final base = Uri.parse(wsBaseUrl);
    return base.replace(
      path: '/ws',
      queryParameters: {'Authorization': authorizationQuery},
    );
  }
}
