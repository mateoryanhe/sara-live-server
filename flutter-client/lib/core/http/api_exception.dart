/// HTTP / 业务错误。
class ApiException implements Exception {
  ApiException({
    required this.code,
    required this.message,
    this.data,
    this.statusCode,
  });

  /// 业务码（Go errercode）；网络层失败可用 -1。
  final int code;
  final String message;
  final Object? data;
  final int? statusCode;

  /// 与 Go `errercode.Token` 对齐。
  bool get isUnauthorized => code == 3;

  @override
  String toString() => 'ApiException(code=$code, message=$message)';
}
