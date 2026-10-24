import '../../core/http/api_client.dart';

/// 对齐 `POST /liveFollow/*`。
class LiveFollowApi {
  LiveFollowApi(this._client);

  final ApiClient _client;

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
