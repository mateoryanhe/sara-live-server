import '../../core/http/api_client.dart';
import 'live_room_models.dart';

class LiveRoomApi {
  LiveRoomApi(this._client);

  final ApiClient _client;

  Future<OnlineUserListPage> fetchOnlineList({
    required int roomId,
    int page = 1,
    int pageSize = 20,
  }) async {
    final res = await _client.post<OnlineUserListPage>(
      '/liveRoom/onlineList',
      body: {
        'roomId': roomId,
        'page': page,
        'pageSize': pageSize,
      },
      parseData: (raw) {
        if (raw is! Map) return null;
        return OnlineUserListPage.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/onlineList empty data');
    }
    return data;
  }

  Future<OneToOneRoomConfig> fetchOneToOneRoom({required int anchorId}) async {
    final res = await _client.post<OneToOneRoomConfig>(
      '/liveRoom/getOneToOne',
      body: {'anchorId': anchorId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return OneToOneRoomConfig.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/getOneToOne empty data');
    }
    return data;
  }

  Future<List<BillingTierItem>> fetchBillingTierList() async {
    final res = await _client.post<List<BillingTierItem>>(
      '/privateRoomBilling/appBillingList',
      body: const {},
      parseData: (raw) {
        if (raw is! Map) return null;
        final list = raw['list'];
        if (list is! List) return <BillingTierItem>[];
        return [
          for (final e in list)
            if (e is Map)
              BillingTierItem.fromJson(Map<String, dynamic>.from(e)),
        ];
      },
    );
    return res.data ?? const [];
  }

  Future<String> createOneToOneRoom({
    required double billing,
    required String title,
    String cover = '',
  }) async {
    final body = <String, dynamic>{
      'billing': billing,
      'title': title,
    };
    if (cover.isNotEmpty) {
      body['cover'] = cover;
    }
    final res = await _client.post<Map<String, dynamic>>(
      '/liveRoom/createOneToOne',
      body: body,
      parseData: (raw) {
        if (raw is! Map) return null;
        return Map<String, dynamic>.from(raw);
      },
    );
    final data = res.data;
    return '${data?['roomId'] ?? ''}';
  }

  Future<LiveRoomDetail> fetchRoom({required int roomId}) async {
    final res = await _client.post<LiveRoomDetail>(
      '/liveRoom/get',
      body: {'roomId': roomId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return LiveRoomDetail.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/get empty data');
    }
    return data;
  }

  Future<String> createShowLiveRoom({
    required String title,
    required String notice,
    List<int>? coverBytes,
    String? coverFilename,
  }) async {
    final fields = {
      'title': title,
      'notice': notice,
      'category': '$liveRoomCategoryHot',
    };
    Map<String, dynamic>? data;
    if (coverBytes != null &&
        coverBytes.isNotEmpty &&
        coverFilename != null &&
        coverFilename.isNotEmpty) {
      final res = await _client.postMultipart<Map<String, dynamic>>(
        '/liveRoom/create',
        fileField: 'cover',
        fileBytes: coverBytes,
        filename: coverFilename,
        fields: fields,
        parseData: (raw) {
          if (raw is! Map) return null;
          return Map<String, dynamic>.from(raw);
        },
      );
      data = res.data;
    } else {
      final res = await _client.postMultipartFields<Map<String, dynamic>>(
        '/liveRoom/create',
        fields: fields,
        parseData: (raw) {
          if (raw is! Map) return null;
          return Map<String, dynamic>.from(raw);
        },
      );
      data = res.data;
    }
    return '${data?['roomId'] ?? ''}';
  }

  /// 修改公告（可传空字符串清空），对齐 `POST /liveRoom/updateNotice`。
  Future<void> updateNotice({required String notice}) async {
    await _client.post<void>(
      '/liveRoom/updateNotice',
      body: {'notice': notice},
      parseData: (_) => null,
    );
  }

  Future<void> startLive() async {
    await _client.post<void>(
      '/liveRoom/startLive',
      body: const {},
      parseData: (_) => null,
    );
  }

  /// 开播/在房心跳（对齐 Go [CommonPeriod=30s]；主播刷新 HeartTime，观众刷新在线心跳）。
  Future<void> reportLiveStartStatus({required int roomId}) async {
    await _client.post<void>(
      '/liveRoom/reportLiveStartStatus',
      body: {'roomId': roomId},
      parseData: (_) => null,
    );
  }

  Future<void> stopLive() async {
    await _client.post<void>(
      '/liveRoom/stopLive',
      body: const {},
      parseData: (_) => null,
    );
  }

  /// 获取全部直播间房间列表。
  ///
  /// `statusFilter=0` 表示不按开播状态过滤，适合观众端直播页展示全部主播。
  Future<LiveRoomListPage> fetchLiveRoomList({
    int page = 1,
    int pageSize = 20,
    int statusFilter = 0,
  }) async {
    final res = await _client.post<LiveRoomListPage>(
      '/liveRoom/roomList',
      body: {
        'page': page,
        'pageSize': pageSize,
        'statusFilter': statusFilter,
      },
      parseData: (raw) {
        if (raw is! Map) return null;
        return LiveRoomListPage.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/roomList empty data');
    }
    return data;
  }

  Future<HotLiveRoomListPage> fetchHotLiveRoomList({
    int pageIndex = 1,
    int pageSize = 20,
  }) async {
    final res = await _client.post<HotLiveRoomListPage>(
      '/liveRoom/hotRoomList',
      body: {
        'pageIndex': pageIndex,
        'pageSize': pageSize,
      },
      parseData: (raw) {
        if (raw is! Map) return null;
        return HotLiveRoomListPage.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/hotRoomList empty data');
    }
    return data;
  }

  Future<void> joinRoom({required int roomId}) async {
    await _client.post<void>(
      '/liveRoom/join',
      body: {'roomId': roomId},
      parseData: (_) => null,
    );
  }

  Future<void> leaveRoom({required int roomId}) async {
    await _client.post<void>(
      '/liveRoom/leave',
      body: {'roomId': roomId},
      parseData: (_) => null,
    );
  }

  Future<void> kickAudience({required int userId}) async {
    await _client.post<void>(
      '/liveRoom/kickAudience',
      body: {'userId': userId},
      parseData: (_) => null,
    );
  }

  Future<void> setAudienceMute({
    required int userId,
    required bool muted,
  }) async {
    await _client.post<void>(
      '/liveRoom/setAudienceMute',
      body: {'userId': userId, 'muted': muted},
      parseData: (_) => null,
    );
  }

  Future<void> cancelAudienceMute({required int userId}) async {
    await _client.post<void>(
      '/liveRoom/cancelAudienceMute',
      body: {'userId': userId},
      parseData: (_) => null,
    );
  }

  Future<void> cancelKickBan({required int userId}) async {
    await _client.post<void>(
      '/liveRoom/cancelKickBan',
      body: {'userId': userId},
      parseData: (_) => null,
    );
  }

  Future<AudienceRestrictStatus> fetchAudienceRestrictStatus({
    required int roomId,
    required int userId,
  }) async {
    final res = await _client.post<AudienceRestrictStatus>(
      '/liveRoom/getAudienceRestrictStatus',
      body: {'roomId': roomId, 'userId': userId},
      parseData: (raw) {
        if (raw is! Map) return null;
        return AudienceRestrictStatus.fromJson(
          Map<String, dynamic>.from(raw),
        );
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/getAudienceRestrictStatus empty data');
    }
    return data;
  }

  Future<void> sendChat({
    required int roomId,
    required String content,
  }) async {
    await _client.post<void>(
      '/liveRoom/sendChat',
      body: {
        'roomId': roomId,
        'content': content,
      },
      parseData: (_) => null,
    );
  }

  Future<SendGiftResult> sendGift({
    required int roomId,
    required int giftId,
    required int count,
  }) async {
    final res = await _client.post<SendGiftResult>(
      '/liveRoom/sendGift',
      body: {
        'roomId': roomId,
        'giftId': giftId,
        'count': count,
      },
      parseData: (raw) {
        if (raw is! Map) return null;
        return SendGiftResult.fromJson(Map<String, dynamic>.from(raw));
      },
    );
    final data = res.data;
    if (data == null) {
      throw StateError('liveRoom/sendGift empty data');
    }
    return data;
  }
}

class SendGiftResult {
  SendGiftResult({required this.cost, required this.diamond});

  final double cost;
  final double diamond;

  factory SendGiftResult.fromJson(Map<String, dynamic> json) {
    return SendGiftResult(
      cost: (json['cost'] as num?)?.toDouble() ?? 0,
      diamond: (json['diamond'] as num?)?.toDouble() ?? 0,
    );
  }
}
