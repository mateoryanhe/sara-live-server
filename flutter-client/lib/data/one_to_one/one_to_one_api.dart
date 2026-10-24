import '../../core/http/api_client.dart';
import 'one_to_one_models.dart';

class OneToOneApi {
  OneToOneApi(this._client);

  final ApiClient _client;

  /// `POST /liveRoom/oneToOneRoomList`
  Future<OneToOneRoomListPage> fetchRoomList({
    int page = 1,
    int pageSize = 20,
    int statusFilter = OneToOneStatusFilter.all,
  }) async {
    final res = await _client.post<OneToOneRoomListPage>(
      '/liveRoom/oneToOneRoomList',
      body: {
        'page': page,
        'pageSize': pageSize,
        'statusFilter': statusFilter,
      },
      parseData: (raw) {
        if (raw is! Map) return null;
        return OneToOneRoomListPage.fromJson(
          Map<String, dynamic>.from(raw),
        );
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/oneToOneRoomList empty data');
    }
    return data;
  }
}
