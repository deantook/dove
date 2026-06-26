import 'package:flutter/material.dart';

/// Hinge 风格的品牌 H 图标。
class HingeHIcon extends StatelessWidget {
  const HingeHIcon({
    super.key,
    required this.color,
    this.size = 28,
    this.strokeWidth = 2.2,
    this.filled = false,
  });

  final Color color;
  final double size;
  final double strokeWidth;
  final bool filled;

  @override
  Widget build(BuildContext context) {
    return CustomPaint(
      size: Size(size, size),
      painter: _HingeHIconPainter(
        color: color,
        strokeWidth: strokeWidth,
        filled: filled,
      ),
    );
  }
}

class _HingeHIconPainter extends CustomPainter {
  _HingeHIconPainter({
    required this.color,
    required this.strokeWidth,
    required this.filled,
  });

  final Color color;
  final double strokeWidth;
  final bool filled;

  @override
  void paint(Canvas canvas, Size size) {
    final w = size.width;
    final h = size.height;

    final leftX = w * 0.28;
    final rightX = w * 0.72;
    final topY = h * 0.18;
    final midY = h * 0.52;
    final bottomY = h * 0.82;
    final barWidth = filled ? w * 0.14 : strokeWidth;

    if (filled) {
      final paint = Paint()
        ..color = color
        ..style = PaintingStyle.fill;

      canvas.drawRRect(
        RRect.fromRectAndRadius(
          Rect.fromLTWH(
            leftX - barWidth / 2,
            topY,
            barWidth,
            bottomY - topY,
          ),
          Radius.circular(barWidth / 2),
        ),
        paint,
      );
      canvas.drawRRect(
        RRect.fromRectAndRadius(
          Rect.fromLTWH(
            rightX - barWidth / 2,
            topY,
            barWidth,
            bottomY - topY,
          ),
          Radius.circular(barWidth / 2),
        ),
        paint,
      );
      canvas.drawRRect(
        RRect.fromRectAndRadius(
          Rect.fromLTWH(
            leftX - barWidth / 2,
            midY - barWidth / 2,
            rightX - leftX + barWidth,
            barWidth,
          ),
          Radius.circular(barWidth / 2),
        ),
        paint,
      );
      return;
    }

    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round
      ..strokeJoin = StrokeJoin.round;

    canvas.drawLine(Offset(leftX, topY), Offset(leftX, bottomY), paint);
    canvas.drawLine(Offset(rightX, topY), Offset(rightX, bottomY), paint);
    canvas.drawLine(Offset(leftX, midY), Offset(rightX, midY), paint);
  }

  @override
  bool shouldRepaint(covariant _HingeHIconPainter oldDelegate) {
    return oldDelegate.color != color ||
        oldDelegate.strokeWidth != strokeWidth ||
        oldDelegate.filled != filled;
  }
}
