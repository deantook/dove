import 'dart:ui';

import 'package:flutter/material.dart';

import '../../../models/discover_user.dart';
import '../../../navigation/user_profile_navigation.dart';
import '../../../theme/app_colors.dart';
import '../../../widgets/moment_image.dart';
import '../../../widgets/user_avatar.dart';
import 'discover_photo_utils.dart';

class DiscoverPhotoOverlay extends StatefulWidget {
  const DiscoverPhotoOverlay({
    super.key,
    required this.user,
    required this.sourceRect,
    required this.resolveSourceRect,
    required this.onClose,
    required this.onDislike,
    required this.onLike,
  });

  final DiscoverUser user;
  final Rect sourceRect;
  final Rect? Function() resolveSourceRect;
  final VoidCallback onClose;
  final VoidCallback onDislike;
  final VoidCallback onLike;

  @override
  State<DiscoverPhotoOverlay> createState() => _DiscoverPhotoOverlayState();
}

class _DiscoverPhotoOverlayState extends State<DiscoverPhotoOverlay>
    with SingleTickerProviderStateMixin {
  static const _duration = Duration(milliseconds: 400);

  late final AnimationController _controller;
  late final Animation<double> _flight;
  late final Animation<double> _blur;
  late final Animation<double> _scrim;
  late final Animation<double> _chromeOpacity;
  late final Animation<Offset> _chromeSlide;
  late final PageController _pageController;
  late Rect _targetRect;

  int _photoIndex = 0;
  bool _isClosing = false;
  bool _galleryEnabled = false;
  Rect? _dismissTargetRect;

  DiscoverUser get user => widget.user;
  List<String> get photos => user.photos;
  List<Color> get _fallbackColors => discoverPhotoFallbackColors(user);

  @override
  void initState() {
    super.initState();
    _pageController = PageController();
    _targetRect = widget.sourceRect;
    _controller = AnimationController(vsync: this, duration: _duration);
    _flight = CurvedAnimation(
      parent: _controller,
      curve: Curves.easeInOutCubic,
      reverseCurve: Curves.easeInOutCubic,
    );
    _blur = CurvedAnimation(
      parent: _controller,
      curve: const Interval(0, 0.9, curve: Curves.easeOutCubic),
      reverseCurve: const Interval(0, 0.9, curve: Curves.easeIn),
    );
    _scrim = CurvedAnimation(
      parent: _controller,
      curve: const Interval(0, 0.9, curve: Curves.easeOut),
      reverseCurve: const Interval(0, 0.9, curve: Curves.easeIn),
    );
    _chromeOpacity = CurvedAnimation(
      parent: _controller,
      curve: const Interval(0.55, 1.0, curve: Curves.easeOutCubic),
      reverseCurve: const Interval(0.55, 1.0, curve: Curves.easeIn),
    );
    _chromeSlide = Tween<Offset>(
      begin: const Offset(0, 0.1),
      end: Offset.zero,
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: const Interval(0.55, 1.0, curve: Curves.easeOutCubic),
      reverseCurve: const Interval(0.55, 1.0, curve: Curves.easeIn),
    ));

    _controller.addStatusListener(_onAnimationStatusChanged);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      setState(() {
        _targetRect = computeDiscoverPhotoTargetRect(context);
      });
      _controller.forward();
    });
  }

  @override
  void dispose() {
    _controller.removeStatusListener(_onAnimationStatusChanged);
    _controller.dispose();
    _pageController.dispose();
    super.dispose();
  }

  void _onAnimationStatusChanged(AnimationStatus status) {
    if (status == AnimationStatus.completed && !_isClosing) {
      setState(() => _galleryEnabled = true);
    }
  }

  Rect _photoRectAt(double t) {
    final source = _isClosing
        ? (_dismissTargetRect ?? widget.sourceRect)
        : widget.sourceRect;
    return Rect.lerp(source, _targetRect, t)!;
  }

  Future<void> _dismiss({VoidCallback? then}) async {
    if (_isClosing) return;
    _isClosing = true;
    _dismissTargetRect = widget.resolveSourceRect() ?? widget.sourceRect;
    await _controller.reverse();
    if (!mounted) return;
    widget.onClose();
    then?.call();
  }

  void _openProfile() {
    UserProfileNavigation.open(
      context,
      name: user.name,
      avatarColors: user.avatarColors,
    );
  }

  Widget _buildPhotoContent({required bool enableGallery}) {
    if (photos.isEmpty) {
      return _PhotoPlaceholder(user: user);
    }

    if (!enableGallery) {
      final index = _photoIndex.clamp(0, photos.length - 1);
      return MomentImage(
        imageUrl: photos[index],
        borderRadius: 20,
        fallbackColors: _fallbackColors,
      );
    }

    return PageView.builder(
      controller: _pageController,
      itemCount: photos.length,
      onPageChanged: (index) => setState(() => _photoIndex = index),
      itemBuilder: (context, index) {
        return MomentImage(
          imageUrl: photos[index],
          borderRadius: 20,
          fallbackColors: _fallbackColors,
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final showGallery = _galleryEnabled && !_isClosing;

    return Material(
      color: Colors.transparent,
      child: AnimatedBuilder(
        animation: _controller,
        builder: (context, child) {
          final photoRect = _photoRectAt(_flight.value);
          final blurSigma = 18 * _blur.value;
          final scrimOpacity = 0.28 * _scrim.value;

          return Stack(
            fit: StackFit.expand,
            children: [
              GestureDetector(
                onTap: _isClosing ? null : () => _dismiss(),
                behavior: HitTestBehavior.opaque,
                child: BackdropFilter(
                  filter: ImageFilter.blur(
                    sigmaX: blurSigma,
                    sigmaY: blurSigma,
                  ),
                  child: Container(
                    color: Colors.black.withValues(alpha: scrimOpacity),
                  ),
                ),
              ),
              Positioned.fromRect(
                rect: photoRect,
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(20),
                  child: _buildPhotoContent(enableGallery: showGallery),
                ),
              ),
              if (photos.length > 1)
                Positioned(
                  top: photoRect.bottom + 12,
                  left: 0,
                  right: 0,
                  child: FadeTransition(
                    opacity: _chromeOpacity,
                    child: SlideTransition(
                      position: _chromeSlide,
                      child: _PhotoIndicator(
                        count: photos.length,
                        index: _photoIndex,
                      ),
                    ),
                  ),
                ),
              Positioned(
                left: 0,
                right: 0,
                bottom: 0,
                child: SafeArea(
                  child: Padding(
                    padding: const EdgeInsets.only(bottom: 24),
                    child: FadeTransition(
                      opacity: _chromeOpacity,
                      child: SlideTransition(
                        position: _chromeSlide,
                        child: _ActionRow(
                          onDislike: () => _dismiss(then: widget.onDislike),
                          onLike: () => _dismiss(then: widget.onLike),
                          onAvatarTap: _openProfile,
                          user: user,
                        ),
                      ),
                    ),
                  ),
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}

class _PhotoPlaceholder extends StatelessWidget {
  const _PhotoPlaceholder({required this.user});

  final DiscoverUser user;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: discoverPhotoFallbackColors(user),
        ),
      ),
      child: const Center(
        child: Icon(
          Icons.image_outlined,
          color: AppColors.textTertiary,
          size: 48,
        ),
      ),
    );
  }
}

class _PhotoIndicator extends StatelessWidget {
  const _PhotoIndicator({
    required this.count,
    required this.index,
  });

  final int count;
  final int index;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(count, (i) {
        final isActive = i == index;
        return AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          margin: const EdgeInsets.symmetric(horizontal: 3),
          width: isActive ? 18 : 6,
          height: 6,
          decoration: BoxDecoration(
            color: isActive
                ? Colors.white
                : Colors.white.withValues(alpha: 0.45),
            borderRadius: BorderRadius.circular(3),
          ),
        );
      }),
    );
  }
}

class _ActionRow extends StatelessWidget {
  const _ActionRow({
    required this.user,
    required this.onDislike,
    required this.onLike,
    required this.onAvatarTap,
  });

  final DiscoverUser user;
  final VoidCallback onDislike;
  final VoidCallback onLike;
  final VoidCallback onAvatarTap;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _CircleActionButton(
          key: const ValueKey('discover_dislike'),
          icon: Icons.close_rounded,
          iconColor: AppColors.textSecondary,
          onTap: onDislike,
        ),
        const SizedBox(width: 28),
        DecoratedBox(
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            border: Border.all(color: Colors.white, width: 3),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.12),
                blurRadius: 12,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: UserAvatar(
            key: const ValueKey('discover_overlay_avatar'),
            name: user.name,
            colors: user.avatarColors,
            size: 56,
            enableProfileNavigation: false,
            onTap: onAvatarTap,
          ),
        ),
        const SizedBox(width: 28),
        _CircleActionButton(
          key: const ValueKey('discover_like'),
          icon: Icons.favorite_rounded,
          iconColor: AppColors.primary,
          onTap: onLike,
        ),
      ],
    );
  }
}

class _CircleActionButton extends StatelessWidget {
  const _CircleActionButton({
    super.key,
    required this.icon,
    required this.iconColor,
    required this.onTap,
  });

  final IconData icon;
  final Color iconColor;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppColors.background,
      shape: const CircleBorder(),
      elevation: 4,
      shadowColor: Colors.black.withValues(alpha: 0.12),
      child: InkWell(
        onTap: onTap,
        customBorder: const CircleBorder(),
        child: SizedBox(
          width: 56,
          height: 56,
          child: Icon(icon, color: iconColor, size: 28),
        ),
      ),
    );
  }
}
