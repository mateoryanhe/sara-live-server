import 'dart:async';
import 'dart:convert';

import 'package:web_socket_channel/status.dart' as ws_status;
import 'package:web_socket_channel/web_socket_channel.dart';

import '../http/auth_token.dart';
import 'push_message.dart';
import 'ws_config.dart';

enum WsConnectionState { disconnected, connecting, connected }

/// 对齐 Go `/ws`：query 鉴权；下行多为 JSON 数组的 Binary/Text 帧。
class WsClient {
  WsClient({
    required this.config,
    required AuthTokenProvider tokenProvider,
    this.onMessages,
    this.onState,
    this.onError,
  }) : _tokenProvider = tokenProvider;

  final WsConfig config;
  final AuthTokenProvider _tokenProvider;
  final PushHandler? onMessages;
  final void Function(WsConnectionState state)? onState;
  final void Function(Object error, StackTrace st)? onError;

  WebSocketChannel? _channel;
  StreamSubscription<dynamic>? _sub;
  Timer? _reconnectTimer;
  bool _manualClose = false;
  int _reconnectAttempt = 0;
  WsConnectionState _state = WsConnectionState.disconnected;

  WsConnectionState get state => _state;

  Future<void> connect() async {
    _manualClose = false;
    await _open();
  }

  Future<void> disconnect() async {
    _manualClose = true;
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
    await _tearDown();
    _setState(WsConnectionState.disconnected);
  }

  Future<void> _open() async {
    final token = _tokenProvider();
    if (token == null) {
      _setState(WsConnectionState.disconnected);
      return;
    }

    await _tearDown();
    _setState(WsConnectionState.connecting);

    try {
      final uri = config.buildUri(token.queryValue);
      final channel = WebSocketChannel.connect(uri);
      _channel = channel;
      await channel.ready;
      _reconnectAttempt = 0;
      _setState(WsConnectionState.connected);
      _sub = channel.stream.listen(
        _onData,
        onError: (Object e, StackTrace st) {
          onError?.call(e, st);
          _scheduleReconnect();
        },
        onDone: () {
          if (!_manualClose) {
            _scheduleReconnect();
          } else {
            _setState(WsConnectionState.disconnected);
          }
        },
        cancelOnError: true,
      );
    } catch (e, st) {
      onError?.call(e, st);
      _scheduleReconnect();
    }
  }

  void _onData(dynamic raw) {
    try {
      late List<PushMessage> list;
      if (raw is List<int>) {
        list = PushMessage.parseBytes(raw);
      } else if (raw is String) {
        list = PushMessage.parseText(raw);
      } else {
        list = PushMessage.parseFrame(raw);
      }
      for (final msg in list) {
        onMessages?.call(msg);
      }
    } catch (e, st) {
      onError?.call(e, st);
    }
  }

  void _scheduleReconnect() {
    _setState(WsConnectionState.disconnected);
    if (_manualClose) {
      return;
    }
    _reconnectTimer?.cancel();
    final delay = _backoff();
    _reconnectAttempt++;
    _reconnectTimer = Timer(delay, () {
      if (!_manualClose) {
        unawaited(_open());
      }
    });
  }

  Duration _backoff() {
    final ms = config.reconnectDelay.inMilliseconds *
        (1 << _reconnectAttempt.clamp(0, 4));
    final capped = ms > config.maxReconnectDelay.inMilliseconds
        ? config.maxReconnectDelay.inMilliseconds
        : ms;
    return Duration(milliseconds: capped);
  }

  Future<void> _tearDown() async {
    await _sub?.cancel();
    _sub = null;
    try {
      await _channel?.sink.close(ws_status.normalClosure);
    } catch (_) {}
    _channel = null;
  }

  void _setState(WsConnectionState next) {
    if (_state == next) {
      return;
    }
    _state = next;
    onState?.call(next);
  }

  /// 预留：若服务端要求客户端上行心跳，可在此发送。
  /// 当前 Go 侧由服务端空闲推 heart，客户端一般只需收。
  void sendJson(Object payload) {
    final ch = _channel;
    if (ch == null || _state != WsConnectionState.connected) {
      return;
    }
    ch.sink.add(jsonEncode(payload));
  }
}
