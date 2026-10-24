/// 1v1 列表在线过滤（对齐 Go `OneToOneRoomStatusFilter*`）。
abstract final class OneToOneStatusFilter {
  static const all = 0;
  static const online = 1;
  static const offline = 2;
}

/// 主播 WebSocket 在线（Go `onlineStatus`）。
abstract final class OneToOneOnlineStatus {
  static const online = 1;
  static const offline = 2;
}

class OneToOneRoomItem {
  OneToOneRoomItem({
    required this.roomId,
    required this.title,
    required this.cover,
    required this.anchorNickname,
    required this.anchorAvatar,
    required this.billing,
    required this.onlineStatus,
    required this.allowCallIcon,
    required this.status,
  });

  final String roomId;
  final String title;
  final String cover;
  final String anchorNickname;
  final String anchorAvatar;
  final double billing;
  final int onlineStatus;
  final bool allowCallIcon;

  /// 直播间开播状态：0 未开播，1 直播中。
  final int status;

  bool get isWsOnline => onlineStatus == OneToOneOnlineStatus.online;

  factory OneToOneRoomItem.fromJson(Map<String, dynamic> json) {
    return OneToOneRoomItem(
      roomId: '${json['roomId'] ?? ''}',
      title: '${json['title'] ?? ''}',
      cover: '${json['cover'] ?? ''}',
      anchorNickname: '${json['anchorNickname'] ?? ''}',
      anchorAvatar: '${json['anchorAvatar'] ?? ''}',
      billing: _jsonDouble(json['billing']),
      onlineStatus: _jsonInt(json['onlineStatus']),
      allowCallIcon: json['allowCallIcon'] == true,
      status: _jsonInt(json['status']),
    );
  }

  Map<String, dynamic> toJson() => {
        'roomId': roomId,
        'title': title,
        'cover': cover,
        'anchorNickname': anchorNickname,
        'anchorAvatar': anchorAvatar,
        'billing': billing,
        'onlineStatus': onlineStatus,
        'allowCallIcon': allowCallIcon,
        'status': status,
      };

}

double _jsonDouble(Object? v) {
  if (v is num) return v.toDouble();
  return double.tryParse('$v') ?? 0;
}

int _jsonInt(Object? v) {
  if (v is num) return v.toInt();
  return int.tryParse('$v') ?? 0;
}

class OneToOneRoomListPage {
  OneToOneRoomListPage({
    required this.total,
    required this.page,
    required this.pageSize,
    required this.list,
  });

  final int total;
  final int page;
  final int pageSize;
  final List<OneToOneRoomItem> list;

  bool get hasMore => page * pageSize < total;

  factory OneToOneRoomListPage.fromJson(Map<String, dynamic> json) {
    final rawList = json['list'];
    final items = <OneToOneRoomItem>[];
    if (rawList is List) {
      for (final e in rawList) {
        if (e is Map) {
          items.add(
            OneToOneRoomItem.fromJson(Map<String, dynamic>.from(e)),
          );
        }
      }
    }
    return OneToOneRoomListPage(
      total: _jsonInt(json['total']),
      page: _jsonInt(json['page']),
      pageSize: _jsonInt(json['pageSize']),
      list: items,
    );
  }

  Map<String, dynamic> toJson() => {
        'total': total,
        'page': page,
        'pageSize': pageSize,
        'list': list.map((e) => e.toJson()).toList(),
      };
}
