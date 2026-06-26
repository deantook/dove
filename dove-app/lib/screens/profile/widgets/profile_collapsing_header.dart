import 'dart:ui' show lerpDouble;

import 'package:flutter/material.dart';

import 'profile_avatar.dart';
import 'profile_header.dart';
import 'profile_user_tags.dart';

class ProfileCollapsingHeader extends StatelessWidget {
  const ProfileCollapsingHeader({
    super.key,
    required this.collapseProgress,
    required this.userName,
    required this.avatarColors,
    required this.tags,
    this.isVerified = true,
    this.showSettings = true,
    this.showBackButton = false,
    this.onSettingsTap,
    this.onBackTap,
  });

  static const expandedAvatarSize = 112.0;
  static const collapsedAvatarSize = 40.0;
  static const horizontalPadding = 20.0;
  static const collapsedLeadingInset = 8.0;
  static const settingsHeight = 48.0;
  static const settingsIconSize = 26.0;
  static const nameSpacing = 12.0;
  static const nameHeight = 34.0;

  final double collapseProgress;
  final String userName;
  final List<Color> avatarColors;
  final List<String> tags;
  final bool isVerified;
  final bool showSettings;
  final bool showBackButton;
  final VoidCallback? onSettingsTap;
  final VoidCallback? onBackTap;

  double get _t => collapseProgress.clamp(0.0, 1.0);

  static double heightForProgress(double progress) {
    final t = progress.clamp(0.0, 1.0);
    const expanded = settingsHeight +
        24 +
        expandedAvatarSize +
        nameSpacing +
        nameHeight +
        24;
    const collapsed = settingsHeight + 8;
    return lerpDouble(expanded, collapsed, t)!;
  }

  @override
  Widget build(BuildContext context) {
    final t = Curves.easeOutCubic.transform(_t);
    final height = heightForProgress(_t);

    return SizedBox(
      height: height,
      width: double.infinity,
      child: LayoutBuilder(
        builder: (context, constraints) {
          final avatarSize =
              lerpDouble(expandedAvatarSize, collapsedAvatarSize, t)!;

          final collapsedOffset = lerpDouble(0, collapsedLeadingInset, t)!;
          final collapsedAvatarLeft =
              (showBackButton ? 48.0 : horizontalPadding) + collapsedOffset;
          final avatarLeft =
              lerpDouble(horizontalPadding, collapsedAvatarLeft, t)!;
          final avatarTopCollapsed = (settingsHeight - avatarSize) / 2;
          final avatarTop =
              lerpDouble(settingsHeight + 24, avatarTopCollapsed, t)!;

          final nameTop = avatarTop + avatarSize + nameSpacing;
          final chromeOpacity = (1 - t * 1.4).clamp(0.0, 1.0);
          final compactLeading = t > 0.5;
          final tagsLeft = avatarLeft + avatarSize + 16;

          return Stack(
            clipBehavior: Clip.none,
            children: [
              ProfileHeader(
                height: settingsHeight,
                iconSize: settingsIconSize,
                onSettingsTap: onSettingsTap,
                onBackTap: onBackTap,
                showSettings: showSettings,
                showBackButton: showBackButton,
                compactLeading: compactLeading,
              ),
              Positioned(
                left: avatarLeft,
                top: avatarTop,
                child: ProfileAvatar(
                  size: avatarSize,
                  userName: userName,
                  avatarColors: avatarColors,
                ),
              ),
              Positioned(
                left: horizontalPadding,
                top: nameTop,
                child: Opacity(
                  opacity: chromeOpacity,
                  child: ProfileNameRow(
                    name: userName,
                    isVerified: isVerified,
                  ),
                ),
              ),
              if (tags.isNotEmpty)
                Positioned(
                  left: tagsLeft,
                  right: horizontalPadding,
                  top: avatarTop,
                  child: Opacity(
                    opacity: chromeOpacity,
                    child: ProfileUserTags(tags: tags),
                  ),
                ),
            ],
          );
        },
      ),
    );
  }
}
