import 'package:flutter/material.dart';

import '../live_match_controller.dart';
import '../live_match_scope.dart';
import 'live_match_floating_bubble.dart';
import 'live_match_sheet.dart';

class LiveMatchHost extends StatefulWidget {
  const LiveMatchHost({super.key, required this.child});

  final Widget child;

  @override
  State<LiveMatchHost> createState() => _LiveMatchHostState();
}

class _LiveMatchHostState extends State<LiveMatchHost> {
  late final LiveMatchController _controller = LiveMatchController();
  var _sheetShowing = false;

  @override
  void initState() {
    super.initState();
    _controller.addListener(_onControllerChanged);
  }

  @override
  void dispose() {
    _controller.removeListener(_onControllerChanged);
    _controller.dispose();
    super.dispose();
  }

  void _onControllerChanged() {
    if (_controller.isExpanded && !_sheetShowing) {
      _presentSheet();
    }

    if (_controller.takeTimeoutNotification()) {
      if (_sheetShowing && mounted) {
        Navigator.of(context).pop();
      }
      _showTimeoutNotice();
    }
  }

  Future<void> _presentSheet() async {
    if (_sheetShowing || !mounted) return;

    _sheetShowing = true;
    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      isDismissible: true,
      enableDrag: false,
      backgroundColor: Colors.transparent,
      builder: (_) => LiveMatchSheet(controller: _controller),
    );

    _sheetShowing = false;

    if (!mounted) return;
    if (_controller.isExpanded) {
      _controller.minimize();
    }
  }

  void _showTimeoutNotice() {
    if (!mounted) return;

    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        behavior: SnackBarBehavior.floating,
        content: Text('当前在线人数不足，请稍后再试'),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return LiveMatchScope(
      controller: _controller,
      child: ListenableBuilder(
        listenable: _controller,
        builder: (context, child) {
          return Stack(
            fit: StackFit.expand,
            children: [
              child!,
              if (_controller.isMinimized)
                LiveMatchFloatingBubble(controller: _controller),
            ],
          );
        },
        child: widget.child,
      ),
    );
  }
}
