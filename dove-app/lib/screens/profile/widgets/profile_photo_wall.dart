import 'package:flutter/material.dart';

import '../../../models/user_profile.dart';
import '../../../theme/app_typography.dart';
import '../../../widgets/moment_image.dart';

class ProfilePhotoWall extends StatelessWidget {
  const ProfilePhotoWall({
    super.key,
    required this.photos,
    required this.fallbackColors,
  });

  static const photoSize = 112.0;
  static const photoSpacing = 8.0;

  final List<String> photos;
  final List<Color> fallbackColors;

  @override
  Widget build(BuildContext context) {
    final displayPhotos = photos.take(kMaxProfilePhotos).toList();
    if (displayPhotos.isEmpty) return const SizedBox.shrink();

    return SizedBox(
      height: photoSize,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        itemCount: displayPhotos.length,
        separatorBuilder: (_, _) => const SizedBox(width: photoSpacing),
        itemBuilder: (context, index) {
          return SizedBox(
            width: photoSize,
            height: photoSize,
            child: MomentImage(
              key: ValueKey(displayPhotos[index]),
              imageUrl: displayPhotos[index],
              borderRadius: 12,
              fallbackColors: fallbackColors,
              aspectRatio: 1,
            ),
          );
        },
      ),
    );
  }
}

class ProfileSectionTitle extends StatelessWidget {
  const ProfileSectionTitle({super.key, required this.title});

  final String title;

  @override
  Widget build(BuildContext context) {
    return Text(title, style: AppTypography.caption);
  }
}
