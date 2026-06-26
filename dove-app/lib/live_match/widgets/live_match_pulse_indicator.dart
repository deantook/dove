import 'package:flutter/material.dart';

import '../../theme/app_colors.dart';

class LiveMatchPulseIndicator extends StatefulWidget {
  const LiveMatchPulseIndicator({
    super.key,
    this.size = 120,
    this.coreSize = 64,
    this.iconSize = 36,
  });

  final double size;
  final double coreSize;
  final double iconSize;

  @override
  State<LiveMatchPulseIndicator> createState() =>
      _LiveMatchPulseIndicatorState();
}

class _LiveMatchPulseIndicatorState extends State<LiveMatchPulseIndicator>
    with SingleTickerProviderStateMixin {
  late final AnimationController _pulseController;

  @override
  void initState() {
    super.initState();
    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1800),
    )..repeat();
  }

  @override
  void dispose() {
    _pulseController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: widget.size,
      height: widget.size,
      child: AnimatedBuilder(
        animation: _pulseController,
        builder: (context, child) {
          return Stack(
            alignment: Alignment.center,
            children: [
              for (var i = 0; i < 3; i++)
                _PulseRing(
                  size: widget.size,
                  progress: (_pulseController.value + i / 3) % 1.0,
                ),
              child!,
            ],
          );
        },
        child: Container(
          width: widget.coreSize,
          height: widget.coreSize,
          decoration: const BoxDecoration(
            color: AppColors.primaryLight,
            shape: BoxShape.circle,
          ),
          child: Icon(
            Icons.bolt_rounded,
            color: AppColors.primary,
            size: widget.iconSize,
          ),
        ),
      ),
    );
  }
}

class _PulseRing extends StatelessWidget {
  const _PulseRing({required this.size, required this.progress});

  final double size;
  final double progress;

  @override
  Widget build(BuildContext context) {
    final scale = 0.55 + progress * 0.85;
    final opacity = (1 - progress).clamp(0.0, 1.0) * 0.35;

    return Transform.scale(
      scale: scale,
      child: Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          border: Border.all(
            color: AppColors.primary.withValues(alpha: opacity),
            width: 2,
          ),
        ),
      ),
    );
  }
}
