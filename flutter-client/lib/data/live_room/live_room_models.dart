class OnlineUserListPage {
  OnlineUserListPage({
    required this.total,
    required this.page,
    required this.pageSize,
    required this.list,
  });

  final int total;
  final int page;
  final int pageSize;
  final List<OnlineUserItem> list;

  bool get hasMore => page * pageSize < total;

  factory OnlineUserListPage.fromJson(Map<String, dynamic> json) {
    final rawList = json['list'];
    final items = <OnlineUserItem>[];
    if (rawList is List) {
      for (final e in rawList) {
        if (e is Map) {
          items.add(OnlineUserItem.fromJson(Map<String, dynamic>.from(e)));
        }
      }
    }
    return OnlineUserListPage(
      total: (json['total'] as num?)?.toInt() ?? items.length,
      page: (json['page'] as num?)?.toInt() ?? 1,
      pageSize: (json['pageSize'] as num?)?.toInt() ?? 20,
      list: items,
    );
  }

  Map<String, dynamic> toJson() => {
        'total': total,
        'page': page,
        'pageSize': pageSize,
        'list': [for (final i in list) i.toJson()],
      };
}

class OnlineUserItem {
  OnlineUserItem({
    required this.userId,
    required this.nickname,
    required this.avatar,
    required this.vipLevel,
    required this.totalReward,
    required this.joinedAt,
    this.muted = false,
    this.gender = 0,
    this.age = 0,
  });

  final String userId;
  final String nickname;
  final String avatar;
  final int vipLevel;
  final double totalReward;
  final String joinedAt;
  final bool muted;
  final int gender;
  final int age;

  factory OnlineUserItem.fromJson(Map<String, dynamic> json) {
    return OnlineUserItem(
      userId: '${json['userId'] ?? ''}',
      nickname: '${json['nickname'] ?? ''}',
      avatar: '${json['avatar'] ?? ''}',
      vipLevel: (json['vipLevel'] as num?)?.toInt() ?? 0,
      totalReward: (json['totalReward'] as num?)?.toDouble() ?? 0,
      joinedAt: '${json['joinedAt'] ?? ''}',
      muted: json['muted'] == true,
      gender: (json['gender'] as num?)?.toInt() ?? 0,
      age: (json['age'] as num?)?.toInt() ?? 0,
    );
  }

  Map<String, dynamic> toJson() => {
        'userId': userId,
        'nickname': nickname,
        'avatar': avatar,
        'vipLevel': vipLevel,
        'totalReward': totalReward,
        'joinedAt': joinedAt,
        'muted': muted,
        'gender': gender,
        'age': age,
      };
}

class OneToOneRoomConfig {
  OneToOneRoomConfig({
    required this.roomId,
    required this.status,
    required this.billing,
    required this.title,
    required this.cover,
    required this.onlineStatus,
  });

  final String roomId;
  final int status;
  final double billing;
  final String title;
  final String cover;
  final int onlineStatus;

  bool get isOnShelf => status == 1;

  factory OneToOneRoomConfig.fromJson(Map<String, dynamic> json) {
    return OneToOneRoomConfig(
      roomId: '${json['roomId'] ?? ''}',
      status: (json['status'] as num?)?.toInt() ?? 0,
      billing: (json['billing'] as num?)?.toDouble() ?? 0,
      title: '${json['title'] ?? ''}',
      cover: '${json['cover'] ?? ''}',
      onlineStatus: (json['onlineStatus'] as num?)?.toInt() ?? 2,
    );
  }

  Map<String, dynamic> toJson() => {
        'roomId': roomId,
        'status': status,
        'billing': billing,
        'title': title,
        'cover': cover,
        'onlineStatus': onlineStatus,
      };
}

class LiveRoomDetail {
  LiveRoomDetail({
    required this.roomId,
    required this.title,
    required this.cover,
    required this.notice,
    required this.status,
    required this.category,
    required this.onlineCount,
  });

  final String roomId;
  final String title;
  final String cover;
  final String notice;
  final int status;
  final int category;
  final int onlineCount;

  bool get isLive => status == 1;

  factory LiveRoomDetail.fromJson(Map<String, dynamic> json) {
    return LiveRoomDetail(
      roomId: '${json['roomId'] ?? ''}',
      title: '${json['title'] ?? ''}',
      cover: '${json['cover'] ?? ''}',
      notice: '${json['notice'] ?? ''}',
      status: (json['status'] as num?)?.toInt() ?? 0,
      category: (json['category'] as num?)?.toInt() ?? 1,
      onlineCount: (json['onlineCount'] as num?)?.toInt() ?? 0,
    );
  }

  Map<String, dynamic> toJson() => {
        'roomId': roomId,
        'title': title,
        'cover': cover,
        'notice': notice,
        'status': status,
        'category': category,
        'onlineCount': onlineCount,
      };
}

/// 秀场分类（对齐 Go category=1）
const liveRoomCategoryHot = 1;

class LiveRoomListItem {
  LiveRoomListItem({
    required this.roomId,
    required this.title,
    required this.cover,
    required this.notice,
    required this.status,
    required this.anchorNickname,
    required this.anchorAvatar,
    required this.onlineCount,
    this.rank = 0,
  });

  final String roomId;
  final String title;
  final String cover;
  final String notice;
  final int status;
  final String anchorNickname;
  final String anchorAvatar;
  final int onlineCount;
  final int rank;

  bool get isLive => status == 1;

  factory LiveRoomListItem.fromJson(Map<String, dynamic> json) {
    return LiveRoomListItem(
      roomId: '${json['roomId'] ?? ''}',
      title: '${json['title'] ?? ''}',
      cover: '${json['cover'] ?? ''}',
      notice: '${json['notice'] ?? ''}',
      status: (json['status'] as num?)?.toInt() ?? 0,
      anchorNickname: '${json['anchorNickname'] ?? ''}',
      anchorAvatar: '${json['anchorAvatar'] ?? ''}',
      onlineCount: (json['onlineCount'] as num?)?.toInt() ?? 0,
      rank: (json['rank'] as num?)?.toInt() ?? 0,
    );
  }
}

/// 观众端房间列表分页结果。
///
/// 直播页使用 `/liveRoom/roomList` 的 `statusFilter=0`，不按开播状态过滤，
/// 以便完整展示主播房间数据。
class LiveRoomListPage {
  LiveRoomListPage({
    required this.total,
    required this.page,
    required this.pageSize,
    required this.list,
  });

  final int total;
  final int page;
  final int pageSize;
  final List<LiveRoomListItem> list;

  bool get hasMore => page * pageSize < total;

  factory LiveRoomListPage.fromJson(Map<String, dynamic> json) {
    final raw = json['list'];
    final items = <LiveRoomListItem>[];
    if (raw is List) {
      for (final e in raw) {
        if (e is Map) {
          items.add(LiveRoomListItem.fromJson(Map<String, dynamic>.from(e)));
        }
      }
    }
    return LiveRoomListPage(
      total: (json['total'] as num?)?.toInt() ?? items.length,
      page: (json['page'] as num?)?.toInt() ?? 1,
      pageSize: (json['pageSize'] as num?)?.toInt() ?? 20,
      list: items,
    );
  }
}

class HotLiveRoomListPage {
  HotLiveRoomListPage({
    required this.total,
    required this.pageIndex,
    required this.pageSize,
    required this.list,
  });

  final int total;
  final int pageIndex;
  final int pageSize;
  final List<LiveRoomListItem> list;

  bool get hasMore => pageIndex * pageSize < total;

  factory HotLiveRoomListPage.fromJson(Map<String, dynamic> json) {
    final raw = json['list'];
    final items = <LiveRoomListItem>[];
    if (raw is List) {
      for (final e in raw) {
        if (e is Map) {
          items.add(
            LiveRoomListItem.fromJson(Map<String, dynamic>.from(e)),
          );
        }
      }
    }
    return HotLiveRoomListPage(
      total: (json['total'] as num?)?.toInt() ?? items.length,
      pageIndex: (json['pageIndex'] as num?)?.toInt() ?? 1,
      pageSize: (json['pageSize'] as num?)?.toInt() ?? 20,
      list: items,
    );
  }
}

class LiveRoomChatLine {
  LiveRoomChatLine({
    required this.senderName,
    required this.content,
    this.senderId = '',
    this.vipLevel = 0,
    this.isGift = false,
    this.giftSummary,
  });

  final String senderName;
  final String senderId;
  final String content;
  final int vipLevel;
  final bool isGift;
  final String? giftSummary;
}

class AudienceRestrictStatus {
  AudienceRestrictStatus({
    required this.muted,
    required this.kickBanned,
    required this.kickRemainSeconds,
  });

  final bool muted;
  final bool kickBanned;
  final int kickRemainSeconds;

  factory AudienceRestrictStatus.fromJson(Map<String, dynamic> json) {
    return AudienceRestrictStatus(
      muted: json['muted'] == true,
      kickBanned: json['kickBanned'] == true,
      kickRemainSeconds: (json['kickRemainSeconds'] as num?)?.toInt() ?? 0,
    );
  }
}

class BillingTierItem {
  BillingTierItem({
    required this.id,
    required this.pricePerMinute,
    required this.sort,
  });

  final String id;
  final double pricePerMinute;
  final int sort;

  factory BillingTierItem.fromJson(Map<String, dynamic> json) {
    return BillingTierItem(
      id: '${json['id'] ?? ''}',
      pricePerMinute: (json['pricePerMinute'] as num?)?.toDouble() ?? 0,
      sort: (json['sort'] as num?)?.toInt() ?? 0,
    );
  }
}
