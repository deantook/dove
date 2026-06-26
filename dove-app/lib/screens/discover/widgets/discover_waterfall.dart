import 'package:flutter/material.dart';

import '../../../models/discover_user.dart';
import '../../../theme/app_colors.dart';
import '../../../theme/app_typography.dart';
import 'discover_photo_utils.dart';
import 'discover_user_card_tile.dart';

class DiscoverWaterfall extends StatelessWidget {
  const DiscoverWaterfall({
    super.key,
    required this.users,
    required this.onUserPhotoTap,
    required this.photoKeyFor,
    required this.onUserRemoved,
    this.previewUserId,
    this.removingUserId,
    this.itemSpacing = 16,
  });

  final List<DiscoverUser> users;
  final DiscoverUserPhotoTap onUserPhotoTap;
  final GlobalKey Function(String userId) photoKeyFor;
  final ValueChanged<String> onUserRemoved;
  final String? previewUserId;
  final String? removingUserId;
  final double itemSpacing;

  @override
  Widget build(BuildContext context) {
    if (users.isEmpty) {
      return SliverFillRemaining(
        hasScrollBody: false,
        child: Center(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 32),
            child: Text(
              '暂无符合条件的推荐',
              textAlign: TextAlign.center,
              style: AppTypography.body.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ),
        ),
      );
    }

    return SliverPadding(
      padding: EdgeInsets.only(bottom: 24),
      sliver: SliverList.separated(
        itemCount: users.length,
        separatorBuilder: (_, _) => SizedBox(height: itemSpacing),
        itemBuilder: (context, index) {
          final user = users[index];
          return DiscoverUserCardTile(
            key: ValueKey(user.id),
            user: user,
            photoKey: photoKeyFor(user.id),
            hidePhoto: previewUserId == user.id,
            isRemoving: removingUserId == user.id,
            onPhotoTap: onUserPhotoTap,
            onRemoveComplete: () => onUserRemoved(user.id),
          );
        },
      ),
    );
  }
}
