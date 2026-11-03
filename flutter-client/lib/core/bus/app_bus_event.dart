import '../ws/push_message.dart';
import '../ws/ws_client.dart';

/// 统一事件总线载荷：生命周期 / WS 连接 / WS 下行 用类型区分。
sealed class AppBusEvent {
  const AppBusEvent();
}

/// App 冷启动与前后台（本机 [AppBus.publish]）。
enum AppLifecycleKind {
  /// [AppScope.create] 依赖就绪（仅一次，sticky 晚订阅补发）。
  started,

  /// 进入后台（[AppLifecycleState.paused]）。
  paused,

  /// 回到前台（[AppLifecycleState.resumed]）。
  resumed,
}

final class AppLifecycleBusEvent extends AppBusEvent {
  AppLifecycleBusEvent(this.kind, {DateTime? at}) : at = at ?? DateTime.now();

  factory AppLifecycleBusEvent.started({DateTime? at}) =>
      AppLifecycleBusEvent(AppLifecycleKind.started, at: at);

  factory AppLifecycleBusEvent.paused({DateTime? at}) =>
      AppLifecycleBusEvent(AppLifecycleKind.paused, at: at);

  factory AppLifecycleBusEvent.resumed({DateTime? at}) =>
      AppLifecycleBusEvent(AppLifecycleKind.resumed, at: at);

  final AppLifecycleKind kind;
  final DateTime at;
}

/// 主壳底部 Tab 被选中（本机 [AppBus.publish]）。
final class ShellTabFocusBusEvent extends AppBusEvent {
  const ShellTabFocusBusEvent(this.tabId);

  final String tabId;
}

@Deprecated('Use AppLifecycleBusEvent')
typedef StartupBusEvent = AppLifecycleBusEvent;

@Deprecated('Use AppLifecycleKind')
typedef StartupKind = AppLifecycleKind;

/// WebSocket 连接状态（本机 [AppScope] / [WsClient.onState] → [AppBus.publish]）。
///
/// 与 Go 侧「上线了 / 离线了」对应：**connected** = 可收 push；**disconnected** = 断线/未连接。
/// 各 Repository / ViewModel **自行订阅**，勿在 AppScope 里集中写业务。
final class WsConnectionBusEvent extends AppBusEvent {
  const WsConnectionBusEvent(this.state, {this.manualDisconnect = false});

  final WsConnectionState state;

  /// [WsClient.disconnect] 主动断开时为 true；被动断线为 false。
  final bool manualDisconnect;
}

/// WS 下行 push（服务端 cmd + data）。
final class WsBusEvent extends AppBusEvent {
  const WsBusEvent({required this.cmd, this.data});

  factory WsBusEvent.fromPush(PushMessage message) =>
      WsBusEvent(cmd: message.cmd, data: message.data);

  final int cmd;
  final Object? data;

  PushMessage toPushMessage() => PushMessage(cmd: cmd, data: data);
}

typedef AppBusHandler = void Function(AppBusEvent event);
