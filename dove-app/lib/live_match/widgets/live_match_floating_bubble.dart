import 'package:flutter/material.dart';

import '../../theme/app_colors.dart';
import '../live_match_controller.dart';
import 'live_match_pulse_indicator.dart';

class LiveMatchFloatingBubble extends StatefulWidget {
  const LiveMatchFloatingBubble({super.key, required this.controller});

  static const size = LiveMatchMetrics.floatingBubbleSize;

  final LiveMatchController controller;

  @override
  State<LiveMatchFloatingBubble> createState() =>
      _LiveMatchFloatingBubbleState();
}

class _LiveMatchFloatingBubbleState extends State<LiveMatchFloatingBubble>
    with SingleTickerProviderStateMixin {
  static const _snapDuration = Duration(milliseconds: 220);
  static const _tapSlop = 8.0;

  Offset? _position;
  var _didDrag = false;
  var _dragDistance = 0.0;
  AnimationController? _snapController;
  Animation<Offset>? _snapAnimation;

  LiveMatchController get _controller => widget.controller;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _syncInitialPosition());
  }

  @override
  void dispose() {
    _snapController?.dispose();
    super.dispose();
  }

  void _syncInitialPosition() {
    if (!mounted) return;
    final media = MediaQuery.of(context);
    final initial = _controller.bubbleOffset ??
        _controller.defaultBubbleOffset(media.size, media.padding);
    setState(() {
      _position = _controller.snapBubbleOffset(
        initial,
        media.size,
        media.padding,
      );
    });
  }

  void _onPanStart(DragStartDetails details) {
    _didDrag = false;
    _dragDistance = 0;
    _snapController?.stop();
    _snapController?.dispose();
    _snapController = null;
    _snapAnimation = null;
  }

  void _onPanUpdate(DragUpdateDetails details) {
    _dragDistance += details.delta.distance;
    if (_dragDistance > _tapSlop) _didDrag = true;
    if (!_didDrag) return;

    final media = MediaQuery.of(context);
    final current = _position;
    if (current == null) return;
    setState(() {
      _position = _controller.clampBubbleOffset(
        current + details.delta,
        media.size,
        media.padding,
      );
    });
  }

  void _openSheet() => _controller.expand();

  void _onPanEnd(DragEndDetails details) {
    if (!_didDrag) {
      _openSheet();
      return;
    }

    final current = _position;
    if (current == null) return;

    final media = MediaQuery.of(context);
    final snapped = _controller.snapBubbleOffset(
      current,
      media.size,
      media.padding,
    );

    _didDrag = false;
    _dragDistance = 0;

    if (snapped == current) {
      _controller.updateBubbleOffset(snapped);
      return;
    }

    _snapController?.dispose();
    _snapController = AnimationController(vsync: this, duration: _snapDuration);
    _snapAnimation = Tween<Offset>(begin: current, end: snapped).animate(
      CurvedAnimation(parent: _snapController!, curve: Curves.easeOutCubic),
    )..addListener(() {
        if (!mounted) return;
        setState(() => _position = _snapAnimation!.value);
      });

    _snapController!.forward().whenComplete(() {
      if (!mounted) return;
      _controller.updateBubbleOffset(snapped);
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_position == null) {
      return const SizedBox.shrink();
    }

    return Positioned(
      key: const ValueKey('live_match_floating_bubble'),
      left: _position!.dx,
      top: _position!.dy,
      child: GestureDetector(
        behavior: HitTestBehavior.opaque,
        onPanStart: _onPanStart,
        onPanUpdate: _onPanUpdate,
        onPanEnd: _onPanEnd,
        child: Material(
          elevation: 8,
          shadowColor: AppColors.primary.withValues(alpha: 0.25),
          color: AppColors.cardBackground,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(LiveMatchFloatingBubble.size / 2),
            side: BorderSide(
              color: AppColors.primary.withValues(alpha: 0.2),
            ),
          ),
          clipBehavior: Clip.antiAlias,
          child: SizedBox(
            width: LiveMatchFloatingBubble.size,
            height: LiveMatchFloatingBubble.size,
            child: const LiveMatchPulseIndicator(
              size: 56,
              coreSize: 40,
              iconSize: 22,
            ),
          ),
        ),
      ),
    );
  }
}
