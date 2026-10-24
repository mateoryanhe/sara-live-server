import '../../core/http/api_client.dart';

/// 对齐 Go `agoradto.RTCRole*`。
abstract final class AgoraRtcRole {
  static const publisher = 1;
  static const subscriber = 2;
}

class LiveRoomRtcCredentials {
  const LiveRoomRtcCredentials({
    required this.token,
    required this.appId,
    required this.channelName,
    required this.userAccount,
    required this.expireAt,
  });

  final String token;
  final String appId;
  final String channelName;
  final String userAccount;
  final int expireAt;

  int get rtcUid => int.tryParse(userAccount) ?? 0;

  factory LiveRoomRtcCredentials.fromJson(Map<String, dynamic> json) {
    return LiveRoomRtcCredentials(
      token: '${json['token'] ?? ''}',
      appId: '${json['appId'] ?? ''}',
      channelName: '${json['channelName'] ?? ''}',
      userAccount: '${json['userAccount'] ?? ''}',
      expireAt: _jsonInt(json['expireAt']),
    );
  }
}

int _jsonInt(Object? v) {
  if (v is int) return v;
  if (v is num) return v.toInt();
  return int.tryParse('$v') ?? 0;
}

class AgoraApi {
  AgoraApi(this._client);

  final ApiClient _client;

  /// `POST /agora/liveRoomToken`
  Future<LiveRoomRtcCredentials> liveRoomToken({
    required String channelName,
    required int role,
  }) async {
    final res = await _client.post<LiveRoomRtcCredentials>(
      '/agora/liveRoomToken',
      body: {
        'channelName': channelName,
        'role': role,
      },
      parseData: (raw) {
        if (raw is! Map) return null;
        return LiveRoomRtcCredentials.fromJson(
          Map<String, dynamic>.from(raw),
        );
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('agora/liveRoomToken empty data');
    }
    return data;
  }
}
