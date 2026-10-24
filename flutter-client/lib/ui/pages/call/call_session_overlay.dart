import 'package:flutter/material.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/core/widgets/cached_avatar.dart';
import 'package:sara_live_app/feature/call/call_controller.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';

/// 全屏通话层（覆盖在主壳之上）。
class CallSessionOverlay extends StatelessWidget {
  const CallSessionOverlay({super.key, required this.scope});

  final AppScope scope;

  @override
  Widget build(BuildContext context) {
    final ctrl = scope.callController;
    final peer = ctrl.peer;
    final theme = Theme.of(context);

    return Material(
      color: Colors.black,
      child: SafeArea(
        child: AppPageFrame(
          child: Stack(
            fit: StackFit.expand,
            children: [
              if (ctrl.phase == CallPhase.inCall ||
                  ctrl.phase == CallPhase.connecting)
                Positioned.fill(
                  child: ctrl.rtcSession.remoteView(),
                ),
              if (ctrl.phase == CallPhase.inCall ||
                  ctrl.phase == CallPhase.connecting)
                Positioned(
                  top: 12,
                  right: 12,
                  width: 96,
                  height: 144,
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(10),
                    child: ctrl.rtcSession.localView(),
                  ),
                ),
              Positioned(
                left: 0,
                right: 0,
                bottom: 0,
                child: Container(
                  padding: const EdgeInsets.fromLTRB(20, 24, 20, 28),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.bottomCenter,
                      end: Alignment.topCenter,
                      colors: [
                        Colors.black.withValues(alpha: 0.85),
                        Colors.transparent,
                      ],
                    ),
                  ),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (peer != null) ...[
                        CachedAvatar(
                          url: peer.avatar,
                          radius: 36,
                          fallbackText: peer.nickname.isNotEmpty
                              ? peer.nickname[0]
                              : '?',
                        ),
                        const SizedBox(height: 10),
                        Text(
                          peer.nickname.isNotEmpty
                              ? peer.nickname
                              : '用户 ${peer.userId}',
                          style: theme.textTheme.titleLarge?.copyWith(
                            color: Colors.white,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ],
                      const SizedBox(height: 8),
                      Text(
                        ctrl.statusLine,
                        textAlign: TextAlign.center,
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: Colors.white70,
                        ),
                      ),
                      const SizedBox(height: 24),
                      _Actions(controller: ctrl),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _Actions extends StatelessWidget {
  const _Actions({required this.controller});

  final CallController controller;

  @override
  Widget build(BuildContext context) {
    switch (controller.phase) {
      case CallPhase.incomingRinging:
        return Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _RoundBtn(
              color: Colors.red,
              icon: Icons.call_end,
              label: '拒接',
              onTap: () => controller.rejectIncoming(),
            ),
            _RoundBtn(
              color: Colors.green,
              icon: Icons.videocam,
              label: '接听',
              onTap: () => controller.acceptIncoming(),
            ),
          ],
        );
      case CallPhase.outgoingRinging:
      case CallPhase.connecting:
      case CallPhase.inCall:
        return _RoundBtn(
          color: Colors.red,
          icon: Icons.call_end,
          label: '挂断',
          onTap: () => controller.hangUp(),
        );
      case CallPhase.idle:
        return const SizedBox.shrink();
    }
  }
}

class _RoundBtn extends StatelessWidget {
  const _RoundBtn({
    required this.color,
    required this.icon,
    required this.label,
    required this.onTap,
  });

  final Color color;
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Material(
          color: color,
          shape: const CircleBorder(),
          clipBehavior: Clip.antiAlias,
          child: InkWell(
            onTap: onTap,
            child: SizedBox(
              width: 64,
              height: 64,
              child: Icon(icon, color: Colors.white, size: 30),
            ),
          ),
        ),
        const SizedBox(height: 8),
        Text(
          label,
          style: const TextStyle(color: Colors.white70),
        ),
      ],
    );
  }
}
