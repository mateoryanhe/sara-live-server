import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';

/// 跨平台通话画面（本地/远端）。
abstract class RtcVideoPlatform {
  /// 远端画面可展示（已收到主播视频或 Web 已进房订阅）。
  ValueListenable<bool> get remoteVideoReady;

  Widget localView();

  Widget remoteView();

  Future<void> disposeViews();
}
