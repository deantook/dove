import 'package:flutter/material.dart';

import '../models/moment.dart';
import '../widgets/segment_tab_bar.dart';
import 'standouts/widgets/moment_feed.dart';

class StandoutsScreen extends StatefulWidget {
  const StandoutsScreen({super.key});

  @override
  State<StandoutsScreen> createState() => _StandoutsScreenState();
}

class _StandoutsScreenState extends State<StandoutsScreen> {
  int _selectedTab = 0;

  List<Moment> _momentsFor(MomentSource source) {
    return mockMoments
        .where((moment) => moment.source == source)
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    final strangerMoments = _momentsFor(MomentSource.stranger);
    final friendMoments = _momentsFor(MomentSource.friend);

    return SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 16, 20, 0),
            child: SegmentTabBar(
              tabs: const ['精选', '好友'],
              selectedIndex: _selectedTab,
              onChanged: (index) => setState(() => _selectedTab = index),
            ),
          ),
          const SizedBox(height: 16),
          Expanded(
            child: IndexedStack(
              index: _selectedTab,
              children: [
                MomentFeedList(
                  moments: strangerMoments,
                  emptyMessage: '暂无精选动态',
                ),
                MomentFeedList(
                  moments: friendMoments,
                  emptyMessage: '暂无好友动态',
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
