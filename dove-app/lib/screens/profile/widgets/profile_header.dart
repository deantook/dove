import 'package:flutter/material.dart';

import '../../../theme/app_colors.dart';
import '../../../theme/app_typography.dart';
import '../../../widgets/nav_icons/dove_nav_icon.dart';
import 'verified_badge.dart';

class ProfileHeader extends StatelessWidget {
  const ProfileHeader({
    super.key,
    this.height = 48,
    this.iconSize = 26,
    this.onSettingsTap,
    this.onBackTap,
    this.showSettings = true,
    this.showBackButton = false,
    this.compactLeading = false,
  });

  final double height;
  final double iconSize;
  final VoidCallback? onSettingsTap;
  final VoidCallback? onBackTap;
  final bool showSettings;
  final bool showBackButton;
  final bool compactLeading;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: height,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(4, 0, 12, 0),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            if (showBackButton)
              _HeaderIconButton(
                onPressed: onBackTap,
                icon: Icon(
                  Icons.arrow_back_ios_new_rounded,
                  size: iconSize * 0.78,
                  color: AppColors.textPrimary,
                ),
              )
            else if (!compactLeading)
              SizedBox(width: height)
            else
              const SizedBox.shrink(),
            const Spacer(),
            if (showSettings)
              _HeaderIconButton(
                key: const ValueKey('profile_settings'),
                onPressed: onSettingsTap,
                icon: DoveNavIcon(
                  assetPath: OtherIconAssets.settings,
                  color: AppColors.textPrimary,
                  size: iconSize,
                ),
              )
            else
              SizedBox(width: height),
          ],
        ),
      ),
    );
  }
}

class _HeaderIconButton extends StatelessWidget {
  const _HeaderIconButton({
    super.key,
    required this.onPressed,
    required this.icon,
  });

  final VoidCallback? onPressed;
  final Widget icon;

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: onPressed,
      padding: EdgeInsets.zero,
      constraints: const BoxConstraints.tightFor(width: 40, height: 40),
      splashRadius: 20,
      icon: icon,
    );
  }
}

class ProfileNameRow extends StatelessWidget {
  const ProfileNameRow({
    super.key,
    required this.name,
    this.isVerified = true,
  });

  final String name;
  final bool isVerified;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Text(
          name,
          style: AppTypography.sectionTitle.copyWith(
            fontSize: 28,
            fontWeight: FontWeight.w700,
            height: 34 / 28,
          ),
        ),
        if (isVerified) ...[
          const SizedBox(width: 8),
          const VerifiedBadge(),
        ],
      ],
    );
  }
}
