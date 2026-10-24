import 'device_tier.dart';

/// API 业务对象磁盘 L2 容量（按设备档位）。
class ApiDiskCapacity {
  const ApiDiskCapacity({
    required this.maxFiles,
    required this.maxBytes,
    required this.defaultTtl,
  });

  final int maxFiles;
  final int maxBytes;
  final Duration defaultTtl;

  static ApiDiskCapacity forTier(DeviceTier tier) {
    switch (tier) {
      case DeviceTier.low:
        return const ApiDiskCapacity(
          maxFiles: 120,
          maxBytes: 12 << 20,
          defaultTtl: Duration(hours: 12),
        );
      case DeviceTier.mid:
        return const ApiDiskCapacity(
          maxFiles: 300,
          maxBytes: 32 << 20,
          defaultTtl: Duration(hours: 24),
        );
      case DeviceTier.high:
        return const ApiDiskCapacity(
          maxFiles: 600,
          maxBytes: 64 << 20,
          defaultTtl: Duration(days: 3),
        );
    }
  }
}
