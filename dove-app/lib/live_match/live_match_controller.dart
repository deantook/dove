import 'dart:async';

import 'package:flutter/material.dart';

import '../widgets/dove_bottom_nav_bar.dart';

enum LiveMatchPresentation { none, sheet, minimized }

abstract final class LiveMatchMetrics {
  static const floatingBubbleSize = 56.0;
  static const bubbleEdgeMargin = 16.0;
  static const bubbleDragMargin = 8.0;
  static const matchTimeout = Duration(seconds: 30);
}

class LiveMatchController extends ChangeNotifier {
  LiveMatchPresentation _presentation = LiveMatchPresentation.none;
  Offset? _bubbleOffset;
  DateTime? _matchStartedAt;
  Timer? _elapsedTicker;
  Timer? _timeoutTimer;
  var _pendingTimeoutNotice = false;

  LiveMatchPresentation get presentation => _presentation;

  bool get isActive => _presentation != LiveMatchPresentation.none;

  bool get isMinimized => _presentation == LiveMatchPresentation.minimized;

  bool get isExpanded => _presentation == LiveMatchPresentation.sheet;

  Offset? get bubbleOffset => _bubbleOffset;

  Duration get elapsed {
    final startedAt = _matchStartedAt;
    if (startedAt == null) return Duration.zero;
    return DateTime.now().difference(startedAt);
  }

  void start() {
    if (_presentation == LiveMatchPresentation.minimized) {
      expand();
      return;
    }
    if (_presentation == LiveMatchPresentation.none) {
      _matchStartedAt = DateTime.now();
      _pendingTimeoutNotice = false;
      _startTimers();
      _presentation = LiveMatchPresentation.sheet;
      notifyListeners();
    }
  }

  void minimize() {
    if (!isActive) return;
    _presentation = LiveMatchPresentation.minimized;
    notifyListeners();
  }

  void expand() {
    if (!isActive) return;
    _presentation = LiveMatchPresentation.sheet;
    notifyListeners();
  }

  void cancel() {
    _stopTimers();
    _presentation = LiveMatchPresentation.none;
    _matchStartedAt = null;
    notifyListeners();
  }

  bool takeTimeoutNotification() {
    if (!_pendingTimeoutNotice) return false;
    _pendingTimeoutNotice = false;
    return true;
  }

  void updateBubbleOffset(Offset offset) {
    _bubbleOffset = offset;
  }

  Offset defaultBubbleOffset(Size screenSize, EdgeInsets padding) {
    const bubbleSize = LiveMatchMetrics.floatingBubbleSize;
    const margin = LiveMatchMetrics.bubbleEdgeMargin;
    final bottomInset = padding.bottom / 2;
    final navBarHeight = DoveBottomNavBar.height + bottomInset;

    return Offset(
      screenSize.width - bubbleSize - margin,
      screenSize.height - navBarHeight - bubbleSize - margin,
    );
  }

  Offset clampBubbleOffset(Offset offset, Size screenSize, EdgeInsets padding) {
    const bubbleSize = LiveMatchMetrics.floatingBubbleSize;
    const margin = LiveMatchMetrics.bubbleDragMargin;
    final bottomInset = padding.bottom / 2;
    final navBarHeight = DoveBottomNavBar.height + bottomInset;

    return Offset(
      offset.dx.clamp(margin, screenSize.width - bubbleSize - margin),
      offset.dy.clamp(
        padding.top + margin,
        screenSize.height - navBarHeight - bubbleSize - margin,
      ),
    );
  }

  Offset snapBubbleOffset(Offset offset, Size screenSize, EdgeInsets padding) {
    const bubbleSize = LiveMatchMetrics.floatingBubbleSize;
    const margin = LiveMatchMetrics.bubbleEdgeMargin;
    final clamped = clampBubbleOffset(offset, screenSize, padding);
    final bubbleCenterX = clamped.dx + bubbleSize / 2;
    final snappedX = bubbleCenterX < screenSize.width / 2
        ? margin
        : screenSize.width - bubbleSize - margin;

    return Offset(snappedX, clamped.dy);
  }

  @override
  void dispose() {
    _stopTimers();
    super.dispose();
  }

  void _startTimers() {
    _stopTimers();
    _elapsedTicker = Timer.periodic(const Duration(seconds: 1), (_) {
      if (!isActive) return;
      notifyListeners();
    });
    _timeoutTimer = Timer(LiveMatchMetrics.matchTimeout, _handleTimeout);
  }

  void _handleTimeout() {
    if (!isActive) return;
    _pendingTimeoutNotice = true;
    cancel();
  }

  void _stopTimers() {
    _elapsedTicker?.cancel();
    _timeoutTimer?.cancel();
    _elapsedTicker = null;
    _timeoutTimer = null;
  }
}

String formatLiveMatchDuration(Duration duration) {
  final minutes = duration.inMinutes.remainder(60).toString().padLeft(2, '0');
  final seconds = duration.inSeconds.remainder(60).toString().padLeft(2, '0');
  return '$minutes:$seconds';
}
