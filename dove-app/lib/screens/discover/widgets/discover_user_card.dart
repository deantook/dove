import 'package:flutter/material.dart';

import '../../../models/discover_user.dart';
import '../../../theme/app_typography.dart';
import 'discover_photo_utils.dart';

class DiscoverUserCard extends StatelessWidget {
  const DiscoverUserCard({
    super.key,
    required this.user,
    required this.photoKey,
    required this.onPhotoTap,
    this.hidePhoto = false,
  });

  final DiscoverUser user;
  final GlobalKey photoKey;
  final DiscoverUserPhotoTap onPhotoTap;
  final bool hidePhoto;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: () {
          final rect = globalRectForKey(photoKey);
          if (rect != null) onPhotoTap(user, rect);
        },
        borderRadius: BorderRadius.circular(20),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(20),
          child: Stack(
            children: [
              KeyedSubtree(
                key: photoKey,
                child: Opacity(
                  opacity: hidePhoto ? 0 : 1,
                  child: DiscoverPhotoFrame(
                    imageUrl: user.resolvedPhotoUrl,
                    fallbackColors: discoverPhotoFallbackColors(user),
                  ),
                ),
              ),
              Positioned.fill(
                child: DecoratedBox(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        Colors.transparent,
                        Colors.black.withValues(alpha: 0.02),
                        Colors.black.withValues(alpha: 0.55),
                      ],
                      stops: const [0.45, 0.7, 1.0],
                    ),
                  ),
                ),
              ),
              if (user.isOnline)
                Positioned(
                  top: 10,
                  right: 10,
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: Colors.black.withValues(alpha: 0.35),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Container(
                          width: 6,
                          height: 6,
                          decoration: const BoxDecoration(
                            color: Color(0xFF34D399),
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 4),
                        Text(
                          '在线',
                          style: AppTypography.smallMeta.copyWith(
                            color: Colors.white,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              Positioned(
                left: 12,
                right: 12,
                bottom: 12,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      '${user.name}, ${user.age}',
                      style: AppTypography.cardHeadline.copyWith(
                        color: Colors.white,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      user.city,
                      style: AppTypography.caption.copyWith(
                        color: Colors.white.withValues(alpha: 0.85),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
