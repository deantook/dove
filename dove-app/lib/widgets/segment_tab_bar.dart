import 'package:flutter/material.dart';

import '../theme/app_colors.dart';
import '../theme/app_typography.dart';

class SegmentTabBar extends StatelessWidget {
  const SegmentTabBar({
    super.key,
    required this.tabs,
    required this.selectedIndex,
    required this.onChanged,
    this.tabSpacing = 24,
    this.showIndicator = true,
  });

  final List<String> tabs;
  final int selectedIndex;
  final ValueChanged<int> onChanged;
  final double tabSpacing;
  final bool showIndicator;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: List.generate(tabs.length, (index) {
        final isSelected = index == selectedIndex;

        return Padding(
          padding: EdgeInsets.only(right: index < tabs.length - 1 ? tabSpacing : 0),
          child: GestureDetector(
            onTap: () => onChanged(index),
            behavior: HitTestBehavior.opaque,
            child: showIndicator
                ? Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      _TabLabel(tabs[index], isSelected: isSelected),
                      const SizedBox(height: 8),
                      AnimatedContainer(
                        duration: const Duration(milliseconds: 200),
                        height: 2,
                        width: isSelected ? 24 : 0,
                        decoration: BoxDecoration(
                          color: AppColors.primary,
                          borderRadius: BorderRadius.circular(1),
                        ),
                      ),
                    ],
                  )
                : _TabLabel(tabs[index], isSelected: isSelected),
          ),
        );
      }),
    );
  }
}

class _TabLabel extends StatelessWidget {
  const _TabLabel(this.label, {required this.isSelected});

  final String label;
  final bool isSelected;

  @override
  Widget build(BuildContext context) {
    return Text(
      label,
      style: AppTypography.cardHeadline.copyWith(
        color: isSelected ? AppColors.textPrimary : AppColors.textTertiary,
        fontWeight: isSelected ? FontWeight.w700 : FontWeight.w500,
      ),
    );
  }
}
