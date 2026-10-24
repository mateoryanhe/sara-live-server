import 'package:permission_handler/permission_handler.dart';

/// Android/iOS：发起视频通话前申请摄像头与麦克风。
Future<bool> ensureCallMediaPermission() async {
  final camera = await Permission.camera.request();
  final mic = await Permission.microphone.request();
  return camera.isGranted && mic.isGranted;
}
