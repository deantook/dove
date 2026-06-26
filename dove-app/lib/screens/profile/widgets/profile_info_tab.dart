import 'package:flutter/material.dart';

import '../../../models/user_profile.dart';
import '../../../theme/app_colors.dart';
import '../../../theme/app_typography.dart';
import 'profile_photo_wall.dart';

class ProfileInfoTab extends StatelessWidget {
  const ProfileInfoTab({
    super.key,
    required this.profileInfo,
    required this.photos,
    required this.avatarColors,
    this.scrollController,
  });

  final ProfileInfo profileInfo;
  final List<String> photos;
  final List<Color> avatarColors;
  final ScrollController? scrollController;

  @override
  Widget build(BuildContext context) {
    return ListView(
      controller: scrollController,
      padding: const EdgeInsets.fromLTRB(20, 8, 20, 24),
      children: [
        if (photos.isNotEmpty) ...[
          const ProfileSectionTitle(title: '照片墙'),
          const SizedBox(height: 12),
          ProfilePhotoWall(
            photos: photos,
            fallbackColors: avatarColors,
          ),
          const SizedBox(height: 24),
        ],
        const ProfileSectionTitle(title: '基本资料'),
        const SizedBox(height: 12),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: AppColors.cardBackground,
            borderRadius: BorderRadius.circular(16),
            border: Border.all(color: AppColors.divider),
          ),
          child: Column(
            children: profileInfo.details.map((detail) {
              final isLast = detail == profileInfo.details.last;
              return Column(
                children: [
                  _DetailRow(label: detail.$1, value: detail.$2),
                  if (!isLast)
                    const Padding(
                      padding: EdgeInsets.symmetric(vertical: 12),
                      child: Divider(height: 1, color: AppColors.divider),
                    ),
                ],
              );
            }).toList(),
          ),
        ),
        if (profileInfo.prompts.isNotEmpty) ...[
          const SizedBox(height: 24),
          const ProfileSectionTitle(title: '自定义文本'),
          const SizedBox(height: 12),
          ...profileInfo.prompts.map(
            (prompt) => _PromptCard(
              key: ValueKey(prompt.title),
              title: prompt.title,
              content: prompt.content,
            ),
          ),
        ],
      ],
    );
  }
}

class _PromptCard extends StatelessWidget {
  const _PromptCard({
    super.key,
    required this.title,
    required this.content,
  });

  final String title;
  final String content;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.cardBackground,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.divider),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: AppTypography.caption.copyWith(
              color: AppColors.textTertiary,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Text(content, style: AppTypography.body),
        ],
      ),
    );
  }
}

class _DetailRow extends StatelessWidget {
  const _DetailRow({
    required this.label,
    required this.value,
  });

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(
          label,
          style: AppTypography.body.copyWith(color: AppColors.textSecondary),
        ),
        const Spacer(),
        Text(value, style: AppTypography.body),
      ],
    );
  }
}
