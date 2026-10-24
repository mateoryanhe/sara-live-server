import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:sara_live_app/app/app_routes.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/ui/contracts/login_ui.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/viewmodel/login_view_model.dart';

/// 启动：自动设备码登录（当前唯一登录方式）。
class SplashPage extends StatefulWidget {
  const SplashPage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> implements LoginUi {
  late final LoginViewModel _vm;
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _vm = LoginViewModel(scope: widget.scope, ui: this);
    WidgetsBinding.instance.addPostFrameCallback((_) => _vm.start());
  }

  @override
  void showLoginLoading() {
    if (!mounted) return;
    setState(() {
      _loading = true;
      _error = null;
    });
  }

  @override
  void hideLoginLoading() {
    if (!mounted) return;
    setState(() => _loading = false);
  }

  @override
  void showLoginError(String message) {
    if (!mounted) return;
    setState(() => _error = message);
  }

  @override
  void onLoginSuccess({required bool isNewUser}) {
    if (!mounted) return;
    context.go(AppRoutes.main);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    return Scaffold(
      body: AppPageFrame(
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(28),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Container(
                  width: 72,
                  height: 72,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: LinearGradient(
                      colors: [
                        cs.primary,
                        cs.primary.withValues(alpha: 0.55),
                      ],
                    ),
                  ),
                  child: Icon(
                    Icons.live_tv_rounded,
                    color: cs.onPrimary,
                    size: 34,
                  ),
                ),
                const SizedBox(height: 20),
                Text(
                  'Sara Live',
                  style: theme.textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.w800,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  '正在进入…',
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: cs.onSurface.withValues(alpha: 0.55),
                  ),
                ),
                const SizedBox(height: 28),
                if (_loading) const CircularProgressIndicator(),
                if (_error != null) ...[
                  Text(_error!, textAlign: TextAlign.center),
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: () => _vm.start(forceRelogin: true),
                    child: const Text('重试登录'),
                  ),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }
}
