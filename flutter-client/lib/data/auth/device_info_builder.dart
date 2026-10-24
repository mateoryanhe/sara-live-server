import 'package:device_info_plus/device_info_plus.dart';
import 'package:flutter/foundation.dart';
import 'package:sara_live_app/app/app_config.dart';
import 'package:sara_live_app/core/platform/platform_check.dart';

/// 组装 Go `entity.DeviceInfo` JSON。
class DeviceInfoBuilder {
  DeviceInfoBuilder({DeviceInfoPlugin? plugin})
      : _plugin = plugin ?? DeviceInfoPlugin();

  final DeviceInfoPlugin _plugin;

  Future<Map<String, dynamic>> build({required String deviceId}) async {
    String deviceType = 'android';
    String deviceModel = '';
    String cpuModel = '';
    String osVersion = 'unknown';

    try {
      if (kIsWeb) {
        deviceType = 'android';
        deviceModel = 'web';
        osVersion = 'web';
      } else if (isAndroidOs) {
        final info = await _plugin.androidInfo;
        deviceType = 'android';
        deviceModel = info.model;
        cpuModel = info.supportedAbis.isNotEmpty
            ? info.supportedAbis.first
            : info.hardware;
        osVersion = info.version.release;
      } else if (isIOSOs) {
        final info = await _plugin.iosInfo;
        deviceType = 'ios';
        deviceModel = info.utsname.machine;
        cpuModel = info.utsname.machine;
        osVersion = info.systemVersion;
      }
    } catch (_) {
      // 兜底字段仍满足服务端 required
    }

    return {
      'deviceType': deviceType,
      'deviceModel': deviceModel,
      'cpuModel': cpuModel,
      'osVersion': osVersion.isEmpty ? 'unknown' : osVersion,
      'appVersion': AppConfig.appVersion,
      'packageName': AppConfig.appPackage,
      'deviceId': deviceId,
    };
  }
}
