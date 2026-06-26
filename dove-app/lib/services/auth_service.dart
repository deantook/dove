import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../models/auth_response.dart';

class AuthException implements Exception {
  final String message;
  final int? statusCode;

  AuthException(this.message, {this.statusCode});

  @override
  String toString() => 'AuthException: $message';
}

class AuthService {
  static final AuthService _instance = AuthService._internal();
  factory AuthService() => _instance;
  AuthService._internal();

  late final Dio _dio;
  String? _token;

  String? get token => _token;
  bool get isAuthenticated => _token != null;

  void initialize({String baseUrl = 'http://localhost:8080'}) {
    _dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 10),
        contentType: 'application/json',
      ),
    );

    // 添加请求/响应拦截器用于调试
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) {
          debugPrint('Request: ${options.method} ${options.path}');
          debugPrint('Data: ${options.data}');
          return handler.next(options);
        },
        onResponse: (response, handler) {
          debugPrint('Response: ${response.statusCode} ${response.data}');
          return handler.next(response);
        },
        onError: (error, handler) {
          debugPrint('Error: ${error.response?.statusCode} ${error.message}');
          return handler.next(error);
        },
      ),
    );
  }

  Future<AuthResponse> login(String username, String password) async {
    try {
      final request = LoginRequest(
        username: username.trim(),
        password: password,
      );

      final response = await _dio.post(
        '/api/auth/login',
        data: request.toJson(),
      );

      if (response.statusCode == 200) {
        final authResponse = AuthResponse.fromJson(response.data);
        _token = authResponse.token;
        return authResponse;
      } else {
        throw AuthException('登录失败: 服务器返回 ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response!.statusCode;
        final data = e.response!.data;
        
        if (statusCode == 401) {
          final message = data['message'] ?? '用户名或密码错误';
          throw AuthException(message, statusCode: 401);
        } else if (statusCode == 400) {
          final message = data['message'] ?? '请求参数无效';
          throw AuthException(message, statusCode: 400);
        } else {
          throw AuthException('登录失败: 服务器错误 ($statusCode)', statusCode: statusCode);
        }
      } else {
        throw AuthException('网络错误: ${e.message}');
      }
    } catch (e) {
      throw AuthException('未知错误: $e');
    }
  }

  void logout() {
    _token = null;
  }

  void setToken(String token) {
    _token = token;
  }
}
