import '../../core/http/api_client.dart';
import 'call_models.dart';

class CallApi {
  CallApi(this._client);

  final ApiClient _client;

  /// `POST /call/oneToOneRoomCall`
  Future<CallRtcCredentials> oneToOneRoomCall({required int targetId}) async {
    final res = await _client.post<CallRtcCredentials>(
      '/call/oneToOneRoomCall',
      body: {'targetId': targetId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return CallRtcCredentials.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('call/oneToOneRoomCall empty data');
    }
    return data;
  }

  /// `POST /call/acceptCall`
  Future<CallRtcCredentials> acceptCall({required String orderId}) async {
    final res = await _client.post<CallRtcCredentials>(
      '/call/acceptCall',
      body: {'orderId': orderId},
      parseData: (raw) {
        if (raw is! Map) return null;
        final map = Map<String, dynamic>.from(raw);
        if (map['success'] == false) return null;
        return CallRtcCredentials.fromJson(map);
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('call/acceptCall empty data');
    }
    return data;
  }

  /// `POST /call/anchorRejectCall`
  Future<void> rejectCall({required String orderId}) async {
    await _client.post<void>(
      '/call/anchorRejectCall',
      body: {'orderId': orderId},
      parseData: (_) => null,
    );
  }

  /// `POST /call/confirmCall`
  Future<void> confirmCall({required String orderId}) async {
    await _client.post<void>(
      '/call/confirmCall',
      body: {'orderId': orderId},
      parseData: (_) => null,
    );
  }

  /// `POST /call/getCallConfirmStatus`
  Future<CallConfirmStatus> getCallConfirmStatus({required String orderId}) async {
    final res = await _client.post<CallConfirmStatus>(
      '/call/getCallConfirmStatus',
      body: {'orderId': orderId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return CallConfirmStatus.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('call/getCallConfirmStatus empty data');
    }
    return data;
  }

  /// `POST /call/callHeart`（约每 10 秒）
  Future<void> callHeart({required String orderId}) async {
    await _client.post<void>(
      '/call/callHeart',
      body: {'orderId': orderId},
      parseData: (_) => null,
    );
  }

  /// `POST /call/callTimeout`
  Future<void> callTimeout({required String orderId}) async {
    await _client.post<void>(
      '/call/callTimeout',
      body: {'orderId': orderId},
      parseData: (_) => null,
    );
  }

  /// `POST /call/endCall`
  Future<CallEndSummary> endCall({required String orderId}) async {
    final res = await _client.post<CallEndSummary>(
      '/call/endCall',
      body: {'orderId': orderId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return CallEndSummary.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('call/endCall empty data');
    }
    return data;
  }
}
