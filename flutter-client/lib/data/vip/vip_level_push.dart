/// VIP 等级推送载荷，对齐 Go `vipdto.VipLevelPushItem`。
class VipLevelPush {
  const VipLevelPush({
    required this.vipLevel,
    required this.levelName,
  });

  final int vipLevel;
  final String levelName;

  factory VipLevelPush.fromData(Object? raw) {
    if (raw is! Map) {
      return const VipLevelPush(vipLevel: 0, levelName: '');
    }
    final m = Map<String, dynamic>.from(raw);
    final levelRaw = m['vipLevel'];
    final level = levelRaw is num
        ? levelRaw.toInt()
        : int.tryParse('${levelRaw ?? ''}') ?? 0;
    return VipLevelPush(
      vipLevel: level,
      levelName: '${m['levelName'] ?? ''}',
    );
  }
}
