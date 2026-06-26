import 'package:flutter/material.dart';

import '../../../models/like_profile.dart';
import '../../../theme/app_colors.dart';
import '../../../theme/app_typography.dart';
import '../../../widgets/user_avatar.dart';

class LikeProfileTile extends StatelessWidget {
  const LikeProfileTile({
    super.key,
    required this.profile,
    this.onTap,
  });

  final LikeProfile profile;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppColors.background,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 12),
          child: Row(
            children: [
              UserAvatar(
                name: profile.name,
                colors: profile.avatarColors,
                size: 56,
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(profile.name, style: AppTypography.cardHeadline),
                    const SizedBox(height: 4),
                    Text(
                      profile.subtitle,
                      style: AppTypography.caption.copyWith(
                        color: AppColors.textSecondary,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 8),
              Text(profile.timeLabel, style: AppTypography.smallMeta),
            ],
          ),
        ),
      ),
    );
  }
}

class LikeProfileList extends StatelessWidget {
  const LikeProfileList({
    super.key,
    required this.profiles,
    required this.emptyMessage,
  });

  final List<LikeProfile> profiles;
  final String emptyMessage;

  @override
  Widget build(BuildContext context) {
    if (profiles.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32),
          child: Text(
            emptyMessage,
            textAlign: TextAlign.center,
            style: AppTypography.body.copyWith(color: AppColors.textSecondary),
          ),
        ),
      );
    }

    return ListView.separated(
      padding: const EdgeInsets.fromLTRB(20, 8, 20, 24),
      itemCount: profiles.length,
      separatorBuilder: (_, _) => const Divider(
        height: 1,
        color: AppColors.divider,
      ),
      itemBuilder: (context, index) {
        final profile = profiles[index];
        return LikeProfileTile(
          key: ValueKey(profile.id),
          profile: profile,
          onTap: () {},
        );
      },
    );
  }
}
