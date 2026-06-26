import 'package:flutter/material.dart';

import 'live_match_controller.dart';

class LiveMatchScope extends InheritedWidget {
  const LiveMatchScope({
    super.key,
    required this.controller,
    required super.child,
  });

  final LiveMatchController controller;

  static LiveMatchController of(BuildContext context) {
    final scope =
        context.dependOnInheritedWidgetOfExactType<LiveMatchScope>();
    assert(scope != null, 'LiveMatchScope not found in widget tree');
    return scope!.controller;
  }

  @override
  bool updateShouldNotify(LiveMatchScope oldWidget) {
    return controller != oldWidget.controller;
  }
}
