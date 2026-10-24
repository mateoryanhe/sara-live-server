import '../../core/http/api_client.dart';
import '../../core/http/auth_token.dart';

class DeviceLoginResult {
  DeviceLoginResult({required this.token, required this.isNewUser});

  final AuthToken token;
  final bool isNewUser;
}

class AuthApi {
  AuthApi(this._client);

  final ApiClient _client;

  /// `POST /auth/deviceLogin`（免鉴权）
  Future<DeviceLoginResult> deviceLogin(Map<String, dynamic> deviceInfo) async {
    final res = await _client.post<DeviceLoginResult>(
      '/auth/deviceLogin',
      requireAuth: false,
      body: {'deviceInfo': deviceInfo},
      parseData: (raw) {
        if (raw is! Map) {
          return null;
        }
        final map = Map<String, dynamic>.from(raw);
        final token = AuthToken.tryParse(map['token']?.toString());
        if (token == null) {
          return null;
        }
        return DeviceLoginResult(
          token: token,
          isNewUser: map['isNewUser'] == true,
        );
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('deviceLogin empty data');
    }
    return data;
  }
}
