import '../../core/http/api_client.dart';
import 'vip_models.dart';

class VipApi {
  VipApi(this._client);

  final ApiClient _client;

  /// `POST /vip/getVipDetail`
  Future<VipDetail> getVipDetail() async {
    final res = await _client.post<VipDetail>(
      '/vip/getVipDetail',
      body: const {},
      parseData: (raw) {
        if (raw is! Map) return null;
        return VipDetail.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('vip/getVipDetail empty data');
    }
    return data;
  }

  /// `POST /vipCfg/vipCfgListForApp`
  Future<List<VipCfgItem>> vipCfgListForApp() async {
    final res = await _client.post<List<VipCfgItem>>(
      '/vipCfg/vipCfgListForApp',
      body: const {},
      parseData: (raw) {
        if (raw is! Map) return null;
        final list = raw['list'];
        if (list is! List) return <VipCfgItem>[];
        return list
            .whereType<Map>()
            .map((e) => VipCfgItem.fromJson(Map<String, dynamic>.from(e)))
            .toList();
      },
    );
    return res.data ?? const [];
  }
}
