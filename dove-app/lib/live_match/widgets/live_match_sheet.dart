import 'package:flutter/material.dart';

import '../../theme/app_colors.dart';
import '../../theme/app_typography.dart';
import '../live_match_controller.dart';
import 'live_match_pulse_indicator.dart';

class LiveMatchSheet extends StatelessWidget {
  const LiveMatchSheet({super.key, required this.controller});

  final LiveMatchController controller;

  @override
  Widget build(BuildContext context) {
    final bottomInset = MediaQuery.paddingOf(context).bottom;

    return Container(
      decoration: const BoxDecoration(
        color: AppColors.cardBackground,
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        boxShadow: [
          BoxShadow(
            color: Color(0x14000000),
            blurRadius: 20,
            offset: Offset(0, -4),
          ),
        ],
      ),
      child: Padding(
        padding: EdgeInsets.fromLTRB(20, 12, 20, 24 + bottomInset),
        child: ListenableBuilder(
          listenable: controller,
          builder: (context, _) {
            final elapsed = formatLiveMatchDuration(controller.elapsed);

            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Row(
                  children: [
                    const Spacer(),
                    IconButton(
                      key: const ValueKey('live_match_minimize'),
                      onPressed: () {
                        controller.minimize();
                        Navigator.of(context).pop();
                      },
                      icon: const Icon(
                        Icons.keyboard_arrow_down_rounded,
                        color: AppColors.textSecondary,
                      ),
                      tooltip: '最小化',
                    ),
                  ],
                ),
                const LiveMatchPulseIndicator(),
                const SizedBox(height: 28),
                Text('正在匹配', style: AppTypography.sectionTitle),
                const SizedBox(height: 8),
                Text(
                  '寻找此刻在线的用户',
                  style: AppTypography.body.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
                const SizedBox(height: 12),
                Text(
                  elapsed,
                  key: const ValueKey('live_match_timer'),
                  style: AppTypography.caption.copyWith(
                    color: AppColors.primary,
                    fontWeight: FontWeight.w600,
                    fontFeatures: const [FontFeature.tabularFigures()],
                  ),
                ),
                const SizedBox(height: 36),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: OutlinedButton(
                    key: const ValueKey('live_match_cancel'),
                    onPressed: () {
                      controller.cancel();
                      Navigator.of(context).pop();
                    },
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppColors.textPrimary,
                      side: const BorderSide(color: AppColors.inputBorder),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(24),
                      ),
                      textStyle: AppTypography.caption.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    child: const Text('取消匹配'),
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }
}
