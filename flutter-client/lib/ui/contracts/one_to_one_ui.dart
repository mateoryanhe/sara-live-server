import '../../data/one_to_one/one_to_one_models.dart';

/// 1v1 列表页 UI 回调（ViewModel 驱动）。
abstract interface class OneToOneUi {
  void showOneToOneLoading({required bool initial});

  void hideOneToOneLoading();

  void renderOneToOneList({
    required List<OneToOneRoomItem> items,
    required int total,
    required int statusFilter,
    required bool hasMore,
    required bool loadingMore,
  });

  void showOneToOneError(String message);
}
