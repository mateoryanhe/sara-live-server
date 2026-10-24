/// 对齐 Go App 登录返回的 token：`{userId}.{tokenSecret}`（不要加 Bearer）。
class AuthToken {
  const AuthToken({required this.userId, required this.secret});

  final String userId;
  final String secret;

  /// 原始 Authorization 头值。
  String get headerValue => '$userId.$secret';

  /// WebSocket：`/ws?Authorization=userId.secret`
  String get queryValue => headerValue;

  static AuthToken? tryParse(String? raw) {
    if (raw == null || raw.isEmpty) {
      return null;
    }
    final parts = raw.split('.');
    if (parts.length != 2 || parts[0].isEmpty || parts[1].isEmpty) {
      return null;
    }
    return AuthToken(userId: parts[0], secret: parts[1]);
  }

  @override
  String toString() => headerValue;
}

/// 由上层（安全存储）提供当前 token；未登录返回 null。
typedef AuthTokenProvider = AuthToken? Function();
