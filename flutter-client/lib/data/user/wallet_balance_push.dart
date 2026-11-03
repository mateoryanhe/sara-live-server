import '../../core/ws/push_cmd.dart';
import '../call/call_models.dart';

/// 对齐 Go `diamonddto.DiamondPushItem` / `golddto.GoldPushItem`。
sealed class WalletBalancePush {
  const WalletBalancePush();

  factory WalletBalancePush.fromCmd(int cmd, Object? raw) {
    final map = pushDataMap(raw);
    if (map == null) {
      return const WalletBalancePushUnknown();
    }
    switch (cmd) {
      case PushCmd.diamond:
        return WalletDiamondPush(
          balance: _toDouble(map['diamond']),
        );
      case PushCmd.gold:
        return WalletGoldPush(
          balance: _toDouble(map['gold']),
        );
      default:
        return const WalletBalancePushUnknown();
    }
  }
}

final class WalletDiamondPush extends WalletBalancePush {
  const WalletDiamondPush({required this.balance});

  final double balance;
}

final class WalletGoldPush extends WalletBalancePush {
  const WalletGoldPush({required this.balance});

  final double balance;
}

final class WalletBalancePushUnknown extends WalletBalancePush {
  const WalletBalancePushUnknown();
}

double _toDouble(Object? v) {
  if (v is num) return v.toDouble();
  return double.tryParse('$v') ?? 0;
}
