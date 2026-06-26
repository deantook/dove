import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../theme/app_colors.dart';
import '../theme/app_typography.dart';

/// 登录页面 - 严格遵守 Hinge 设计系统
/// - 极简主义，大量留白
/// - 白色背景占比 > 75%
/// - 品牌紫只出现在 CTA 按钮
class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  final _usernameFocusNode = FocusNode();
  final _passwordFocusNode = FocusNode();

  bool _isPasswordVisible = false;

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    _usernameFocusNode.dispose();
    _passwordFocusNode.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    final username = _usernameController.text.trim();
    final password = _passwordController.text;

    if (username.isEmpty) {
      _showError('请输入用户名');
      _usernameFocusNode.requestFocus();
      return;
    }

    if (password.isEmpty) {
      _showError('请输入密码');
      _passwordFocusNode.requestFocus();
      return;
    }

    final authProvider = context.read<AuthProvider>();
    final success = await authProvider.login(username, password);

    if (success && mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            '欢迎回来',
            style: AppTypography.caption.copyWith(color: Colors.white),
          ),
          backgroundColor: AppColors.primary,
          duration: const Duration(seconds: 2),
          behavior: SnackBarBehavior.floating,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          margin: const EdgeInsets.all(20),
        ),
      );
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(
          message,
          style: AppTypography.caption.copyWith(color: Colors.white),
        ),
        backgroundColor: const Color(0xFFE57373),
        duration: const Duration(seconds: 2),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        margin: const EdgeInsets.all(20),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // 1. 白色背景占比极高
      backgroundColor: AppColors.cardBackground,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 20),
          child: ConstrainedBox(
            constraints: BoxConstraints(
              minHeight: MediaQuery.of(context).size.height -
                  MediaQuery.of(context).padding.top -
                  MediaQuery.of(context).padding.bottom,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 2. 大量留白 - 顶部空白区（视觉上 1:0.6 比例）
                const SizedBox(height: 120),

                // 3. Page Title 风格：32px, w700, 大留白
                Text(
                  "What's your name?",
                  style: AppTypography.pageTitle,
                ),
                const SizedBox(height: 12),

                // 4. 次级说明文字：16px, #666666
                Text(
                  '输入你的用户名和密码开始',
                  style: AppTypography.body.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),

                // 5. 大区块间距 32px + 额外留白
                const SizedBox(height: 64),

                // 6. 输入区域
                _buildUsernameField(),
                const SizedBox(height: 24),
                _buildPasswordField(),

                // 7. 错误提示（如果有）
                Consumer<AuthProvider>(
                  builder: (context, authProvider, child) {
                    if (authProvider.hasError) {
                      return Padding(
                        padding: const EdgeInsets.only(top: 24),
                        child: Text(
                          authProvider.errorMessage ?? '登录失败',
                          style: AppTypography.caption.copyWith(
                            color: const Color(0xFFE57373),
                          ),
                        ),
                      );
                    }
                    return const SizedBox.shrink();
                  },
                ),

                // 8. 弹性空白区，保持内容靠上，CTA 在底部舒适位置
                // 使用 MediaQuery 获取屏幕高度，计算适当留白
                Builder(
                  builder: (context) {
                    final screenHeight = MediaQuery.of(context).size.height;
                    // 小屏幕使用较小留白，大屏幕使用较大留白
                    // 保持内容:留白比例约 1:0.6
                    final breathingRoom = (screenHeight * 0.15).clamp(40.0, 160.0);
                    return SizedBox(height: breathingRoom);
                  },
                ),

                // 9. 底部 CTA 区域 - 胶囊按钮，品牌紫只在此处出现
                Padding(
                  padding: const EdgeInsets.only(bottom: 40),
                  child: _buildLoginButton(),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  /// 用户名输入框
  /// - 高度 48px，圆角 12px
  /// - 无重边框，使用 #F7F7F7 背景
  /// - 无阴影
  Widget _buildUsernameField() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Caption 标签：14px, #666666
        Text(
          '用户名',
          style: AppTypography.caption.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 12),

        // Text Input: 48px 高，12px 圆角，#F7F7F7 背景
        TextField(
          controller: _usernameController,
          focusNode: _usernameFocusNode,
          style: AppTypography.body,
          textInputAction: TextInputAction.next,
          onSubmitted: (_) => _passwordFocusNode.requestFocus(),
          decoration: InputDecoration(
            hintText: '输入用户名',
            hintStyle: AppTypography.body.copyWith(
              color: AppColors.textTertiary,
            ),
            filled: true,
            fillColor: AppColors.backgroundSecondary,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide.none,
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide.none,
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(
                color: AppColors.primary,
                width: 1.5,
              ),
            ),
            contentPadding: const EdgeInsets.symmetric(
              horizontal: 20,
              vertical: 12,
            ),
          ),
        ),
      ],
    );
  }

  /// 密码输入框
  Widget _buildPasswordField() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '密码',
          style: AppTypography.caption.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 12),

        TextField(
          controller: _passwordController,
          focusNode: _passwordFocusNode,
          obscureText: !_isPasswordVisible,
          style: AppTypography.body,
          textInputAction: TextInputAction.done,
          onSubmitted: (_) => _handleLogin(),
          decoration: InputDecoration(
            hintText: '输入密码（任意位数）',
            hintStyle: AppTypography.body.copyWith(
              color: AppColors.textTertiary,
            ),
            filled: true,
            fillColor: AppColors.backgroundSecondary,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide.none,
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide.none,
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(
                color: AppColors.primary,
                width: 1.5,
              ),
            ),
            contentPadding: const EdgeInsets.symmetric(
              horizontal: 20,
              vertical: 12,
            ),
            // 密码可见性切换按钮
            suffixIcon: IconButton(
              icon: Icon(
                _isPasswordVisible
                    ? Icons.visibility_off_outlined
                    : Icons.visibility_outlined,
                color: AppColors.textTertiary,
                size: 20,
              ),
              onPressed: () {
                setState(() {
                  _isPasswordVisible = !_isPasswordVisible;
                });
              },
            ),
          ),
        ),
      ],
    );
  }

  /// 登录按钮 - 品牌紫只在此处出现
  /// - 胶囊按钮：48px 高，24px 圆角
  /// - #7C3AED 背景，白色文字
  /// - font-weight: 600
  /// - 无阴影
  Widget _buildLoginButton() {
    return Consumer<AuthProvider>(
      builder: (context, authProvider, child) {
        return SizedBox(
          width: double.infinity,
          height: 48,
          child: ElevatedButton(
            onPressed: authProvider.isLoading ? null : _handleLogin,
            style: ElevatedButton.styleFrom(
              // 品牌紫 - 只用于关键 CTA
              backgroundColor: AppColors.primary,
              foregroundColor: Colors.white,
              disabledBackgroundColor: AppColors.primary.withValues(alpha: 0.5),
              // 无阴影
              elevation: 0,
              shadowColor: Colors.transparent,
              // 胶囊按钮：24px 圆角
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(24),
              ),
            ),
            child: authProvider.isLoading
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      color: Colors.white,
                      strokeWidth: 2,
                    ),
                  )
                : const Text(
                    '继续',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      height: 24 / 16,
                    ),
                  ),
          ),
        );
      },
    );
  }
}
