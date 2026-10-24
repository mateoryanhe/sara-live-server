class UserProfile {
  UserProfile({
    required this.userId,
    required this.nickname,
    required this.avatar,
    required this.gold,
    required this.diamond,
    required this.shareCode,
    required this.vipLevel,
    required this.isAnchor,
    required this.prettyId,
    required this.gender,
    required this.followCount,
    required this.followerCount,
  });

  final String userId;
  final String nickname;
  final String avatar;
  final double gold;
  final double diamond;
  final String shareCode;
  final int vipLevel;
  final bool isAnchor;
  final String prettyId;
  final int gender;
  final int followCount;
  final int followerCount;

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      userId: '${json['userId'] ?? ''}',
      nickname: '${json['nickname'] ?? ''}',
      avatar: '${json['avatar'] ?? ''}',
      gold: _toDouble(json['gold']),
      diamond: _toDouble(json['diamond']),
      shareCode: '${json['shareCode'] ?? ''}',
      vipLevel: _toInt(json['vipLevel']),
      isAnchor: json['isAnchor'] == true,
      prettyId: '${json['prettyId'] ?? ''}',
      gender: _toInt(json['gender']),
      followCount: _toInt(json['followCount']),
      followerCount: _toInt(json['followerCount']),
    );
  }

  UserProfile copyWith({
    String? nickname,
    String? avatar,
  }) {
    return UserProfile(
      userId: userId,
      nickname: nickname ?? this.nickname,
      avatar: avatar ?? this.avatar,
      gold: gold,
      diamond: diamond,
      shareCode: shareCode,
      vipLevel: vipLevel,
      isAnchor: isAnchor,
      prettyId: prettyId,
      gender: gender,
      followCount: followCount,
      followerCount: followerCount,
    );
  }

  Map<String, dynamic> toJson() => {
        'userId': userId,
        'nickname': nickname,
        'avatar': avatar,
        'gold': gold,
        'diamond': diamond,
        'shareCode': shareCode,
        'vipLevel': vipLevel,
        'isAnchor': isAnchor,
        'prettyId': prettyId,
        'gender': gender,
        'followCount': followCount,
        'followerCount': followerCount,
      };

  static double _toDouble(Object? v) {
    if (v is num) return v.toDouble();
    return double.tryParse('$v') ?? 0;
  }

  static int _toInt(Object? v) {
    if (v is num) return v.toInt();
    return int.tryParse('$v') ?? 0;
  }
}
