import 'package:flutter/foundation.dart';
import '../models/auth_response.dart';
import '../services/auth_service.dart';

enum AuthState {
  initial,
  loading,
  authenticated,
  unauthenticated,
  error,
}

class AuthProvider extends ChangeNotifier {
  AuthState _state = AuthState.initial;
  AuthState get state => _state;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  String? _username;
  String? get username => _username;

  AuthResponse? _authResponse;
  AuthResponse? get authResponse => _authResponse;

  bool get isLoading => _state == AuthState.loading;
  bool get isAuthenticated => _state == AuthState.authenticated;
  bool get hasError => _state == AuthState.error;

  void initialize() {
    AuthService().initialize();
    // 初始化完成后，如果没有登录状态，显示登录页面
    _state = AuthState.unauthenticated;
    notifyListeners();
  }

  Future<bool> login(String username, String password) async {
    if (username.isEmpty || password.isEmpty) {
      _state = AuthState.error;
      _errorMessage = '请输入用户名和密码';
      notifyListeners();
      return false;
    }

    _state = AuthState.loading;
    _errorMessage = null;
    notifyListeners();

    try {
      final response = await AuthService().login(username, password);
      _authResponse = response;
      _username = response.username;
      _state = AuthState.authenticated;
      notifyListeners();
      return true;
    } on AuthException catch (e) {
      _state = AuthState.error;
      _errorMessage = e.message;
      notifyListeners();
      return false;
    } catch (e) {
      _state = AuthState.error;
      _errorMessage = '登录失败: $e';
      notifyListeners();
      return false;
    }
  }

  void logout() {
    AuthService().logout();
    _state = AuthState.unauthenticated;
    _username = null;
    _authResponse = null;
    _errorMessage = null;
    notifyListeners();
  }

  void clearError() {
    if (_state == AuthState.error) {
      _state = AuthState.unauthenticated;
      _errorMessage = null;
      notifyListeners();
    }
  }

  void reset() {
    _state = AuthState.initial;
    _errorMessage = null;
    notifyListeners();
  }
}
