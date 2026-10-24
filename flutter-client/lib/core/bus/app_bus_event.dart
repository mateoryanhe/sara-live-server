import '../ws/push_message.dart';

/// 统一事件总线载荷：生命周期 / WS 用类型区分。
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

@Deprecated('Use AppLifecycleBusEvent')
typedef StartupBusEvent = AppLifecycleBusEvent;

@Deprecated('Use AppLifecycleKind')
typedef StartupKind = AppLifecycleKind;

/// WS 下行或本机 WS 相关事件（含 [PushCmd.wsConnected]）。
final class WsBusEvent extends AppBusEvent {
  const WsBusEvent({required this.cmd, this.data});

  factory WsBusEvent.fromPush(PushMessage message) =>
      WsBusEvent(cmd: message.cmd, data: message.data);

  final int cmd;
  final Object? data;

  PushMessage toPushMessage() => PushMessage(cmd: cmd, data: data);
}

typedef AppBusHandler = void Function(AppBusEvent event);
