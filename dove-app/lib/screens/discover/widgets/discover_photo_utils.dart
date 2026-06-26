import 'package:flutter/material.dart';

import '../../../models/discover_user.dart';
import '../../../widgets/moment_image.dart';

/// 读取 [key] 对应组件在屏幕上的位置和尺寸。
Rect? globalRectForKey(GlobalKey key) {
  final context = key.currentContext;
  if (context == null) return null;

  final box = context.findRenderObject();
  if (box is! RenderBox || !box.hasSize || !box.attached) return null;

  return box.localToGlobal(Offset.zero) & box.size;
}

/// 计算弹出层中 4:3 大图的目标位置（与操作栏布局对齐）。
Rect computeDiscoverPhotoTargetRect(BuildContext context) {
  final media = MediaQuery.of(context);
  const horizontalPadding = 20.0;
  const chromeBelow = 118.0;

  final width = media.size.width - horizontalPadding * 2;
  final height = width * 3 / 4;
  final safeTop = media.padding.top;
  final safeBottom = media.padding.bottom;
  final availableHeight = media.size.height - safeTop - safeBottom;
  final totalHeight = height + chromeBelow;
  final top = safeTop + (availableHeight - totalHeight) / 2;

  return Rect.fromLTWH(horizontalPadding, top, width, height);
}

/// 列表卡片与弹出层共用的照片容器。
class DiscoverPhotoFrame extends StatelessWidget {
  const DiscoverPhotoFrame({
    super.key,
    required this.imageUrl,
    required this.fallbackColors,
    this.aspectRatio = 4 / 3,
    this.borderRadius = 20,
    this.fit = BoxFit.cover,
  });

  final String imageUrl;
  final List<Color> fallbackColors;
  final double aspectRatio;
  final double borderRadius;
  final BoxFit fit;

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(borderRadius),
      child: MomentImage(
        imageUrl: imageUrl,
        aspectRatio: aspectRatio,
        borderRadius: borderRadius,
        fallbackColors: fallbackColors,
      ),
    );
  }
}

List<Color> discoverPhotoFallbackColors(DiscoverUser user) {
  return [
    user.avatarColors.first.withValues(alpha: 0.5),
    user.avatarColors.last.withValues(alpha: 0.8),
  ];
}

typedef DiscoverUserPhotoTap = void Function(DiscoverUser user, Rect sourceRect);
