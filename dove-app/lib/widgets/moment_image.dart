import 'package:flutter/material.dart';

import '../theme/app_colors.dart';

class MomentImage extends StatelessWidget {
  const MomentImage({
    super.key,
    required this.imageUrl,
    this.height,
    this.aspectRatio,
    this.borderRadius = 16,
    this.fallbackColors = const [AppColors.backgroundSecondary, AppColors.divider],
  });

  final String imageUrl;
  final double? height;
  final double? aspectRatio;
  final double borderRadius;
  final List<Color> fallbackColors;

  @override
  Widget build(BuildContext context) {
    final image = ClipRRect(
      borderRadius: BorderRadius.circular(borderRadius),
      child: Image.network(
        imageUrl,
        width: double.infinity,
        height: aspectRatio == null ? height : null,
        fit: BoxFit.cover,
        loadingBuilder: (context, child, progress) {
          if (progress == null) return child;
          return _placeholder();
        },
        errorBuilder: (_, _, _) => _placeholder(),
      ),
    );

    if (aspectRatio != null) {
      return AspectRatio(aspectRatio: aspectRatio!, child: image);
    }

    return SizedBox(width: double.infinity, height: height, child: image);
  }

  Widget _placeholder() {
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(borderRadius),
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: fallbackColors,
        ),
      ),
      child: const Center(
        child: Icon(
          Icons.image_outlined,
          color: AppColors.textTertiary,
          size: 32,
        ),
      ),
    );
  }
}
