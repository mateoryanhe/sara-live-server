/// Go `apiResponseMiddleware` 信封：`{code,message,data}`。
class ApiResponse<T> {
  const ApiResponse({
    required this.code,
    required this.message,
    this.data,
  });

  final int code;
  final String message;
  final T? data;

  bool get isSuccess => code == 0;

  factory ApiResponse.fromJson(
    Map<String, dynamic> json, {
    T? Function(Object? raw)? parseData,
  }) {
    final rawData = json['data'];
    return ApiResponse<T>(
      code: (json['code'] as num?)?.toInt() ?? -1,
      message: json['message']?.toString() ?? '',
      data: parseData != null ? parseData(rawData) : rawData as T?,
    );
  }
}
