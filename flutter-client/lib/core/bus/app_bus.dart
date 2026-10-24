import 'app_bus_event.dart';

/// 进程内统一事件总线：生命周期 + WS 共用 `add` / `remove` / `publish`。
///
/// ```dart
/// bus.add((e) {
///   switch (e) {
///     case AppLifecycleBusEvent(:final kind):
///       switch (kind) {
///         case AppLifecycleKind.started:
///         case AppLifecycleKind.paused:
///         case AppLifecycleKind.resumed:
///       }
///     case WsBusEvent(:final cmd, :final data):
///       break;
///     default:
///       break;
///   }
/// });
/// ```
///
/// 仅 [AppLifecycleKind.started] 为 sticky：晚订阅会补发一次。
class AppBus {
  final List<AppBusHandler> _handlers = <AppBusHandler>[];
  AppLifecycleBusEvent? _stickyStarted;

  void add(AppBusHandler handler) {
    if (!_handlers.contains(handler)) {
      _handlers.add(handler);
    }
    final sticky = _stickyStarted;
    if (sticky != null) {
      handler(sticky);
    }
  }

  void remove(AppBusHandler handler) {
    _handlers.remove(handler);
  }

  void publish(AppBusEvent event) {
    if (event is AppLifecycleBusEvent &&
        event.kind == AppLifecycleKind.started) {
      _stickyStarted = event;
    }
    for (final h in List<AppBusHandler>.of(_handlers)) {
      h(event);
    }
  }

  bool get hasStarted => _stickyStarted != null;
}
