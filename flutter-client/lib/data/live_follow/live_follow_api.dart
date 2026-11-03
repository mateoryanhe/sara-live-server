import '../../core/http/api_client.dart';

class LiveFollowMutationResult {
  const LiveFollowMutationResult({
    required this.following,
    required this.followerCount,
  });

  final bool following;
  final int followerCount;

  factory LiveFollowMutationResult.fromJson(Map<String, dynamic> json) {
    return LiveFollowMutationResult(
      following: json['following'] == true,
      followerCount: _toInt(json['followerCount']),
    );
  }

  static int _toInt(Object? v) {
    if (v is num) return v.toInt();
    return int.tryParse('$v') ?? 0;
  }
}

/// 对齐 `POST /liveFollow/*`。
class LiveFollowApi {
  LiveFollowApi(this._client);

  final ApiClient _client;

  Future<bool> isFollowing({required int anchorId}) async {
    final res = await _client.post<Map<String, dynamic>>(
      '/liveFollow/isFollowing',
      body: {'anchorId': anchorId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return Map<String, dynamic>.from(raw);
      },
    );
    return res.data?['following'] == true;
  }

  Future<LiveFollowMutationResult> follow({required int anchorId}) async {
    final res = await _client.post<LiveFollowMutationResult>(
      '/liveFollow/follow',
      body: {'anchorId': anchorId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return LiveFollowMutationResult.fromJson(
          Map<String, dynamic>.from(raw),
        );
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveFollow/follow empty data');
    }
    return data;
  }

  Future<LiveFollowMutationResult> unfollow({required int anchorId}) async {
    final res = await _client.post<LiveFollowMutationResult>(
      '/liveFollow/unfollow',
      body: {'anchorId': anchorId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return LiveFollowMutationResult.fromJson(
          Map<String, dynamic>.from(raw),
        );
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveFollow/unfollow empty data');
    }
    return data;
  }

  Future<void> blockUser({required int targetId}) async {
    await _client.post<void>(
      '/liveFollow/block',
      body: {'targetId': targetId},
      parseData: (_) => null,
    );
  }

  Future<void> unblockUser({required int targetId}) async {
    await _client.post<void>(
      '/liveFollow/unblock',
      body: {'targetId': targetId},
      parseData: (_) => null,
    );
  }
}
