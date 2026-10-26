import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:sara_live_app/app/app_routes.dart';
import 'package:sara_live_app/app/app_scope.dart';
import 'package:sara_live_app/l10n/app_localizations.dart';
import 'package:sara_live_app/ui/contracts/login_ui.dart';
import 'package:sara_live_app/ui/widgets/app_page_frame.dart';
import 'package:sara_live_app/viewmodel/login_view_model.dart';

/// Explicit account-switch login page.
///
/// This page is intentionally separate from [SplashPage]. Splash performs an
/// automatic login with the persisted device code, while this page waits for
/// the user to enter (and confirm) a device code after the previous session
/// has been cleared.
class DeviceCodeLoginPage extends StatefulWidget {
  const DeviceCodeLoginPage({super.key, required this.scope});

  final AppScope scope;

  @override
  State<DeviceCodeLoginPage> createState() => _DeviceCodeLoginPageState();
}

class _DeviceCodeLoginPageState extends State<DeviceCodeLoginPage>
    implements LoginUi {
  late final LoginViewModel _vm;
  late final TextEditingController _controller;
  bool _loading = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _vm = LoginViewModel(scope: widget.scope, ui: this);
    _controller = TextEditingController(
      text: widget.scope.session.storedDeviceId ?? '',
    );
    if (_controller.text.isEmpty) {
      widget.scope.session.getOrCreateDeviceId().then((id) {
        if (!mounted || _controller.text.isNotEmpty) return;
        _controller.text = id;
        _controller.selection = TextSelection.collapsed(offset: id.length);
      });
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
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
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }

  @override
  void onLoginSuccess({required bool isNewUser}) {
    if (!mounted) return;
    context.go(AppRoutes.main);
  }

  Future<void> _submit() async {
    if (_loading) return;
    final l10n = AppLocalizations.of(context);
    final value = _controller.text.trim();
    if (value.isEmpty) {
      showLoginError(l10n.deviceCodeEmpty);
      return;
    }
    if (value.length > 128) {
      showLoginError(l10n.deviceCodeTooLong);
      return;
    }
    FocusManager.instance.primaryFocus?.unfocus();
    await _vm.loginWithDeviceId(value);
  }

  Future<void> _cancel() async {
    if (_loading) return;
    await widget.scope.cancelDeviceCodeLogin();
    if (!mounted) return;
    context.go(AppRoutes.splash);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);
    final cs = theme.colorScheme;

    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) {
        if (!didPop) {
          unawaited(_cancel());
        }
      },
      child: Scaffold(
        appBar: AppBar(
          title: Text(l10n.deviceCodeEditTitle),
          leading: IconButton(
            tooltip: l10n.commonCancel,
            onPressed: _loading ? null : _cancel,
            icon: const Icon(Icons.arrow_back_rounded),
          ),
        ),
        body: AppPageFrame(
          child: SafeArea(
            top: false,
            child: SingleChildScrollView(
              padding: const EdgeInsets.fromLTRB(24, 28, 24, 32),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Container(
                    width: 68,
                    height: 68,
                    alignment: Alignment.center,
                    decoration: BoxDecoration(
                      color: cs.primary.withValues(alpha: 0.12),
                      shape: BoxShape.circle,
                    ),
                    child: Icon(
                      Icons.phonelink_setup_rounded,
                      size: 34,
                      color: cs.primary,
                    ),
                  ),
                  const SizedBox(height: 20),
                  Text(
                    l10n.deviceCodeEditTitle,
                    style: theme.textTheme.headlineSmall?.copyWith(
                      fontWeight: FontWeight.w800,
                    ),
                  ),
                  const SizedBox(height: 10),
                  Text(
                    l10n.deviceCodeWarning,
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: cs.onSurface.withValues(alpha: 0.65),
                    ),
                  ),
                  const SizedBox(height: 24),
                  TextField(
                    controller: _controller,
                    autofocus: true,
                    enabled: !_loading,
                    maxLength: 128,
                    textInputAction: TextInputAction.done,
                    decoration: InputDecoration(
                      labelText: l10n.deviceCodeTitle,
                      hintText: l10n.deviceCodeHint,
                      prefixIcon: const Icon(Icons.key_rounded),
                      border: const OutlineInputBorder(),
                    ),
                    onSubmitted: (_) => _submit(),
                  ),
                  if (_error != null) ...[
                    const SizedBox(height: 8),
                    Text(
                      _error!,
                      style: TextStyle(color: cs.error),
                    ),
                  ],
                  const SizedBox(height: 16),
                  SizedBox(
                    height: 48,
                    child: FilledButton.icon(
                      onPressed: _loading ? null : _submit,
                      icon: _loading
                          ? const SizedBox(
                              width: 18,
                              height: 18,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                          : const Icon(Icons.login_rounded),
                      label: Text(l10n.commonSave),
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextButton(
                    onPressed: _loading ? null : _cancel,
                    child: Text(l10n.commonCancel),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
