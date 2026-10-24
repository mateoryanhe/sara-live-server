import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:http/http.dart' as http;

import 'api_exception.dart';
import 'api_response.dart';
import 'auth_token.dart';
import 'http_logger.dart';

/// 对齐 Go App：`POST` + JSON，勿传 `X-H5-Client`。
class ApiClient {
  ApiClient({
    required String Function() baseUrlProvider,
    AuthTokenProvider? tokenProvider,
    this.acceptLanguage,
    this.appPackage,
    this.timeout = const Duration(seconds: 20),
    /// 失败后再试次数（不含首次）；默认 2 ⇒ 最多共 3 次。
    this.maxRetries = 2,
    /// 首次重试前等待；之后按次数线性递增。
    this.retryDelay = const Duration(milliseconds: 400),
    http.Client? httpClient,
    void Function(ApiException error)? onUnauthorized,
  })  : _baseUrlProvider = baseUrlProvider,
        _tokenProvider = tokenProvider,
        _http = httpClient ?? http.Client(),
        _onUnauthorized = onUnauthorized;

  /// 与 Go `httpserver.ReqId` 对齐。
  static const reqIdHeader = 'reqId';

  final String Function() _baseUrlProvider;

  /// 如 `https://api.example.com`，不要以 `/` 结尾。
  String get baseUrl => _baseUrlProvider();
  final AuthTokenProvider? _tokenProvider;
  final String? acceptLanguage;
  final String? appPackage;
  final Duration timeout;
  final int maxRetries;
  final Duration retryDelay;
  final http.Client _http;
  final void Function(ApiException error)? _onUnauthorized;

  /// [path] 如 `/auth/emailLogin`。
  Future<ApiResponse<T>> post<T>(
    String path, {
    Map<String, dynamic>? body,
    bool requireAuth = true,
    T? Function(Object? raw)? parseData,
    Map<String, String>? extraHeaders,
  }) {
    final reqId = _newReqId();
    return _withRetry(() async {
      final uri = Uri.parse('$baseUrl$path');
      final headers = <String, String>{
        'Content-Type': 'application/json; charset=utf-8',
        'Accept': 'application/json',
        reqIdHeader: reqId,
        if (acceptLanguage != null && acceptLanguage!.isNotEmpty)
          'Accept-Language': acceptLanguage!,
        if (appPackage != null && appPackage!.isNotEmpty)
          'X-App-Package': appPackage!,
        ...?extraHeaders,
      };
      // extraHeaders 若覆盖成空，回填本次 reqId
      if (headers[reqIdHeader] == null || headers[reqIdHeader]!.isEmpty) {
        headers[reqIdHeader] = reqId;
      }

      if (requireAuth) {
        _attachAuth(headers);
      }

      final encodedBody = jsonEncode(body ?? const {});
      final sw = Stopwatch()..start();
      late http.Response resp;
      try {
        resp = await _http
            .post(uri, headers: headers, body: encodedBody)
            .timeout(timeout);
      } catch (e) {
        sw.stop();
        await HttpLogger.log(
          method: 'POST',
          uri: uri,
          headers: headers,
          requestBody: encodedBody,
          statusCode: null,
          responseBody: null,
          elapsed: sw.elapsed,
          error: e,
        );
        throw ApiException(code: -1, message: e.toString());
      }

      return _parseAndLog(
        resp: resp,
        method: 'POST',
        uri: uri,
        headers: headers,
        requestBody: encodedBody,
        elapsed: sw.elapsed,
        parseData: parseData,
      );
    });
  }

  /// multipart 仅表单字段（无文件），如 `liveRoom/create` 不传封面。
  Future<ApiResponse<T>> postMultipartFields<T>(
    String path, {
    Map<String, String>? fields,
    bool requireAuth = true,
    T? Function(Object? raw)? parseData,
  }) {
    final reqId = _newReqId();
    return _withRetry(() async {
      final uri = Uri.parse('$baseUrl$path');
      final request = http.MultipartRequest('POST', uri);
      request.headers['Accept'] = 'application/json';
      request.headers[reqIdHeader] = reqId;
      if (acceptLanguage != null && acceptLanguage!.isNotEmpty) {
        request.headers['Accept-Language'] = acceptLanguage!;
      }
      if (appPackage != null && appPackage!.isNotEmpty) {
        request.headers['X-App-Package'] = appPackage!;
      }
      if (requireAuth) {
        _attachAuth(request.headers);
      }
      if (fields != null) {
        request.fields.addAll(fields);
      }

      final reqHeaders = Map<String, String>.from(request.headers);
      final sw = Stopwatch()..start();
      late http.Response resp;
      try {
        final streamed = await _http.send(request).timeout(timeout);
        resp = await http.Response.fromStream(streamed).timeout(timeout);
      } catch (e) {
        sw.stop();
        await HttpLogger.log(
          method: 'POST',
          uri: uri,
          headers: reqHeaders,
          requestBody: 'multipart fields-only ${fields ?? {}}',
          statusCode: null,
          responseBody: null,
          elapsed: sw.elapsed,
          error: e,
          note: 'multipart-fields',
        );
        throw ApiException(code: -1, message: e.toString());
      }

      return _parseAndLog(
        resp: resp,
        method: 'POST',
        uri: uri,
        headers: reqHeaders,
        requestBody: 'multipart fields-only ${fields ?? {}}',
        elapsed: sw.elapsed,
        parseData: parseData,
        note: 'multipart-fields',
      );
    });
  }

  /// multipart 上传，如头像字段名 `file`（对齐 Go UploadAvatarReq）。
  Future<ApiResponse<T>> postMultipart<T>(
    String path, {
    required String fileField,
    required List<int> fileBytes,
    required String filename,
    Map<String, String>? fields,
    bool requireAuth = true,
    T? Function(Object? raw)? parseData,
  }) {
    final reqId = _newReqId();
    return _withRetry(() async {
      final uri = Uri.parse('$baseUrl$path');
      final request = http.MultipartRequest('POST', uri);
      request.headers['Accept'] = 'application/json';
      request.headers[reqIdHeader] = reqId;
      if (acceptLanguage != null && acceptLanguage!.isNotEmpty) {
        request.headers['Accept-Language'] = acceptLanguage!;
      }
      if (appPackage != null && appPackage!.isNotEmpty) {
        request.headers['X-App-Package'] = appPackage!;
      }
      if (requireAuth) {
        _attachAuth(request.headers);
      }
      if (fields != null) {
        request.fields.addAll(fields);
      }
      request.files.add(
        http.MultipartFile.fromBytes(
          fileField,
          fileBytes,
          filename: filename,
        ),
      );

      final reqHeaders = Map<String, String>.from(request.headers);
      final sw = Stopwatch()..start();
      late http.Response resp;
      try {
        final streamed = await _http.send(request).timeout(timeout);
        resp = await http.Response.fromStream(streamed).timeout(timeout);
      } catch (e) {
        sw.stop();
        await HttpLogger.log(
          method: 'POST',
          uri: uri,
          headers: reqHeaders,
          requestBody: 'multipart field=$fileField file=$filename bytes=${fileBytes.length}',
          statusCode: null,
          responseBody: null,
          elapsed: sw.elapsed,
          error: e,
          note: 'multipart',
        );
        throw ApiException(code: -1, message: e.toString());
      }

      return _parseAndLog(
        resp: resp,
        method: 'POST',
        uri: uri,
        headers: reqHeaders,
        requestBody:
            'multipart field=$fileField file=$filename bytes=${fileBytes.length}',
        elapsed: sw.elapsed,
        parseData: parseData,
        note: 'multipart',
      );
    });
  }

  /// Go 日志按数字解析 reqId；用微秒时间戳（可转 uint64）。
  String _newReqId() => DateTime.now().microsecondsSinceEpoch.toString();

  Future<T> _withRetry<T>(Future<T> Function() action) async {
    final retries = maxRetries < 0 ? 0 : maxRetries;
    var attempt = 0;
    while (true) {
      try {
        return await action();
      } catch (e) {
        if (!_shouldRetry(e) || attempt >= retries) {
          rethrow;
        }
        attempt++;
        final wait = retryDelay * attempt;
        await Future<void>.delayed(wait);
      }
    }
  }

  bool _shouldRetry(Object error) {
    if (error is! ApiException) {
      return error is TimeoutException ||
          error is SocketException ||
          error is HttpException ||
          error is IOException;
    }
    // 业务失败 / 鉴权失败不重试
    if (error.isUnauthorized || error.code > 0) {
      return false;
    }
    final status = error.statusCode;
    if (status == null) {
      // 网络超时等包装成 code=-1、无 statusCode
      return true;
    }
    return status == 408 || status == 429 || status >= 500;
  }

  void _attachAuth(Map<String, String> headers) {
    final token = _tokenProvider?.call();
    if (token == null) {
      throw ApiException(code: 3, message: 'not logged in');
    }
    headers['Authorization'] = token.headerValue;
  }

  Future<ApiResponse<T>> _parseAndLog<T>({
    required http.Response resp,
    required String method,
    required Uri uri,
    required Map<String, String> headers,
    required String? requestBody,
    required Duration elapsed,
    T? Function(Object? raw)? parseData,
    String? note,
  }) async {
    final bodyStr = utf8.decode(resp.bodyBytes);
    try {
      final api = _parseResponseBody<T>(
        statusCode: resp.statusCode,
        bodyStr: bodyStr,
        parseData: parseData,
      );
      await HttpLogger.log(
        method: method,
        uri: uri,
        headers: headers,
        requestBody: requestBody,
        statusCode: resp.statusCode,
        responseBody: bodyStr,
        elapsed: elapsed,
        note: note,
      );
      return api;
    } catch (e) {
      await HttpLogger.log(
        method: method,
        uri: uri,
        headers: headers,
        requestBody: requestBody,
        statusCode: resp.statusCode,
        responseBody: bodyStr,
        elapsed: elapsed,
        error: e,
        note: note,
      );
      rethrow;
    }
  }

  ApiResponse<T> _parseResponseBody<T>({
    required int statusCode,
    required String bodyStr,
    T? Function(Object? raw)? parseData,
  }) {
    if (statusCode < 200 || statusCode >= 300) {
      throw ApiException(
        code: -1,
        message: 'http $statusCode',
        statusCode: statusCode,
      );
    }

    final decoded = jsonDecode(bodyStr);
    if (decoded is! Map<String, dynamic>) {
      throw ApiException(code: -1, message: 'invalid response body');
    }

    final api = ApiResponse<T>.fromJson(decoded, parseData: parseData);
    if (!api.isSuccess) {
      final err = ApiException(
        code: api.code,
        message: api.message,
        data: api.data,
        statusCode: statusCode,
      );
      if (err.isUnauthorized) {
        _onUnauthorized?.call(err);
      }
      throw err;
    }
    return api;
  }

  void close() => _http.close();
}
