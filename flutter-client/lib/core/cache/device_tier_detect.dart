import 'package:device_info_plus/device_info_plus.dart';
import 'package:flutter/foundation.dart';

import '../platform/platform_check.dart';
import 'device_tier.dart';

/// 探测设备档位（启动一次即可）。
Future<DeviceTier> detectDeviceTier({DeviceInfoPlugin? plugin}) async {
  if (kIsWeb) {
    return DeviceTier.mid;
  }
  final p = plugin ?? DeviceInfoPlugin();
  try {
    if (isAndroidOs) {
      final info = await p.androidInfo;
      final mb = info.physicalRamSize;
      if (mb > 0) {
        return deviceTierFromTotalMemoryMb(mb);
      }
    }
  } catch (_) {
    // fall through
  }
  return DeviceTier.mid;
}
