import 'package:flutter/material.dart';

/// 底部导航「我的」Tab 圆形头像。
class NavProfileAvatar extends StatelessWidget {
  const NavProfileAvatar({
    super.key,
    this.size = 28,
    this.isActive = false,
  });

  final double size;
  final bool isActive;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        border: isActive
            ? Border.all(color: Colors.white, width: 1.5)
            : null,
        gradient: const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFF8B7355),
            Color(0xFF5C7A6B),
          ],
        ),
      ),
      child: Icon(
        Icons.person,
        size: size * 0.55,
        color: Colors.white.withValues(alpha: 0.85),
      ),
    );
  }
}
