import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

abstract final class OtherIconAssets {
  static const settings = 'assets/icon/other/齿轮.svg';
  static const verified = 'assets/icon/other/认证.svg';
}

abstract final class NavIconAssets {
  static const discover = 'assets/icon/nav/首页.svg';
  static const standouts = 'assets/icon/nav/星星.svg';
  static const likes = 'assets/icon/nav/爱心.svg';
  static const messages = 'assets/icon/nav/消息.svg';
  static const profile = 'assets/icon/nav/用户.svg';
}

class DoveNavIcon extends StatelessWidget {
  const DoveNavIcon({
    super.key,
    required this.assetPath,
    required this.color,
    this.size = 30,
  });

  final String assetPath;
  final Color color;
  final double size;

  @override
  Widget build(BuildContext context) {
    return SvgPicture.asset(
      assetPath,
      width: size,
      height: size,
      colorFilter: ColorFilter.mode(color, BlendMode.srcIn),
    );
  }
}
