import '../../core/http/api_client.dart';
import 'user_profile.dart';

class UserApi {
  UserApi(this._client);

  final ApiClient _client;

  /// `POST /userInfo/get`；`userId` 省略/0 为当前登录用户。
  Future<UserProfile> getUserInfo({String? userId}) async {
    final body = <String, dynamic>{};
    if (userId != null && userId.isNotEmpty && userId != '0') {
      body['userId'] = int.tryParse(userId) ?? userId;
    }
    final res = await _client.post<UserProfile>(
      '/userInfo/get',
      body: body,
      parseData: (raw) {
        if (raw is! Map) return null;
        return UserProfile.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('userInfo/get empty data');
    }
    return data;
  }

  /// `POST /userInfo/updateNickname`
  Future<String> updateNickname(String nickname) async {
    final res = await _client.post<String>(
      '/userInfo/updateNickname',
      body: {'nickname': nickname},
      parseData: (raw) {
        if (raw is Map && raw['nickname'] != null) {
          return '${raw['nickname']}';
        }
        return null;
      },
    );
    final data = res.data;
    if (data == null || data.isEmpty) {
      throw StateError('updateNickname empty data');
    }
    return data;
  }

  /// `POST /userInfo/uploadAvatar` multipart field `file`
  Future<String> uploadAvatar({
    required List<int> bytes,
    required String filename,
  }) async {
    final res = await _client.postMultipart<String>(
      '/userInfo/uploadAvatar',
      fileField: 'file',
      fileBytes: bytes,
      filename: filename,
      parseData: (raw) {
        if (raw is Map && raw['avatar'] != null) {
          return '${raw['avatar']}';
        }
        return null;
      },
    );
    final data = res.data;
    if (data == null || data.isEmpty) {
      throw StateError('uploadAvatar empty data');
    }
    return data;
  }
}
