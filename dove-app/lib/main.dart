import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'live_match/widgets/live_match_host.dart';
import 'providers/auth_provider.dart';
import 'screens/login_screen.dart';
import 'screens/main_shell.dart';
import 'theme/app_colors.dart';
import 'theme/app_theme.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const DoveApp());
}

class DoveApp extends StatelessWidget {
  const DoveApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) {
        final authProvider = AuthProvider();
        authProvider.initialize();
        return authProvider;
      },
      child: MaterialApp(
        title: 'Dove',
        debugShowCheckedModeBanner: false,
        theme: AppTheme.light,
        home: const AuthGate(),
      ),
    );
  }
}

/// 认证入口组件
/// 根据登录状态决定显示登录页面还是主应用
class AuthGate extends StatelessWidget {
  const AuthGate({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<AuthProvider>(
      builder: (context, authProvider, child) {
        // 根据认证状态决定显示哪个页面
        switch (authProvider.state) {
          case AuthState.initial:
            // 初始化状态显示加载页面
            return const Scaffold(
              body: Center(
                child: CircularProgressIndicator(
                  color: AppColors.primary,
                ),
              ),
            );
          case AuthState.authenticated:
            // 已登录，显示主应用
            return const LiveMatchHost(
              child: MainShell(),
            );
          case AuthState.loading:
          case AuthState.unauthenticated:
          case AuthState.error:
            // 未登录或登录中或出错，显示登录页面
            return const LoginScreen();
        }
      },
    );
  }
}
