/// 缓存 key 约定：用户数据必须以 userId 为命名空间。
///
/// 格式：
/// - 用户：`user:{userId}:{scope}`，例如 `user:10001:profile`
/// - 全局：`global:{scope}`，仅设备/包级、与登录无关的数据
class CacheKey {
  CacheKey._();

  static String user(int userId, String scope) {
    assert(userId > 0, 'userId must be > 0');
    assert(scope.isNotEmpty && !scope.contains(':'), 'scope must be non-empty and not contain ":"');
    return 'user:$userId:$scope';
  }

  /// 该用户下所有 key 的前缀，用于登出/踢人整清。
  static String userPrefix(int userId) {
    assert(userId > 0, 'userId must be > 0');
    return 'user:$userId:';
  }

  static String global(String scope) {
    assert(scope.isNotEmpty && !scope.contains(':'), 'scope must be non-empty and not contain ":"');
    return 'global:$scope';
  }

  static bool isUserKey(String key, int userId) =>
      key.startsWith(userPrefix(userId));
}
