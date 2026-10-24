/// 对齐 Go `vipdto.AppVipDetailRes` / `vipcfgdto.AppVipCfgItem`。
class VipDetail {
  const VipDetail({
    required this.vipLevel,
    required this.levelName,
    required this.totalRecharge,
    required this.isMaxLevel,
    required this.nextLevel,
    required this.nextUpgradeLimit,
    required this.rechargeToNextLevel,
    this.currentCfg,
    this.nextCfg,
  });

  final int vipLevel;
  final String levelName;
  final double totalRecharge;
  final bool isMaxLevel;
  final int nextLevel;
  final double nextUpgradeLimit;
  final double rechargeToNextLevel;
  final VipCfgItem? currentCfg;
  final VipCfgItem? nextCfg;

  bool get hasVip => vipLevel > 0;

  /// 「我的」列表副文案（进度；不含等级名，等级在 value）。
  String get mineSummary {
    if (!hasVip) {
      if (nextLevel > 0) {
        return '距 Lv.$nextLevel 还需 ${_fmt(rechargeToNextLevel)} 金币';
      }
      return '';
    }
    if (isMaxLevel) return '已满级';
    if (nextLevel > 0) {
      return '距下一级 ${_fmt(rechargeToNextLevel)} 金币';
    }
    return '';
  }

  String get mineTitleValue {
    if (!hasVip) return '未开通';
    if (levelName.isNotEmpty) return levelName;
    return 'Lv.$vipLevel';
  }

  factory VipDetail.fromJson(Map<String, dynamic> json) {
    return VipDetail(
      vipLevel: _toInt(json['vipLevel']),
      levelName: '${json['levelName'] ?? ''}',
      totalRecharge: _toDouble(json['totalRecharge']),
      isMaxLevel: json['isMaxLevel'] == true,
      nextLevel: _toInt(json['nextLevel']),
      nextUpgradeLimit: _toDouble(json['nextUpgradeLimit']),
      rechargeToNextLevel: _toDouble(json['rechargeToNextLevel']),
      currentCfg: _cfg(json['currentCfg']),
      nextCfg: _cfg(json['nextCfg']),
    );
  }

  Map<String, dynamic> toJson() => {
        'vipLevel': vipLevel,
        'levelName': levelName,
        'totalRecharge': totalRecharge,
        'isMaxLevel': isMaxLevel,
        'nextLevel': nextLevel,
        'nextUpgradeLimit': nextUpgradeLimit,
        'rechargeToNextLevel': rechargeToNextLevel,
        if (currentCfg != null) 'currentCfg': currentCfg!.toJson(),
        if (nextCfg != null) 'nextCfg': nextCfg!.toJson(),
      };

  static VipCfgItem? _cfg(Object? raw) {
    if (raw is! Map) return null;
    return VipCfgItem.fromJson(Map<String, dynamic>.from(raw));
  }

  static String _fmt(double v) {
    if (v == v.roundToDouble()) return '${v.toInt()}';
    return v.toStringAsFixed(2);
  }
}

class VipCfgItem {
  const VipCfgItem({
    required this.level,
    required this.levelName,
    required this.levelIcon,
    required this.upgradeRechargeLimit,
    required this.privilegeList,
  });

  final int level;
  final String levelName;
  final String levelIcon;
  final double upgradeRechargeLimit;
  final List<VipPrivilegeItem> privilegeList;

  factory VipCfgItem.fromJson(Map<String, dynamic> json) {
    final rawList = json['privilegeList'];
    final list = <VipPrivilegeItem>[];
    if (rawList is List) {
      for (final item in rawList) {
        if (item is Map) {
          list.add(VipPrivilegeItem.fromJson(Map<String, dynamic>.from(item)));
        }
      }
    }
    return VipCfgItem(
      level: _toInt(json['level']),
      levelName: '${json['levelName'] ?? ''}',
      levelIcon: '${json['levelIcon'] ?? ''}',
      upgradeRechargeLimit: _toDouble(json['upgradeRechargeLimit']),
      privilegeList: list,
    );
  }

  Map<String, dynamic> toJson() => {
        'level': level,
        'levelName': levelName,
        'levelIcon': levelIcon,
        'upgradeRechargeLimit': upgradeRechargeLimit,
        'privilegeList': privilegeList.map((e) => e.toJson()).toList(),
      };
}

class VipPrivilegeItem {
  const VipPrivilegeItem({
    required this.privilegeType,
    required this.icon,
    required this.title,
    required this.desc,
    required this.animation,
  });

  /// 2=进场特效 3=公屏评论特效 4=客服优先
  final int privilegeType;
  final String icon;
  final VipI18nText title;
  final VipI18nText desc;
  final String animation;

  factory VipPrivilegeItem.fromJson(Map<String, dynamic> json) {
    return VipPrivilegeItem(
      privilegeType: _toInt(json['privilegeType']),
      icon: '${json['icon'] ?? ''}',
      title: VipI18nText.fromJson(json['title']),
      desc: VipI18nText.fromJson(json['desc']),
      animation: '${json['animation'] ?? ''}',
    );
  }

  Map<String, dynamic> toJson() => {
        'privilegeType': privilegeType,
        'icon': icon,
        'title': title.toJson(),
        'desc': desc.toJson(),
        'animation': animation,
      };

  String get typeLabel => switch (privilegeType) {
        2 => '进场特效',
        3 => '公屏特效',
        4 => '客服优先',
        _ => '特权',
      };
}

class VipI18nText {
  const VipI18nText({
    this.en = '',
    this.es = '',
    this.pt = '',
    this.hi = '',
    this.id = '',
  });

  final String en;
  final String es;
  final String pt;
  final String hi;
  final String id;

  factory VipI18nText.fromJson(Object? raw) {
    if (raw is! Map) return const VipI18nText();
    final m = Map<String, dynamic>.from(raw);
    return VipI18nText(
      en: '${m['en'] ?? ''}',
      es: '${m['es'] ?? ''}',
      pt: '${m['pt'] ?? ''}',
      hi: '${m['hi'] ?? ''}',
      id: '${m['id'] ?? ''}',
    );
  }

  Map<String, dynamic> toJson() => {
        'en': en,
        'es': es,
        'pt': pt,
        'hi': hi,
        'id': id,
      };

  /// VIP 特权文案暂无中文；优先英文，再印尼语等。
  String resolve([String? prefer]) {
    final order = <String>[
      if (prefer != null && prefer.isNotEmpty) prefer,
      'en',
      'id',
      'es',
      'pt',
      'hi',
    ];
    for (final key in order) {
      final v = switch (key) {
        'en' => en,
        'es' => es,
        'pt' => pt,
        'hi' => hi,
        'id' => id,
        _ => '',
      };
      if (v.trim().isNotEmpty) return v.trim();
    }
    return '';
  }
}

double _toDouble(Object? v) {
  if (v is num) return v.toDouble();
  return double.tryParse('$v') ?? 0;
}

int _toInt(Object? v) {
  if (v is num) return v.toInt();
  return int.tryParse('$v') ?? 0;
}
