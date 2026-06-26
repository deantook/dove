import 'package:flutter/material.dart';

import '../../../models/moment.dart';
import '../../../theme/app_colors.dart';
import '../../../theme/app_typography.dart';
import '../../../widgets/moment_image.dart';
import '../moment_detail_screen.dart';
import 'moment_avatar.dart';

class MomentCard extends StatelessWidget {
  const MomentCard({
    super.key,
    required this.moment,
    this.onTap,
  });

  final Moment moment;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppColors.cardBackground,
      borderRadius: BorderRadius.circular(16),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(16),
            border: Border.all(color: AppColors.divider),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  MomentAvatar(
                    name: moment.authorName,
                    colors: moment.avatarColors,
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          moment.authorName,
                          style: AppTypography.cardHeadline,
                        ),
                        Text(moment.timeLabel, style: AppTypography.smallMeta),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              Text(moment.content, style: AppTypography.body),
              if (moment.resolvedImageUrl != null) ...[
                const SizedBox(height: 12),
                MomentImage(
                  imageUrl: moment.resolvedImageUrl!,
                  height: 180,
                  fallbackColors: [
                    moment.avatarColors.first.withValues(alpha: 0.6),
                    moment.avatarColors.last.withValues(alpha: 0.8),
                  ],
                ),
              ],
              const SizedBox(height: 12),
              Row(
                children: [
                  Icon(
                    Icons.favorite_border_rounded,
                    size: 18,
                    color: AppColors.textTertiary,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    '${moment.likeCount}',
                    style: AppTypography.smallMeta,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class MomentFeedList extends StatelessWidget {
  const MomentFeedList({
    super.key,
    required this.moments,
    required this.emptyMessage,
  });

  final List<Moment> moments;
  final String emptyMessage;

  @override
  Widget build(BuildContext context) {
    if (moments.isEmpty) {
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
      itemCount: moments.length,
      separatorBuilder: (_, _) => const SizedBox(height: 16),
      itemBuilder: (context, index) {
        final moment = moments[index];
        return MomentCard(
          key: ValueKey(moment.id),
          moment: moment,
          onTap: () {
            Navigator.of(context).push(MomentDetailScreen.route(moment));
          },
        );
      },
    );
  }
}
