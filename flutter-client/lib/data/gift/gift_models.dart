class AppGiftItem {
  AppGiftItem({
    required this.id,
    required this.name,
    required this.icon,
    required this.price,
    required this.sort,
  });

  final String id;
  final String name;
  final String icon;
  final double price;
  final int sort;

  factory AppGiftItem.fromJson(Map<String, dynamic> json) {
    return AppGiftItem(
      id: '${json['id'] ?? ''}',
      name: '${json['name'] ?? ''}',
      icon: '${json['icon'] ?? ''}',
      price: (json['price'] as num?)?.toDouble() ?? 0,
      sort: (json['sort'] as num?)?.toInt() ?? 0,
    );
  }
}
