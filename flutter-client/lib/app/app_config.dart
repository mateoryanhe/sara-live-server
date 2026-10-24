/// App 运行配置。可用 `--dart-define=API_BASE=https://...` 覆盖内置默认。
class AppConfig {
  const AppConfig._();

  /// 编译期默认 API（设置里未填自定义地址时使用）。
  static const String defaultApiBaseUrl = String.fromEnvironment(
    'API_BASE',
    defaultValue: 'https://www.bigtktool.shop',
  );

  @Deprecated('Use defaultApiBaseUrl or ServerEndpointStore.apiBaseUrl')
  static String get apiBaseUrl => defaultApiBaseUrl;

  static const String appPackage = 'net.saralive.sara_live_app';
  static const String appVersion = '0.1.0';

  /// 由 API 根地址推导 WS：`https`→`wss`，`http`→`ws`。
  static String wsBaseUrlFor(String apiBaseUrl) {
    final u = Uri.parse(apiBaseUrl);
    final scheme = u.scheme == 'https'
        ? 'wss'
        : (u.scheme == 'http' ? 'ws' : u.scheme);
    return u.replace(scheme: scheme).toString().replaceAll(RegExp(r'/$'), '');
  }

  @Deprecated('Use ServerEndpointStore.wsBaseUrl')
  static String get wsBaseUrl => wsBaseUrlFor(defaultApiBaseUrl);
}
