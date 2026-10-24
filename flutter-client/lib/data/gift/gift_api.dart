import '../../core/http/api_client.dart';
import 'gift_models.dart';

class GiftApi {
  GiftApi(this._client);

  final ApiClient _client;

  /// `POST /gift/appGiftList`
  Future<List<AppGiftItem>> fetchAppGiftList() async {
    final res = await _client.post<List<AppGiftItem>>(
      '/gift/appGiftList',
      body: const {},
      parseData: (raw) {
        if (raw is! Map) return null;
        final list = raw['list'];
        if (list is! List) return <AppGiftItem>[];
        return [
          for (final e in list)
            if (e is Map)
              AppGiftItem.fromJson(Map<String, dynamic>.from(e)),
        ];
      },
    );
    return res.data ?? const [];
  }
}
