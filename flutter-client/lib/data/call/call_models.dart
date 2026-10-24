/// 业务来源（对齐 Go `CallOrderSource*`）。
abstract final class CallSource {
  static const liveRoom = 1;
  static const privateMessage = 2;
  static const oneToOneRoom = 3;
}

class CallRtcCredentials {
  const CallRtcCredentials({
    required this.orderId,
    required this.appId,
    required this.channelName,
    required this.token,
    required this.userAccount,
    required this.expireAt,
  });

  final String orderId;
  final String appId;
  final String channelName;
  final String token;
  final String userAccount;
  final int expireAt;

  int get rtcUid => int.tryParse(userAccount) ?? 0;

  factory CallRtcCredentials.fromJson(Map<String, dynamic> json) {
    return CallRtcCredentials(
      orderId: '${json['orderId'] ?? ''}',
      appId: '${json['appId'] ?? ''}',
      channelName: '${json['channelName'] ?? ''}',
      token: '${json['token'] ?? ''}',
      userAccount: '${json['userAccount'] ?? ''}',
      expireAt: _jsonInt(json['expireAt']),
    );
  }
}

class CallConfirmStatus {
  CallConfirmStatus({
    required this.success,
    required this.orderId,
    required this.selfConfirmed,
    required this.peerConfirmed,
    required this.peerConfirmAt,
    required this.answerConfirmCount,
    required this.callStarted,
  });

  final bool success;
  final String orderId;
  final bool selfConfirmed;
  final bool peerConfirmed;
  final int peerConfirmAt;
  final int answerConfirmCount;
  final bool callStarted;

  factory CallConfirmStatus.fromJson(Map<String, dynamic> json) {
    return CallConfirmStatus(
      success: json['success'] == true,
      orderId: '${json['orderId'] ?? ''}',
      selfConfirmed: json['selfConfirmed'] == true,
      peerConfirmed: json['peerConfirmed'] == true,
      peerConfirmAt: _jsonInt(json['peerConfirmAt']),
      answerConfirmCount: _jsonInt(json['answerConfirmCount']),
      callStarted: json['callStarted'] == true,
    );
  }
}

class CallEndSummary {
  CallEndSummary({
    required this.success,
    required this.orderId,
    required this.callDuration,
    required this.billingDuration,
    required this.totalCost,
  });

  final bool success;
  final String orderId;
  final int callDuration;
  final int billingDuration;
  final double totalCost;

  factory CallEndSummary.fromJson(Map<String, dynamic> json) {
    return CallEndSummary(
      success: json['success'] == true,
      orderId: '${json['orderId'] ?? ''}',
      callDuration: _jsonInt(json['callDuration']),
      billingDuration: _jsonInt(json['billingDuration']),
      totalCost: _jsonDouble(json['totalCost']),
    );
  }
}

class CallRequestPush {
  CallRequestPush({
    required this.orderId,
    required this.source,
    required this.callerId,
    required this.channelName,
    required this.callType,
    required this.callerNickname,
    required this.callerAvatar,
    required this.receiverToken,
    required this.message,
  });

  final String orderId;
  final int source;
  final String callerId;
  final String channelName;
  final int callType;
  final String callerNickname;
  final String callerAvatar;
  final String receiverToken;
  final String message;

  factory CallRequestPush.fromJson(Map<String, dynamic> json) {
    return CallRequestPush(
      orderId: '${json['orderId'] ?? ''}',
      source: _jsonInt(json['source']),
      callerId: '${json['callerId'] ?? ''}',
      channelName: '${json['channelName'] ?? ''}',
      callType: _jsonInt(json['callType']),
      callerNickname: '${json['callerNickname'] ?? ''}',
      callerAvatar: '${json['callerAvatar'] ?? ''}',
      receiverToken: '${json['receiverToken'] ?? ''}',
      message: '${json['message'] ?? ''}',
    );
  }
}

Map<String, dynamic>? pushDataMap(Object? data) {
  if (data is Map) {
    return Map<String, dynamic>.from(data);
  }
  return null;
}

int _jsonInt(Object? v) {
  if (v is num) return v.toInt();
  return int.tryParse('$v') ?? 0;
}

double _jsonDouble(Object? v) {
  if (v is num) return v.toDouble();
  return double.tryParse('$v') ?? 0;
}
