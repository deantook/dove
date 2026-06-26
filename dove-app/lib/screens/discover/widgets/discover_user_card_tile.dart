import 'package:flutter/material.dart';

import '../../../models/discover_user.dart';
import 'discover_photo_utils.dart';
import 'discover_user_card.dart';

class DiscoverUserCardTile extends StatefulWidget {
  const DiscoverUserCardTile({
    super.key,
    required this.user,
    required this.photoKey,
    required this.onPhotoTap,
    required this.onRemoveComplete,
    this.hidePhoto = false,
    this.isRemoving = false,
  });

  final DiscoverUser user;
  final GlobalKey photoKey;
  final DiscoverUserPhotoTap onPhotoTap;
  final VoidCallback onRemoveComplete;
  final bool hidePhoto;
  final bool isRemoving;

  @override
  State<DiscoverUserCardTile> createState() => _DiscoverUserCardTileState();
}

class _DiscoverUserCardTileState extends State<DiscoverUserCardTile>
    with SingleTickerProviderStateMixin {
  static const _duration = Duration(milliseconds: 320);

  late final AnimationController _controller;
  late final Animation<double> _fade;
  late final Animation<double> _heightFactor;
  late final Animation<Offset> _slide;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(vsync: this, duration: _duration);
    _fade = CurvedAnimation(
      parent: _controller,
      curve: Curves.easeInCubic,
    );
    _heightFactor = Tween<double>(begin: 1, end: 0).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.easeInOutCubic,
    ));
    _slide = Tween<Offset>(
      begin: Offset.zero,
      end: const Offset(0, -0.06),
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.easeInCubic,
    ));

    if (widget.isRemoving) {
      _startRemoving();
    }
  }

  @override
  void didUpdateWidget(DiscoverUserCardTile oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.isRemoving && !oldWidget.isRemoving) {
      _startRemoving();
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _startRemoving() async {
    await _controller.forward();
    if (mounted) widget.onRemoveComplete();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        final heightFactor =
            widget.isRemoving ? _heightFactor.value.clamp(0.0, 1.0) : 1.0;
        final opacity = widget.isRemoving ? 1 - _fade.value : 1.0;
        final slide = widget.isRemoving ? _slide.value : Offset.zero;

        return ClipRect(
          child: Align(
            alignment: Alignment.topCenter,
            heightFactor: heightFactor,
            child: Opacity(
              opacity: opacity,
              child: Transform.translate(
                offset: Offset(0, slide.dy * 48),
                child: child,
              ),
            ),
          ),
        );
      },
      child: DiscoverUserCard(
        user: widget.user,
        photoKey: widget.photoKey,
        hidePhoto: widget.hidePhoto,
        onPhotoTap: widget.onPhotoTap,
      ),
    );
  }
}
