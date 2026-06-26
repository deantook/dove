import 'package:flutter/material.dart';

import '../theme/app_colors.dart';
import '../theme/app_typography.dart';

class BaseTabScreen extends StatelessWidget {
  const BaseTabScreen({
    super.key,
    required this.subtitle,
  });

  final String subtitle;

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 16),
            Text(
              subtitle,
              style: AppTypography.body.copyWith(color: AppColors.textSecondary),
            ),
            const Spacer(),
          ],
        ),
      ),
    );
  }
}
