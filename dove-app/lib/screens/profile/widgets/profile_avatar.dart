import 'package:flutter/material.dart';

import '../../../data/mock_image_urls.dart';

/// 个人页大头像：紫色描边。
class ProfileAvatar extends StatelessWidget {
  const ProfileAvatar({
    super.key,
    this.size = 112,
    this.userName,
    this.avatarUrl,
    this.avatarColors = const [
      Color(0xFFD4C4B0),
      Color(0xFF8B7355),
      Color(0xFF5C7A6B),
    ],
  });

  final double size;
  final String? userName;
  final String? avatarUrl;
  final List<Color> avatarColors;

  static const _ringColor = Color(0xFF5E2652);
  static const _ringWidth = 3.0;
  static const _photoInset = 4.0;

  String? get _resolvedAvatarUrl =>
      avatarUrl ?? (userName != null ? MockImageUrls.avatarFor(userName!) : null);

  @override
  Widget build(BuildContext context) {
    final photoSize = size - (_ringWidth + _photoInset) * 2;

    return SizedBox(
      width: size,
      height: size,
      child: Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          border: Border.all(color: _ringColor, width: _ringWidth),
        ),
        child: Center(
          child: ClipOval(
            child: SizedBox(
              width: photoSize,
              height: photoSize,
              child: _buildPhoto(photoSize),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildPhoto(double photoSize) {
    final url = _resolvedAvatarUrl;
    if (url != null) {
      return Image.network(
        url,
        width: photoSize,
        height: photoSize,
        fit: BoxFit.cover,
        loadingBuilder: (context, child, progress) {
          if (progress == null) return child;
          return _placeholderPhoto(photoSize);
        },
        errorBuilder: (_, _, _) => _placeholderPhoto(photoSize),
      );
    }
    return _placeholderPhoto(photoSize);
  }

  Widget _placeholderPhoto(double photoSize) {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: avatarColors.length >= 2
              ? avatarColors
              : [avatarColors.first, avatarColors.first],
        ),
      ),
      child: Icon(
        Icons.person_outline_rounded,
        size: photoSize * 0.45,
        color: Colors.white.withValues(alpha: 0.9),
      ),
    );
  }
}
