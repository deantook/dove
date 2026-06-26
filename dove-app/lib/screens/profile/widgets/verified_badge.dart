import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../../../widgets/nav_icons/dove_nav_icon.dart';

class VerifiedBadge extends StatelessWidget {
  const VerifiedBadge({
    super.key,
    this.size = 28,
  });

  final double size;

  @override
  Widget build(BuildContext context) {
    return SvgPicture.asset(
      OtherIconAssets.verified,
      width: size,
      height: size,
    );
  }
}
