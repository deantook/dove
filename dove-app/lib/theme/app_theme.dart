import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'app_colors.dart';
import 'app_typography.dart';

abstract final class AppTheme {
  static ThemeData get light {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      scaffoldBackgroundColor: AppColors.background,
      colorScheme: const ColorScheme.light(
        primary: AppColors.primary,
        onPrimary: Colors.white,
        surface: AppColors.background,
        onSurface: AppColors.textPrimary,
      ),
      dividerColor: AppColors.divider,
      appBarTheme: const AppBarTheme(
        elevation: 0,
        scrolledUnderElevation: 0,
        backgroundColor: AppColors.background,
        foregroundColor: AppColors.textPrimary,
        systemOverlayStyle: SystemUiOverlayStyle.dark,
      ),
      textTheme: const TextTheme(
        headlineLarge: AppTypography.pageTitle,
        headlineMedium: AppTypography.sectionTitle,
        titleMedium: AppTypography.cardHeadline,
        bodyLarge: AppTypography.body,
        bodyMedium: AppTypography.caption,
        labelSmall: AppTypography.smallMeta,
      ),
    );
  }
}
