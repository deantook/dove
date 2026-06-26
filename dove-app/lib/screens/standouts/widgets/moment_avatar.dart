import 'package:flutter/material.dart';

import '../../../widgets/user_avatar.dart';

class MomentAvatar extends StatelessWidget {
  const MomentAvatar({
    super.key,
    required this.name,
    required this.colors,
    this.size = 40,
    this.enableProfileNavigation = true,
  });

  final String name;
  final List<Color> colors;
  final double size;
  final bool enableProfileNavigation;

  @override
  Widget build(BuildContext context) {
    return UserAvatar(
      name: name,
      colors: colors,
      size: size,
      enableProfileNavigation: enableProfileNavigation,
    );
  }
}
