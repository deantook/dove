import 'package:flutter/material.dart';

import '../theme/app_colors.dart';

typedef DoveNavIconBuilder = Widget Function(bool isActive, Color color);

class DoveNavItem {
  const DoveNavItem({required this.iconBuilder, this.semanticLabel});

  final DoveNavIconBuilder iconBuilder;
  final String? semanticLabel;
}

class DoveBottomNavBar extends StatelessWidget {
  const DoveBottomNavBar({
    super.key,
    required this.currentIndex,
    required this.items,
    required this.onTap,
  });

  static const height = 55.0;

  final int currentIndex;
  final List<DoveNavItem> items;
  final ValueChanged<int> onTap;

  @override
  Widget build(BuildContext context) {
    final bottomInset = MediaQuery.paddingOf(context).bottom / 2;

    return ColoredBox(
      color: AppColors.navBarBackground,
      child: Padding(
        padding: EdgeInsets.only(bottom: bottomInset),
        child: SizedBox(
          height: height,
          child: Row(
            children: List.generate(items.length, (index) {
              final item = items[index];
              final isActive = index == currentIndex;
              final color = isActive
                  ? AppColors.navBarActive
                  : AppColors.navBarInactive;

              return Expanded(
                child: Semantics(
                  label: item.semanticLabel,
                  button: true,
                  selected: isActive,
                  child: GestureDetector(
                    key: ValueKey('nav_tab_$index'),
                    onTap: () => onTap(index),
                    behavior: HitTestBehavior.opaque,
                    child: Center(child: item.iconBuilder(isActive, color)),
                  ),
                ),
              );
            }),
          ),
        ),
      ),
    );
  }
}
