import 'package:flutter/material.dart';

import '../data/mock_image_urls.dart';
import '../navigation/user_profile_navigation.dart';
import '../theme/app_typography.dart';

class UserAvatar extends StatelessWidget {
  const UserAvatar({
    super.key,
    required this.name,
    required this.colors,
    this.size = 40,
    this.imageUrl,
    this.onTap,
    this.enableProfileNavigation = true,
  });

  final String name;
  final List<Color> colors;
  final double size;
  final String? imageUrl;
  final VoidCallback? onTap;
  final bool enableProfileNavigation;

  String? get _resolvedImageUrl => imageUrl ?? MockImageUrls.avatarFor(name);

  @override
  Widget build(BuildContext context) {
    final avatar = _resolvedImageUrl != null
        ? _NetworkAvatar(
            url: _resolvedImageUrl!,
            size: size,
            fallback: _gradientAvatar(),
          )
        : _gradientAvatar();

    final handler = onTap ??
        (enableProfileNavigation
            ? () => UserProfileNavigation.open(
                  context,
                  name: name,
                  avatarColors: colors,
                )
            : null);

    if (handler == null) return avatar;

    return GestureDetector(
      onTap: handler,
      behavior: HitTestBehavior.opaque,
      child: avatar,
    );
  }

  Widget _gradientAvatar() {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: colors,
        ),
      ),
      alignment: Alignment.center,
      child: Text(
        name.characters.first.toUpperCase(),
        style: AppTypography.caption.copyWith(
          color: Colors.white,
          fontWeight: FontWeight.w600,
          fontSize: size * 0.38,
        ),
      ),
    );
  }
}

class _NetworkAvatar extends StatelessWidget {
  const _NetworkAvatar({
    required this.url,
    required this.size,
    required this.fallback,
  });

  final String url;
  final double size;
  final Widget fallback;

  @override
  Widget build(BuildContext context) {
    return ClipOval(
      child: SizedBox(
        width: size,
        height: size,
        child: Image.network(
          url,
          width: size,
          height: size,
          fit: BoxFit.cover,
          loadingBuilder: (context, child, progress) {
            if (progress == null) return child;
            return fallback;
          },
          errorBuilder: (_, _, _) => fallback,
        ),
      ),
    );
  }
}
